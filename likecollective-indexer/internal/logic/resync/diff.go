package resync

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"likecollective-indexer/ent"
	entaccount "likecollective-indexer/ent/account"
	entnftclass "likecollective-indexer/ent/nftclass"
	entstaking "likecollective-indexer/ent/staking"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/util/ordered"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// missing is what a Change reports as the old value of a row the resync has to
// create.
const missing = "(none)"

// Table names as the report prints them. The ent Label constants are singular
// and name the entity, not the table, so there is nothing to reuse here.
const (
	tableAccounts   = "accounts"
	tableNFTClasses = "nft_classes"
	tableStakings   = "stakings"
)

// compareAddress orders addresses by their raw bytes. ordered.Normalize cannot:
// common.Address is a [20]byte array, which is not cmp.Ordered.
func compareAddress(a common.Address, b common.Address) int {
	return bytes.Compare(a[:], b[:])
}

var cmpStakingKeyAccount ordered.Comparator[stakingKey] = func(a, b stakingKey) int {
	return compareAddress(a.Account, b.Account)
}

var cmpStakingKeyBookNFT ordered.Comparator[stakingKey] = func(a, b stakingKey) int {
	return compareAddress(a.BookNFT, b.BookNFT)
}

var compareStakingKey = ordered.CombineComparators(
	cmpStakingKeyAccount,
	cmpStakingKeyBookNFT,
)

func toUint256(value *big.Int) (*uint256.Int, error) {
	converted, overflow := uint256.FromBig(value)
	if overflow {
		return nil, fmt.Errorf("value %s overflows uint256", value)
	}
	return converted, nil
}

func appendChange(changes []Change, table string, key string, field string, old string, new string) []Change {
	if old == new {
		return changes
	}
	return append(changes, Change{
		Table: table,
		Key:   key,
		Field: field,
		Old:   old,
		New:   new,
	})
}

// stakingsByKey indexes stored stakings by the (account, book NFT) pair they
// are unique on, normalising the mixed-case hex the database holds.
func stakingsByKey(dbStakings []*ent.Staking) map[stakingKey]*ent.Staking {
	byKey := make(map[stakingKey]*ent.Staking, len(dbStakings))
	for _, dbStaking := range dbStakings {
		byKey[stakingKey{
			Account: common.HexToAddress(dbStaking.Edges.Account.EvmAddress),
			BookNFT: common.HexToAddress(dbStaking.Edges.NftClass.Address),
		}] = dbStaking
	}
	return byKey
}

// diff reports every column the snapshot would rewrite. pool_share and
// last_staked_at are left out: the first is derived from the amounts and
// recomputed on persist, the second is not part of a snapshot at all.
func diff(
	accounts []*model.Account,
	nftClasses []*model.NFTClass,
	stakings []*model.Staking,
	dbAccounts []*ent.Account,
	dbNFTClasses []*ent.NFTClass,
	dbStakingsByKey map[stakingKey]*ent.Staking,
) []Change {
	changes := []Change{}

	dbAccountsByAddress := map[common.Address]*ent.Account{}
	for _, dbAccount := range dbAccounts {
		dbAccountsByAddress[common.HexToAddress(dbAccount.EvmAddress)] = dbAccount
	}
	dbNFTClassesByAddress := map[common.Address]*ent.NFTClass{}
	for _, dbNFTClass := range dbNFTClasses {
		dbNFTClassesByAddress[common.HexToAddress(dbNFTClass.Address)] = dbNFTClass
	}
	// number_of_stakers is recomputed on persist from the stakings that carry a
	// non-zero staked_amount; mirror that here so the report covers it.
	stakerCounts := map[common.Address]int{}
	for _, staking := range stakings {
		if !staking.StakedAmount.IsZero() {
			stakerCounts[staking.BookNFTEvmAddress]++
		}
	}

	for _, account := range accounts {
		key := account.EVMAddress.String()
		dbAccount, ok := dbAccountsByAddress[account.EVMAddress]
		if !ok {
			changes = appendChange(changes, tableAccounts, key, "*", missing, "created")
			continue
		}
		changes = appendChange(changes, tableAccounts, key, entaccount.FieldStakedAmount,
			(*uint256.Int)(dbAccount.StakedAmount).String(), account.StakedAmount.String())
		changes = appendChange(changes, tableAccounts, key, entaccount.FieldPendingRewardAmount,
			(*uint256.Int)(dbAccount.PendingRewardAmount).String(), account.PendingRewardAmount.String())
		changes = appendChange(changes, tableAccounts, key, entaccount.FieldClaimedRewardAmount,
			(*uint256.Int)(dbAccount.ClaimedRewardAmount).String(), account.ClaimedRewardAmount.String())
	}

	for _, nftClass := range nftClasses {
		key := nftClass.EVMAddress.String()
		dbNFTClass, ok := dbNFTClassesByAddress[nftClass.EVMAddress]
		if !ok {
			changes = appendChange(changes, tableNFTClasses, key, "*", missing, "created")
			continue
		}
		changes = appendChange(changes, tableNFTClasses, key, entnftclass.FieldStakedAmount,
			(*uint256.Int)(dbNFTClass.StakedAmount).String(), nftClass.StakedAmount.String())
		changes = appendChange(changes, tableNFTClasses, key, entnftclass.FieldNumberOfStakers,
			strconv.FormatUint(uint64(dbNFTClass.NumberOfStakers), 10),
			strconv.Itoa(stakerCounts[nftClass.EVMAddress]))
	}

	for _, staking := range stakings {
		key := fmt.Sprintf("%s@%s", staking.AccountEVMAddress, staking.BookNFTEvmAddress)
		dbStaking, ok := dbStakingsByKey[stakingKey{staking.AccountEVMAddress, staking.BookNFTEvmAddress}]
		if !ok {
			changes = appendChange(changes, tableStakings, key, "*", missing, "created")
			continue
		}
		changes = appendChange(changes, tableStakings, key, entstaking.FieldStakedAmount,
			(*uint256.Int)(dbStaking.StakedAmount).String(), staking.StakedAmount.String())
		changes = appendChange(changes, tableStakings, key, entstaking.FieldPendingRewardAmount,
			(*uint256.Int)(dbStaking.PendingRewardAmount).String(), staking.PendingRewardAmount.String())
		changes = appendChange(changes, tableStakings, key, entstaking.FieldClaimedRewardAmount,
			(*uint256.Int)(dbStaking.ClaimedRewardAmount).String(), staking.ClaimedRewardAmount.String())
	}

	return changes
}
