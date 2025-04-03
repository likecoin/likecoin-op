// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"likenft-indexer/ent/predicate"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/ent/transactionmemo"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TransactionMemoUpdate is the builder for updating TransactionMemo entities.
type TransactionMemoUpdate struct {
	config
	hooks    []Hook
	mutation *TransactionMemoMutation
}

// Where appends a list predicates to the TransactionMemoUpdate builder.
func (tmu *TransactionMemoUpdate) Where(ps ...predicate.TransactionMemo) *TransactionMemoUpdate {
	tmu.mutation.Where(ps...)
	return tmu
}

// SetTransactionHash sets the "transaction_hash" field.
func (tmu *TransactionMemoUpdate) SetTransactionHash(s string) *TransactionMemoUpdate {
	tmu.mutation.SetTransactionHash(s)
	return tmu
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableTransactionHash(s *string) *TransactionMemoUpdate {
	if s != nil {
		tmu.SetTransactionHash(*s)
	}
	return tmu
}

// SetBookNftID sets the "book_nft_id" field.
func (tmu *TransactionMemoUpdate) SetBookNftID(s string) *TransactionMemoUpdate {
	tmu.mutation.SetBookNftID(s)
	return tmu
}

// SetNillableBookNftID sets the "book_nft_id" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableBookNftID(s *string) *TransactionMemoUpdate {
	if s != nil {
		tmu.SetBookNftID(*s)
	}
	return tmu
}

// SetFrom sets the "from" field.
func (tmu *TransactionMemoUpdate) SetFrom(s string) *TransactionMemoUpdate {
	tmu.mutation.SetFrom(s)
	return tmu
}

// SetNillableFrom sets the "from" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableFrom(s *string) *TransactionMemoUpdate {
	if s != nil {
		tmu.SetFrom(*s)
	}
	return tmu
}

// SetTo sets the "to" field.
func (tmu *TransactionMemoUpdate) SetTo(s string) *TransactionMemoUpdate {
	tmu.mutation.SetTo(s)
	return tmu
}

// SetNillableTo sets the "to" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableTo(s *string) *TransactionMemoUpdate {
	if s != nil {
		tmu.SetTo(*s)
	}
	return tmu
}

// SetTokenID sets the "token_id" field.
func (tmu *TransactionMemoUpdate) SetTokenID(t typeutil.Uint64) *TransactionMemoUpdate {
	tmu.mutation.ResetTokenID()
	tmu.mutation.SetTokenID(t)
	return tmu
}

// SetNillableTokenID sets the "token_id" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableTokenID(t *typeutil.Uint64) *TransactionMemoUpdate {
	if t != nil {
		tmu.SetTokenID(*t)
	}
	return tmu
}

// AddTokenID adds t to the "token_id" field.
func (tmu *TransactionMemoUpdate) AddTokenID(t typeutil.Uint64) *TransactionMemoUpdate {
	tmu.mutation.AddTokenID(t)
	return tmu
}

// SetMemo sets the "memo" field.
func (tmu *TransactionMemoUpdate) SetMemo(s string) *TransactionMemoUpdate {
	tmu.mutation.SetMemo(s)
	return tmu
}

// SetNillableMemo sets the "memo" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableMemo(s *string) *TransactionMemoUpdate {
	if s != nil {
		tmu.SetMemo(*s)
	}
	return tmu
}

// SetBlockNumber sets the "block_number" field.
func (tmu *TransactionMemoUpdate) SetBlockNumber(t typeutil.Uint64) *TransactionMemoUpdate {
	tmu.mutation.ResetBlockNumber()
	tmu.mutation.SetBlockNumber(t)
	return tmu
}

// SetNillableBlockNumber sets the "block_number" field if the given value is not nil.
func (tmu *TransactionMemoUpdate) SetNillableBlockNumber(t *typeutil.Uint64) *TransactionMemoUpdate {
	if t != nil {
		tmu.SetBlockNumber(*t)
	}
	return tmu
}

// AddBlockNumber adds t to the "block_number" field.
func (tmu *TransactionMemoUpdate) AddBlockNumber(t typeutil.Uint64) *TransactionMemoUpdate {
	tmu.mutation.AddBlockNumber(t)
	return tmu
}

