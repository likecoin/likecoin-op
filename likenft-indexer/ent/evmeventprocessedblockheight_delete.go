// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"likenft-indexer/ent/evmeventprocessedblockheight"
	"likenft-indexer/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EVMEventProcessedBlockHeightDelete is the builder for deleting a EVMEventProcessedBlockHeight entity.
type EVMEventProcessedBlockHeightDelete struct {
	config
	hooks    []Hook
	mutation *EVMEventProcessedBlockHeightMutation
}

// Where appends a list predicates to the EVMEventProcessedBlockHeightDelete builder.
func (eepbhd *EVMEventProcessedBlockHeightDelete) Where(ps ...predicate.EVMEventProcessedBlockHeight) *EVMEventProcessedBlockHeightDelete {
	eepbhd.mutation.Where(ps...)
	return eepbhd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (eepbhd *EVMEventProcessedBlockHeightDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, eepbhd.sqlExec, eepbhd.mutation, eepbhd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (eepbhd *EVMEventProcessedBlockHeightDelete) ExecX(ctx context.Context) int {
	n, err := eepbhd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (eepbhd *EVMEventProcessedBlockHeightDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(evmeventprocessedblockheight.Table, sqlgraph.NewFieldSpec(evmeventprocessedblockheight.FieldID, field.TypeInt))
	if ps := eepbhd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, eepbhd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	eepbhd.mutation.done = true
	return affected, err
}

// EVMEventProcessedBlockHeightDeleteOne is the builder for deleting a single EVMEventProcessedBlockHeight entity.
type EVMEventProcessedBlockHeightDeleteOne struct {
	eepbhd *EVMEventProcessedBlockHeightDelete
}

// Where appends a list predicates to the EVMEventProcessedBlockHeightDelete builder.
func (eepbhdo *EVMEventProcessedBlockHeightDeleteOne) Where(ps ...predicate.EVMEventProcessedBlockHeight) *EVMEventProcessedBlockHeightDeleteOne {
	eepbhdo.eepbhd.mutation.Where(ps...)
	return eepbhdo
}

// Exec executes the deletion query.
func (eepbhdo *EVMEventProcessedBlockHeightDeleteOne) Exec(ctx context.Context) error {
	n, err := eepbhdo.eepbhd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{evmeventprocessedblockheight.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (eepbhdo *EVMEventProcessedBlockHeightDeleteOne) ExecX(ctx context.Context) {
	if err := eepbhdo.Exec(ctx); err != nil {
		panic(err)
	}
}
