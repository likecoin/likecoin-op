package database

import (
	"context"
	"database/sql"
	"fmt"

	"likecollective-indexer/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

// stakingStateLockKey is the Postgres advisory lock key serialising every
// writer of accounts, nft_classes and stakings. Any constant works as long as
// nothing else in the database takes the same key.
const stakingStateLockKey int64 = 0x6c696b655f737473 // "like_sts"

// stakingStateLockTxKey carries the transaction holding the staking state lock
// on the context fn runs with.
type stakingStateLockTxKey struct{}

// WithStakingStateLock runs fn holding the advisory lock that every writer of
// the staking state takes around its load, chain read and persist.
//
// Holding it only around the persist would not be enough. Account totals are
// written as the loaded value moved by a correction, so two writers that loaded
// the same account before either committed would each overwrite the other's
// correction.
//
// The lock is transaction scoped, so Postgres releases it when the transaction
// ends or its connection drops, and a crashed writer cannot leave it held. fn
// does its reads on other connections, but its writes go through
// WithStakingStateTx onto the transaction holding the lock, and are committed
// only when fn returns nil. Were they committed on a connection of their own,
// losing the lock mid-pass -- the connection dropped, or the idle transaction
// terminated during a long chain read -- would go unnoticed, and the write
// could land after another writer's. On the lock's own transaction they fail
// with it instead. A caller must therefore return any error the writes
// returned, or the part written before it would be committed. Each waiter
// holds a connection while it waits, so a pool capped at or below the number
// of waiters would leave the holder unable to get one.
func WithStakingStateLock(
	ctx context.Context,
	dbService Service,
	fn func(ctx context.Context) error,
) error {
	tx, err := dbService.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin staking state lock transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock($1)", stakingStateLockKey); err != nil {
		return fmt.Errorf("failed to take staking state lock: %w", err)
	}

	if err := fn(context.WithValue(ctx, stakingStateLockTxKey{}, tx)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit staking state lock transaction: %w", err)
	}
	return nil
}

// WithStakingStateTx runs fn on the transaction holding the staking state lock
// if ctx is inside WithStakingStateLock, where it is committed with the lock,
// and on a transaction of its own otherwise.
func WithStakingStateTx(
	ctx context.Context,
	dbService Service,
	fn func(tx *ent.Tx) error,
) error {
	lockTx, ok := ctx.Value(stakingStateLockTxKey{}).(*sql.Tx)
	if !ok {
		return WithTx(ctx, dbService.Client(), fn)
	}
	// Built like the client in New, whose hooks and options it does not share.
	drv := &lockTxDriver{entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: lockTx})}
	client := ent.NewClient(ent.Driver(drv))
	if debug {
		client = client.Debug()
	}
	return WithTx(ctx, client, fn)
}

// lockTxDriver runs every statement on the lock's transaction, and leaves
// committing and rolling it back to WithStakingStateLock.
type lockTxDriver struct {
	*entsql.Driver
}

func (d *lockTxDriver) Tx(context.Context) (dialect.Tx, error) {
	return dialect.NopTx(d), nil
}
