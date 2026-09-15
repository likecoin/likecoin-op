package database

import (
	"context"
	"errors"
	"fmt"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/schema/typeutil"
	"likecollective-indexer/ent/stakingstatehead"
)

// ErrStakingStateHeadBehind is returned when a write was read from a head
// older than one the staking state has already been committed at.
var ErrStakingStateHeadBehind = errors.New("staking state head is behind the committed head")

type StakingStateHeadRepository interface {
	// GetBlockNumber returns the newest head the staking state has been
	// committed at, or 0 if none has been recorded.
	GetBlockNumber(ctx context.Context) (uint64, error)

	// AdvanceBlockNumber records blockNumber as the committed head, failing
	// with ErrStakingStateHeadBehind if a newer one is already recorded.
	AdvanceBlockNumber(ctx context.Context, tx *ent.Tx, blockNumber uint64) error
}

type stakingStateHeadRepository struct {
	dbService Service
}

func MakeStakingStateHeadRepository(
	dbService Service,
) StakingStateHeadRepository {
	return &stakingStateHeadRepository{
		dbService: dbService,
	}
}

func (r *stakingStateHeadRepository) GetBlockNumber(ctx context.Context) (uint64, error) {
	head, err := r.dbService.Client().StakingStateHead.Query().
		Order(ent.Asc(stakingstatehead.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(head.BlockNumber), nil
}

func (r *stakingStateHeadRepository) AdvanceBlockNumber(
	ctx context.Context,
	tx *ent.Tx,
	blockNumber uint64,
) error {
	// The comparison and the write are one statement, but creating the first
	// row is not: two callers that both find no row would both insert one.
	// Callers hold the staking state lock, which rules that out.
	advanced, err := tx.StakingStateHead.Update().
		Where(stakingstatehead.BlockNumberLTE(typeutil.Uint64(blockNumber))).
		SetBlockNumber(typeutil.Uint64(blockNumber)).
		Save(ctx)
	if err != nil {
		return err
	}
	if advanced > 0 {
		return nil
	}

	head, err := tx.StakingStateHead.Query().
		Order(ent.Asc(stakingstatehead.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_, err = tx.StakingStateHead.Create().
				SetBlockNumber(typeutil.Uint64(blockNumber)).
				Save(ctx)
			return err
		}
		return err
	}
	return fmt.Errorf(
		"%w: read at %d, committed at %d",
		ErrStakingStateHeadBehind, blockNumber, uint64(head.BlockNumber),
	)
}
