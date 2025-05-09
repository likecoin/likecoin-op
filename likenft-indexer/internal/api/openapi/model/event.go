package model

import (
	"strconv"

	"likenft-indexer/ent"
	"likenft-indexer/openapi/api"
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
