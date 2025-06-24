package openapi

import (
	"context"

	"likenft-indexer/internal/api/openapi/model"
	"likenft-indexer/openapi/api"
)

func (h *OpenAPIHandler) TokenBookNFTsByAccount(ctx context.Context, params api.TokenBookNFTsByAccountParams) (*api.TokenBookNFTsByAccountOK, error) {
	pagination := model.NFTClassPagination{
		PaginationLimit: params.PaginationLimit,
		PaginationKey:   params.PaginationKey,
		Reverse:         params.Reverse,
	}

	nftClasses, count, nextKey, err := h.nftClassRepository.QueryNFTClassesByAccountTokens(
		ctx,
		params.EvmAddress,
		pagination.ToEntPagination(),
	)

	if err != nil {
		return nil, err
	}

	apiBookNFTs := make([]api.BookNFT, len(nftClasses))

	for i, nftClass := range nftClasses {
		apiNFTClass, err := model.MakeNFTClass(nftClass)
		if err != nil {
			return nil, err
		}
		apiBookNFTs[i] = *apiNFTClass
	}

	return &api.TokenBookNFTsByAccountOK{
		Pagination: api.PaginationResponse{
			NextKey: nextKey,
			Count:   count,
		},
		Data: apiBookNFTs,
	}, nil
}
