package stakingstate

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var testOtherAccount = common.HexToAddress("0x2000000000000000000000000000000000000002")

func rewardDepositedEvent(amount uint64) *ent.StakingEvent {
	return &ent.StakingEvent{
		EventType:                stakingevent.EventTypeRewardDeposited,
		BlockNumber:              typeutil.Uint64(testEventBlock),
		AccountEvmAddress:        testAccount.Hex(),
		NftClassAddress:          testBookNFT.Hex(),
		PendingRewardAmountAdded: typeutil.Uint256(uint256.NewInt(amount)),
	}
}

// poolAtHead is chain-backed state for a pool of two stakers, holding the
// amounts an earlier event re-read at a head past the deposit.
func poolAtHead(staked, otherStaked uint64) *stakingState {
	return &stakingState{
		nftClasses: []*model.NFTClass{{
			EVMAddress:   testBookNFT,
			StakedAmount: uint256.NewInt(staked + otherStaked),
		}},
		stakings: []*model.Staking{
			{
				AccountEVMAddress:   testAccount,
				BookNFTEvmAddress:   testBookNFT,
				StakedAmount:        uint256.NewInt(staked),
				PendingRewardAmount: uint256.NewInt(0),
				ClaimedRewardAmount: uint256.NewInt(0),
			},
			{
				AccountEVMAddress:   testOtherAccount,
				BookNFTEvmAddress:   testBookNFT,
				StakedAmount:        uint256.NewInt(otherStaked),
				PendingRewardAmount: uint256.NewInt(0),
				ClaimedRewardAmount: uint256.NewInt(0),
			},
		},
		chainBacked: true,
	}
}

func distributedTo(events []*ent.StakingEvent) map[common.Address]uint64 {
	distributed := make(map[common.Address]uint64)
	for _, event := range events {
		if event.EventType != stakingevent.EventTypeRewardDepositDistributed {
			continue
		}
		distributed[common.HexToAddress(event.AccountEvmAddress)] =
			(*uint256.Int)(event.PendingRewardAmountAdded).Uint64()
	}
	return distributed
}

func TestChainBackedDepositIsSplitByTheStakeAtItsBlock(t *testing.T) {
	// The queue is behind: at the deposit the two stakers held 50 each, but by
	// the head an earlier event read, the other staker had left. Splitting by
	// the loaded amounts would hand the whole deposit to the one who stayed.
	state := poolAtHead(50, 0)
	event := rewardDepositedEvent(10)
	client := &stubEVMClient{
		staked: map[common.Address]*big.Int{
			testAccount:      big.NewInt(50),
			testOtherAccount: big.NewInt(50),
		},
		total: map[common.Address]*big.Int{testBookNFT: big.NewInt(100)},
	}

	if err := readDepositStakes(
		context.Background(), client, state, []*ent.StakingEvent{event},
	); err != nil {
		t.Fatalf("readDepositStakes: %v", err)
	}
	_, events, err := state.Process([]*ent.StakingEvent{event})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	distributed := distributedTo(events)
	if distributed[testAccount] != 5 || distributed[testOtherAccount] != 5 {
		t.Fatalf("distributed = %v, want 5 each", distributed)
	}
	for _, block := range client.readBlocks {
		if block.Uint64() != testEventBlock {
			t.Fatalf("read at block %s, want every read at the deposit's block %d", block, testEventBlock)
		}
	}
	if len(client.readBlocks) != 3 {
		t.Fatalf("reads = %d, want the total and both stakers, the empty one included", len(client.readBlocks))
	}
}

func TestChainBackedDepositRefusesToSplitWithoutTheStakeAtItsBlock(t *testing.T) {
	// Falling back to the loaded amounts would quietly write the wrong history.
	state := poolAtHead(50, 50)

	_, _, err := state.Process([]*ent.StakingEvent{rewardDepositedEvent(10)})

	if !errors.Is(err, ErrRewardDepositedEventApplication) {
		t.Fatalf("err = %v, want %v", err, ErrRewardDepositedEventApplication)
	}
}

func TestChainBackedDepositRefusesToSplitWithAHolderNotIndexed(t *testing.T) {
	// A third holder's Staked was never processed, so no row of theirs is
	// loaded. Splitting would leave them out of the history for good.
	state := poolAtHead(50, 50)
	client := &stubEVMClient{
		staked: map[common.Address]*big.Int{
			testAccount:      big.NewInt(50),
			testOtherAccount: big.NewInt(50),
		},
		total: map[common.Address]*big.Int{testBookNFT: big.NewInt(150)},
	}

	err := readDepositStakes(
		context.Background(), client, state, []*ent.StakingEvent{rewardDepositedEvent(10)},
	)

	if !errors.Is(err, ErrHolderNotIndexed) {
		t.Fatalf("err = %v, want %v", err, ErrHolderNotIndexed)
	}
}

func TestChainBackedDepositSplitsAPoolWhoseTotalDriftedBelowItsPositions(t *testing.T) {
	// A restake before the contract fix left pending rewards out of the pool
	// total. The reads at that block will never change, so refusing would fail
	// the deposit on every retry.
	state := poolAtHead(50, 50)
	event := rewardDepositedEvent(10)
	client := &stubEVMClient{
		staked: map[common.Address]*big.Int{
			testAccount:      big.NewInt(50),
			testOtherAccount: big.NewInt(50),
		},
		total: map[common.Address]*big.Int{testBookNFT: big.NewInt(90)},
	}

	if err := readDepositStakes(
		context.Background(), client, state, []*ent.StakingEvent{event},
	); err != nil {
		t.Fatalf("readDepositStakes: %v", err)
	}
	if _, _, err := state.Process([]*ent.StakingEvent{event}); err != nil {
		t.Fatalf("Process: %v", err)
	}
}
