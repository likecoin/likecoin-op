package stakingstate_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/logic/stakingstate"

	"github.com/holiman/uint256"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	goyaml "gopkg.in/yaml.v2"
)

type StakingStateTestCaseEventType string

const (
	StakingStateTestCaseEventTypeStaked                   StakingStateTestCaseEventType = "staked"
	StakingStateTestCaseEventTypeUnstaked                 StakingStateTestCaseEventType = "unstaked"
	StakingStateTestCaseEventTypeRewardDeposited          StakingStateTestCaseEventType = "rewardDeposited"
	StakingStateTestCaseEventTypeRewardDepositDistributed StakingStateTestCaseEventType = "rewardDepositDistributed"
	StakingStateTestCaseEventTypeRewardClaimed            StakingStateTestCaseEventType = "rewardClaimed"
	StakingStateTestCaseEventTypeAllRewardsClaimed        StakingStateTestCaseEventType = "allRewardsClaimed"
	StakingStateTestCaseEventTypeStakePositionTransferred StakingStateTestCaseEventType = "stakePositionTransferred"
	StakingStateTestCaseEventTypeStakePositionReceived    StakingStateTestCaseEventType = "stakePositionReceived"
)

type StakingStateTestCaseEventStaked struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	StakedAmount string `yaml:"stakedAmount"`
}

func (e *StakingStateTestCaseEventStaked) ToEnt() *ent.StakingEvent {
	stakedAmount, err := uint256.FromDecimal(e.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:         stakingevent.EventTypeStaked,
		AccountEvmAddress: e.Account,
		NftClassAddress:   e.BookNFT,
		StakedAmountAdded: typeutil.Uint256(stakedAmount),
	}
}

type StakingStateTestCaseEventUnstaked struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	StakedAmount string `yaml:"stakedAmount"`
}

func (e *StakingStateTestCaseEventUnstaked) ToEnt() *ent.StakingEvent {
	stakedAmount, err := uint256.FromDecimal(e.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:           stakingevent.EventTypeUnstaked,
		AccountEvmAddress:   e.Account,
		NftClassAddress:     e.BookNFT,
		StakedAmountRemoved: typeutil.Uint256(stakedAmount),
	}
}

type StakingStateTestCaseEventRewardDeposited struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	RewardAmount string `yaml:"rewardAmount"`
}

func (e *StakingStateTestCaseEventRewardDeposited) ToEnt() *ent.StakingEvent {
	rewardAmount, err := uint256.FromDecimal(e.RewardAmount)
	if err != nil {
		panic("failed to convert reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                stakingevent.EventTypeRewardDeposited,
		AccountEvmAddress:        e.Account,
		NftClassAddress:          e.BookNFT,
		PendingRewardAmountAdded: typeutil.Uint256(rewardAmount),
	}
}

type StakingStateTestCaseEventRewardDepositDistributed struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	RewardAmount string `yaml:"rewardAmount"`
}

func (e *StakingStateTestCaseEventRewardDepositDistributed) ToEnt() *ent.StakingEvent {
	rewardAmount, err := uint256.FromDecimal(e.RewardAmount)
	if err != nil {
		panic("failed to convert reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                stakingevent.EventTypeRewardDepositDistributed,
		AccountEvmAddress:        e.Account,
		NftClassAddress:          e.BookNFT,
		PendingRewardAmountAdded: typeutil.Uint256(rewardAmount),
	}
}

type StakingStateTestCaseEventRewardClaimed struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	RewardAmount string `yaml:"rewardAmount"`
}

