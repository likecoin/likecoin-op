package database

import (
	"context"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/nftclass"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/staking"

	"entgo.io/ent/dialect/sql"
	"github.com/holiman/uint256"
)

type NFTClassRepository interface {
	QueryNFTClasses(
		ctx context.Context,
		filter QueryNFTClassesFilter,
		pagination NFTClassPagination,
	) (
		bookNFTs []*ent.NFTClass,
		count int,
		nextKey int,
		err error,
	)
	QueryNFTClass(
		ctx context.Context,
		address string,
	) (*ent.NFTClass, error)

	QueryNFTClassesByAddresses(
		ctx context.Context,
		addresses []string,
	) ([]*ent.NFTClass, error)

	GetOrCreateNFTClass(
		ctx context.Context,
		tx *ent.Tx,
		address string,
	) (*ent.NFTClass, error)

	CreateOrUpdateNFTClass(
		ctx context.Context,
		tx *ent.Tx,
		address string,
		stakedAmount typeutil.Uint256,
	) (*ent.NFTClass, error)

	RecomputeNumberOfStakersByNFTClassAddress(
		ctx context.Context,
		tx *ent.Tx,
		nftClassAddress string,
	) (*ent.NFTClass, error)
}

type nftClassRepository struct {
	dbService Service
}

func MakeNFTClassRepository(dbService Service) NFTClassRepository {
	return &nftClassRepository{dbService: dbService}
}

func (r *nftClassRepository) QueryNFTClasses(
	ctx context.Context,
	filter QueryNFTClassesFilter,
	pagination NFTClassPagination,
) (
	nftClasses []*ent.NFTClass,
	count int,
	nextKey int,
	err error,
) {
	q := r.dbService.Client().NFTClass.Query()
	q = filter.HandleFilter(q)

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
		usingCustomSort := pagination.SortBy != nil && pagination.SortOrder != nil
		if usingCustomSort {
			currentOffset := 0
			if pagination.Key != nil {
				currentOffset = *pagination.Key
			}
			if currentOffset+len(nftClasses) < count {
				nextKey = currentOffset + len(nftClasses)
			}
		} else {
			nextKey = nftClasses[len(nftClasses)-1].ID
		}
	}

	return nftClasses, count, nextKey, nil
}

func (r *nftClassRepository) QueryNFTClass(
	ctx context.Context,
	address string,
) (*ent.NFTClass, error) {
	nftClass, err := r.dbService.Client().NFTClass.Query().Where(
		nftclass.AddressEqualFold(address),
	).Only(ctx)
	if err != nil {
		return nil, err
	}

	return nftClass, nil
}

func (r *nftClassRepository) QueryNFTClassesByAddresses(
	ctx context.Context,
	addresses []string,
) ([]*ent.NFTClass, error) {
	return r.dbService.Client().NFTClass.Query().Where(nftclass.AddressIn(addresses...)).All(ctx)
}

func (r *nftClassRepository) GetOrCreateNFTClass(
	ctx context.Context,
	tx *ent.Tx,
	address string,
) (*ent.NFTClass, error) {
	nftClass, err := tx.NFTClass.Query().Where(nftclass.AddressEqualFold(address)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return tx.NFTClass.Create().
				SetAddress(address).
				SetStakedAmount((typeutil.Uint256)(uint256.NewInt(0))).
				SetNumberOfStakers(0).
				SetNillableLastStakedAt(nil).
				Save(ctx)
		}
		return nil, err
	}
	return nftClass, nil
}

func (r *nftClassRepository) CreateOrUpdateNFTClass(
	ctx context.Context,
	tx *ent.Tx,
	address string,
	stakedAmount typeutil.Uint256,
) (*ent.NFTClass, error) {
	nftClass, err := tx.NFTClass.Query().Where(nftclass.AddressEqualFold(address)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return tx.NFTClass.Create().
				SetAddress(address).
				SetStakedAmount(stakedAmount).
				Save(ctx)
		}
		return nil, err
	}
	return tx.NFTClass.UpdateOne(nftClass).
		SetStakedAmount(stakedAmount).
		Save(ctx)
}

func (r *nftClassRepository) RecomputeNumberOfStakersByNFTClassAddress(
	ctx context.Context,
	tx *ent.Tx,
	nftClassAddress string,
) (*ent.NFTClass, error) {
	nftClass, err := tx.NFTClass.Query().Where(nftclass.AddressEqualFold(nftClassAddress)).Only(ctx)
	if err != nil {
		return nil, err
	}

	c, err := tx.Staking.Query().Where(
		staking.NftClassIDEQ(nftClass.ID),
		func(s *sql.Selector) {
			s.Where(sql.GT(staking.FieldStakedAmount, 0))
		},
	).Count(ctx)
	if err != nil {
		return nil, err
	}

	return tx.NFTClass.UpdateOne(nftClass).SetNumberOfStakers(uint64(c)).Save(ctx)
}
