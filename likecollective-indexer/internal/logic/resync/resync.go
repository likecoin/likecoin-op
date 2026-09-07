// Package resync rebuilds the materialised staking state -- accounts,
// nft_classes and stakings -- from a single on-chain snapshot, instead of from
// the event log the webhook feeds. The webhook is the indexer's only ingest
// path, so a delivery that never arrives leaves those running totals wrong
// forever; this is the repair.
package resync

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"math/big"
	"slices"

	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/util/parallel"
	"likecollective-indexer/internal/util/retry"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// stakingKey identifies one stakings row. common.Address is a comparable
// [20]byte, so keying on it also normalises the mixed-case hex the database
// stores.
type stakingKey struct {
	Account common.Address
	BookNFT common.Address
}

// Change is one column the snapshot disagrees with the database about.
type Change struct {
	Table string
	Key   string
	Field string
	Old   string
	New   string
}

// Snapshot is the on-chain state at BlockNumber, shaped for
// persistor.StakingStatePersistor.Persist, together with how it differs from
// what is currently stored.
type Snapshot struct {
	BlockNumber *big.Int
	Positions   int

	Accounts   []*model.Account
	NFTClasses []*model.NFTClass
	Stakings   []*model.Staking

	Changes []Change
}

// Build reads the whole staking state off chain at blockNumber and compares it
// with the database. It writes nothing.
//
// claimed_reward_amount has no on-chain source -- the contract keeps no
// lifetime-claimed counter -- so it is carried over from the database as is,
// and account totals are recomputed as the sum over that account's stakings,
// which is the invariant the event applications maintain.
func Build(
	ctx context.Context,
	logger *slog.Logger,
	evmClient evm.EVMClient,
	dbService database.Service,
	blockNumber *big.Int,
	concurrency int,
) (*Snapshot, error) {
	dbStakings, err := database.MakeStakingRepository(dbService).QueryAllStakings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query stakings: %w", err)
	}
	dbNFTClasses, err := database.MakeNFTClassRepository(dbService).QueryAllNFTClasses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query nft classes: %w", err)
	}
	dbAccounts, err := database.MakeAccountRepository(dbService).QueryAllAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	positions, err := evmClient.ListStakePositions(ctx, blockNumber, concurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to list stake positions: %w", err)
	}
	logger.Info("enumerated stake positions",
		"block_number", blockNumber,
		"positions", len(positions),
	)

	// The union of what the chain has and what the database has. A row the
	// chain no longer knows about -- a position burned on a full unstake --
	// has to be zeroed, not left sitting at its last value, which is why the
	// database side is folded in rather than only iterated over on chain.
	stakingKeys := map[stakingKey]struct{}{}
	bookNFTs := map[common.Address]struct{}{}
	accountAddresses := map[common.Address]struct{}{}

	for _, position := range positions {
		stakingKeys[stakingKey{position.Owner, position.BookNFT}] = struct{}{}
		bookNFTs[position.BookNFT] = struct{}{}
		accountAddresses[position.Owner] = struct{}{}
	}
	for _, dbStaking := range dbStakings {
		key := stakingKey{
			Account: common.HexToAddress(dbStaking.Edges.Account.EvmAddress),
			BookNFT: common.HexToAddress(dbStaking.Edges.NftClass.Address),
		}
		stakingKeys[key] = struct{}{}
		bookNFTs[key.BookNFT] = struct{}{}
		accountAddresses[key.Account] = struct{}{}
	}
	for _, dbNFTClass := range dbNFTClasses {
		bookNFTs[common.HexToAddress(dbNFTClass.Address)] = struct{}{}
	}
	for _, dbAccount := range dbAccounts {
		accountAddresses[common.HexToAddress(dbAccount.EvmAddress)] = struct{}{}
	}

	sortedStakingKeys := slices.SortedFunc(maps.Keys(stakingKeys), compareStakingKey)
	sortedBookNFTs := slices.SortedFunc(maps.Keys(bookNFTs), compareAddress)
	sortedAccountAddresses := slices.SortedFunc(maps.Keys(accountAddresses), compareAddress)

	logger.Info("reading staking state from chain",
		"stakings", len(sortedStakingKeys),
		"nft_classes", len(sortedBookNFTs),
		"accounts", len(sortedAccountAddresses),
	)

	stakings, err := parallel.MapWithLimit(ctx, concurrency, sortedStakingKeys, func(
		ctx context.Context,
		key stakingKey,
	) (*model.Staking, error) {
		return retry.Value(ctx, retry.DefaultAttempts, retry.DefaultBase, func(ctx context.Context) (*model.Staking, error) {
			return readStaking(ctx, evmClient, blockNumber, key)
		})
	})
	if err != nil {
		return nil, err
	}

	nftClasses, err := parallel.MapWithLimit(ctx, concurrency, sortedBookNFTs, func(
		ctx context.Context,
		bookNFT common.Address,
	) (*model.NFTClass, error) {
		return retry.Value(ctx, retry.DefaultAttempts, retry.DefaultBase, func(ctx context.Context) (*model.NFTClass, error) {
			return readNFTClass(ctx, evmClient, blockNumber, bookNFT)
		})
	})
	if err != nil {
		return nil, err
	}

	// pool.totalStaked is maintained by stake/unstake, independently of the
	// positions it is meant to summarise. The two are read here through
	// different accessors, so disagreement means the contract's own state is
	// internally inconsistent -- worth saying out loud, but both halves are
	// still the best available truth, so it does not stop the snapshot.
	stakedByBookNFT := map[common.Address]*uint256.Int{}
	for _, staking := range stakings {
		total, ok := stakedByBookNFT[staking.BookNFTEvmAddress]
		if !ok {
			total = uint256.NewInt(0)
		}
		stakedByBookNFT[staking.BookNFTEvmAddress] = new(uint256.Int).Add(total, staking.StakedAmount)
	}
	for _, nftClass := range nftClasses {
		summed, ok := stakedByBookNFT[nftClass.EVMAddress]
		if !ok {
			summed = uint256.NewInt(0)
		}
		if !summed.Eq(nftClass.StakedAmount) {
			logger.Warn("total stake disagrees with the sum of its stakings",
				"book_nft", nftClass.EVMAddress,
				"sum_of_stakings", summed,
				"from_get_total_stake", nftClass.StakedAmount,
			)
		}
	}

	// claimed_reward_amount is not readable from the contract; carry the stored
	// value across so the resync does not wipe it.
	dbStakingsByKey := stakingsByKey(dbStakings)
	for _, staking := range stakings {
		key := stakingKey{staking.AccountEVMAddress, staking.BookNFTEvmAddress}
		if dbStaking, ok := dbStakingsByKey[key]; ok {
			staking.ClaimedRewardAmount = (*uint256.Int)(dbStaking.ClaimedRewardAmount)
		}
	}

	// getStakeForUser sums the contract's own per-user position index, which is
	// maintained separately from ERC721 ownership. Enumerating by ownerOf and
	// summing the same positions has to agree; if it does not, one of the two
	// indexes is corrupt and the snapshot is not trustworthy.
	positionStakedByKey := map[stakingKey]*uint256.Int{}
	for _, position := range positions {
		key := stakingKey{position.Owner, position.BookNFT}
		staked, err := toUint256(position.StakedAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to convert staked amount of position %s: %w", position.TokenID, err)
		}
		if existing, ok := positionStakedByKey[key]; ok {
			staked = new(uint256.Int).Add(existing, staked)
		}
		positionStakedByKey[key] = staked
	}
	mismatches := 0
	for _, staking := range stakings {
		key := stakingKey{staking.AccountEVMAddress, staking.BookNFTEvmAddress}
		fromPositions, ok := positionStakedByKey[key]
		if !ok {
			fromPositions = uint256.NewInt(0)
		}
		if !fromPositions.Eq(staking.StakedAmount) {
			mismatches++
			logger.Warn("staked amount disagrees between position enumeration and getStakeForUser",
				"account", staking.AccountEVMAddress,
				"book_nft", staking.BookNFTEvmAddress,
				"from_positions", fromPositions,
				"from_get_stake_for_user", staking.StakedAmount,
			)
		}
	}
	if mismatches > 0 {
		return nil, fmt.Errorf("%d stakings disagree between position enumeration and getStakeForUser", mismatches)
	}

	// Account rows are the per-account totals of their stakings.
	accountsByAddress := map[common.Address]*model.Account{}
	for _, address := range sortedAccountAddresses {
		accountsByAddress[address] = model.NewAccount(address.String())
	}
	for _, staking := range stakings {
		account := accountsByAddress[staking.AccountEVMAddress]
		account.StakedAmount = new(uint256.Int).Add(account.StakedAmount, staking.StakedAmount)
		account.PendingRewardAmount = new(uint256.Int).Add(account.PendingRewardAmount, staking.PendingRewardAmount)
		account.ClaimedRewardAmount = new(uint256.Int).Add(account.ClaimedRewardAmount, staking.ClaimedRewardAmount)
	}
	accounts := make([]*model.Account, 0, len(sortedAccountAddresses))
	for _, address := range sortedAccountAddresses {
		accounts = append(accounts, accountsByAddress[address])
	}

	return &Snapshot{
		BlockNumber: blockNumber,
		Positions:   len(positions),
		Accounts:    accounts,
		NFTClasses:  nftClasses,
		Stakings:    stakings,
		Changes: diff(
			accounts, nftClasses, stakings,
			dbAccounts, dbNFTClasses, dbStakingsByKey,
		),
	}, nil
}

