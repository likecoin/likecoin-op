package parallel

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestMapWithLimitKeepsInputOrder(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}

	results, err := MapWithLimit(context.Background(), 3, items, func(ctx context.Context, item int) (int, error) {
		return item * 10, nil
	})
	if err != nil {
		t.Fatalf("MapWithLimit: %v", err)
	}
	for i, item := range items {
		if results[i] != item*10 {
			t.Fatalf("results[%d] = %d, want %d", i, results[i], item*10)
		}
	}
}

func TestMapWithLimitReturnsTheFirstError(t *testing.T) {
	wantErr := errors.New("boom")

	results, err := MapWithLimit(context.Background(), 2, []int{1, 2, 3, 4}, func(ctx context.Context, item int) (int, error) {
		if item == 2 {
			return 0, wantErr
		}
		return item, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
	if results != nil {
		t.Fatalf("results = %v, want nil on error", results)
	}
}

func TestMapWithLimitReturnsCancellationInsteadOfAPartialSlice(t *testing.T) {
	// The caller's context is cancelled mid-run. Whether the producer stops
	// short is up to the select, so what is pinned here is the part that
	// matters either way: a cancelled run reports the cancellation rather than
	// handing back a slice whose untouched entries are at their zero value.
	ctx, cancel := context.WithCancel(context.Background())

	var once sync.Once
	results, err := MapWithLimit(ctx, 1, []int{1, 2, 3, 4}, func(ctx context.Context, item int) (int, error) {
		once.Do(cancel)
		return item, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want %v", err, context.Canceled)
	}
	if results != nil {
		t.Fatalf("results = %v, want nil on cancellation", results)
	}
}

func TestMapWithLimitStartsNoMoreWorkersThanItems(t *testing.T) {
	// limit comes from a cli flag. Starting that many goroutines regardless of
	// the items would exhaust memory before any work began.
	results, err := MapWithLimit(context.Background(), 1<<40, []int{1, 2}, func(ctx context.Context, item int) (int, error) {
		return item, nil
	})
	if err != nil {
		t.Fatalf("MapWithLimit: %v", err)
	}
	if len(results) != 2 || results[0] != 1 || results[1] != 2 {
		t.Fatalf("results = %v, want [1 2]", results)
	}

	empty, err := MapWithLimit(context.Background(), 1<<40, []int{}, func(ctx context.Context, item int) (int, error) {
		return item, nil
	})
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty = %v, %v; want no results and no error", empty, err)
	}
}
