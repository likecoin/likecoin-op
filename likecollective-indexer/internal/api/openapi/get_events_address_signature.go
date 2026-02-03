package openapi

import (
	"context"
	"math"

	"likecollective-indexer/internal/api/openapi/model"
	"likecollective-indexer/openapi/api"
)

func (h *openAPIHandler) EventsAddressSignatureGet(
	ctx context.Context,
	params api.EventsAddressSignatureGetParams,
) (*api.EventsAddressSignatureGetOK, error) {
	ps := model.OpenAPIEventParameters{
		Address:                 &params.Address,
		Signature:               &params.Signature,
		Limit:                   params.Limit,
		Page:                    params.Page,
		SortBy:                  params.SortBy,
		SortOrder:               params.SortOrder,
		FilterBlockTimestamp:    params.FilterBlockTimestamp,
		FilterBlockTimestampGte: params.FilterBlockTimestampGte,
		FilterBlockTimestampGt:  params.FilterBlockTimestampGt,
		FilterBlockTimestampLte: params.FilterBlockTimestampLte,
		FilterBlockTimestampLt:  params.FilterBlockTimestampLt,
		FilterTopic1:            params.FilterTopic1,
		FilterTopic2:            params.FilterTopic2,
		FilterTopic3:            params.FilterTopic3,
		FilterTopic0:            params.FilterTopic0,
	}
	dbFilter, err := ps.ToEntFilter()
	if err != nil {
		return nil, err
	}

	events, count, err := h.evmEventRepository.GetEvmEvents(ctx, dbFilter)
	if err != nil {
		return nil, err
	}

	apiEvents := make([]api.Event, len(events))

	for i, n := range events {
		apiEvent, err := model.MakeEvent(n)
		if err != nil {
			return nil, err
		}
		apiEvents[i] = apiEvent
	}

	return &api.EventsAddressSignatureGetOK{
		Meta: api.EventQueryMetadata{
			ChainIds:      []int{},
			Address:       params.Address,
			Signature:     params.Signature,
			Page:          params.Page.Value,
			LimitPerChain: params.Limit.Value,
			TotalItems:    count,
			TotalPages:    int(math.Ceil(float64(count) / float64(params.Limit.Value))),
		},
		Data: apiEvents,
	}, nil
}
