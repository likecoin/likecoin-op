package resync

import (
	"context"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"testing"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

const bookNFTB = "0x1000000000000000000000000000000000000002"

// stubEVMClient serves a fixed chain state. Embedding the interface leaves
// every method build does not call unimplemented, so reaching one panics.
type stubEVMClient struct {
	evm.EVMClient

	positions []*evm.StakePosition
	stakes    map[stakingKey]int64
	pending   map[stakingKey]int64
	totals    map[common.Address]int64
}

func (c *stubEVMClient) ListStakePositions(context.Context, *big.Int, int) ([]*evm.StakePosition, error) {
	return c.positions, nil
}

func (c *stubEVMClient) GetStakeForUser(_ context.Context, _ *big.Int, user common.Address, bookNFT common.Address) (*big.Int, error) {
	return big.NewInt(c.stakes[stakingKey{user, bookNFT}]), nil
}

func (c *stubEVMClient) GetPendingRewardsForUser(_ context.Context, _ *big.Int, user common.Address, bookNFT common.Address) (*big.Int, error) {
	return big.NewInt(c.pending[stakingKey{user, bookNFT}]), nil
}

func (c *stubEVMClient) GetTotalStake(_ context.Context, _ *big.Int, bookNFT common.Address) (*big.Int, error) {
	return big.NewInt(c.totals[bookNFT]), nil
}

func key(account string, bookNFT string) stakingKey {
	return stakingKey{common.HexToAddress(account), common.HexToAddress(bookNFT)}
}

func position(tokenID int64, account string, bookNFT string, staked int64) *evm.StakePosition {
	return &evm.StakePosition{
		TokenID:      big.NewInt(tokenID),
		Owner:        common.HexToAddress(account),
		BookNFT:      common.HexToAddress(bookNFT),
		StakedAmount: big.NewInt(staked),
	}
}

func runBuild(
	t *testing.T,
	client *stubEVMClient,
	dbAccounts []*ent.Account,
	dbNFTClasses []*ent.NFTClass,
	dbStakings []*ent.Staking,
) (*Snapshot, error) {
	t.Helper()
	return build(
		context.Background(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		client,
		big.NewInt(1),
		2,
		dbAccounts,
		dbNFTClasses,
		dbStakings,
	)
}

func storedAccount(account string, staked uint64) *ent.Account {
	return &ent.Account{
		EvmAddress:          account,
		StakedAmount:        amount(staked),
		PendingRewardAmount: amount(0),
		ClaimedRewardAmount: amount(0),
	}
}

func storedNFTClass(bookNFT string, staked uint64) *ent.NFTClass {
	return &ent.NFTClass{Address: bookNFT, StakedAmount: amount(staked)}
}

func findStaking(snapshot *Snapshot, k stakingKey) *model.Staking {
	for _, staking := range snapshot.Stakings {
		if staking.AccountEVMAddress == k.Account && staking.BookNFTEvmAddress == k.BookNFT {
			return staking
		}
	}
	return nil
}

func findAccount(snapshot *Snapshot, account string) *model.Account {
	for _, a := range snapshot.Accounts {
		if a.EVMAddress == common.HexToAddress(account) {
			return a
		}
	}
	return nil
}

func expectAmount(t *testing.T, name string, got *uint256.Int, want uint64) {
	t.Helper()
	if !got.Eq(uint256.NewInt(want)) {
		t.Fatalf("%s: got %s, want %d", name, got, want)
	}
}

func TestBuildZeroesAStakingTheChainNoLongerHas(t *testing.T) {
	// A position burned on a full unstake is gone from enumeration, but its
	// row is still in the database and must be zeroed rather than skipped.
	stored := storedStaking(accountA, bookNFTA, 100, 5)
	stored.ClaimedRewardAmount = amount(7)

	snapshot, err := runBuild(t, &stubEVMClient{},
		[]*ent.Account{storedAccount(accountA, 100)},
		[]*ent.NFTClass{storedNFTClass(bookNFTA, 100)},
		[]*ent.Staking{stored},
	)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	staking := findStaking(snapshot, key(accountA, bookNFTA))
	if staking == nil {
		t.Fatalf("database-only staking dropped from snapshot: %+v", snapshot.Stakings)
	}
	expectAmount(t, "staked", staking.StakedAmount, 0)
	expectAmount(t, "pending", staking.PendingRewardAmount, 0)
	expectAmount(t, "claimed", staking.ClaimedRewardAmount, 7)

	if len(snapshot.NFTClasses) != 1 {
		t.Fatalf("expected the nft class to be carried, got %+v", snapshot.NFTClasses)
	}
	expectAmount(t, "nft class staked", snapshot.NFTClasses[0].StakedAmount, 0)

	if findChange(snapshot.Changes, tableStakings, "staked_amount") == nil {
		t.Fatalf("zeroing not reported, got %+v", snapshot.Changes)
	}
}

func TestBuildSumsAccountTotalsAndCarriesClaimedRewards(t *testing.T) {
	client := &stubEVMClient{
		positions: []*evm.StakePosition{
			position(1, accountA, bookNFTA, 60),
			position(2, accountA, bookNFTA, 40),
			position(3, accountA, bookNFTB, 50),
		},
		stakes: map[stakingKey]int64{
			key(accountA, bookNFTA): 100,
			key(accountA, bookNFTB): 50,
		},
		pending: map[stakingKey]int64{
			key(accountA, bookNFTA): 3,
			key(accountA, bookNFTB): 4,
		},
		totals: map[common.Address]int64{
			common.HexToAddress(bookNFTA): 100,
			common.HexToAddress(bookNFTB): 50,
		},
	}
	stored := storedStaking(accountA, bookNFTA, 100, 3)
	stored.ClaimedRewardAmount = amount(9)

	snapshot, err := runBuild(t, client,
		// accountB has no stakings anywhere and must still come out zeroed.
		[]*ent.Account{storedAccount(accountA, 100), storedAccount(accountB, 1)},
		nil,
		[]*ent.Staking{stored},
	)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	a := findAccount(snapshot, accountA)
	if a == nil {
		t.Fatalf("account A missing: %+v", snapshot.Accounts)
	}
	expectAmount(t, "account A staked", a.StakedAmount, 150)
	expectAmount(t, "account A pending", a.PendingRewardAmount, 7)
	expectAmount(t, "account A claimed", a.ClaimedRewardAmount, 9)

	b := findAccount(snapshot, accountB)
	if b == nil {
		t.Fatalf("database-only account B dropped: %+v", snapshot.Accounts)
	}
	expectAmount(t, "account B staked", b.StakedAmount, 0)

	onB := findStaking(snapshot, key(accountA, bookNFTB))
	if onB == nil {
		t.Fatalf("chain-only staking missing: %+v", snapshot.Stakings)
	}
	expectAmount(t, "claimed without a stored row", onB.ClaimedRewardAmount, 0)
}

func TestBuildRejectsPositionsThatDisagreeWithGetStakeForUser(t *testing.T) {
	client := &stubEVMClient{
		positions: []*evm.StakePosition{position(1, accountA, bookNFTA, 100)},
		stakes:    map[stakingKey]int64{key(accountA, bookNFTA): 90},
		totals:    map[common.Address]int64{common.HexToAddress(bookNFTA): 90},
	}

	snapshot, err := runBuild(t, client, nil, nil, nil)
	if err == nil {
		t.Fatalf("expected a mismatch error, got snapshot %+v", snapshot)
	}
	if !strings.Contains(err.Error(), "disagree") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffReportsPoolShareRecomputedOnPersist(t *testing.T) {
	nftClasses := []*model.NFTClass{{
		EVMAddress:   common.HexToAddress(bookNFTA),
		StakedAmount: uint256.NewInt(200),
	}}
	stakings := []*model.Staking{snapshotStaking(accountA, bookNFTA, 50, 0)}
	stored := storedStaking(accountA, bookNFTA, 50, 0)
	stored.PoolShare = "0.5"

	changes := diff(nil, nftClasses, stakings, nil, nil, stakingsByKey([]*ent.Staking{stored}))

	change := findChange(changes, tableStakings, "pool_share")
	if change == nil {
		t.Fatalf("pool_share drift not reported, got %+v", changes)
	}
	if change.Old != "0.5" || change.New != "0.25" {
		t.Fatalf("unexpected change %+v", change)
	}
}