func readStaking(
	ctx context.Context,
	evmClient evm.EVMClient,
	blockNumber *big.Int,
	key stakingKey,
) (*model.Staking, error) {
	stakedAmount, err := evmClient.GetStakeForUser(ctx, blockNumber, key.Account, key.BookNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to get stake for user %s on %s: %w", key.Account, key.BookNFT, err)
	}
	pendingRewardAmount, err := evmClient.GetPendingRewardsForUser(ctx, blockNumber, key.Account, key.BookNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending rewards for user %s on %s: %w", key.Account, key.BookNFT, err)
	}

	staked, err := toUint256(stakedAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to convert staked amount of %s on %s: %w", key.Account, key.BookNFT, err)
	}
	pending, err := toUint256(pendingRewardAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pending reward amount of %s on %s: %w", key.Account, key.BookNFT, err)
	}

	// claimed_reward_amount is filled in from the database by the caller.
	return &model.Staking{
		AccountEVMAddress:   key.Account,
		BookNFTEvmAddress:   key.BookNFT,
		StakedAmount:        staked,
		PendingRewardAmount: pending,
		ClaimedRewardAmount: uint256.NewInt(0),
	}, nil
}

func readNFTClass(
	ctx context.Context,
	evmClient evm.EVMClient,
	blockNumber *big.Int,
	bookNFT common.Address,
) (*model.NFTClass, error) {
	totalStake, err := evmClient.GetTotalStake(ctx, blockNumber, bookNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stake of %s: %w", bookNFT, err)
	}

	staked, err := toUint256(totalStake)
	if err != nil {
		return nil, fmt.Errorf("failed to convert total stake of %s: %w", bookNFT, err)
	}

	return &model.NFTClass{
		EVMAddress:   bookNFT,
		StakedAmount: staked,
	}, nil
}
