package parallel

import (
	"context"
	"sync"
)

// MapWithLimit applies f to every item, running at most limit calls at a time,
// and returns the results in input order. The first error cancels the context
// handed to the remaining calls and is the one returned.
func MapWithLimit[T any, R any](
	ctx context.Context,
	limit int,
	items []T,
	f func(ctx context.Context, item T) (R, error),
) ([]R, error) {
	if limit <= 0 {
		limit = 1
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]R, len(items))

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)

	indexes := make(chan int)

	for range limit {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range indexes {
				result, err := f(ctx, items[index])
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						cancel()
					}
					mu.Unlock()
					return
				}
				results[index] = result
			}
		}()
	}

producing:
	for index := range items {
		select {
		case indexes <- index:
		case <-ctx.Done():
			break producing
		}
	}
	close(indexes)
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}
