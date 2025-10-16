package stakingstate

import (
	"errors"
	"fmt"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/shopspring/decimal"
)

type StakingEventApplication interface {
	GetStakingEvent() *ent.StakingEvent
	Apply(state *stakingState) (*stakingState, []StakingEventApplication, error)
}

func MakeStakingEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	switch stakingEvent.EventType {
	case stakingevent.EventTypeStaked:
		return makeStakedEventApplication(stakingEvent)
	case stakingevent.EventTypeUnstaked:
		return makeUnstakedEventApplication(stakingEvent)
	case stakingevent.EventTypeRewardClaimed:
		return makeRewardClaimedEventApplication(stakingEvent)
	case stakingevent.EventTypeRewardDeposited:
		return makeRewardDepositedEventApplication(stakingEvent)
	case stakingevent.EventTypeRewardDepositDistributed:
		return makeRewardDepositDistributedEventApplication(stakingEvent)
	case stakingevent.EventTypeAllRewardsClaimed:
		return makeAllRewardsClaimedEventApplication(stakingEvent)
	case stakingevent.EventTypeStakePositionTransferred:
		return makeStakePositionTransferredEventApplication(stakingEvent)
	case stakingevent.EventTypeStakePositionReceived:
		return makeStakePositionReceivedEventApplication(stakingEvent)
	default:
		return nil, errors.New("invalid staking event type")
	}
}

var ErrStakedEventApplication = errors.New("err staked event application")

type stakedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	stakedAmountAdded *uint256.Int
}

func makeStakedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	stakedAmountAdded := (*uint256.Int)(stakingEvent.StakedAmountAdded)
	if stakedAmountAdded == nil {
		return nil, errors.Join(
			ErrStakedEventApplication,
			errors.New("failed to convert staked amount to uint256"),
		)
	}

	return &stakedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		stakedAmountAdded,
	}, nil
}

func (s *stakedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *stakedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account := state.GetOrCreateAccount(s.accountEVMAddress)
	nftClass := state.GetOrCreateNFTClass(s.nftClassAddress)
	staking := state.GetOrCreateStaking(s.accountEVMAddress, s.nftClassAddress)

	account.StakedAmount = uint256.NewInt(0).Add(account.StakedAmount, s.stakedAmountAdded)
	nftClass.StakedAmount = uint256.NewInt(0).Add(nftClass.StakedAmount, s.stakedAmountAdded)
	staking.StakedAmount = uint256.NewInt(0).Add(staking.StakedAmount, s.stakedAmountAdded)

	return state, []StakingEventApplication{}, nil
}

var ErrUnstakedEventApplication = errors.New("err unstaked event application")

type unstakedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress   common.Address
	nftClassAddress     common.Address
	stakedAmountRemoved *uint256.Int
}

func makeUnstakedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	stakedAmountRemoved := (*uint256.Int)(stakingEvent.StakedAmountRemoved)
	if stakedAmountRemoved == nil {
		return nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("failed to convert unstaked amount to uint256"),
		)
	}

	return &unstakedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		stakedAmountRemoved,
	}, nil
}

func (s *unstakedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *unstakedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account, ok := state.GetAccountByAddress(s.accountEVMAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("account not found"),
		)
	}
	nftClass, ok := state.GetNFTClassByAddress(s.nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("nft class not found"),
		)
	}
	staking, ok := state.GetStakingByAddress(s.accountEVMAddress, s.nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("staking not found"),
		)
	}

	if account.StakedAmount.Cmp(s.stakedAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("staked amount is less than staked amount removed"),
		)
	}
	if nftClass.StakedAmount.Cmp(s.stakedAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("nft class staked amount is less than staked amount removed"),
		)
	}
	if staking.StakedAmount.Cmp(s.stakedAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrUnstakedEventApplication,
			errors.New("staking staked amount is less than staked amount removed"),
		)
	}

	account.StakedAmount = uint256.NewInt(0).Sub(account.StakedAmount, s.stakedAmountRemoved)
	nftClass.StakedAmount = uint256.NewInt(0).Sub(nftClass.StakedAmount, s.stakedAmountRemoved)
	staking.StakedAmount = uint256.NewInt(0).Sub(staking.StakedAmount, s.stakedAmountRemoved)

	return state, []StakingEventApplication{}, nil
}

