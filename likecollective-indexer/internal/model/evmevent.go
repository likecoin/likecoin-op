package model

import (
	"likecollective-indexer/ent"
	"likecollective-indexer/internal/util/ordered"
)

var cmpBlockNumber ordered.Comparator[*ent.EVMEvent] = func(a, b *ent.EVMEvent) int {
	return ordered.Normalize(a.BlockNumber, b.BlockNumber)
}
var cmpTransactionIndex ordered.Comparator[*ent.EVMEvent] = func(a, b *ent.EVMEvent) int {
	return ordered.Normalize(a.TransactionIndex, b.TransactionIndex)
}

var topic0OrderMap map[string]int = map[string]int{}

// cmpTopic0 puts ranked topics before unranked ones. Two unranked topics are
// equal, so they fall through to the log index rather than comparing as
// greater both ways.
var cmpTopic0 ordered.Comparator[*ent.EVMEvent] = func(a, b *ent.EVMEvent) int {
	orderA, rankedA := topic0OrderMap[a.Topic0]
	orderB, rankedB := topic0OrderMap[b.Topic0]
	switch {
	case rankedA && rankedB:
		return ordered.Normalize(orderA, orderB)
	case rankedA:
		return -1
	case rankedB:
		return 1
	}
	return 0
}

var cmpLogIndex ordered.Comparator[*ent.EVMEvent] = func(a, b *ent.EVMEvent) int {
	return ordered.Normalize(a.LogIndex, b.LogIndex)
}

var EvmEventsProcessingComparator = ordered.CombineComparators(
	cmpBlockNumber,
	cmpTransactionIndex,
	cmpTopic0,
	cmpLogIndex,
)
