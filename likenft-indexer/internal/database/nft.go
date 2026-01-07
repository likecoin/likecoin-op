package database

import (
	"context"
	"math/big"
	"time"

	"likenft-indexer/ent"
	"likenft-indexer/ent/nft"
	"likenft-indexer/ent/nftclass"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/evm/model"
)

type NFTRepository interface {
	QueryNFTsByEvmAddress(
		ctx context.Context,
		accountEvmAddress string,
		contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
		contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
		pagination NFTPagination,
	) (nfts []*ent.NFT, count int, nextKey int, err error)
	QueryNFTsByBookNFTAndEvmAddress(
		ctx context.Context,
		bookNFTId string,
		accountEvmAddress string,
		pagination NFTPagination,
	) (nfts []*ent.NFT, count int, nextKey int, err error)
	GetOrCreate(
		ctx context.Context,
		ContractAddress string,
		TokenID *big.Int,
		TokenURI string,
		Image string,
		ImageData *string,
		ExternalURL *string,
		Description string,
		Name string,
		Attributes []model.ERC721MetadataAttribute,
		BackgroundColor *string,
		AnimationURL *string,
		YoutubeURL *string,
		OwnerAddress string,
		Owner *ent.Account,
		nftClass *ent.NFTClass,
	) (*ent.NFT, error)

	UpdateOwner(
		ctx context.Context,
		ContractAddress string,
		TokenID *big.Int,
		NewOwnerAddress string,
		NewOwner *ent.Account,
	) error
}

type nftRepository struct {
	dbService Service
}

func MakeNFTRepository(
	dbService Service,
) NFTRepository {
	return &nftRepository{
		dbService: dbService,
	}
}

func (r *nftRepository) QueryNFTsByEvmAddress(
	ctx context.Context,
	accountEvmAddress string,
	contractLevelMetadataEQ ContractLevelMetadataFilterEquatable,
	contractLevelMetadataNEQ ContractLevelMetadataFilterEquatable,
	pagination NFTPagination,
) (nfts []*ent.NFT, count int, nextKey int, err error) {
	bookNFTQ := r.dbService.Client().NFTClass.Query()

	bookNFTQ = contractLevelMetadataEQ.ApplyEQ(bookNFTQ)
	bookNFTQ = contractLevelMetadataNEQ.ApplyNEQ(bookNFTQ)

	q := bookNFTQ.QueryNfts().Where(
		nft.OwnerAddressEqualFold(accountEvmAddress),
	)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	nfts, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(nfts) > 0 {
		nextKey = nfts[len(nfts)-1].ID
	}

	return nfts, count, nextKey, nil
}

func (r *nftRepository) QueryNFTsByBookNFTAndEvmAddress(
	ctx context.Context,
	bookNFTId string,
	accountEvmAddress string,
	pagination NFTPagination,
) (nfts []*ent.NFT, count int, nextKey int, err error) {
	q := r.dbService.Client().NFT.Query().
		Where(
			nft.HasClassWith(nftclass.AddressEqualFold(bookNFTId)),
			nft.OwnerAddressEqualFold(accountEvmAddress),
		)

	count, err = q.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	q = pagination.HandlePagination(q)

	nfts, err = q.All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	nextKey = 0
	if len(nfts) > 0 {
		nextKey = nfts[len(nfts)-1].ID
	}

	return nfts, count, nextKey, nil
}

func (r *nftRepository) GetOrCreate(
	ctx context.Context,
	contractAddress string,
	tokenID *big.Int,
	tokenURI string,
	image string,
	imageData *string,
	externalURL *string,
	description string,
	name string,
	attributes []model.ERC721MetadataAttribute,
	backgroundColor *string,
	animationURL *string,
	youtubeURL *string,
	ownerAddress string,
	owner *ent.Account,
	nftClass *ent.NFTClass,
) (*ent.NFT, error) {
	resChan := make(chan *ent.NFT, 1)

	err := WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		n, err := tx.NFT.Query().
			Where(
				nft.ContractAddressEqualFold(contractAddress),
				nft.TokenIDEQ(typeutil.Uint64(tokenID.Uint64()))).
			Only(ctx)

		if err == nil {
			resChan <- n
			return nil
		}

		if !ent.IsNotFound(err) {
			return err
		}

		err = tx.NFT.Create().
			SetContractAddress(contractAddress).
			SetTokenID(typeutil.Uint64(tokenID.Uint64())).
			SetTokenURI(tokenURI).
			SetImage(image).
			SetNillableImageData(imageData).
			SetNillableExternalURL(externalURL).
			SetDescription(description).
			SetName(name).
			SetAttributes(attributes).
			SetNillableBackgroundColor(backgroundColor).
			SetNillableAnimationURL(animationURL).
			SetNillableYoutubeURL(youtubeURL).
			SetOwnerAddress(ownerAddress).
			SetOwner(owner).
			SetMintedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetClass(nftClass).Exec(ctx)

		if err != nil {
			return err
		}

		n, err = tx.NFT.Query().
			Where(
				nft.ContractAddressEqualFold(contractAddress),
				nft.TokenIDEQ(typeutil.Uint64(tokenID.Uint64()))).
			Only(ctx)

		if err != nil {
			return err
		}

		resChan <- n
		return nil

	})

	if err != nil {
		return nil, err
	}
	return <-resChan, nil
}

func (r *nftRepository) UpdateOwner(
	ctx context.Context,
	contractAddress string,
	tokenID *big.Int,
	newOwnerAddress string,
	newOwner *ent.Account,
) error {
	err := WithTx(ctx, r.dbService.Client(), func(tx *ent.Tx) error {
		_, err := tx.NFT.Query().
			Where(
				nft.ContractAddressEqualFold(contractAddress),
				nft.TokenIDEQ(typeutil.Uint64(tokenID.Uint64())),
			).
			Only(ctx)

		if err != nil {
			return err
		}

		return tx.NFT.Update().
			SetOwnerAddress(newOwnerAddress).
			SetOwner(newOwner).
			Where(
				nft.ContractAddressEqualFold(contractAddress),
				nft.TokenIDEQ(typeutil.Uint64(tokenID.Uint64())),
			).Exec(ctx)
	})

	if err != nil {
		return err
	}
	return nil
}