var ErrRewardClaimedEventApplication = errors.New("err reward claimed event application")

type rewardClaimedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	pendingRewardAmountRemoved *uint256.Int
	claimedRewardAmountAdded   *uint256.Int
}

func makeRewardClaimedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	pendingRewardAmountRemoved := (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved)
	if pendingRewardAmountRemoved == nil {
		return nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("failed to convert pending reward amount removed to uint256"),
		)
	}
	claimedRewardAmountAdded := (*uint256.Int)(stakingEvent.ClaimedRewardAmountAdded)
	if claimedRewardAmountAdded == nil {
		return nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("failed to convert reward claimed amount to uint256"),
		)
	}

	return &rewardClaimedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		pendingRewardAmountRemoved,
		claimedRewardAmountAdded,
	}, nil
}

func (s *rewardClaimedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *rewardClaimedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account, ok := state.GetAccountByAddress(s.accountEVMAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("account not found"),
		)
	}
	staking, ok := state.GetStakingByAddress(s.accountEVMAddress, s.nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("staking not found"),
		)
	}

	if account.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("pending reward amount is less than pending reward amount removed"),
		)
	}
	if staking.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrRewardClaimedEventApplication,
			errors.New("staking pending reward amount is less than pending reward amount removed"),
		)
	}

	account.PendingRewardAmount = uint256.NewInt(0).Sub(account.PendingRewardAmount, s.pendingRewardAmountRemoved)
	account.ClaimedRewardAmount = uint256.NewInt(0).Add(account.ClaimedRewardAmount, s.claimedRewardAmountAdded)
	staking.PendingRewardAmount = uint256.NewInt(0).Sub(staking.PendingRewardAmount, s.pendingRewardAmountRemoved)
	staking.ClaimedRewardAmount = uint256.NewInt(0).Add(staking.ClaimedRewardAmount, s.claimedRewardAmountAdded)

	return state, []StakingEventApplication{}, nil
}

var ErrRewardDepositedEventApplication = errors.New("err reward deposited event application")

type rewardDepositedEventApplication struct {
	stakingEvent *ent.StakingEvent
}

func makeRewardDepositedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	return &rewardDepositedEventApplication{
		stakingEvent,
	}, nil
}

func (s *rewardDepositedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *rewardDepositedEventApplication) Apply(
	state *stakingState,
) (
	*stakingState, []StakingEventApplication, error,
) {
	nftClassAddress := common.HexToAddress(s.stakingEvent.NftClassAddress)
	rewardAmount := (*uint256.Int)(s.stakingEvent.PendingRewardAmountAdded)
	if rewardAmount == nil {
		return nil, nil, errors.Join(
			ErrRewardDepositedEventApplication,
			errors.New("failed to convert reward amount to uint256"),
		)
	}

	nftClass, ok := state.GetNFTClassByAddress(nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrRewardDepositedEventApplication,
			errors.New("nft class not found"),
		)
	}

	stakings := state.GetStakingsByNFTClassAddress(nftClassAddress)
	distributedApplications := make([]StakingEventApplication, 0)

	nonZeroStakings := make([]*model.Staking, 0)
	for _, staking := range stakings {
		if staking.StakedAmount.IsZero() {
			continue
		}
		nonZeroStakings = append(nonZeroStakings, staking)
	}

	for _, staking := range nonZeroStakings {
		var poolShares *big.Rat
		if (*uint256.Int)(nftClass.StakedAmount).IsZero() {
			poolShares = big.NewRat(0, 1)
		} else {
			poolShares = big.NewRat(staking.StakedAmount.ToBig().Int64(), nftClass.StakedAmount.ToBig().Int64())
		}
		pendingRewardAmountRat := big.NewRat(0, 1).
			Mul(poolShares,
				big.NewRat(
					rewardAmount.ToBig().Int64(),
					big.NewInt(1).Int64(),
				))
		pendingRewardAmount, err := uint256.FromDecimal(
			decimal.NewFromBigRat(pendingRewardAmountRat, 18).
				Floor().
				String(),
		)
		if err != nil {
			return nil, nil, errors.Join(
				ErrRewardDepositedEventApplication,
				fmt.Errorf("failed to convert pending reward amount to uint256: %w", err),
			)
		}
		distributedApplication, err := makeRewardDepositDistributedEventApplication(
			MakeRewardDepositDistributedEvent(
				s.stakingEvent,
				staking.AccountEVMAddress.String(),
				pendingRewardAmount,
			),
		)
		if err != nil {
			return nil, nil, errors.Join(
				ErrRewardDepositedEventApplication,
				fmt.Errorf("failed to make reward deposit distributed event application: %w", err),
			)
		}
		distributedApplications = append(
			distributedApplications,
			distributedApplication,
		)
	}

	return state, distributedApplications, nil
}

type rewardDepositDistributedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	pendingRewardAmountAdded *uint256.Int
}

func makeRewardDepositDistributedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	pendingRewardAmountAdded := (*uint256.Int)(stakingEvent.PendingRewardAmountAdded)
	if pendingRewardAmountAdded == nil {
		return nil, errors.New("failed to convert pending reward amount added to uint256")
	}

	return &rewardDepositDistributedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		pendingRewardAmountAdded,
	}, nil
}

func (s *rewardDepositDistributedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *rewardDepositDistributedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account := state.GetOrCreateAccount(s.accountEVMAddress)
	staking := state.GetOrCreateStaking(s.accountEVMAddress, s.nftClassAddress)

	account.PendingRewardAmount = uint256.NewInt(0).Add(account.PendingRewardAmount, s.pendingRewardAmountAdded)
	staking.PendingRewardAmount = uint256.NewInt(0).Add(staking.PendingRewardAmount, s.pendingRewardAmountAdded)

	return state, []StakingEventApplication{}, nil
}

var ErrAllRewardsClaimedEventApplication = errors.New("err all rewards claimed event application")

type allRewardsClaimedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	pendingRewardAmountRemoved *uint256.Int
	claimedRewardAmountAdded   *uint256.Int
}

func makeAllRewardsClaimedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	pendingRewardAmountRemoved := (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved)
	if pendingRewardAmountRemoved == nil {
		return nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("failed to convert pending reward amount removed to uint256"),
		)
	}
	claimedRewardAmountAdded := (*uint256.Int)(stakingEvent.ClaimedRewardAmountAdded)
	if claimedRewardAmountAdded == nil {
		return nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("failed to convert claimed reward amount added to uint256"),
		)
	}

	return &allRewardsClaimedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		pendingRewardAmountRemoved,
		claimedRewardAmountAdded,
	}, nil
}

func (s *allRewardsClaimedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *allRewardsClaimedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account, ok := state.GetAccountByAddress(s.accountEVMAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("account not found"),
		)
	}
	staking, ok := state.GetStakingByAddress(s.accountEVMAddress, s.nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("staking not found"),
		)
	}

	if account.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("pending reward amount is less than pending reward amount removed"),
		)
	}
	if staking.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrAllRewardsClaimedEventApplication,
			errors.New("staking pending reward amount is less than pending reward amount removed"),
		)
	}

	account.PendingRewardAmount = uint256.NewInt(0).Sub(account.PendingRewardAmount, s.pendingRewardAmountRemoved)
	account.ClaimedRewardAmount = uint256.NewInt(0).Add(account.ClaimedRewardAmount, s.claimedRewardAmountAdded)
	staking.PendingRewardAmount = uint256.NewInt(0).Sub(staking.PendingRewardAmount, s.pendingRewardAmountRemoved)
	staking.ClaimedRewardAmount = uint256.NewInt(0).Add(staking.ClaimedRewardAmount, s.claimedRewardAmountAdded)

	return state, []StakingEventApplication{}, nil
}

var ErrStakePositionTransferredEventApplication = errors.New("err stake position transferred event application")

type stakePositionTransferredEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	stakedAmountRemoved        *uint256.Int
	pendingRewardAmountRemoved *uint256.Int
}

func makeStakePositionTransferredEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	stakedAmountRemoved := (*uint256.Int)(stakingEvent.StakedAmountRemoved)
	if stakedAmountRemoved == nil {
		return nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("failed to convert staked amount removed to uint256"),
		)
	}
	pendingRewardAmountRemoved := (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved)
	if pendingRewardAmountRemoved == nil {
		return nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("failed to convert pending reward amount removed to uint256"),
		)
	}

	return &stakePositionTransferredEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		stakedAmountRemoved,
		pendingRewardAmountRemoved,
	}, nil
}

