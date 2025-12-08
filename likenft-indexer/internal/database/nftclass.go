package database

import (
	"context"
	"math/big"
	"slices"
	"time"

	"likenft-indexer/ent"
	"likenft-indexer/ent/nft"
	"likenft-indexer/ent/nftclass"
	"likenft-indexer/ent/predicate"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/evm/model"

	"entgo.io/ent/dialect/sql"
)

type NFTClassRepository interface {
	QueryAllNFTClasses(ctx context.Context) ([]*ent.NFTClass, error)
	QueryAllNFTClassesOfLowestEventBlockHeight(
		ctx context.Context,
		filterDisabledForIndexing bool,
	) ([]*ent.NFTClass, error)
	QueryNFTClassByAddress(ctx context.Context, address string) (*ent.NFTClass, error)
	QueryNFTClassesByAddressesExact(
		ctx context.Context,
		addresses []string,
	) ([]*ent.NFTClass, error)
	QueryNFTClasses(
		ctx context.Context,
		contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
		contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
		pagination NFTClassPagination,
	) (nftClasses []*ent.NFTClass, count int, nextKey int, err error)
	QueryNFTClassesByAccountTokens(
		ctx context.Context,
		accountEvmAddress string,
		contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
		contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
		pagination NFTClassPagination,
	) (nftClasses []*ent.NFTClass, count int, nextKey int, err error)
	QueryNFTClassesByEvmAddress(
		ctx context.Context,
		accountEvmAddress string,
		contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
		contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
		pagination NFTClassPagination,
	) (nftClasses []*ent.NFTClass, count int, nextKey int, err error)
	InsertNFTClass(
		ctx context.Context,
		address string,
		name string,
		symbol string,
		ownerAddress *string,
		minterAddresses []string,
		totalSupply *big.Int,
		maxSupply typeutil.Uint64,
		metadata *model.ContractLevelMetadata,
		bannerImage string,
		featuredImage string,
		deployerAddress string,
		deployedBlockNumber typeutil.Uint64,
		latestEventBlockNumber typeutil.Uint64,
		mintedAt time.Time,
		owner *ent.Account,
	) error
	UpdateMetadata(ctx context.Context, address string, metadata *model.ContractLevelMetadata) error
	UpdateOwner(ctx context.Context, address string, newOwner *ent.Account) error
	UpdateMinterAddresses(ctx context.Context, address string, minterAddresses []string) error
	UpdateUpdaterAddresses(ctx context.Context, address string, updaterAddresses []string) error
	UpdateTotalSupply(ctx context.Context, address string, newTotalSupply *big.Int) error
	UpdateNFTClassesLatestEventBlockNumber(
		ctx context.Context,
		addresses []string,
		latestEventBlockNumber typeutil.Uint64,
	) error
	DisableForIndexing(
		ctx context.Context,
		address string,
		reason string,
	) error
}

type nftClassRepository struct {
	dbService Service
}

func MakeNFTClassRepository(
	dbService Service,
) NFTClassRepository {
	return &nftClassRepository{
		dbService: dbService,
	}
}

func (r *nftClassRepository) QueryAllNFTClasses(ctx context.Context) ([]*ent.NFTClass, error) {
	return r.dbService.Client().NFTClass.Query().All(ctx)
}

func (r *nftClassRepository) QueryAllNFTClassesOfLowestEventBlockHeight(
	ctx context.Context,
	filterDisabledForIndexing bool,
) ([]*ent.NFTClass, error) {
	q := r.dbService.Client().NFTClass.Query()

	if filterDisabledForIndexing {
		q = q.Where(nftclass.DisabledForIndexingEQ(false))
	}

	return q.
		Order(nftclass.ByLatestEventBlockNumber(sql.OrderAsc())).
		All(ctx)
}

func (r *nftClassRepository) QueryNFTClassByAddress(ctx context.Context, address string) (*ent.NFTClass, error) {
	return r.dbService.Client().NFTClass.Query().
		Where(nftclass.AddressEqualFold(address)).
		Only(ctx)
}

