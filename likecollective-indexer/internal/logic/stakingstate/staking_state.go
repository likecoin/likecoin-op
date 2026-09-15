package stakingstate

import (
	"context"
	"fmt"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/logic/stakingstate/loader"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/logic/stakingstate/persistor"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type StakingState interface {
	Process(stakingEvents []*ent.StakingEvent) (*stakingState, []*ent.StakingEvent, error)
	Persist(ctx context.Context, headBlockNumber *big.Int, stakingEvents []*ent.StakingEvent, persistor persistor.StakingStatePersistor) error
}

type stakingState struct {
	accounts   []*model.Account
	nftClasses []*model.NFTClass
	stakings   []*model.Staking

	// chainBacked marks a state whose staked and pending amounts are replaced
	// from the chain once the events are applied. The applications then write
	// only the history and leave those amounts as loaded: see add, lacks and
	// subtract.
	chainBacked bool

	// touched holds the stakings an event named, which the chain is re-read
	// for even when they hold nothing as loaded.
	touched map[*model.Staking]struct{}
}

func LoadStakingState(
	ctx context.Context,
	stakingStateLoader loader.StakingStateLoader,
	stakingEvents []*ent.StakingEvent,
) (StakingState, error) {
	return loadStakingState(ctx, stakingStateLoader, stakingEvents, false)
}

func loadStakingState(
	ctx context.Context,
	stakingStateLoader loader.StakingStateLoader,
	stakingEvents []*ent.StakingEvent,
	chainBacked bool,
) (*stakingState, error) {
	loadStates := make([]*loader.LoadState, 0)
	for _, stakingEvent := range stakingEvents {
		loadStateFactory, err := loader.MakeLoadStateFactory(stakingEvent)
		if err != nil {
			return nil, err
		}
		loadStates = append(loadStates, loadStateFactory.MakeLoadState())
	}

	accounts, nftClasses, stakings, err := stakingStateLoader.Load(ctx, loader.CombineLoadStates(loadStates))
	if err != nil {
		return nil, err
	}

	return &stakingState{
		accounts:    accounts,
		nftClasses:  nftClasses,
		stakings:    stakings,
		chainBacked: chainBacked,
	}, nil
}

func (s *stakingState) Process(stakingEvents []*ent.StakingEvent) (*stakingState, []*ent.StakingEvent, error) {
	applications := make([]StakingEventApplication, 0)
	for _, stakingEvent := range stakingEvents {
		application, err := MakeStakingEventApplication(stakingEvent)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to make staking event application: %w", err)
		}
		applications = append(applications, application)
	}

	return s.run(applications)
}

func (s *stakingState) run(applications []StakingEventApplication) (*stakingState, []*ent.StakingEvent, error) {
	processedStakingEvents := make([]*ent.StakingEvent, 0)
	newState := s
	for _, application := range applications {
		var (
			furtherApplications []StakingEventApplication
			err                 error
		)
		newState, furtherApplications, err = application.Apply(newState)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to apply staking event application: %w", err)
		}
		processedStakingEvents = append(
			processedStakingEvents, application.GetStakingEvent(),
		)

		if len(furtherApplications) > 0 {
			var (
				furtherProcessedStakingEvents []*ent.StakingEvent
			)
			newState, furtherProcessedStakingEvents, err = s.run(furtherApplications)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to run further applications: %w", err)
			}
			processedStakingEvents = append(processedStakingEvents, furtherProcessedStakingEvents...)
		}
	}
	return newState, processedStakingEvents, nil
}

func (s *stakingState) Persist(
	ctx context.Context,
	headBlockNumber *big.Int,
	stakingEvents []*ent.StakingEvent,
	persistor persistor.StakingStatePersistor,
) error {
	return persistor.Persist(
		ctx,
		headBlockNumber,
		stakingEvents,
		s.accounts,
		s.nftClasses,
		s.stakings,
	)
}