func (s *stakePositionTransferredEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *stakePositionTransferredEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account, ok := state.GetAccountByAddress(s.accountEVMAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("account not found"),
		)
	}
	staking, ok := state.GetStakingByAddress(s.accountEVMAddress, s.nftClassAddress)
	if !ok {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("staking not found"),
		)
	}

	if account.StakedAmount.Cmp(s.stakedAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("stake position transferred staked amount is less than staked amount removed"),
		)
	}
	if account.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("stake position transferred pending reward amount is less than pending reward amount removed"),
		)
	}
	if staking.StakedAmount.Cmp(s.stakedAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("stake position transferred staking staked amount is less than staked amount removed"),
		)
	}
	if staking.PendingRewardAmount.Cmp(s.pendingRewardAmountRemoved) < 0 {
		return nil, nil, errors.Join(
			ErrStakePositionTransferredEventApplication,
			errors.New("stake position transferred staking pending reward amount is less than pending reward amount removed"),
		)
	}

	account.StakedAmount = uint256.NewInt(0).Sub(account.StakedAmount, s.stakedAmountRemoved)
	account.PendingRewardAmount = uint256.NewInt(0).Sub(account.PendingRewardAmount, s.pendingRewardAmountRemoved)
	staking.StakedAmount = uint256.NewInt(0).Sub(staking.StakedAmount, s.stakedAmountRemoved)
	staking.PendingRewardAmount = uint256.NewInt(0).Sub(staking.PendingRewardAmount, s.pendingRewardAmountRemoved)

	return state, []StakingEventApplication{}, nil
}

var ErrStakePositionReceivedEventApplication = errors.New("err stake position received event application")

type stakePositionReceivedEventApplication struct {
	stakingEvent *ent.StakingEvent

	accountEVMAddress common.Address
	nftClassAddress   common.Address

	stakedAmountAdded        *uint256.Int
	pendingRewardAmountAdded *uint256.Int
}

func makeStakePositionReceivedEventApplication(stakingEvent *ent.StakingEvent) (StakingEventApplication, error) {
	accountEVMAddress := common.HexToAddress(stakingEvent.AccountEvmAddress)
	nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)
	stakedAmountAdded := (*uint256.Int)(stakingEvent.StakedAmountAdded)
	if stakedAmountAdded == nil {
		return nil, errors.Join(
			ErrStakePositionReceivedEventApplication,
			errors.New("failed to convert staked amount added to uint256"),
		)
	}
	pendingRewardAmountAdded := (*uint256.Int)(stakingEvent.PendingRewardAmountAdded)
	if pendingRewardAmountAdded == nil {
		return nil, errors.Join(
			ErrStakePositionReceivedEventApplication,
			errors.New("failed to convert pending reward amount added to uint256"),
		)
	}

	return &stakePositionReceivedEventApplication{
		stakingEvent,
		accountEVMAddress,
		nftClassAddress,
		stakedAmountAdded,
		pendingRewardAmountAdded,
	}, nil
}

func (s *stakePositionReceivedEventApplication) GetStakingEvent() *ent.StakingEvent {
	return s.stakingEvent
}

func (s *stakePositionReceivedEventApplication) Apply(state *stakingState) (*stakingState, []StakingEventApplication, error) {
	account := state.GetOrCreateAccount(s.accountEVMAddress)
	staking := state.GetOrCreateStaking(s.accountEVMAddress, s.nftClassAddress)

	account.StakedAmount = uint256.NewInt(0).Add(account.StakedAmount, s.stakedAmountAdded)
	account.PendingRewardAmount = uint256.NewInt(0).Add(account.PendingRewardAmount, s.pendingRewardAmountAdded)
	staking.StakedAmount = uint256.NewInt(0).Add(staking.StakedAmount, s.stakedAmountAdded)
	staking.PendingRewardAmount = uint256.NewInt(0).Add(staking.PendingRewardAmount, s.pendingRewardAmountAdded)

	return state, []StakingEventApplication{}, nil
}
