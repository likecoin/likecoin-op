package stakingstate

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"sync"
	"testing"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	testAccount = common.HexToAddress("0x2000000000000000000000000000000000000001")
	testBookNFT = common.HexToAddress("0x1000000000000000000000000000000000000001")
)

// stubEVMClient embeds the interface so only the head lookup and the three
// reads the reconciler makes need implementing; anything else panics, which is
// the point. Every read records the block it was asked for, and err fails the
// reads only -- the head lookup always succeeds, so a test setting it is
// exercising a read failure rather than a failure to resolve the head.
type stubEVMClient struct {
	evm.EVMClient
	head    *big.Int
	staked  map[common.Address]*big.Int
	pending map[common.Address]*big.Int
	total   map[common.Address]*big.Int
	err     error

	mu         sync.Mutex
	readBlocks []*big.Int
}

func (s *stubEVMClient) LatestBlockNumber(ctx context.Context) (*big.Int, error) {
	return s.head, nil
}

func (s *stubEVMClient) recordRead(blockNumber *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.readBlocks = append(s.readBlocks, blockNumber)
}

func (s *stubEVMClient) GetStakeForUser(
	ctx context.Context, blockNumber *big.Int, user common.Address, bookNFT common.Address,
) (*big.Int, error) {
	s.recordRead(blockNumber)
	if s.err != nil {
		return nil, s.err
	}
	return s.staked[user], nil
}

func (s *stubEVMClient) GetPendingRewardsForUser(
	ctx context.Context, blockNumber *big.Int, user common.Address, bookNFT common.Address,
) (*big.Int, error) {
	s.recordRead(blockNumber)
	if s.err != nil {
		return nil, s.err
	}
	return s.pending[user], nil
}

func (s *stubEVMClient) GetTotalStake(
	ctx context.Context, blockNumber *big.Int, bookNFT common.Address,
) (*big.Int, error) {
	s.recordRead(blockNumber)
	if s.err != nil {
		return nil, s.err
	}
	return s.total[bookNFT], nil
}

func stateWith(accountStaked, accountPending, stakingStaked, stakingPending uint64) *stakingState {
	return &stakingState{
		accounts: []*model.Account{{
			EVMAddress:          testAccount,
			StakedAmount:        uint256.NewInt(accountStaked),
			PendingRewardAmount: uint256.NewInt(accountPending),
			ClaimedRewardAmount: uint256.NewInt(7),
		}},
		nftClasses: []*model.NFTClass{{
			EVMAddress:   testBookNFT,
			StakedAmount: uint256.NewInt(0),
		}},
		stakings: []*model.Staking{{
			AccountEVMAddress:   testAccount,
			BookNFTEvmAddress:   testBookNFT,
			StakedAmount:        uint256.NewInt(stakingStaked),
			PendingRewardAmount: uint256.NewInt(stakingPending),
			ClaimedRewardAmount: uint256.NewInt(7),
		}},
	}
}

// testEventBlock is the block the event under test was emitted in; the stub
// chain's head sits on it unless a test moves one of them.
const testEventBlock = 1000

