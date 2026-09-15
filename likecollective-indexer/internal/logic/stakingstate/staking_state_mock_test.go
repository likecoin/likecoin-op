package stakingstate_test

import (
	"context"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/logic/stakingstate/loader"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type stakingStateTestMockLoader struct {
	snapshottedAccounts   []StakingStateTestCaseSnapshottedAccounts
	snapshottedNftClasses []StakingStateTestCaseSnapshottedNftClasses
	snapshottedStakings   []StakingStateTestCaseSnapshottedStakings
}

func MakeStakingStateTestMockLoader(
	snapshottedAccounts []StakingStateTestCaseSnapshottedAccounts,
	snapshottedNftClasses []StakingStateTestCaseSnapshottedNftClasses,
	snapshottedStakings []StakingStateTestCaseSnapshottedStakings,
) *stakingStateTestMockLoader {
	return &stakingStateTestMockLoader{
		snapshottedAccounts,
		snapshottedNftClasses,
		snapshottedStakings,
	}
}

func (l *stakingStateTestMockLoader) Load(
	ctx context.Context,
	loadState *loader.LoadState,
) ([]*model.Account, []*model.NFTClass, []*model.Staking, error) {
	accounts := make([]*model.Account, 0)
	nftClasses := make([]*model.NFTClass, 0)
	stakings := make([]*model.Staking, 0)
	for _, account := range l.snapshottedAccounts {
		a := &model.Account{
			EVMAddress:          common.HexToAddress(account.Address),
			StakedAmount:        uint256.MustFromDecimal(account.StakedAmount),
			PendingRewardAmount: uint256.MustFromDecimal(account.PendingRewardAmount),
			ClaimedRewardAmount: uint256.MustFromDecimal(account.ClaimedRewardAmount),
		}
		accounts = append(accounts, a)
	}
	for _, nftClass := range l.snapshottedNftClasses {
		n := &model.NFTClass{
			EVMAddress:   common.HexToAddress(nftClass.Address),
			StakedAmount: uint256.MustFromDecimal(nftClass.StakedAmount),
		}
		nftClasses = append(nftClasses, n)
	}
	for _, staking := range l.snapshottedStakings {
		s := &model.Staking{
			AccountEVMAddress:   common.HexToAddress(staking.AccountAddress),
			BookNFTEvmAddress:   common.HexToAddress(staking.BookNFTAddress),
			StakedAmount:        uint256.MustFromDecimal(staking.StakedAmount),
			PendingRewardAmount: uint256.MustFromDecimal(staking.PendingRewardAmount),
			ClaimedRewardAmount: uint256.MustFromDecimal(staking.ClaimedRewardAmount),
		}
		stakings = append(stakings, s)
	}
	return accounts, nftClasses, stakings, nil
}

type MockStorage struct {
	Accounts   map[string]StakingStateAccount            `yaml:"accounts"`
	NftClasses map[string]StakingStateBookNFT            `yaml:"nftClasses"`
	Stakings   map[string]map[string]StakingStateStaking `yaml:"stakings"`
}

func MakeMockStorage() *MockStorage {
	return &MockStorage{
		Accounts:   make(map[string]StakingStateAccount),
		NftClasses: make(map[string]StakingStateBookNFT),
		Stakings:   make(map[string]map[string]StakingStateStaking),
	}
}

type stakingStateTestMockPersistor struct {
	mockStorage *MockStorage
}

func MakeStakingStateTestMockPersistor() *stakingStateTestMockPersistor {
	mockStorage := MakeMockStorage()
	return &stakingStateTestMockPersistor{
		mockStorage,
	}
}

func (p *stakingStateTestMockPersistor) AlreadyApplied(
	ctx context.Context,
	transactionHash string,
	transactionIndex uint,
	logIndex uint,
) (bool, error) {
	return false, nil
}

func (p *stakingStateTestMockPersistor) WithLock(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func (p *stakingStateTestMockPersistor) Persist(
	ctx context.Context,
	headBlockNumber *big.Int,
	stakingEvents []*ent.StakingEvent,
	accounts []*model.Account,
	nftClasses []*model.NFTClass,
	stakings []*model.Staking,
) error {
	stakingStateAccounts := make(map[string]StakingStateAccount)
	for _, account := range accounts {
		stakingStateAccounts[account.EVMAddress.String()] = StakingStateAccount{
			EVMAddress:          account.EVMAddress.String(),
			StakedAmount:        account.StakedAmount.String(),
			PendingRewardAmount: account.PendingRewardAmount.String(),
			ClaimedRewardAmount: account.ClaimedRewardAmount.String(),
		}
	}
	p.mockStorage.Accounts = stakingStateAccounts
	stakingStateNftClasses := make(map[string]StakingStateBookNFT)
	for _, nftClass := range nftClasses {
		stakingStateNftClasses[nftClass.EVMAddress.String()] = StakingStateBookNFT{
			EVMAddress:   nftClass.EVMAddress.String(),
			StakedAmount: nftClass.StakedAmount.String(),
		}
	}
	p.mockStorage.NftClasses = stakingStateNftClasses
	stakingStateStakings := make(map[string]map[string]StakingStateStaking)
	for _, staking := range stakings {
		if _, ok := stakingStateStakings[staking.AccountEVMAddress.String()]; !ok {
			stakingStateStakings[staking.AccountEVMAddress.String()] = make(map[string]StakingStateStaking)
		}
		stakingStateStakings[staking.AccountEVMAddress.String()][staking.BookNFTEvmAddress.String()] = StakingStateStaking{
			AccountEVMAddress:   staking.AccountEVMAddress.String(),
			BookNFTEvmAddress:   staking.BookNFTEvmAddress.String(),
			StakedAmount:        staking.StakedAmount.String(),
			PendingRewardAmount: staking.PendingRewardAmount.String(),
			ClaimedRewardAmount: staking.ClaimedRewardAmount.String(),
		}
	}
	p.mockStorage.Stakings = stakingStateStakings
	return nil
}
