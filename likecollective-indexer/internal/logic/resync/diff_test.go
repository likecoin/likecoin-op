package resync

import (
	"testing"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/internal/logic/stakingstate/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

const (
	accountA = "0x2000000000000000000000000000000000000001"
	accountB = "0x2000000000000000000000000000000000000002"
	bookNFTA = "0x1000000000000000000000000000000000000001"
)

func amount(v uint64) typeutil.Uint256 {
	return typeutil.Uint256(uint256.NewInt(v))
}

func snapshotStaking(account string, bookNFT string, staked uint64, pending uint64) *model.Staking {
	return &model.Staking{
		AccountEVMAddress:   common.HexToAddress(account),
		BookNFTEvmAddress:   common.HexToAddress(bookNFT),
		StakedAmount:        uint256.NewInt(staked),
		PendingRewardAmount: uint256.NewInt(pending),
		ClaimedRewardAmount: uint256.NewInt(0),
	}
}

func storedStaking(account string, bookNFT string, staked uint64, pending uint64) *ent.Staking {
	s := &ent.Staking{
		StakedAmount:        amount(staked),
		PendingRewardAmount: amount(pending),
		ClaimedRewardAmount: amount(0),
	}
	s.Edges.Account = &ent.Account{EvmAddress: account}
	s.Edges.NftClass = &ent.NFTClass{Address: bookNFT}
	return s
}

func findChange(changes []Change, table string, field string) *Change {
	for i := range changes {
		if changes[i].Table == table && changes[i].Field == field {
			return &changes[i]
		}
	}
	return nil
}

func TestDiffReportsNothingWhenTheChainAgrees(t *testing.T) {
	stakings := []*model.Staking{snapshotStaking(accountA, bookNFTA, 100, 5)}
	stored := []*ent.Staking{storedStaking(accountA, bookNFTA, 100, 5)}

	changes := diff(nil, nil, stakings, nil, nil, stakingsByKey(stored))

	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestDiffReportsADriftedStakedAmount(t *testing.T) {
	stakings := []*model.Staking{snapshotStaking(accountA, bookNFTA, 60, 0)}
	stored := []*ent.Staking{storedStaking(accountA, bookNFTA, 190074458108, 0)}

	changes := diff(nil, nil, stakings, nil, nil, stakingsByKey(stored))

	change := findChange(changes, "stakings", "staked_amount")
	if change == nil {
		t.Fatalf("staked_amount drift not reported, got %+v", changes)
	}
	if change.Old != "190074458108" || change.New != "60" {
		t.Fatalf("unexpected change %+v", change)
	}
}

func TestDiffMatchesRowsAcrossAddressCase(t *testing.T) {
	// The database holds mixed-case hex from the webhook; keying on the raw
	// string would report every row as needing creation.
	stakings := []*model.Staking{snapshotStaking(accountA, bookNFTA, 100, 0)}
	stored := []*ent.Staking{storedStaking(
		common.HexToAddress(accountA).Hex(),
		common.HexToAddress(bookNFTA).Hex(),
		100, 0,
	)}

	changes := diff(nil, nil, stakings, nil, nil, stakingsByKey(stored))

	if len(changes) != 0 {
		t.Fatalf("case difference reported as drift: %+v", changes)
	}
}

func TestDiffReportsARowThatHasToBeCreated(t *testing.T) {
	stakings := []*model.Staking{snapshotStaking(accountB, bookNFTA, 7, 0)}

	changes := diff(nil, nil, stakings, nil, nil, stakingsByKey(nil))

	change := findChange(changes, "stakings", "*")
	if change == nil || change.Old != missing || change.New != "created" {
		t.Fatalf("row creation not reported, got %+v", changes)
	}
}

func TestDiffReportsNumberOfStakersFromNonZeroStakingsOnly(t *testing.T) {
	// number_of_stakers is recomputed on persist from stakings that still
	// carry stake, so a fully unstaked position must not be counted.
	nftClasses := []*model.NFTClass{{
		EVMAddress:   common.HexToAddress(bookNFTA),
		StakedAmount: uint256.NewInt(100),
	}}
	stakings := []*model.Staking{
		snapshotStaking(accountA, bookNFTA, 100, 0),
		snapshotStaking(accountB, bookNFTA, 0, 0),
	}
	storedNFTClasses := []*ent.NFTClass{{
		Address:         bookNFTA,
		StakedAmount:    amount(100),
		NumberOfStakers: 2,
	}}

	changes := diff(nil, nftClasses, stakings, nil, storedNFTClasses, stakingsByKey(nil))

	change := findChange(changes, "nft_classes", "number_of_stakers")
	if change == nil {
		t.Fatalf("number_of_stakers drift not reported, got %+v", changes)
	}
	if change.Old != "2" || change.New != "1" {
		t.Fatalf("unexpected change %+v", change)
	}
}

func TestCompareStakingKeyOrdersByAccountBeforeBookNFT(t *testing.T) {
	lowAccountHighBook := stakingKey{common.HexToAddress(accountA), common.HexToAddress(bookNFTA)}
	highAccountLowBook := stakingKey{common.HexToAddress(accountB), common.HexToAddress("0x0")}

	if compareStakingKey(lowAccountHighBook, highAccountLowBook) >= 0 {
		t.Fatal("account is not the primary ordering key")
	}
	if compareStakingKey(lowAccountHighBook, lowAccountHighBook) != 0 {
		t.Fatal("a key does not compare equal to itself")
	}
}