func (e *StakingStateTestCaseEventRewardClaimed) ToEnt() *ent.StakingEvent {
	rewardAmount, err := uint256.FromDecimal(e.RewardAmount)
	if err != nil {
		panic("failed to convert reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                  stakingevent.EventTypeRewardClaimed,
		AccountEvmAddress:          e.Account,
		NftClassAddress:            e.BookNFT,
		PendingRewardAmountRemoved: typeutil.Uint256(rewardAmount),
		ClaimedRewardAmountAdded:   typeutil.Uint256(rewardAmount),
	}
}

type StakingStateTestCaseEventAllRewardsClaimed struct {
	BookNFT      string `yaml:"bookNFT"`
	Account      string `yaml:"account"`
	RewardAmount string `yaml:"rewardAmount"`
}

func (e *StakingStateTestCaseEventAllRewardsClaimed) ToEnt() *ent.StakingEvent {
	rewardAmount, err := uint256.FromDecimal(e.RewardAmount)
	if err != nil {
		panic("failed to convert reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                  stakingevent.EventTypeAllRewardsClaimed,
		AccountEvmAddress:          e.Account,
		NftClassAddress:            e.BookNFT,
		PendingRewardAmountRemoved: typeutil.Uint256(rewardAmount),
		ClaimedRewardAmountAdded:   typeutil.Uint256(rewardAmount),
	}
}

type StakingStateTestCaseEventStakePositionTransferred struct {
	BookNFT             string `yaml:"bookNFT"`
	Account             string `yaml:"account"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
}

func (e *StakingStateTestCaseEventStakePositionTransferred) ToEnt() *ent.StakingEvent {
	stakedAmount, err := uint256.FromDecimal(e.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	pendingRewardAmount, err := uint256.FromDecimal(e.PendingRewardAmount)
	if err != nil {
		panic("failed to convert pending reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                  stakingevent.EventTypeStakePositionTransferred,
		AccountEvmAddress:          e.Account,
		NftClassAddress:            e.BookNFT,
		StakedAmountRemoved:        typeutil.Uint256(stakedAmount),
		PendingRewardAmountRemoved: typeutil.Uint256(pendingRewardAmount),
	}
}

type StakingStateTestCaseEventStakePositionReceived struct {
	BookNFT             string `yaml:"bookNFT"`
	Account             string `yaml:"account"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
}

func (e *StakingStateTestCaseEventStakePositionReceived) ToEnt() *ent.StakingEvent {
	stakedAmount, err := uint256.FromDecimal(e.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	pendingRewardAmount, err := uint256.FromDecimal(e.PendingRewardAmount)
	if err != nil {
		panic("failed to convert pending reward amount to big.Int")
	}
	return &ent.StakingEvent{
		EventType:                stakingevent.EventTypeStakePositionReceived,
		AccountEvmAddress:        e.Account,
		NftClassAddress:          e.BookNFT,
		StakedAmountAdded:        typeutil.Uint256(stakedAmount),
		PendingRewardAmountAdded: typeutil.Uint256(pendingRewardAmount),
	}
}

type StakingStateTestCaseEvent struct {
	Type                              StakingStateTestCaseEventType                      `yaml:"type"`
	StakedEventData                   *StakingStateTestCaseEventStaked                   `yaml:"stakedEventData"`
	UnstakedEventData                 *StakingStateTestCaseEventUnstaked                 `yaml:"unstakedEventData"`
	RewardDepositedEventData          *StakingStateTestCaseEventRewardDeposited          `yaml:"rewardDepositedEventData"`
	RewardDepositDistributedEventData *StakingStateTestCaseEventRewardDepositDistributed `yaml:"rewardDepositDistributedEventData"`
	RewardClaimedEventData            *StakingStateTestCaseEventRewardClaimed            `yaml:"rewardClaimedEventData"`
	AllRewardsClaimedEventData        *StakingStateTestCaseEventAllRewardsClaimed        `yaml:"allRewardsClaimedEventData"`
	StakePositionTransferredEventData *StakingStateTestCaseEventStakePositionTransferred `yaml:"stakePositionTransferredEventData"`
	StakePositionReceivedEventData    *StakingStateTestCaseEventStakePositionReceived    `yaml:"stakePositionReceivedEventData"`
}

func (e *StakingStateTestCaseEvent) ToStakingEvent() *ent.StakingEvent {
	switch e.Type {
	case StakingStateTestCaseEventTypeStaked:
		return e.StakedEventData.ToEnt()
	case StakingStateTestCaseEventTypeUnstaked:
		return e.UnstakedEventData.ToEnt()
	case StakingStateTestCaseEventTypeRewardDeposited:
		return e.RewardDepositedEventData.ToEnt()
	case StakingStateTestCaseEventTypeRewardDepositDistributed:
		return e.RewardDepositDistributedEventData.ToEnt()
	case StakingStateTestCaseEventTypeRewardClaimed:
		return e.RewardClaimedEventData.ToEnt()
	case StakingStateTestCaseEventTypeAllRewardsClaimed:
		return e.AllRewardsClaimedEventData.ToEnt()
	case StakingStateTestCaseEventTypeStakePositionTransferred:
		return e.StakePositionTransferredEventData.ToEnt()
	case StakingStateTestCaseEventTypeStakePositionReceived:
		return e.StakePositionReceivedEventData.ToEnt()
	}
	panic("unknown event type")
}

func FromEnt(stakingEvent *ent.StakingEvent) *StakingStateTestCaseEvent {
	switch stakingEvent.EventType {
	case stakingevent.EventTypeStaked:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeStaked,
			StakedEventData: &StakingStateTestCaseEventStaked{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				StakedAmount: (*uint256.Int)(stakingEvent.StakedAmountAdded).String(),
			},
		}
	case stakingevent.EventTypeUnstaked:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeUnstaked,
			UnstakedEventData: &StakingStateTestCaseEventUnstaked{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				StakedAmount: (*uint256.Int)(stakingEvent.StakedAmountRemoved).String(),
			},
		}
	case stakingevent.EventTypeRewardDeposited:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeRewardDeposited,
			RewardDepositedEventData: &StakingStateTestCaseEventRewardDeposited{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				RewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountAdded).String(),
			},
		}
	case stakingevent.EventTypeRewardDepositDistributed:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeRewardDepositDistributed,
			RewardDepositDistributedEventData: &StakingStateTestCaseEventRewardDepositDistributed{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				RewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountAdded).String(),
			},
		}
	case stakingevent.EventTypeRewardClaimed:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeRewardClaimed,
			RewardClaimedEventData: &StakingStateTestCaseEventRewardClaimed{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				RewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved).String(),
			},
		}
	case stakingevent.EventTypeAllRewardsClaimed:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeAllRewardsClaimed,
			AllRewardsClaimedEventData: &StakingStateTestCaseEventAllRewardsClaimed{
				BookNFT:      stakingEvent.NftClassAddress,
				Account:      stakingEvent.AccountEvmAddress,
				RewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved).String(),
			},
		}
	case stakingevent.EventTypeStakePositionTransferred:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeStakePositionTransferred,
			StakePositionTransferredEventData: &StakingStateTestCaseEventStakePositionTransferred{
				BookNFT:             stakingEvent.NftClassAddress,
				Account:             stakingEvent.AccountEvmAddress,
				StakedAmount:        (*uint256.Int)(stakingEvent.StakedAmountRemoved).String(),
				PendingRewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountRemoved).String(),
			},
		}
	case stakingevent.EventTypeStakePositionReceived:
		return &StakingStateTestCaseEvent{
			Type: StakingStateTestCaseEventTypeStakePositionReceived,
			StakePositionReceivedEventData: &StakingStateTestCaseEventStakePositionReceived{
				BookNFT:             stakingEvent.NftClassAddress,
				Account:             stakingEvent.AccountEvmAddress,
				StakedAmount:        (*uint256.Int)(stakingEvent.StakedAmountAdded).String(),
				PendingRewardAmount: (*uint256.Int)(stakingEvent.PendingRewardAmountAdded).String(),
			},
		}
	}
	panic("unknown event type")
}

