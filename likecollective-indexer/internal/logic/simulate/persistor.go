package simulate

import (
	"context"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/logic/stakingstate/persistor"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type echo func(
	state *State,
) error

type simulationPersistor struct {
	echo echo
}

func MakeSimulationPersistor(echo echo) persistor.StakingStatePersistor {
	return &simulationPersistor{
		echo,
	}
}

// AlreadyApplied is always false: a simulation persists nothing.
func (p *simulationPersistor) AlreadyApplied(
	ctx context.Context,
	transactionHash string,
	transactionIndex uint,
	logIndex uint,
) (bool, error) {
	return false, nil
}

// WithLock runs fn directly: a simulation has no other writer.
func (p *simulationPersistor) WithLock(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func (p *simulationPersistor) Persist(
	ctx context.Context,
	headBlockNumber *big.Int,
	stakingEvents []*ent.StakingEvent,
	accounts []*model.Account,
	nftClasses []*model.NFTClass,
	stakings []*model.Staking,
) error {
	state := &State{
		BookPendingRewards: make(map[common.Address]*uint256.Int),
		BookStakedAmounts:  make(map[common.Address]*uint256.Int),
		UserPendingRewards: make(map[common.Address]map[common.Address]*uint256.Int),
		UserStakedAmounts:  make(map[common.Address]map[common.Address]*uint256.Int),
	}

	for _, nftClass := range nftClasses {
		state.BookStakedAmounts[nftClass.EVMAddress] = nftClass.StakedAmount
	}

	for _, staking := range stakings {
		if _, ok := state.BookPendingRewards[staking.BookNFTEvmAddress]; !ok {
			state.BookPendingRewards[staking.BookNFTEvmAddress] = uint256.NewInt(0)
		}
		state.BookPendingRewards[staking.BookNFTEvmAddress] = uint256.NewInt(0).Add(state.BookPendingRewards[staking.BookNFTEvmAddress], staking.PendingRewardAmount)

		if _, ok := state.UserPendingRewards[staking.AccountEVMAddress]; !ok {
			state.UserPendingRewards[staking.AccountEVMAddress] = make(map[common.Address]*uint256.Int)
		}
		state.UserPendingRewards[staking.AccountEVMAddress][staking.BookNFTEvmAddress] = staking.PendingRewardAmount

		if _, ok := state.UserStakedAmounts[staking.AccountEVMAddress]; !ok {
			state.UserStakedAmounts[staking.AccountEVMAddress] = make(map[common.Address]*uint256.Int)
		}
		state.UserStakedAmounts[staking.AccountEVMAddress][staking.BookNFTEvmAddress] = staking.StakedAmount

	}

	return p.echo(state)
}
