package database

import (
	"time"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/evmevent"

	"entgo.io/ent/dialect/sql"
)

type EvmEventsFilterSortBy string

const (
	EvmEventsFilterSortByBlockNumber    EvmEventsFilterSortBy = "block_number"
	EvmEventsFilterSortByBlockTimestamp EvmEventsFilterSortBy = "block_timestamp"
)

func (s EvmEventsFilterSortBy) AsEntOrderOption(opts ...sql.OrderTermOption) evmevent.OrderOption {
	switch s {
	case EvmEventsFilterSortByBlockNumber:
		return evmevent.ByBlockNumber(opts...)
	case EvmEventsFilterSortByBlockTimestamp:
		return evmevent.ByTimestamp(opts...)
	}
	return evmevent.ByID()
}

type EvmEventsFilterSortOrder string

const (
	EvmEventsFilterSortOrderAsc  EvmEventsFilterSortOrder = "asc"
	EvmEventsFilterSortOrderDesc EvmEventsFilterSortOrder = "desc"
)

func (o EvmEventsFilterSortOrder) AsSql() sql.OrderTermOption {
	switch o {
	case EvmEventsFilterSortOrderAsc:
		return sql.OrderAsc()
	case EvmEventsFilterSortOrderDesc:
		return sql.OrderDesc()
	}
	return sql.OrderAsc()
}

type EvmEventsFilter struct {
	// Contract address.
	Address *string
	// Event signature.
	Signature *string

	// Limit.
	Limit int
	// Page.
	Page int
	// Sort_by.
	SortBy *EvmEventsFilterSortBy
	// Sort_order.
	SortOrder *EvmEventsFilterSortOrder
	// Filter_block_timestamp.
	FilterBlockTimestamp *time.Time
	// Filter_block_timestamp_gte.
	FilterBlockTimestampGte *time.Time
	// Filter_block_timestamp_gt.
	FilterBlockTimestampGt *time.Time
	// Filter_block_timestamp_lte.
	FilterBlockTimestampLte *time.Time
	// Filter_block_timestamp_lt.
	FilterBlockTimestampLt *time.Time
	// Filter_topic_1.
	FilterTopic1 *string
	// Filter_topic_2.
	FilterTopic2 *string
	// Filter_topic_3.
	FilterTopic3 *string
	// Filter_topic_0.
	FilterTopic0 *string
}

func MakeEvmEventsFilter(
	address *string,
	signature *string,

	limit int,
	page int,
	sortBy *EvmEventsFilterSortBy,
	sortOrder *EvmEventsFilterSortOrder,
	filterBlockTimestamp *time.Time,
	filterBlockTimestampGte *time.Time,
	filterBlockTimestampGt *time.Time,
	filterBlockTimestampLte *time.Time,
	filterBlockTimestampLt *time.Time,
	filterTopic1 *string,
	filterTopic2 *string,
	filterTopic3 *string,
	filterTopic0 *string,
) *EvmEventsFilter {
	return &EvmEventsFilter{
		Address:                 address,
		Signature:               signature,
		Limit:                   limit,
		Page:                    page,
		SortBy:                  sortBy,
		SortOrder:               sortOrder,
		FilterBlockTimestamp:    filterBlockTimestamp,
		FilterBlockTimestampGte: filterBlockTimestampGte,
		FilterBlockTimestampGt:  filterBlockTimestampGt,
		FilterBlockTimestampLte: filterBlockTimestampLte,
		FilterBlockTimestampLt:  filterBlockTimestampLt,
		FilterTopic1:            filterTopic1,
		FilterTopic2:            filterTopic2,
		FilterTopic3:            filterTopic3,
		FilterTopic0:            filterTopic0,
	}
}

func (f *EvmEventsFilter) HandleSort(q *ent.EVMEventQuery) *ent.EVMEventQuery {
	if f.SortBy != nil {
		q = q.Order(f.SortBy.AsEntOrderOption(f.SortOrder.AsSql()))
	}
	return q
}

func (f *EvmEventsFilter) HandleFilter(q *ent.EVMEventQuery) *ent.EVMEventQuery {
	if f.Address != nil {
		q = q.Where(evmevent.AddressEqualFold(*f.Address))
	}
	if f.Signature != nil {
		q = q.Where(evmevent.SignatureEqualFold(*f.Signature))
	}
	if f.FilterBlockTimestamp != nil {
		q = q.Where(evmevent.TimestampEQ(*f.FilterBlockTimestamp))
	}
	if f.FilterBlockTimestampGte != nil {
		q = q.Where(evmevent.TimestampGTE(*f.FilterBlockTimestampGte))
	}
	if f.FilterBlockTimestampGt != nil {
		q = q.Where(evmevent.TimestampGT(*f.FilterBlockTimestampGt))
	}
	if f.FilterBlockTimestampLte != nil {
		q = q.Where(evmevent.TimestampLTE(*f.FilterBlockTimestampLte))
	}
	if f.FilterBlockTimestampLt != nil {
		q = q.Where(evmevent.TimestampLT(*f.FilterBlockTimestampLt))
	}
	if f.FilterTopic1 != nil {
		q = q.Where(evmevent.Topic1EQ(*f.FilterTopic1))
	}
	if f.FilterTopic2 != nil {
		q = q.Where(evmevent.Topic2EQ(*f.FilterTopic2))
	}
	if f.FilterTopic3 != nil {
		q = q.Where(evmevent.Topic3EQ(*f.FilterTopic3))
	}
	if f.FilterTopic0 != nil {
		q = q.Where(evmevent.Topic0EQ(*f.FilterTopic0))
	}
	return q
}

func (f *EvmEventsFilter) HandlePagination(q *ent.EVMEventQuery) *ent.EVMEventQuery {
	return q.Limit(f.Limit).Offset(f.Page * f.Limit)
}
