package model

import (
	"strconv"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/openapi/api"
)

func MakeEvent(e *ent.EVMEvent) (api.Event, error) {
	data := ""

	if e.Data != nil {
		data = *e.DataHex
	}

	topics := make([]string, 1, 4)
	topics[0] = e.Topic0Hex

	if e.Topic1Hex != nil {
		topics = append(topics, *e.Topic1Hex)
		if e.Topic2Hex != nil {
			topics = append(topics, *e.Topic2Hex)
			if e.Topic3Hex != nil {
				topics = append(topics, *e.Topic3Hex)
			}
		}
	}

	indexedParams, err := MakeAPIAdditionalProps(e.IndexedParams)
	if err != nil {
		return api.Event{}, err
	}

	nonIndexedParams, err := MakeAPIAdditionalProps(e.NonIndexedParams)
	if err != nil {
		return api.Event{}, err
	}

	return api.Event{
		ChainID:          int(e.ChainID),
		BlockNumber:      strconv.FormatUint(uint64(e.BlockNumber), 10),
		BlockHash:        e.BlockHash,
		BlockTimestamp:   strconv.FormatInt(e.Timestamp.Unix(), 10),
		TransactionHash:  e.TransactionHash,
		TransactionIndex: int(e.TransactionIndex),
		LogIndex:         int(e.LogIndex),
		Address:          e.Address,
		Data:             data,
		Topics:           topics,
		Decoded: api.EventDecoded{
			Name:             e.Name,
			Signature:        e.Signature,
			IndexedParams:    api.EventDecodedIndexedParams(indexedParams),
			NonIndexedParams: api.EventDecodedNonIndexedParams(nonIndexedParams),
		},
	}, nil
}

type OpenAPIEventParameters struct {
	// Contract address.
	Address *string
	// Event signature.
	Signature *string

	// Limit.
	Limit api.OptInt
	// Page.
	Page api.OptInt
	// Sort_by.
	SortBy api.OptEventSortRequestSortBy
	// Sort_order.
	SortOrder api.OptEventSortRequestSortOrder
	// Filter_block_timestamp.
	FilterBlockTimestamp api.OptString
	// Filter_block_timestamp_gte.
	FilterBlockTimestampGte api.OptString
	// Filter_block_timestamp_gt.
	FilterBlockTimestampGt api.OptString
	// Filter_block_timestamp_lte.
	FilterBlockTimestampLte api.OptString
	// Filter_block_timestamp_lt.
	FilterBlockTimestampLt api.OptString
	// Filter_topic_1.
	FilterTopic1 api.OptString
	// Filter_topic_2.
	FilterTopic2 api.OptString
	// Filter_topic_3.
	FilterTopic3 api.OptString
	// Filter_topic_0.
	FilterTopic0 api.OptString
}

func (e *OpenAPIEventParameters) ToEntFilter() (*database.EvmEventsFilter, error) {
	sortBy := FromOpt(e.SortBy)
	var filterSortBy *database.EvmEventsFilterSortBy = nil
	if sortBy != nil {
		f := database.EvmEventsFilterSortBy(*sortBy)
		filterSortBy = &f
	}

	sortOrder := FromOpt(e.SortOrder)
	var filterSortOrder *database.EvmEventsFilterSortOrder = nil
	if sortOrder != nil {
		f := database.EvmEventsFilterSortOrder(*sortOrder)
		filterSortOrder = &f
	}

	filterBlockTimeStamp, err := TimeFromOptString(e.FilterBlockTimestamp)
	if err != nil {
		return nil, err
	}
	filterBlockTimestampGte, err := TimeFromOptString(e.FilterBlockTimestampGte)
	if err != nil {
		return nil, err
	}
	filterBlockTimestampGt, err := TimeFromOptString(e.FilterBlockTimestampGt)
	if err != nil {
		return nil, err
	}
	filterBlockTimestampLte, err := TimeFromOptString(e.FilterBlockTimestampLte)
	if err != nil {
		return nil, err
	}
	filterBlockTimestampLt, err := TimeFromOptString(e.FilterBlockTimestampLt)
	if err != nil {
		return nil, err
	}
	filterTopic1 := FromOpt(e.FilterTopic1)
	filterTopic2 := FromOpt(e.FilterTopic2)
	filterTopic3 := FromOpt(e.FilterTopic3)
	filterTopic0 := FromOpt(e.FilterTopic0)

	return database.MakeEvmEventsFilter(
		e.Address,
		e.Signature,
		e.Limit.Value,
		e.Page.Value,
		filterSortBy,
		filterSortOrder,
		filterBlockTimeStamp,
		filterBlockTimestampGte,
		filterBlockTimestampGt,
		filterBlockTimestampLte,
		filterBlockTimestampLt,
		filterTopic1,
		filterTopic2,
		filterTopic3,
		filterTopic0,
	), nil
}