func chainSaying(staked, pending, total uint64) *stubEVMClient {
	return &stubEVMClient{
		head:    big.NewInt(testEventBlock),
		staked:  map[common.Address]*big.Int{testAccount: big.NewInt(int64(staked))},
		pending: map[common.Address]*big.Int{testAccount: big.NewInt(int64(pending))},
		total:   map[common.Address]*big.Int{testBookNFT: big.NewInt(int64(total))},
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func reconcile(t *testing.T, state *stakingState, client *stubEVMClient) {
	t.Helper()
	if err := reconcileFromChain(
		context.Background(), testLogger(), client, state, big.NewInt(testEventBlock),
	); err != nil {
		t.Fatalf("reconcileFromChain: %v", err)
	}
}

func TestReconcileReplacesDriftedAmountsWithTheChainsOwn(t *testing.T) {
	// The drift this exists for: the stored staked amount is wildly wrong
	// because a delta was wrong, and accumulating more deltas never fixes it.
	state := stateWith(190074458108, 0, 190074458108, 0)

	reconcile(t, state, chainSaying(60, 5, 60))

	if state.stakings[0].StakedAmount.Uint64() != 60 {
		t.Fatalf("staking staked = %s, want 60", state.stakings[0].StakedAmount)
	}
	if state.stakings[0].PendingRewardAmount.Uint64() != 5 {
		t.Fatalf("staking pending = %s, want 5", state.stakings[0].PendingRewardAmount)
	}
	if state.nftClasses[0].StakedAmount.Uint64() != 60 {
		t.Fatalf("nft class staked = %s, want 60", state.nftClasses[0].StakedAmount)
	}
}

func TestReconcileIsIdempotent(t *testing.T) {
	// An eth_call at block N returns the state after all of block N, so two
	// logs from the same block re-read the same answer. Applying it twice must
	// not double count -- that is the whole reason the amount is replaced
	// rather than added.
	state := stateWith(100, 2, 100, 2)
	client := chainSaying(140, 9, 140)

	reconcile(t, state, client)
	firstAccount := state.accounts[0].StakedAmount.Clone()
	firstPending := state.accounts[0].PendingRewardAmount.Clone()

	reconcile(t, state, client)

	if !state.accounts[0].StakedAmount.Eq(firstAccount) {
		t.Fatalf("account staked moved on the second pass: %s then %s",
			firstAccount, state.accounts[0].StakedAmount)
	}
	if !state.accounts[0].PendingRewardAmount.Eq(firstPending) {
		t.Fatalf("account pending moved on the second pass: %s then %s",
			firstPending, state.accounts[0].PendingRewardAmount)
	}
}

func TestReconcileMovesTheAccountByTheCorrectionOnly(t *testing.T) {
	// The account totals stakings that are not all loaded, so it cannot be
	// recomputed as a sum here; it has to move by exactly what the staking
	// moved. This account holds 250 across two book NFTs, only one loaded.
	state := stateWith(250, 10, 100, 4)

	reconcile(t, state, chainSaying(130, 6, 130))

	// staking went 100 -> 130, so the account goes 250 -> 280, keeping the
	// 150 held in the book NFT that was never loaded.
	if state.accounts[0].StakedAmount.Uint64() != 280 {
		t.Fatalf("account staked = %s, want 280", state.accounts[0].StakedAmount)
	}
	if state.accounts[0].PendingRewardAmount.Uint64() != 12 {
		t.Fatalf("account pending = %s, want 12", state.accounts[0].PendingRewardAmount)
	}
}

func TestReconcileHandlesAnAmountThatFell(t *testing.T) {
	state := stateWith(250, 10, 100, 4)

	reconcile(t, state, chainSaying(40, 1, 40))

	if state.accounts[0].StakedAmount.Uint64() != 190 {
		t.Fatalf("account staked = %s, want 190", state.accounts[0].StakedAmount)
	}
	if state.accounts[0].PendingRewardAmount.Uint64() != 7 {
		t.Fatalf("account pending = %s, want 7", state.accounts[0].PendingRewardAmount)
	}
}

func TestReconcileSkipsRowsHoldingNothing(t *testing.T) {
	// A pool keeps its fully unstaked rows forever. Re-reading them would make
	// a deposit cost grow with the pool's whole history rather than with its
	// current stakers.
	state := stateWith(0, 0, 0, 0)
	client := chainSaying(0, 0, 0)

	reconcile(t, state, client)

	// GetTotalStake still runs for the nft class; the two per-staking reads
	// must not.
	if got := len(client.readBlocks); got != 1 {
		t.Fatalf("chain was read %d times, want 1 (the nft class total only)", got)
	}
	if !state.stakings[0].StakedAmount.IsZero() {
		t.Fatal("an empty staking should have been left alone")
	}
}

func TestReconcileClampsInsteadOfWrappingTheAccountTotal(t *testing.T) {
	// The account total may already be too low from drift banked earlier. A
	// uint256 subtraction would wrap that into an astronomical balance.
	state := stateWith(5, 0, 100, 0)

	reconcile(t, state, chainSaying(0, 0, 0))

	if !state.accounts[0].StakedAmount.IsZero() {
		t.Fatalf("account staked = %s, want 0", state.accounts[0].StakedAmount)
	}
}

func TestReconcileLeavesClaimedRewardAlone(t *testing.T) {
	// Lifetime cumulative, and the contract keeps no counter for it, so a
	// snapshot has nothing to read it from.
	state := stateWith(100, 0, 100, 0)

	reconcile(t, state, chainSaying(1, 1, 1))

	if state.stakings[0].ClaimedRewardAmount.Uint64() != 7 {
		t.Fatalf("staking claimed = %s, want 7 untouched", state.stakings[0].ClaimedRewardAmount)
	}
	if state.accounts[0].ClaimedRewardAmount.Uint64() != 7 {
		t.Fatalf("account claimed = %s, want 7 untouched", state.accounts[0].ClaimedRewardAmount)
	}
}

func TestReconcileFailsRatherThanPersistPartialState(t *testing.T) {
	// A failed read must not leave some rows re-read and others accumulated;
	// the caller aborts and the event is retried.
	state := stateWith(100, 0, 100, 0)
	client := chainSaying(1, 1, 1)
	client.err = errors.New("rpc is down")

	err := reconcileFromChain(
		context.Background(), testLogger(), client, state, big.NewInt(testEventBlock),
	)

	if err == nil {
		t.Fatal("expected a read failure to be returned")
	}
	if len(client.readBlocks) == 0 {
		t.Fatal("the failure should have come from a read, not from resolving the head")
	}
}

func TestReconcilePinsEveryReadToOneHead(t *testing.T) {
	// A pool-wide pass is many reads. Each taking whatever block the node had
	// at that moment would mix rows from before and after a block that landed
	// mid-pass into a set of totals no block ever held.
	state := stateWith(100, 0, 100, 0)
	client := chainSaying(1, 1, 1)
	client.head = big.NewInt(testEventBlock + 5)

	reconcile(t, state, client)

	if len(client.readBlocks) == 0 {
		t.Fatal("expected the chain to be read")
	}
	for _, blockNumber := range client.readBlocks {
		if blockNumber == nil || blockNumber.Cmp(client.head) != 0 {
			t.Fatalf("a read was pinned to %s, want the head %s", blockNumber, client.head)
		}
	}
}

func TestReconcileRefusesAHeadBehindTheEvent(t *testing.T) {
	// A node that has not caught up to the event's own block would return a
	// truth older than the event, and nothing would re-drive the event to
	// correct it. Failing here hands the event back to asynq to retry.
	state := stateWith(100, 0, 100, 0)
	client := chainSaying(1, 1, 1)
	client.head = big.NewInt(testEventBlock - 1)

	err := reconcileFromChain(
		context.Background(), testLogger(), client, state, big.NewInt(testEventBlock),
	)

	if !errors.Is(err, ErrChainBehindEvent) {
		t.Fatalf("err = %v, want %v", err, ErrChainBehindEvent)
	}
	if len(client.readBlocks) != 0 {
		t.Fatal("nothing should have been read from a lagging head")
	}
	if state.stakings[0].StakedAmount.Uint64() != 100 {
		t.Fatalf("staking staked = %s, want 100 untouched", state.stakings[0].StakedAmount)
	}
}

func unstakedEvent(amount uint64) *ent.StakingEvent {
	return &ent.StakingEvent{
		EventType:           stakingevent.EventTypeUnstaked,
		AccountEvmAddress:   testAccount.Hex(),
		NftClassAddress:     testBookNFT.Hex(),
		StakedAmountRemoved: typeutil.Uint256(uint256.NewInt(amount)),
	}
}

func stakedEvent(amount uint64) *ent.StakingEvent {
	return &ent.StakingEvent{
		EventType:         stakingevent.EventTypeStaked,
		AccountEvmAddress: testAccount.Hex(),
		NftClassAddress:   testBookNFT.Hex(),
		StakedAmountAdded: typeutil.Uint256(uint256.NewInt(amount)),
	}
}

func TestChainBackedStateAppliesARemovalTheLoadedAmountAlreadyReflects(t *testing.T) {
	// Two unstakes of a 100 position are queued. The first was re-read at a
	// head past both, so the staking loads as 0; the account holds 30 in
	// another pool. Refusing the second removal against that 0 would fail it
	// on every retry, and resync would write the same 0 back.
	state := stateWith(30, 0, 0, 0)
	state.chainBacked = true

	processed, events, err := state.Process([]*ent.StakingEvent{unstakedEvent(50)})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("history rows = %d, want 1", len(events))
	}

	reconcile(t, processed, chainSaying(0, 0, 0))

	if !processed.stakings[0].StakedAmount.IsZero() {
		t.Fatalf("staking staked = %s, want 0", processed.stakings[0].StakedAmount)
	}
	if processed.accounts[0].StakedAmount.Uint64() != 30 {
		t.Fatalf("account staked = %s, want the 30 held elsewhere kept", processed.accounts[0].StakedAmount)
	}
}

func TestChainBackedStateReReadsAnEmptyRowAnEventNamed(t *testing.T) {
	// A stake into a row holding nothing leaves it holding nothing as loaded,
	// since the delta is not applied; the row must still be read, or the stake
	// is never seen.
	state := stateWith(30, 0, 0, 0)
	state.chainBacked = true

	processed, _, err := state.Process([]*ent.StakingEvent{stakedEvent(100)})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	reconcile(t, processed, chainSaying(100, 0, 100))

	if processed.stakings[0].StakedAmount.Uint64() != 100 {
		t.Fatalf("staking staked = %s, want 100", processed.stakings[0].StakedAmount)
	}
	if processed.accounts[0].StakedAmount.Uint64() != 130 {
		t.Fatalf("account staked = %s, want 130", processed.accounts[0].StakedAmount)
	}
}

func TestDeltaStateStillRefusesARemovalLargerThanTheAmount(t *testing.T) {
	// Simulation checks the delta arithmetic itself, so an underflow there is
	// the bug it exists to catch.
	state := stateWith(30, 0, 0, 0)

	_, _, err := state.Process([]*ent.StakingEvent{unstakedEvent(50)})

	if !errors.Is(err, ErrUnstakedEventApplication) {
		t.Fatalf("err = %v, want %v", err, ErrUnstakedEventApplication)
	}
}
