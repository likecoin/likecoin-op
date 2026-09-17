package model

import (
	"slices"
	"testing"

	"likecollective-indexer/ent"
)

func TestEvmEventsProcessingComparatorOrdersLogsWithinATransaction(t *testing.T) {
	events := []*ent.EVMEvent{
		{ID: 3, BlockNumber: 10, TransactionIndex: 1, LogIndex: 7, Topic0: "0xb"},
		{ID: 4, BlockNumber: 11, TransactionIndex: 0, LogIndex: 0, Topic0: "0xa"},
		{ID: 2, BlockNumber: 10, TransactionIndex: 1, LogIndex: 5, Topic0: "0xa"},
		{ID: 1, BlockNumber: 10, TransactionIndex: 0, LogIndex: 9, Topic0: "0xc"},
	}

	slices.SortFunc(events, EvmEventsProcessingComparator)

	for i, event := range events {
		if event.ID != i+1 {
			t.Fatalf("position %d holds event %d, want %d", i, event.ID, i+1)
		}
	}
	if EvmEventsProcessingComparator(events[0], events[0]) != 0 {
		t.Fatal("an event does not compare equal to itself")
	}
}