// Mutation returns the TransactionMemoMutation object of the builder.
func (tmu *TransactionMemoUpdate) Mutation() *TransactionMemoMutation {
	return tmu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tmu *TransactionMemoUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tmu.sqlSave, tmu.mutation, tmu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tmu *TransactionMemoUpdate) SaveX(ctx context.Context) int {
	affected, err := tmu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tmu *TransactionMemoUpdate) Exec(ctx context.Context) error {
	_, err := tmu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tmu *TransactionMemoUpdate) ExecX(ctx context.Context) {
	if err := tmu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tmu *TransactionMemoUpdate) check() error {
	if v, ok := tmu.mutation.TransactionHash(); ok {
		if err := transactionmemo.TransactionHashValidator(v); err != nil {
			return &ValidationError{Name: "transaction_hash", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.transaction_hash": %w`, err)}
		}
	}
	if v, ok := tmu.mutation.BookNftID(); ok {
		if err := transactionmemo.BookNftIDValidator(v); err != nil {
			return &ValidationError{Name: "book_nft_id", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.book_nft_id": %w`, err)}
		}
	}
	if v, ok := tmu.mutation.From(); ok {
		if err := transactionmemo.FromValidator(v); err != nil {
			return &ValidationError{Name: "from", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.from": %w`, err)}
		}
	}
	if v, ok := tmu.mutation.To(); ok {
		if err := transactionmemo.ToValidator(v); err != nil {
			return &ValidationError{Name: "to", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.to": %w`, err)}
		}
	}
	return nil
}

func (tmu *TransactionMemoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tmu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(transactionmemo.Table, transactionmemo.Columns, sqlgraph.NewFieldSpec(transactionmemo.FieldID, field.TypeInt))
	if ps := tmu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tmu.mutation.TransactionHash(); ok {
		_spec.SetField(transactionmemo.FieldTransactionHash, field.TypeString, value)
	}
	if value, ok := tmu.mutation.BookNftID(); ok {
		_spec.SetField(transactionmemo.FieldBookNftID, field.TypeString, value)
	}
	if value, ok := tmu.mutation.From(); ok {
		_spec.SetField(transactionmemo.FieldFrom, field.TypeString, value)
	}
	if value, ok := tmu.mutation.To(); ok {
		_spec.SetField(transactionmemo.FieldTo, field.TypeString, value)
	}
	if value, ok := tmu.mutation.TokenID(); ok {
		vv, err := transactionmemo.ValueScanner.TokenID.Value(value)
		if err != nil {
			return 0, err
		}
		_spec.SetField(transactionmemo.FieldTokenID, field.TypeUint64, vv)
	}
	if value, ok := tmu.mutation.AddedTokenID(); ok {
		vv, err := transactionmemo.ValueScanner.TokenID.Value(value)
		if err != nil {
			return 0, err
		}
		_spec.AddField(transactionmemo.FieldTokenID, field.TypeUint64, vv)
	}
	if value, ok := tmu.mutation.Memo(); ok {
		_spec.SetField(transactionmemo.FieldMemo, field.TypeString, value)
	}
	if value, ok := tmu.mutation.BlockNumber(); ok {
		vv, err := transactionmemo.ValueScanner.BlockNumber.Value(value)
		if err != nil {
			return 0, err
		}
		_spec.SetField(transactionmemo.FieldBlockNumber, field.TypeUint64, vv)
	}
	if value, ok := tmu.mutation.AddedBlockNumber(); ok {
		vv, err := transactionmemo.ValueScanner.BlockNumber.Value(value)
		if err != nil {
			return 0, err
		}
		_spec.AddField(transactionmemo.FieldBlockNumber, field.TypeUint64, vv)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tmu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactionmemo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tmu.mutation.done = true
	return n, nil
}

// TransactionMemoUpdateOne is the builder for updating a single TransactionMemo entity.
type TransactionMemoUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TransactionMemoMutation
}

// SetTransactionHash sets the "transaction_hash" field.
func (tmuo *TransactionMemoUpdateOne) SetTransactionHash(s string) *TransactionMemoUpdateOne {
	tmuo.mutation.SetTransactionHash(s)
	return tmuo
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableTransactionHash(s *string) *TransactionMemoUpdateOne {
	if s != nil {
		tmuo.SetTransactionHash(*s)
	}
	return tmuo
}

// SetBookNftID sets the "book_nft_id" field.
func (tmuo *TransactionMemoUpdateOne) SetBookNftID(s string) *TransactionMemoUpdateOne {
	tmuo.mutation.SetBookNftID(s)
	return tmuo
}

// SetNillableBookNftID sets the "book_nft_id" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableBookNftID(s *string) *TransactionMemoUpdateOne {
	if s != nil {
		tmuo.SetBookNftID(*s)
	}
	return tmuo
}

// SetFrom sets the "from" field.
func (tmuo *TransactionMemoUpdateOne) SetFrom(s string) *TransactionMemoUpdateOne {
	tmuo.mutation.SetFrom(s)
	return tmuo
}

// SetNillableFrom sets the "from" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableFrom(s *string) *TransactionMemoUpdateOne {
	if s != nil {
		tmuo.SetFrom(*s)
	}
	return tmuo
}

// SetTo sets the "to" field.
func (tmuo *TransactionMemoUpdateOne) SetTo(s string) *TransactionMemoUpdateOne {
	tmuo.mutation.SetTo(s)
	return tmuo
}

// SetNillableTo sets the "to" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableTo(s *string) *TransactionMemoUpdateOne {
	if s != nil {
		tmuo.SetTo(*s)
	}
	return tmuo
}

// SetTokenID sets the "token_id" field.
func (tmuo *TransactionMemoUpdateOne) SetTokenID(t typeutil.Uint64) *TransactionMemoUpdateOne {
	tmuo.mutation.ResetTokenID()
	tmuo.mutation.SetTokenID(t)
	return tmuo
}

// SetNillableTokenID sets the "token_id" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableTokenID(t *typeutil.Uint64) *TransactionMemoUpdateOne {
	if t != nil {
		tmuo.SetTokenID(*t)
	}
	return tmuo
}

// AddTokenID adds t to the "token_id" field.
func (tmuo *TransactionMemoUpdateOne) AddTokenID(t typeutil.Uint64) *TransactionMemoUpdateOne {
	tmuo.mutation.AddTokenID(t)
	return tmuo
}

// SetMemo sets the "memo" field.
func (tmuo *TransactionMemoUpdateOne) SetMemo(s string) *TransactionMemoUpdateOne {
	tmuo.mutation.SetMemo(s)
	return tmuo
}

// SetNillableMemo sets the "memo" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableMemo(s *string) *TransactionMemoUpdateOne {
	if s != nil {
		tmuo.SetMemo(*s)
	}
	return tmuo
}