// add returns total + amount, or total as it is on chain-backed state.
//
// reconcileFromChain moves an account by how far each of its stakings moved
// from what was loaded, so on chain-backed state nothing may move them first.
func (s *stakingState) add(total *uint256.Int, amount *uint256.Int) *uint256.Int {
	if s.chainBacked {
		return total
	}
	return new(uint256.Int).Add(total, amount)
}

// lacks reports whether total is too small for amount to be removed from it.
//
// Never on chain-backed state. The loaded amount there may already have been
// re-read at a head past this event -- an earlier event in the queue read it
// after this one's removal landed on chain -- so an amount smaller than the
// removal is expected rather than a sign of corruption. Refusing it would fail
// the event on every retry, and nothing after it could correct that.
func (s *stakingState) lacks(total *uint256.Int, amount *uint256.Int) bool {
	if s.chainBacked {
		return false
	}
	return total.Lt(amount)
}

// subtract returns total - amount, or total as it is on chain-backed state. On
// delta state the caller must have checked lacks first.
func (s *stakingState) subtract(total *uint256.Int, amount *uint256.Int) *uint256.Int {
	if s.chainBacked {
		return total
	}
	return new(uint256.Int).Sub(total, amount)
}

func (s *stakingState) touch(staking *model.Staking) {
	if s.touched == nil {
		s.touched = make(map[*model.Staking]struct{})
	}
	s.touched[staking] = struct{}{}
}

func (s *stakingState) isTouched(staking *model.Staking) bool {
	_, ok := s.touched[staking]
	return ok
}

func (s *stakingState) GetAccountByAddress(evmAddress common.Address) (*model.Account, bool) {
	for _, account := range s.accounts {
		if account.EVMAddress.String() == evmAddress.String() {
			return account, true
		}
	}
	return nil, false
}

func (s *stakingState) GetNFTClassByAddress(evmAddress common.Address) (*model.NFTClass, bool) {
	for _, nftClass := range s.nftClasses {
		if nftClass.EVMAddress.String() == evmAddress.String() {
			return nftClass, true
		}
	}
	return nil, false
}

func (s *stakingState) GetStakingByAddress(accountEVMAddress common.Address, bookNFTEVMAddress common.Address) (*model.Staking, bool) {
	for _, staking := range s.stakings {
		if staking.AccountEVMAddress.String() == accountEVMAddress.String() && staking.BookNFTEvmAddress.String() == bookNFTEVMAddress.String() {
			return staking, true
		}
	}
	return nil, false
}

func (s *stakingState) GetOrCreateAccount(evmAddress common.Address) *model.Account {
	for _, account := range s.accounts {
		if account.EVMAddress.String() == evmAddress.String() {
			return account
		}
	}
	account := model.NewAccount(evmAddress.String())
	s.accounts = append(s.accounts, account)
	return account
}

func (s *stakingState) GetOrCreateNFTClass(evmAddress common.Address) *model.NFTClass {
	for _, nftClass := range s.nftClasses {
		if nftClass.EVMAddress.String() == evmAddress.String() {
			return nftClass
		}
	}
	nftClass := model.NewNFTClass(evmAddress.String())
	s.nftClasses = append(s.nftClasses, nftClass)
	return nftClass
}

func (s *stakingState) GetOrCreateStaking(
	accountEVMAddress common.Address,
	bookNFTEVMAddress common.Address,
) *model.Staking {
	for _, staking := range s.stakings {
		if staking.AccountEVMAddress.String() == accountEVMAddress.String() && staking.BookNFTEvmAddress.String() == bookNFTEVMAddress.String() {
			return staking
		}
	}
	staking := model.NewStakingFromStakingKey(database.NewStakingKey(accountEVMAddress.String(), bookNFTEVMAddress.String()))
	s.stakings = append(s.stakings, staking)
	return staking
}

func (s *stakingState) GetStakingsByNFTClassAddress(nftClassAddress common.Address) []*model.Staking {
	stakings := make([]*model.Staking, 0)
	for _, staking := range s.stakings {
		if staking.BookNFTEvmAddress.String() == nftClassAddress.String() {
			stakings = append(stakings, staking)
		}
	}
	return stakings
}
