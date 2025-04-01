// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"likenft-indexer/ent/evmevent"
	"likenft-indexer/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EVMEventDelete is the builder for deleting a EVMEvent entity.
type EVMEventDelete struct {
	config
	hooks    []Hook
	mutation *EVMEventMutation
}

// Where appends a list predicates to the EVMEventDelete builder.
func (eed *EVMEventDelete) Where(ps ...predicate.EVMEvent) *EVMEventDelete {
	eed.mutation.Where(ps...)
	return eed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (eed *EVMEventDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, eed.sqlExec, eed.mutation, eed.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (eed *EVMEventDelete) ExecX(ctx context.Context) int {
	n, err := eed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (eed *EVMEventDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(evmevent.Table, sqlgraph.NewFieldSpec(evmevent.FieldID, field.TypeInt))
	if ps := eed.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, eed.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	eed.mutation.done = true
	return affected, err
}

// EVMEventDeleteOne is the builder for deleting a single EVMEvent entity.
type EVMEventDeleteOne struct {
	eed *EVMEventDelete
}

// Where appends a list predicates to the EVMEventDelete builder.
func (eedo *EVMEventDeleteOne) Where(ps ...predicate.EVMEvent) *EVMEventDeleteOne {
	eedo.eed.mutation.Where(ps...)
	return eedo
}

// Exec executes the deletion query.
func (eedo *EVMEventDeleteOne) Exec(ctx context.Context) error {
	n, err := eedo.eed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{evmevent.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (eedo *EVMEventDeleteOne) ExecX(ctx context.Context) {
	if err := eedo.Exec(ctx); err != nil {
		panic(err)
	}
}