// SetBlockNumber sets the "block_number" field.
func (tmuo *TransactionMemoUpdateOne) SetBlockNumber(t typeutil.Uint64) *TransactionMemoUpdateOne {
	tmuo.mutation.ResetBlockNumber()
	tmuo.mutation.SetBlockNumber(t)
	return tmuo
}

// SetNillableBlockNumber sets the "block_number" field if the given value is not nil.
func (tmuo *TransactionMemoUpdateOne) SetNillableBlockNumber(t *typeutil.Uint64) *TransactionMemoUpdateOne {
	if t != nil {
		tmuo.SetBlockNumber(*t)
	}
	return tmuo
}

// AddBlockNumber adds t to the "block_number" field.
func (tmuo *TransactionMemoUpdateOne) AddBlockNumber(t typeutil.Uint64) *TransactionMemoUpdateOne {
	tmuo.mutation.AddBlockNumber(t)
	return tmuo
}

// Mutation returns the TransactionMemoMutation object of the builder.
func (tmuo *TransactionMemoUpdateOne) Mutation() *TransactionMemoMutation {
	return tmuo.mutation
}

// Where appends a list predicates to the TransactionMemoUpdate builder.
func (tmuo *TransactionMemoUpdateOne) Where(ps ...predicate.TransactionMemo) *TransactionMemoUpdateOne {
	tmuo.mutation.Where(ps...)
	return tmuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tmuo *TransactionMemoUpdateOne) Select(field string, fields ...string) *TransactionMemoUpdateOne {
	tmuo.fields = append([]string{field}, fields...)
	return tmuo
}

