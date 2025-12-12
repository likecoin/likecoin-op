package database

import (
	"likecollective-indexer/ent"
	"likecollective-indexer/ent/account"
	"likecollective-indexer/ent/nftclass"

	"entgo.io/ent/dialect/sql"
)

type NFTClassPagination struct {
	Limit     *int
	Key       *int
	Reverse   *bool
	SortBy    *string
	SortOrder *string
}

func (f *NFTClassPagination) HandlePagination(q *ent.NFTClassQuery) *ent.NFTClassQuery {
	// Handle sorting
	usingCustomSort := f.SortBy != nil && f.SortOrder != nil
	if usingCustomSort {
		// Use custom sorting
		sortField := *f.SortBy
		var orderFunc func(*sql.Selector)
		if *f.SortOrder == "desc" {
			orderFunc = sql.OrderByField(sortField, sql.OrderDesc()).ToFunc()
		} else {
			orderFunc = sql.OrderByField(sortField, sql.OrderAsc()).ToFunc()
		}
		q = q.Order(orderFunc)
		// Add secondary sort by id for consistent pagination
		q = q.Order(sql.OrderByField("id", sql.OrderAsc()).ToFunc())
	} else {
		// Use default id-based sorting
		if f.Reverse != nil && *f.Reverse {
			q = q.Order(sql.OrderByField("id", sql.OrderDesc()).ToFunc())
		} else {
			q = q.Order(sql.OrderByField("id", sql.OrderAsc()).ToFunc())
		}
	}

	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	} else {
		q = q.Limit(100)
	}

	if f.Key != nil {
		if usingCustomSort {
			// When using custom sort, treat Key as offset
			q = q.Offset(*f.Key)
		} else {
			// When using default id-based sort, treat Key as cursor
			if f.Reverse != nil && *f.Reverse {
				q = q.Where(sql.FieldLT("id", *f.Key))
			} else {
				q = q.Where(sql.FieldGT("id", *f.Key))
			}
		}
	}

	return q
}

type NFTClassesRequestTimeFrameSortBy string

const (
	NFTClasssRequestTimeFrameSortByStakedAmount    NFTClassesRequestTimeFrameSortBy = "staked_amount"
	NFTClasssRequestTimeFrameSortByLastStakedAt    NFTClassesRequestTimeFrameSortBy = "last_staked_at"
	NFTClasssRequestTimeFrameSortByNumberOfStakers NFTClassesRequestTimeFrameSortBy = "number_of_stakers"
)

type NFTClassesRequestTimeFrameSortOrder string

const (
	NFTClasssRequestTimeFrameSortOrderAsc  NFTClassesRequestTimeFrameSortOrder = "asc"
	NFTClasssRequestTimeFrameSortOrderDesc NFTClassesRequestTimeFrameSortOrder = "desc"
)

type QueryNFTClassesFilter struct {
	filterBookNftIn *[]string
	filterAccountIn *[]string
}

func NewQueryNFTClassesFilter(
	filterBookNftIn *[]string,
	filterAccountIn *[]string,
) QueryNFTClassesFilter {
	return QueryNFTClassesFilter{
		filterBookNftIn,
		filterAccountIn,
	}
}

func (f *QueryNFTClassesFilter) HandleFilter(q *ent.NFTClassQuery) *ent.NFTClassQuery {
	if f.filterBookNftIn != nil {
		q = q.Where(nftclass.AddressIn(*f.filterBookNftIn...))
	}
	if f.filterAccountIn != nil {
		q = q.Where(nftclass.HasAccountsWith(account.EvmAddressIn(*f.filterAccountIn...)))
	}
	return q
}
