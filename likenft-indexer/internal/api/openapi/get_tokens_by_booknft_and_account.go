package openapi

import (
	"context"

	"likenft-indexer/internal/api/openapi/model"
	"likenft-indexer/openapi/api"
)

func (h *OpenAPIHandler) TokensByBookNFTAndAccount(ctx context.Context, params api.TokensByBookNFTAndAccountParams) (*api.TokensByBookNFTAndAccountOK, error) {
	ps := model.NFTPagination{
		PaginationLimit: params.PaginationLimit,
		PaginationKey:   params.PaginationKey,
		Reverse:         params.Reverse,
	}

	nfts, count, nextKey, err := h.nftRepository.QueryNFTsByBookNFTAndEvmAddress(
		ctx,
		params.ID,
		params.EvmAddress,
		ps.ToEntPagination(),
	)

	if err != nil {
		return nil, err
	}

	apiNFTs := make([]api.NFT, len(nfts))

	for i, n := range nfts {
		apiNFTs[i] = model.MakeNFT(n)
	}

	return &api.TokensByBookNFTAndAccountOK{
		Pagination: api.PaginationResponse{
			NextKey: nextKey,
			Count:   count,
		},
		Data: apiNFTs,
	}, nil
}
