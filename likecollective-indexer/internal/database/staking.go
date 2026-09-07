package database

import (
	"context"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/account"
	"likecollective-indexer/ent/nftclass"
	"likecollective-indexer/ent/predicate"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/staking"

	"github.com/holiman/uint256"
	"github.com/shopspring/decimal"
)

type StakingKey struct {
	AccountEVMAddress string
	BookNFTEVMAddress string
}

func NewStakingKey(
	accountEVMAddress string,
	bookNFTEVMAddress string,
) StakingKey {
	return StakingKey{
		AccountEVMAddress: accountEVMAddress,
		BookNFTEVMAddress: bookNFTEVMAddress,
	}
}

func (k *StakingKey) ToPredicate() predicate.Staking {
	return staking.And(
		staking.HasAccountWith(account.EvmAddressEqualFold(k.AccountEVMAddress)),
		staking.HasNftClassWith(nftclass.AddressEqualFold(k.BookNFTEVMAddress)),
	)
}

type StakingRepository interface {
	QueryStakings(
		ctx context.Context,
		filter QueryStakingsFilter,
		pagination StakingPagination,
	) (stakings []*ent.Staking, count int, nextKey int, err error)

	QueryStakingsByKeys(
		ctx context.Context,
		keys []StakingKey,
	) ([]*ent.Staking, error)

	QueryStakingsByNFTClassAddresses(
		ctx context.Context,
		bookNFTAddresses []string,
	) ([]*ent.Staking, error)

	// QueryAllStakings returns every staking with its account and nft class
	// edges loaded, unfiltered and unpaginated.
	QueryAllStakings(ctx context.Context) ([]*ent.Staking, error)

	GetOrCreateStaking(
		ctx context.Context,
		tx *ent.Tx,
		nftClassID int,
		accountID int,
	) (*ent.Staking, error)

	CreateOrUpdateStaking(
		ctx context.Context,
		tx *ent.Tx,
		bookNFTEvmAddress string,
		accountEvmAddress string,
		stakedAmount typeutil.Uint256,
		pendingRewardAmount typeutil.Uint256,
		claimedRewardAmount typeutil.Uint256,
	) (*ent.Staking, error)

	RecomputePoolSharesByNFTClassAddress(
		ctx context.Context,
		tx *ent.Tx,
		nftClassAddress string,
	) error
}

type stakingRepository struct {
	dbService Service
}

func MakeStakingRepository(
	dbService Service,
) StakingRepository {
	return &stakingRepository{
		dbService: dbService,
	}
}

func (r *stakingRepository) QueryStakings(
	ctx context.Context,
	filter QueryStakingsFilter,
	pagination StakingPagination,
) (
	stakings []*ent.Staking,
	count int,
	nextKey int,
	err error,
) {
	q := r.dbService.Client().Staking.Query().
		WithAccount().
		WithNftClass()
	q = filter.HandleFilter(q)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	q = pagination.HandlePagination(q)

	stakings, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(stakings) > 0 {
		nextKey = stakings[len(stakings)-1].ID
	}

	return stakings, count, nextKey, nil
}

func (r *stakingRepository) QueryStakingsByKeys(
	ctx context.Context,
	keys []StakingKey,
) ([]*ent.Staking, error) {
	keysPredicates := make([]predicate.Staking, 0)
	for _, key := range keys {
		keysPredicates = append(
			keysPredicates,
			key.ToPredicate(),
		)
	}
	return r.dbService.Client().Staking.Query().WithAccount().WithNftClass().Where(staking.Or(keysPredicates...)).All(ctx)
}

func (r *stakingRepository) QueryStakingsByNFTClassAddresses(
	ctx context.Context,
	nftClassAddresses []string,
) ([]*ent.Staking, error) {
	if len(nftClassAddresses) == 0 {
		return []*ent.Staking{}, nil
	}

	predicates := make([]predicate.Staking, 0)
	for _, bookNFTAddress := range nftClassAddresses {
		predicates = append(predicates, staking.HasNftClassWith(nftclass.AddressEqualFold(bookNFTAddress)))
	}

	return r.dbService.Client().Staking.Query().WithAccount().WithNftClass().Where(staking.Or(predicates...)).All(ctx)
}

