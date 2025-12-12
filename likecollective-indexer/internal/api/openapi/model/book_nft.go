package model

import (
	"likecollective-indexer/ent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/openapi/api"

	"github.com/holiman/uint256"
)

func MakeBookNFT(nftClass *ent.NFTClass) api.BookNFT {
	return api.BookNFT{
		EvmAddress:      api.EvmAddress(nftClass.Address),
		StakedAmount:    api.Uint256((*uint256.Int)(nftClass.StakedAmount).String()),
		LastStakedAt:    api.NewNilDateTime(nftClass.LastStakedAt),
		NumberOfStakers: int(nftClass.NumberOfStakers),
	}
}

type NFTClassFilterParams struct {
	FilterBookNftIn []api.EvmAddress
	FilterAccountIn []api.EvmAddress
}

func (p *NFTClassFilterParams) ToEntFilter() database.QueryNFTClassesFilter {

	var filterBookNftIn *[]string
	var filterAccountIn *[]string
	if len(p.FilterBookNftIn) > 0 {
		_filterBookNftIn := make([]string, len(p.FilterBookNftIn))
		for i, bookNft := range p.FilterBookNftIn {
			_filterBookNftIn[i] = string(bookNft)
		}
		filterBookNftIn = &_filterBookNftIn
	}
	if len(p.FilterAccountIn) > 0 {
		_filterAccountIn := make([]string, len(p.FilterAccountIn))
		for i, account := range p.FilterAccountIn {
			_filterAccountIn[i] = string(account)
		}
		filterAccountIn = &_filterAccountIn
	}

	return database.NewQueryNFTClassesFilter(filterBookNftIn, filterAccountIn)
}

type NFTClassPagination struct {
	// Pagination.limit.
	PaginationLimit api.OptInt
	// Pagination.key.
	PaginationKey api.OptInt
	// Reverse.
	Reverse api.OptBool
	// Sort by.
	SortBy api.OptBookNftsGetSortBy
	// Sort order.
	SortOrder api.OptBookNftsGetSortOrder
}

func (p *NFTClassPagination) ToEntPagination() database.NFTClassPagination {
	limit := FromOpt(p.PaginationLimit)
	if limit != nil && *limit == 0 {
		limit = nil
	}

	key := FromOpt(p.PaginationKey)
	if key != nil && *key == 0 {
		key = nil
	}

	reverse := FromOpt(p.Reverse)

	var sortBy *string
	if p.SortBy.Set {
		s := string(p.SortBy.Value)
		sortBy = &s
	}

	var sortOrder *string
	if p.SortOrder.Set {
		s := string(p.SortOrder.Value)
		sortOrder = &s
	}

	return database.NFTClassPagination{
		Limit:     limit,
		Key:       key,
		Reverse:   reverse,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}