// Save executes the query and returns the updated TransactionMemo entity.
func (tmuo *TransactionMemoUpdateOne) Save(ctx context.Context) (*TransactionMemo, error) {
	return withHooks(ctx, tmuo.sqlSave, tmuo.mutation, tmuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tmuo *TransactionMemoUpdateOne) SaveX(ctx context.Context) *TransactionMemo {
	node, err := tmuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tmuo *TransactionMemoUpdateOne) Exec(ctx context.Context) error {
	_, err := tmuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tmuo *TransactionMemoUpdateOne) ExecX(ctx context.Context) {
	if err := tmuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tmuo *TransactionMemoUpdateOne) check() error {
	if v, ok := tmuo.mutation.TransactionHash(); ok {
		if err := transactionmemo.TransactionHashValidator(v); err != nil {
			return &ValidationError{Name: "transaction_hash", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.transaction_hash": %w`, err)}
		}
	}
	if v, ok := tmuo.mutation.BookNftID(); ok {
		if err := transactionmemo.BookNftIDValidator(v); err != nil {
			return &ValidationError{Name: "book_nft_id", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.book_nft_id": %w`, err)}
		}
	}
	if v, ok := tmuo.mutation.From(); ok {
		if err := transactionmemo.FromValidator(v); err != nil {
			return &ValidationError{Name: "from", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.from": %w`, err)}
		}
	}
	if v, ok := tmuo.mutation.To(); ok {
		if err := transactionmemo.ToValidator(v); err != nil {
			return &ValidationError{Name: "to", err: fmt.Errorf(`ent: validator failed for field "TransactionMemo.to": %w`, err)}
		}
	}
	return nil
}

func (tmuo *TransactionMemoUpdateOne) sqlSave(ctx context.Context) (_node *TransactionMemo, err error) {
	if err := tmuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(transactionmemo.Table, transactionmemo.Columns, sqlgraph.NewFieldSpec(transactionmemo.FieldID, field.TypeInt))
	id, ok := tmuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TransactionMemo.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tmuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, transactionmemo.FieldID)
		for _, f := range fields {
			if !transactionmemo.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != transactionmemo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tmuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tmuo.mutation.TransactionHash(); ok {
		_spec.SetField(transactionmemo.FieldTransactionHash, field.TypeString, value)
	}
	if value, ok := tmuo.mutation.BookNftID(); ok {
		_spec.SetField(transactionmemo.FieldBookNftID, field.TypeString, value)
	}
	if value, ok := tmuo.mutation.From(); ok {
		_spec.SetField(transactionmemo.FieldFrom, field.TypeString, value)
	}
	if value, ok := tmuo.mutation.To(); ok {
		_spec.SetField(transactionmemo.FieldTo, field.TypeString, value)
	}
	if value, ok := tmuo.mutation.TokenID(); ok {
		vv, err := transactionmemo.ValueScanner.TokenID.Value(value)
		if err != nil {
			return nil, err
		}
		_spec.SetField(transactionmemo.FieldTokenID, field.TypeUint64, vv)
	}
	if value, ok := tmuo.mutation.AddedTokenID(); ok {
		vv, err := transactionmemo.ValueScanner.TokenID.Value(value)
		if err != nil {
			return nil, err
		}
		_spec.AddField(transactionmemo.FieldTokenID, field.TypeUint64, vv)
	}
	if value, ok := tmuo.mutation.Memo(); ok {
		_spec.SetField(transactionmemo.FieldMemo, field.TypeString, value)
	}
	if value, ok := tmuo.mutation.BlockNumber(); ok {
		vv, err := transactionmemo.ValueScanner.BlockNumber.Value(value)
		if err != nil {
			return nil, err
		}
		_spec.SetField(transactionmemo.FieldBlockNumber, field.TypeUint64, vv)
	}
	if value, ok := tmuo.mutation.AddedBlockNumber(); ok {
		vv, err := transactionmemo.ValueScanner.BlockNumber.Value(value)
		if err != nil {
			return nil, err
		}
		_spec.AddField(transactionmemo.FieldBlockNumber, field.TypeUint64, vv)
	}
	_node = &TransactionMemo{config: tmuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tmuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactionmemo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tmuo.mutation.done = true
	return _node, nil
}