func (r *nftClassRepository) QueryNFTClassesByAddressesExact(
	ctx context.Context,
	addresses []string,
) ([]*ent.NFTClass, error) {
	res, err := r.dbService.Client().NFTClass.Query().
		Where(nftclass.AddressIn(addresses...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		if !slices.ContainsFunc(res, func(nftClass *ent.NFTClass) bool {
			return nftClass.Address == address
		}) {
			return nil, &ent.NotFoundError{}
		}
	}

	return res, nil
}

func (r *nftClassRepository) QueryNFTClasses(
	ctx context.Context,
	contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
	contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
	pagination NFTClassPagination,
) (nftClasses []*ent.NFTClass, count int, nextKey int, err error) {
	q := r.dbService.Client().NFTClass.Query()

	q = contractLevelMetadataEQ.ApplyEQ(q)
	q = contractLevelMetadataNEQ.ApplyNEQ(q)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	nftClasses, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(nftClasses) > 0 {
		nextKey = nftClasses[len(nftClasses)-1].ID
	}

	return nftClasses, count, nextKey, nil
}

func (r *nftClassRepository) QueryNFTClassesByAccountTokens(
	ctx context.Context,
	accountEvmAddress string,
	contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
	contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
	pagination NFTClassPagination,
) (nftClasses []*ent.NFTClass, count int, nextKey int, err error) {
	q := r.dbService.Client().NFTClass.Query().
		Where(
			nftclass.HasNftsWith(
				nft.OwnerAddressEqualFold(accountEvmAddress),
			),
		)

	q = contractLevelMetadataEQ.ApplyEQ(q)
	q = contractLevelMetadataNEQ.ApplyNEQ(q)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	nftClasses, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(nftClasses) > 0 {
		nextKey = nftClasses[len(nftClasses)-1].ID
	}

	return nftClasses, count, nextKey, nil
}

func (r *nftClassRepository) QueryNFTClassesByEvmAddress(
	ctx context.Context,
	accountEvmAddress string,
	contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
	contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
	pagination NFTClassPagination,
) (
	nftClasses []*ent.NFTClass,
	count int,
	nextKey int,
	err error,
) {
	q := r.dbService.Client().NFTClass.Query().
		Where(
			nftclass.OwnerAddressEqualFold(accountEvmAddress),
		)

	q = contractLevelMetadataEQ.ApplyEQ(q)
	q = contractLevelMetadataNEQ.ApplyNEQ(q)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	nftClasses, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(nftClasses) > 0 {
		nextKey = nftClasses[len(nftClasses)-1].ID
	}

	return nftClasses, count, nextKey, nil
}

func (r *nftClassRepository) InsertNFTClass(
	ctx context.Context,
	address string,
	name string,
	symbol string,
	ownerAddress *string,
	minterAddresses []string,
	totalSupply *big.Int,
	maxSupply typeutil.Uint64,
	metadata *model.ContractLevelMetadata,
	bannerImage string,
	featuredImage string,
	deployerAddress string,
	deployedBlockNumber typeutil.Uint64,
	latestEventBlockNumber typeutil.Uint64,
	mintedAt time.Time,
	owner *ent.Account,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		builder := tx.NFTClass.Create().
			SetAddress(address).
			SetName(name).
			SetSymbol(symbol).
			SetNillableOwnerAddress(ownerAddress).
			SetMinterAddresses(minterAddresses).
			SetTotalSupply(totalSupply).
			SetMaxSupply(maxSupply).
			SetMetadata(metadata).
			SetBannerImage(bannerImage).
			SetFeaturedImage(featuredImage).
			SetDeployerAddress(deployerAddress).
			SetDeployedBlockNumber(deployedBlockNumber).
			SetLatestEventBlockNumber(latestEventBlockNumber).
			SetMintedAt(mintedAt).
			SetOwner(owner).
			SetUpdatedAt(time.Now())

		return builder.Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateMetadata(
	ctx context.Context,
	address string,
	metadata *model.ContractLevelMetadata,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		_, err := tx.NFTClass.Query().
			Where(nftclass.AddressEqualFold(address)).
			Only(ctx)
		if err != nil {
			return err
		}
		return r.dbService.Client().NFTClass.Update().
			SetMetadata(metadata).
			Where(nftclass.AddressEqualFold(address)).
			Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateOwner(
	ctx context.Context,
	address string,
	newOwner *ent.Account,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		_, err := tx.NFTClass.Query().
			Where(nftclass.AddressEqualFold(address)).
			Only(ctx)
		if err != nil {
			return err
		}
		return r.dbService.Client().NFTClass.Update().
			SetOwner(newOwner).
			SetOwnerAddress(newOwner.EvmAddress).
			Where(nftclass.AddressEqualFold(address)).
			Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateMinterAddresses(
	ctx context.Context,
	address string,
	minterAddresses []string,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		return tx.NFTClass.Update().
			SetMinterAddresses(minterAddresses).
			Where(nftclass.AddressEqualFold(address)).
			Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateUpdaterAddresses(
	ctx context.Context,
	address string,
	updaterAddresses []string,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		return tx.NFTClass.Update().
			SetUpdaterAddresses(updaterAddresses).
			Where(nftclass.AddressEqualFold(address)).
			Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateTotalSupply(
	ctx context.Context,
	address string,
	newTotalSupply *big.Int,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		_, err := tx.NFTClass.Query().
			Where(nftclass.AddressEqualFold(address)).
			Only(ctx)
		if err != nil {
			return err
		}
		return r.dbService.Client().NFTClass.Update().
			SetTotalSupply(newTotalSupply).
			Where(nftclass.AddressEqualFold(address)).
			Exec(ctx)
	})
}

func (r *nftClassRepository) UpdateNFTClassesLatestEventBlockNumber(
	ctx context.Context,
	addresses []string,
	latestEventBlockNumber typeutil.Uint64,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		predicates := make([]predicate.NFTClass, 0, len(addresses))
		for _, address := range addresses {
			predicates = append(predicates, nftclass.AddressEqualFold(address))
		}
		return tx.NFTClass.
			Update().
			Where(nftclass.Or(predicates...)).
			SetLatestEventBlockNumber(latestEventBlockNumber).Exec(ctx)
	})
}

func (r *nftClassRepository) DisableForIndexing(
	ctx context.Context,
	address string,
	reason string,
) error {
	return WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		return tx.NFTClass.Update().
			Where(nftclass.AddressEqualFold(address)).
			SetDisabledForIndexing(true).
			SetDisabledForIndexingReason(reason).
			Exec(ctx)
	})
}