func (r *stakingRepository) GetOrCreateStaking(
	ctx context.Context,
	tx *ent.Tx,
	nftClassID int,
	accountID int,
) (*ent.Staking, error) {
	staking, err := tx.Staking.Query().
		Where(staking.NftClassID(nftClassID), staking.AccountID(accountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return tx.Staking.Create().
				SetNftClassID(nftClassID).
				SetAccountID(accountID).
				SetPoolShare("0").
				SetStakedAmount((typeutil.Uint256)(uint256.NewInt(0))).
				SetPendingRewardAmount((typeutil.Uint256)(uint256.NewInt(0))).
				SetClaimedRewardAmount((typeutil.Uint256)(uint256.NewInt(0))).
				Save(ctx)
		}
		return nil, err
	}
	return staking, nil
}

func (r *stakingRepository) CreateOrUpdateStaking(
	ctx context.Context,
	tx *ent.Tx,
	bookNFTEvmAddress string,
	accountEvmAddress string,
	stakedAmount typeutil.Uint256,
	pendingRewardAmount typeutil.Uint256,
	claimedRewardAmount typeutil.Uint256,
) (*ent.Staking, error) {
	account, err := tx.Account.Query().Where(account.EvmAddressEqualFold(accountEvmAddress)).Only(ctx)
	if err != nil {
		return nil, err
	}
	nftClass, err := tx.NFTClass.Query().Where(nftclass.AddressEqualFold(bookNFTEvmAddress)).Only(ctx)
	if err != nil {
		return nil, err
	}
	s, err := tx.Staking.Query().WithAccount().WithNftClass().
		Where(
			staking.AccountIDEQ(account.ID),
			staking.NftClassIDEQ(nftClass.ID),
		).
		Only(ctx)

	var poolShares *big.Rat
	if (*uint256.Int)(nftClass.StakedAmount).IsZero() {
		poolShares = big.NewRat(0, 1)
	} else {
		poolShares = big.NewRat((*uint256.Int)(stakedAmount).ToBig().Int64(), (*uint256.Int)(nftClass.StakedAmount).ToBig().Int64())
	}
	poolSharesPercentage := decimal.NewFromBigRat(poolShares, 2)

	if err != nil {
		if ent.IsNotFound(err) {
			return tx.Staking.Create().
				SetAccountID(account.ID).
				SetNftClassID(nftClass.ID).
				SetPoolShare(poolSharesPercentage.String()).
				SetStakedAmount(stakedAmount).
				SetPendingRewardAmount(pendingRewardAmount).
				SetClaimedRewardAmount(claimedRewardAmount).
				Save(ctx)
		}
		return nil, err
	}

	return tx.Staking.UpdateOne(s).
		SetStakedAmount(stakedAmount).
		SetPoolShare(poolSharesPercentage.String()).
		SetPendingRewardAmount(pendingRewardAmount).
		SetClaimedRewardAmount(claimedRewardAmount).
		Save(ctx)
}

func (r *stakingRepository) RecomputePoolSharesByNFTClassAddress(
	ctx context.Context,
	tx *ent.Tx,
	nftClassAddress string,
) error {
	nftClass, err := tx.NFTClass.Query().Where(nftclass.AddressEqualFold(nftClassAddress)).Only(ctx)
	if err != nil {
		return err
	}

	stakings, err := tx.Staking.Query().Where(staking.NftClassIDEQ(nftClass.ID)).All(ctx)
	if err != nil {
		return err
	}

	for _, staking := range stakings {
		var poolShares *big.Rat
		if (*uint256.Int)(nftClass.StakedAmount).IsZero() {
			poolShares = big.NewRat(0, 1)
		} else {
			poolShares = big.NewRat((*uint256.Int)(staking.StakedAmount).ToBig().Int64(), (*uint256.Int)(nftClass.StakedAmount).ToBig().Int64())
		}
		poolSharePercentage := decimal.NewFromBigRat(poolShares, 2)
		if _, err = tx.Staking.UpdateOne(staking).SetPoolShare(poolSharePercentage.String()).Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *stakingRepository) QueryAllStakings(ctx context.Context) ([]*ent.Staking, error) {
	return r.dbService.Client().Staking.Query().
		WithAccount().
		WithNftClass().
		All(ctx)
}
