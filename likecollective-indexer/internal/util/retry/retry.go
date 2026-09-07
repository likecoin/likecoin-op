package retry

import (
	"context"
	"math/rand/v2"
	"time"
)

// Defaults for bulk on-chain reads. Walking thousands of tokens is enough to
// trip a provider's rate limit however modest the concurrency, so callers retry
// rather than lose the whole run.
const (
	DefaultAttempts = 6
	DefaultBase     = 500 * time.Millisecond
)

// Value calls fn until it succeeds, up to attempts times, backing off
// base, 2*base, 4*base ... between tries, jittered. A cancelled context stops
// it early. The error returned is the one from the last attempt.
//
// It retries on any error rather than inspecting it: the callers are read-only
// RPC reads, where the failures worth surviving -- a provider's rate limit, a
// dropped connection -- are not reliably distinguishable from the response.
//
// The jitter matters more than the backoff curve. The failure this exists for
// is a rate limit, which rejects every in-flight caller within the same
// millisecond; sleeping them all for exactly the same interval would reproduce
// the burst that caused it, once per attempt.
func Value[R any](
	ctx context.Context,
	attempts int,
	base time.Duration,
	fn func(ctx context.Context) (R, error),
) (R, error) {
	var (
		result R
		err    error
	)

	if attempts < 1 {
		attempts = 1
	}

	backoff := base
	for attempt := range attempts {
		result, err = fn(ctx)
		if err == nil {
			return result, nil
		}
		if attempt == attempts-1 {
			break
		}

		select {
		case <-time.After(jitter(backoff)):
			backoff *= 2
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}

	return result, err
}

// jitter spreads a delay over [d/2, d), so concurrent retriers of the same
// rate limit stop waking together.
func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	half := d / 2
	return half + rand.N(half+1)
}