type StakingStateAccount struct {
	EVMAddress          string `yaml:"address"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
	ClaimedRewardAmount string `yaml:"claimedRewardAmount"`
}

type StakingStateBookNFT struct {
	EVMAddress   string `yaml:"address"`
	StakedAmount string `yaml:"stakedAmount"`
}

type StakingStateStaking struct {
	AccountEVMAddress   string `yaml:"accountAddress"`
	BookNFTEvmAddress   string `yaml:"bookNFTAddress"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
	ClaimedRewardAmount string `yaml:"claimedRewardAmount"`
}

type StakingState struct {
	Accounts   map[string]StakingStateAccount            `yaml:"accounts"`
	NftClasses map[string]StakingStateBookNFT            `yaml:"nftClasses"`
	Stakings   map[string]map[string]StakingStateStaking `yaml:"stakings"`
}

type StakingStateTestCaseStep struct {
	Description            string                      `yaml:"description"`
	Events                 []StakingStateTestCaseEvent `yaml:"events"`
	StakingState           *MockStorage                `yaml:"stakingstate"`
	ProcessedStakingEvents []StakingStateTestCaseEvent `yaml:"processedStakingEvents"`
	Error                  string                      `yaml:"error"`
}

type StakingStateTestCaseSnapshottedAccounts struct {
	Address             string `yaml:"address"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
	ClaimedRewardAmount string `yaml:"claimedRewardAmount"`
}

func (s *StakingStateTestCaseSnapshottedAccounts) ToEnt() *ent.Account {
	stakedAmount, err := uint256.FromDecimal(s.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	pendingRewardAmount, err := uint256.FromDecimal(s.PendingRewardAmount)
	if err != nil {
		panic("failed to convert pending reward amount to big.Int")
	}
	claimedRewardAmount, err := uint256.FromDecimal(s.ClaimedRewardAmount)
	if err != nil {
		panic("failed to convert claimed reward amount to big.Int")
	}
	return &ent.Account{
		EvmAddress:          s.Address,
		StakedAmount:        typeutil.Uint256(stakedAmount),
		PendingRewardAmount: typeutil.Uint256(pendingRewardAmount),
		ClaimedRewardAmount: typeutil.Uint256(claimedRewardAmount),
	}
}

type StakingStateTestCaseSnapshottedNftClasses struct {
	Address      string `yaml:"address"`
	StakedAmount string `yaml:"stakedAmount"`
}

func (s *StakingStateTestCaseSnapshottedNftClasses) ToEnt() *ent.NFTClass {
	stakedAmount, err := uint256.FromDecimal(s.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	return &ent.NFTClass{
		Address:      s.Address,
		StakedAmount: typeutil.Uint256(stakedAmount),
	}
}

type StakingStateTestCaseSnapshottedStakings struct {
	AccountAddress      string `yaml:"accountAddress"`
	BookNFTAddress      string `yaml:"bookNFTAddress"`
	StakedAmount        string `yaml:"stakedAmount"`
	PendingRewardAmount string `yaml:"pendingRewardAmount"`
	ClaimedRewardAmount string `yaml:"claimedRewardAmount"`
}

func (s *StakingStateTestCaseSnapshottedStakings) ToEnt() *ent.Staking {
	stakedAmount, err := uint256.FromDecimal(s.StakedAmount)
	if err != nil {
		panic("failed to convert staked amount to big.Int")
	}
	pendingRewardAmount, err := uint256.FromDecimal(s.PendingRewardAmount)
	if err != nil {
		panic("failed to convert pending reward amount to big.Int")
	}
	claimedRewardAmount, err := uint256.FromDecimal(s.ClaimedRewardAmount)
	if err != nil {
		panic("failed to convert claimed reward amount to big.Int")
	}
	return &ent.Staking{
		Edges: ent.StakingEdges{
			Account: &ent.Account{
				EvmAddress: s.AccountAddress,
			},
			NftClass: &ent.NFTClass{
				Address: s.BookNFTAddress,
			},
		},
		StakedAmount:        typeutil.Uint256(stakedAmount),
		PendingRewardAmount: typeutil.Uint256(pendingRewardAmount),
		ClaimedRewardAmount: typeutil.Uint256(claimedRewardAmount),
	}
}

type StakingStateTestCaseParameters struct {
	SnapshottedAccounts   []StakingStateTestCaseSnapshottedAccounts   `yaml:"snapshottedAccounts"`
	SnapshottedNftClasses []StakingStateTestCaseSnapshottedNftClasses `yaml:"snapshottedNftClasses"`
	SnapshottedStakings   []StakingStateTestCaseSnapshottedStakings   `yaml:"snapshottedStakings"`
}

type StakingStateTestCase struct {
	Name       string                          `yaml:"name"`
	Parameters *StakingStateTestCaseParameters `yaml:"parameters"`
	Steps      []StakingStateTestCaseStep      `yaml:"steps"`
}

func TestStakingStateFromTestData(t *testing.T) {
	Convey("Test Staking State From Test Data", t, func() {
		rootDir := "staking_state_testdata/"
		entries, err := os.ReadDir(rootDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			f, err := os.Open(rootDir + entry.Name())
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			decoder := goyaml.NewDecoder(f)

			Convey(entry.Name(), func() {
				for {
					var testCase StakingStateTestCase
					err := decoder.Decode(&testCase)
					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						t.Fatal(err)
					}

					Convey(testCase.Name, func() {
						loader := MakeStakingStateTestMockLoader(
							testCase.Parameters.SnapshottedAccounts,
							testCase.Parameters.SnapshottedNftClasses,
							testCase.Parameters.SnapshottedStakings,
						)
						persistor := MakeStakingStateTestMockPersistor()

						stakingState, err := stakingstate.LoadStakingState(context.Background(), loader, []*ent.StakingEvent{})
						if err != nil {
							t.Fatal(err)
						}

						for _, step := range testCase.Steps {
							stakingEvents := make([]*ent.StakingEvent, 0)
							for _, event := range step.Events {
								stakingEvents = append(stakingEvents, event.ToStakingEvent())
							}
							var err error
							var processedStakingEvents []*ent.StakingEvent
							stakingState, processedStakingEvents, err = stakingState.Process(stakingEvents)
							if err != nil {
								if step.Error != "" {
									So(err.Error(), ShouldContainSubstring, step.Error)
									return
								}
								t.Fatal(err)
							}

							err = stakingState.Persist(context.Background(), nil, processedStakingEvents, persistor)
							if err != nil {
								t.Fatal(err)
							}

							actualStakingStateBytes, err := json.Marshal(persistor.mockStorage)
							if err != nil {
								t.Fatal(err)
							}

							expectedStakingStateBytes, err := json.Marshal(step.StakingState)
							if err != nil {
								t.Fatal(err)
							}

							processedTestCaseEvents := make([]*StakingStateTestCaseEvent, 0)
							for _, processedStakingEvent := range processedStakingEvents {
								processedTestCaseEvents = append(processedTestCaseEvents, FromEnt(processedStakingEvent))
							}

							actualTestCaseEventsBytes, err := json.Marshal(processedTestCaseEvents)
							if err != nil {
								t.Fatal(err)
							}

							expectedTestCaseEventsBytes, err := json.Marshal(step.ProcessedStakingEvents)
							if err != nil {
								t.Fatal(err)
							}

							if step.Description != "" {
								Convey(step.Description, func() {
									require.JSONEq(t, string(expectedStakingStateBytes), string(actualStakingStateBytes))
									require.JSONEq(t, string(expectedTestCaseEventsBytes), string(actualTestCaseEventsBytes))
								})
							} else {
								require.JSONEq(t, string(expectedStakingStateBytes), string(actualStakingStateBytes))
								require.JSONEq(t, string(expectedTestCaseEventsBytes), string(actualTestCaseEventsBytes))
							}
						}
					})
				}
			})
		}

	})
}
