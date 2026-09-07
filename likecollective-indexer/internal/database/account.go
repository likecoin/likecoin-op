package database

import (
	"context"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/account"
	"likecollective-indexer/ent/nftclass"
	"likecollective-indexer/ent/predicate"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/staking"

	"github.com/holiman/uint256"
)

type AccountRepository interface {
	QueryAccounts(ctx context.Context, pagination AccountPagination, filter QueryAccountsFilter) (accounts []*ent.Account, count int, nextKey int, err error)
	QueryAccount(ctx context.Context, evmAddress string) (*ent.Account, error)

	QueryAccountsByEvmAddresses(ctx context.Context, evmAddresses []string) ([]*ent.Account, error)
	QueryAccountsByNFTClassAddresses(ctx context.Context, nftClassAddresses []string) ([]*ent.Account, error)

	// QueryAllAccounts returns every account, unfiltered and unpaginated.
	QueryAllAccounts(ctx context.Context) ([]*ent.Account, error)

	GetOrCreateAccount(
		ctx context.Context,
		tx *ent.Tx,
		evmAddress string,
	) (*ent.Account, error)

	CreateOrUpdateAccount(
		ctx context.Context,
		tx *ent.Tx,
		evmAddress string,
		stakedAmount typeutil.Uint256,
		pendingRewardAmount typeutil.Uint256,
		claimedRewardAmount typeutil.Uint256,
	) (*ent.Account, error)
}

type accountRepository struct {
	dbService Service
}

func MakeAccountRepository(dbService Service) AccountRepository {
	return &accountRepository{dbService: dbService}
}

func (r *accountRepository) QueryAccounts(
	ctx context.Context,
	pagination AccountPagination,
	filter QueryAccountsFilter,
) (
	accounts []*ent.Account,
	count int,
	nextKey int,
	err error,
) {
	q := r.dbService.Client().Account.Query()
	q = filter.HandleFilter(q)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	accounts, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(accounts) > 0 {
		nextKey = accounts[len(accounts)-1].ID
	}

	return accounts, count, 0, nil
}

func (r *accountRepository) QueryAccount(
	ctx context.Context,
	evmAddress string,
) (*ent.Account, error) {
	account, err := r.dbService.Client().Account.Query().Where(
		account.EvmAddressEqualFold(evmAddress),
	).Only(ctx)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) QueryAccountsByEvmAddresses(
	ctx context.Context,
	evmAddresses []string,
) ([]*ent.Account, error) {
	return r.dbService.Client().Account.Query().Where(
		account.EvmAddressIn(evmAddresses...),
	).All(ctx)
}

func (r *accountRepository) QueryAccountsByNFTClassAddresses(
	ctx context.Context,
	nftClassAddresses []string,
) ([]*ent.Account, error) {
	if len(nftClassAddresses) == 0 {
		return []*ent.Account{}, nil
	}

	addressPredicates := make([]predicate.NFTClass, 0)
	for _, bookNFTAddress := range nftClassAddresses {
		addressPredicates = append(addressPredicates, nftclass.AddressEqualFold(bookNFTAddress))
	}
	dbStakings, err := r.dbService.Client().Staking.Query().WithAccount().WithNftClass().Where(
		staking.HasNftClassWith(nftclass.Or(addressPredicates...)),
	).All(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make([]*ent.Account, 0)
	for _, staking := range dbStakings {
		accounts = append(accounts, staking.Edges.Account)
	}

	return accounts, nil
}

func (r *accountRepository) GetOrCreateAccount(
	ctx context.Context,
	tx *ent.Tx,
	evmAddress string,
) (*ent.Account, error) {
	account, err := tx.Account.Query().
		Where(account.EvmAddressEqualFold(evmAddress)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return tx.Account.Create().
				SetEvmAddress(evmAddress).
				SetStakedAmount((typeutil.Uint256)(uint256.NewInt(0))).
				SetPendingRewardAmount((typeutil.Uint256)(uint256.NewInt(0))).
				SetClaimedRewardAmount((typeutil.Uint256)(uint256.NewInt(0))).
				Save(ctx)
		}
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) CreateOrUpdateAccount(
	ctx context.Context,
	tx *ent.Tx,
	evmAddress string,
	stakedAmount typeutil.Uint256,
	pendingRewardAmount typeutil.Uint256,
	claimedRewardAmount typeutil.Uint256,
) (*ent.Account, error) {
	account, err := tx.Account.Query().
		Where(account.EvmAddressEqualFold(evmAddress)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return tx.Account.Create().
				SetEvmAddress(evmAddress).
				SetStakedAmount(stakedAmount).
				SetPendingRewardAmount(pendingRewardAmount).
				SetClaimedRewardAmount(claimedRewardAmount).
				Save(ctx)
		}
		return nil, err
	}
	return tx.Account.UpdateOne(account).
		SetStakedAmount(stakedAmount).
		SetPendingRewardAmount(pendingRewardAmount).
		SetClaimedRewardAmount(claimedRewardAmount).
		Save(ctx)
}

func (r *accountRepository) QueryAllAccounts(ctx context.Context) ([]*ent.Account, error) {
	return r.dbService.Client().Account.Query().All(ctx)
}
