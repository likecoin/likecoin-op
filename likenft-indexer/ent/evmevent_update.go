// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"likenft-indexer/ent/evmevent"
	"likenft-indexer/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EVMEventUpdate is the builder for updating EVMEvent entities.
type EVMEventUpdate struct {
	config
	hooks    []Hook
	mutation *EVMEventMutation
}

// Where appends a list predicates to the EVMEventUpdate builder.
func (eeu *EVMEventUpdate) Where(ps ...predicate.EVMEvent) *EVMEventUpdate {
	eeu.mutation.Where(ps...)
	return eeu
}

// SetTransactionHash sets the "transaction_hash" field.
func (eeu *EVMEventUpdate) SetTransactionHash(s string) *EVMEventUpdate {
	eeu.mutation.SetTransactionHash(s)
	return eeu
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTransactionHash(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTransactionHash(*s)
	}
	return eeu
}

// SetTransactionIndex sets the "transaction_index" field.
func (eeu *EVMEventUpdate) SetTransactionIndex(u uint) *EVMEventUpdate {
	eeu.mutation.ResetTransactionIndex()
	eeu.mutation.SetTransactionIndex(u)
	return eeu
}

// SetNillableTransactionIndex sets the "transaction_index" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTransactionIndex(u *uint) *EVMEventUpdate {
	if u != nil {
		eeu.SetTransactionIndex(*u)
	}
	return eeu
}

// AddTransactionIndex adds u to the "transaction_index" field.
func (eeu *EVMEventUpdate) AddTransactionIndex(u int) *EVMEventUpdate {
	eeu.mutation.AddTransactionIndex(u)
	return eeu
}

// SetBlockHash sets the "block_hash" field.
func (eeu *EVMEventUpdate) SetBlockHash(s string) *EVMEventUpdate {
	eeu.mutation.SetBlockHash(s)
	return eeu
}

// SetNillableBlockHash sets the "block_hash" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableBlockHash(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetBlockHash(*s)
	}
	return eeu
}

// SetBlockNumber sets the "block_number" field.
func (eeu *EVMEventUpdate) SetBlockNumber(u uint64) *EVMEventUpdate {
	eeu.mutation.ResetBlockNumber()
	eeu.mutation.SetBlockNumber(u)
	return eeu
}

// SetNillableBlockNumber sets the "block_number" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableBlockNumber(u *uint64) *EVMEventUpdate {
	if u != nil {
		eeu.SetBlockNumber(*u)
	}
	return eeu
}

// AddBlockNumber adds u to the "block_number" field.
func (eeu *EVMEventUpdate) AddBlockNumber(u int64) *EVMEventUpdate {
	eeu.mutation.AddBlockNumber(u)
	return eeu
}

// SetLogIndex sets the "log_index" field.
func (eeu *EVMEventUpdate) SetLogIndex(u uint) *EVMEventUpdate {
	eeu.mutation.ResetLogIndex()
	eeu.mutation.SetLogIndex(u)
	return eeu
}

// SetNillableLogIndex sets the "log_index" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableLogIndex(u *uint) *EVMEventUpdate {
	if u != nil {
		eeu.SetLogIndex(*u)
	}
	return eeu
}

// AddLogIndex adds u to the "log_index" field.
func (eeu *EVMEventUpdate) AddLogIndex(u int) *EVMEventUpdate {
	eeu.mutation.AddLogIndex(u)
	return eeu
}

// SetAddress sets the "address" field.
func (eeu *EVMEventUpdate) SetAddress(s string) *EVMEventUpdate {
	eeu.mutation.SetAddress(s)
	return eeu
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableAddress(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetAddress(*s)
	}
	return eeu
}

// SetTopic0 sets the "topic0" field.
func (eeu *EVMEventUpdate) SetTopic0(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic0(s)
	return eeu
}

// SetNillableTopic0 sets the "topic0" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic0(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic0(*s)
	}
	return eeu
}

// SetTopic0Hex sets the "topic0_hex" field.
func (eeu *EVMEventUpdate) SetTopic0Hex(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic0Hex(s)
	return eeu
}

// SetNillableTopic0Hex sets the "topic0_hex" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic0Hex(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic0Hex(*s)
	}
	return eeu
}

// SetTopic1 sets the "topic1" field.
func (eeu *EVMEventUpdate) SetTopic1(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic1(s)
	return eeu
}

// SetNillableTopic1 sets the "topic1" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic1(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic1(*s)
	}
	return eeu
}

// ClearTopic1 clears the value of the "topic1" field.
func (eeu *EVMEventUpdate) ClearTopic1() *EVMEventUpdate {
	eeu.mutation.ClearTopic1()
	return eeu
}

// SetTopic1Hex sets the "topic1_hex" field.
func (eeu *EVMEventUpdate) SetTopic1Hex(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic1Hex(s)
	return eeu
}

// SetNillableTopic1Hex sets the "topic1_hex" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic1Hex(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic1Hex(*s)
	}
	return eeu
}

// ClearTopic1Hex clears the value of the "topic1_hex" field.
func (eeu *EVMEventUpdate) ClearTopic1Hex() *EVMEventUpdate {
	eeu.mutation.ClearTopic1Hex()
	return eeu
}

// SetTopic2 sets the "topic2" field.
func (eeu *EVMEventUpdate) SetTopic2(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic2(s)
	return eeu
}

// SetNillableTopic2 sets the "topic2" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic2(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic2(*s)
	}
	return eeu
}

// ClearTopic2 clears the value of the "topic2" field.
func (eeu *EVMEventUpdate) ClearTopic2() *EVMEventUpdate {
	eeu.mutation.ClearTopic2()
	return eeu
}

// SetTopic2Hex sets the "topic2_hex" field.
func (eeu *EVMEventUpdate) SetTopic2Hex(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic2Hex(s)
	return eeu
}

// SetNillableTopic2Hex sets the "topic2_hex" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic2Hex(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic2Hex(*s)
	}
	return eeu
}

// ClearTopic2Hex clears the value of the "topic2_hex" field.
func (eeu *EVMEventUpdate) ClearTopic2Hex() *EVMEventUpdate {
	eeu.mutation.ClearTopic2Hex()
	return eeu
}

// SetTopic3 sets the "topic3" field.
func (eeu *EVMEventUpdate) SetTopic3(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic3(s)
	return eeu
}

// SetNillableTopic3 sets the "topic3" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic3(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic3(*s)
	}
	return eeu
}

// ClearTopic3 clears the value of the "topic3" field.
func (eeu *EVMEventUpdate) ClearTopic3() *EVMEventUpdate {
	eeu.mutation.ClearTopic3()
	return eeu
}

// SetTopic3Hex sets the "topic3_hex" field.
func (eeu *EVMEventUpdate) SetTopic3Hex(s string) *EVMEventUpdate {
	eeu.mutation.SetTopic3Hex(s)
	return eeu
}

// SetNillableTopic3Hex sets the "topic3_hex" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTopic3Hex(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetTopic3Hex(*s)
	}
	return eeu
}

// ClearTopic3Hex clears the value of the "topic3_hex" field.
func (eeu *EVMEventUpdate) ClearTopic3Hex() *EVMEventUpdate {
	eeu.mutation.ClearTopic3Hex()
	return eeu
}

// SetData sets the "data" field.
func (eeu *EVMEventUpdate) SetData(s string) *EVMEventUpdate {
	eeu.mutation.SetData(s)
	return eeu
}

// SetNillableData sets the "data" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableData(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetData(*s)
	}
	return eeu
}

// ClearData clears the value of the "data" field.
func (eeu *EVMEventUpdate) ClearData() *EVMEventUpdate {
	eeu.mutation.ClearData()
	return eeu
}

// SetDataHex sets the "data_hex" field.
func (eeu *EVMEventUpdate) SetDataHex(s string) *EVMEventUpdate {
	eeu.mutation.SetDataHex(s)
	return eeu
}

// SetNillableDataHex sets the "data_hex" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableDataHex(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetDataHex(*s)
	}
	return eeu
}

// ClearDataHex clears the value of the "data_hex" field.
func (eeu *EVMEventUpdate) ClearDataHex() *EVMEventUpdate {
	eeu.mutation.ClearDataHex()
	return eeu
}

// SetRemoved sets the "removed" field.
func (eeu *EVMEventUpdate) SetRemoved(b bool) *EVMEventUpdate {
	eeu.mutation.SetRemoved(b)
	return eeu
}

// SetNillableRemoved sets the "removed" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableRemoved(b *bool) *EVMEventUpdate {
	if b != nil {
		eeu.SetRemoved(*b)
	}
	return eeu
}

// SetStatus sets the "status" field.
func (eeu *EVMEventUpdate) SetStatus(e evmevent.Status) *EVMEventUpdate {
	eeu.mutation.SetStatus(e)
	return eeu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableStatus(e *evmevent.Status) *EVMEventUpdate {
	if e != nil {
		eeu.SetStatus(*e)
	}
	return eeu
}

// SetFailedReason sets the "failed_reason" field.
func (eeu *EVMEventUpdate) SetFailedReason(s string) *EVMEventUpdate {
	eeu.mutation.SetFailedReason(s)
	return eeu
}

// SetNillableFailedReason sets the "failed_reason" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableFailedReason(s *string) *EVMEventUpdate {
	if s != nil {
		eeu.SetFailedReason(*s)
	}
	return eeu
}

// ClearFailedReason clears the value of the "failed_reason" field.
func (eeu *EVMEventUpdate) ClearFailedReason() *EVMEventUpdate {
	eeu.mutation.ClearFailedReason()
	return eeu
}

// SetTimestamp sets the "timestamp" field.
func (eeu *EVMEventUpdate) SetTimestamp(t time.Time) *EVMEventUpdate {
	eeu.mutation.SetTimestamp(t)
	return eeu
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (eeu *EVMEventUpdate) SetNillableTimestamp(t *time.Time) *EVMEventUpdate {
	if t != nil {
		eeu.SetTimestamp(*t)
	}
	return eeu
}

// Mutation returns the EVMEventMutation object of the builder.
func (eeu *EVMEventUpdate) Mutation() *EVMEventMutation {
	return eeu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eeu *EVMEventUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, eeu.sqlSave, eeu.mutation, eeu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eeu *EVMEventUpdate) SaveX(ctx context.Context) int {
	affected, err := eeu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eeu *EVMEventUpdate) Exec(ctx context.Context) error {
	_, err := eeu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eeu *EVMEventUpdate) ExecX(ctx context.Context) {
	if err := eeu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eeu *EVMEventUpdate) check() error {
	if v, ok := eeu.mutation.TransactionHash(); ok {
		if err := evmevent.TransactionHashValidator(v); err != nil {
			return &ValidationError{Name: "transaction_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.transaction_hash": %w`, err)}
		}
	}
	if v, ok := eeu.mutation.BlockHash(); ok {
		if err := evmevent.BlockHashValidator(v); err != nil {
			return &ValidationError{Name: "block_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.block_hash": %w`, err)}
		}
	}
	if v, ok := eeu.mutation.Address(); ok {
		if err := evmevent.AddressValidator(v); err != nil {
			return &ValidationError{Name: "address", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.address": %w`, err)}
		}
	}
	if v, ok := eeu.mutation.Topic0(); ok {
		if err := evmevent.Topic0Validator(v); err != nil {
			return &ValidationError{Name: "topic0", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0": %w`, err)}
		}
	}
	if v, ok := eeu.mutation.Topic0Hex(); ok {
		if err := evmevent.Topic0HexValidator(v); err != nil {
			return &ValidationError{Name: "topic0_hex", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0_hex": %w`, err)}
		}
	}
	if v, ok := eeu.mutation.Status(); ok {
		if err := evmevent.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.status": %w`, err)}
		}
	}
	return nil
}

func (eeu *EVMEventUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := eeu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(evmevent.Table, evmevent.Columns, sqlgraph.NewFieldSpec(evmevent.FieldID, field.TypeInt))
	if ps := eeu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eeu.mutation.TransactionHash(); ok {
		_spec.SetField(evmevent.FieldTransactionHash, field.TypeString, value)
	}
	if value, ok := eeu.mutation.TransactionIndex(); ok {
		_spec.SetField(evmevent.FieldTransactionIndex, field.TypeUint, value)
	}
	if value, ok := eeu.mutation.AddedTransactionIndex(); ok {
		_spec.AddField(evmevent.FieldTransactionIndex, field.TypeUint, value)
	}
	if value, ok := eeu.mutation.BlockHash(); ok {
		_spec.SetField(evmevent.FieldBlockHash, field.TypeString, value)
	}
	if value, ok := eeu.mutation.BlockNumber(); ok {
		_spec.SetField(evmevent.FieldBlockNumber, field.TypeUint64, value)
	}
	if value, ok := eeu.mutation.AddedBlockNumber(); ok {
		_spec.AddField(evmevent.FieldBlockNumber, field.TypeUint64, value)
	}
	if value, ok := eeu.mutation.LogIndex(); ok {
		_spec.SetField(evmevent.FieldLogIndex, field.TypeUint, value)
	}
	if value, ok := eeu.mutation.AddedLogIndex(); ok {
		_spec.AddField(evmevent.FieldLogIndex, field.TypeUint, value)
	}
	if value, ok := eeu.mutation.Address(); ok {
		_spec.SetField(evmevent.FieldAddress, field.TypeString, value)
	}
	if value, ok := eeu.mutation.Topic0(); ok {
		_spec.SetField(evmevent.FieldTopic0, field.TypeString, value)
	}
	if value, ok := eeu.mutation.Topic0Hex(); ok {
		_spec.SetField(evmevent.FieldTopic0Hex, field.TypeString, value)
	}
	if value, ok := eeu.mutation.Topic1(); ok {
		_spec.SetField(evmevent.FieldTopic1, field.TypeString, value)
	}
	if eeu.mutation.Topic1Cleared() {
		_spec.ClearField(evmevent.FieldTopic1, field.TypeString)
	}
	if value, ok := eeu.mutation.Topic1Hex(); ok {
		_spec.SetField(evmevent.FieldTopic1Hex, field.TypeString, value)
	}
	if eeu.mutation.Topic1HexCleared() {
		_spec.ClearField(evmevent.FieldTopic1Hex, field.TypeString)
	}
	if value, ok := eeu.mutation.Topic2(); ok {
		_spec.SetField(evmevent.FieldTopic2, field.TypeString, value)
	}
	if eeu.mutation.Topic2Cleared() {
		_spec.ClearField(evmevent.FieldTopic2, field.TypeString)
	}
	if value, ok := eeu.mutation.Topic2Hex(); ok {
		_spec.SetField(evmevent.FieldTopic2Hex, field.TypeString, value)
	}
	if eeu.mutation.Topic2HexCleared() {
		_spec.ClearField(evmevent.FieldTopic2Hex, field.TypeString)
	}
	if value, ok := eeu.mutation.Topic3(); ok {
		_spec.SetField(evmevent.FieldTopic3, field.TypeString, value)
	}
	if eeu.mutation.Topic3Cleared() {
		_spec.ClearField(evmevent.FieldTopic3, field.TypeString)
	}
	if value, ok := eeu.mutation.Topic3Hex(); ok {
		_spec.SetField(evmevent.FieldTopic3Hex, field.TypeString, value)
	}
	if eeu.mutation.Topic3HexCleared() {
		_spec.ClearField(evmevent.FieldTopic3Hex, field.TypeString)
	}
	if value, ok := eeu.mutation.Data(); ok {
		_spec.SetField(evmevent.FieldData, field.TypeString, value)
	}
	if eeu.mutation.DataCleared() {
		_spec.ClearField(evmevent.FieldData, field.TypeString)
	}
	if value, ok := eeu.mutation.DataHex(); ok {
		_spec.SetField(evmevent.FieldDataHex, field.TypeString, value)
	}
	if eeu.mutation.DataHexCleared() {
		_spec.ClearField(evmevent.FieldDataHex, field.TypeString)
	}
	if value, ok := eeu.mutation.Removed(); ok {
		_spec.SetField(evmevent.FieldRemoved, field.TypeBool, value)
	}
	if value, ok := eeu.mutation.Status(); ok {
		_spec.SetField(evmevent.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := eeu.mutation.FailedReason(); ok {
		_spec.SetField(evmevent.FieldFailedReason, field.TypeString, value)
	}
	if eeu.mutation.FailedReasonCleared() {
		_spec.ClearField(evmevent.FieldFailedReason, field.TypeString)
	}
	if value, ok := eeu.mutation.Timestamp(); ok {
		_spec.SetField(evmevent.FieldTimestamp, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eeu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{evmevent.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eeu.mutation.done = true
	return n, nil
}

// EVMEventUpdateOne is the builder for updating a single EVMEvent entity.
type EVMEventUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EVMEventMutation
}

// SetTransactionHash sets the "transaction_hash" field.
func (eeuo *EVMEventUpdateOne) SetTransactionHash(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTransactionHash(s)
	return eeuo
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTransactionHash(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTransactionHash(*s)
	}
	return eeuo
}

// SetTransactionIndex sets the "transaction_index" field.
func (eeuo *EVMEventUpdateOne) SetTransactionIndex(u uint) *EVMEventUpdateOne {
	eeuo.mutation.ResetTransactionIndex()
	eeuo.mutation.SetTransactionIndex(u)
	return eeuo
}

// SetNillableTransactionIndex sets the "transaction_index" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTransactionIndex(u *uint) *EVMEventUpdateOne {
	if u != nil {
		eeuo.SetTransactionIndex(*u)
	}
	return eeuo
}

// AddTransactionIndex adds u to the "transaction_index" field.
func (eeuo *EVMEventUpdateOne) AddTransactionIndex(u int) *EVMEventUpdateOne {
	eeuo.mutation.AddTransactionIndex(u)
	return eeuo
}

// SetBlockHash sets the "block_hash" field.
func (eeuo *EVMEventUpdateOne) SetBlockHash(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetBlockHash(s)
	return eeuo
}

// SetNillableBlockHash sets the "block_hash" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableBlockHash(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetBlockHash(*s)
	}
	return eeuo
}

// SetBlockNumber sets the "block_number" field.
func (eeuo *EVMEventUpdateOne) SetBlockNumber(u uint64) *EVMEventUpdateOne {
	eeuo.mutation.ResetBlockNumber()
	eeuo.mutation.SetBlockNumber(u)
	return eeuo
}

// SetNillableBlockNumber sets the "block_number" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableBlockNumber(u *uint64) *EVMEventUpdateOne {
	if u != nil {
		eeuo.SetBlockNumber(*u)
	}
	return eeuo
}

// AddBlockNumber adds u to the "block_number" field.
func (eeuo *EVMEventUpdateOne) AddBlockNumber(u int64) *EVMEventUpdateOne {
	eeuo.mutation.AddBlockNumber(u)
	return eeuo
}

// SetLogIndex sets the "log_index" field.
func (eeuo *EVMEventUpdateOne) SetLogIndex(u uint) *EVMEventUpdateOne {
	eeuo.mutation.ResetLogIndex()
	eeuo.mutation.SetLogIndex(u)
	return eeuo
}

// SetNillableLogIndex sets the "log_index" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableLogIndex(u *uint) *EVMEventUpdateOne {
	if u != nil {
		eeuo.SetLogIndex(*u)
	}
	return eeuo
}

// AddLogIndex adds u to the "log_index" field.
func (eeuo *EVMEventUpdateOne) AddLogIndex(u int) *EVMEventUpdateOne {
	eeuo.mutation.AddLogIndex(u)
	return eeuo
}

// SetAddress sets the "address" field.
func (eeuo *EVMEventUpdateOne) SetAddress(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetAddress(s)
	return eeuo
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableAddress(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetAddress(*s)
	}
	return eeuo
}

// SetTopic0 sets the "topic0" field.
func (eeuo *EVMEventUpdateOne) SetTopic0(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic0(s)
	return eeuo
}

// SetNillableTopic0 sets the "topic0" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic0(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic0(*s)
	}
	return eeuo
}

// SetTopic0Hex sets the "topic0_hex" field.
func (eeuo *EVMEventUpdateOne) SetTopic0Hex(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic0Hex(s)
	return eeuo
}

// SetNillableTopic0Hex sets the "topic0_hex" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic0Hex(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic0Hex(*s)
	}
	return eeuo
}

// SetTopic1 sets the "topic1" field.
func (eeuo *EVMEventUpdateOne) SetTopic1(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic1(s)
	return eeuo
}

// SetNillableTopic1 sets the "topic1" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic1(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic1(*s)
	}
	return eeuo
}

// ClearTopic1 clears the value of the "topic1" field.
func (eeuo *EVMEventUpdateOne) ClearTopic1() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic1()
	return eeuo
}

// SetTopic1Hex sets the "topic1_hex" field.
func (eeuo *EVMEventUpdateOne) SetTopic1Hex(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic1Hex(s)
	return eeuo
}

// SetNillableTopic1Hex sets the "topic1_hex" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic1Hex(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic1Hex(*s)
	}
	return eeuo
}

// ClearTopic1Hex clears the value of the "topic1_hex" field.
func (eeuo *EVMEventUpdateOne) ClearTopic1Hex() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic1Hex()
	return eeuo
}

// SetTopic2 sets the "topic2" field.
func (eeuo *EVMEventUpdateOne) SetTopic2(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic2(s)
	return eeuo
}

// SetNillableTopic2 sets the "topic2" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic2(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic2(*s)
	}
	return eeuo
}

// ClearTopic2 clears the value of the "topic2" field.
func (eeuo *EVMEventUpdateOne) ClearTopic2() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic2()
	return eeuo
}

// SetTopic2Hex sets the "topic2_hex" field.
func (eeuo *EVMEventUpdateOne) SetTopic2Hex(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic2Hex(s)
	return eeuo
}

// SetNillableTopic2Hex sets the "topic2_hex" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic2Hex(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic2Hex(*s)
	}
	return eeuo
}

// ClearTopic2Hex clears the value of the "topic2_hex" field.
func (eeuo *EVMEventUpdateOne) ClearTopic2Hex() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic2Hex()
	return eeuo
}

// SetTopic3 sets the "topic3" field.
func (eeuo *EVMEventUpdateOne) SetTopic3(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic3(s)
	return eeuo
}

// SetNillableTopic3 sets the "topic3" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic3(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic3(*s)
	}
	return eeuo
}

// ClearTopic3 clears the value of the "topic3" field.
func (eeuo *EVMEventUpdateOne) ClearTopic3() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic3()
	return eeuo
}

// SetTopic3Hex sets the "topic3_hex" field.
func (eeuo *EVMEventUpdateOne) SetTopic3Hex(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetTopic3Hex(s)
	return eeuo
}

// SetNillableTopic3Hex sets the "topic3_hex" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTopic3Hex(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetTopic3Hex(*s)
	}
	return eeuo
}

// ClearTopic3Hex clears the value of the "topic3_hex" field.
func (eeuo *EVMEventUpdateOne) ClearTopic3Hex() *EVMEventUpdateOne {
	eeuo.mutation.ClearTopic3Hex()
	return eeuo
}

// SetData sets the "data" field.
func (eeuo *EVMEventUpdateOne) SetData(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetData(s)
	return eeuo
}

// SetNillableData sets the "data" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableData(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetData(*s)
	}
	return eeuo
}

// ClearData clears the value of the "data" field.
func (eeuo *EVMEventUpdateOne) ClearData() *EVMEventUpdateOne {
	eeuo.mutation.ClearData()
	return eeuo
}

// SetDataHex sets the "data_hex" field.
func (eeuo *EVMEventUpdateOne) SetDataHex(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetDataHex(s)
	return eeuo
}

// SetNillableDataHex sets the "data_hex" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableDataHex(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetDataHex(*s)
	}
	return eeuo
}

// ClearDataHex clears the value of the "data_hex" field.
func (eeuo *EVMEventUpdateOne) ClearDataHex() *EVMEventUpdateOne {
	eeuo.mutation.ClearDataHex()
	return eeuo
}

// SetRemoved sets the "removed" field.
func (eeuo *EVMEventUpdateOne) SetRemoved(b bool) *EVMEventUpdateOne {
	eeuo.mutation.SetRemoved(b)
	return eeuo
}

// SetNillableRemoved sets the "removed" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableRemoved(b *bool) *EVMEventUpdateOne {
	if b != nil {
		eeuo.SetRemoved(*b)
	}
	return eeuo
}

// SetStatus sets the "status" field.
func (eeuo *EVMEventUpdateOne) SetStatus(e evmevent.Status) *EVMEventUpdateOne {
	eeuo.mutation.SetStatus(e)
	return eeuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableStatus(e *evmevent.Status) *EVMEventUpdateOne {
	if e != nil {
		eeuo.SetStatus(*e)
	}
	return eeuo
}

// SetFailedReason sets the "failed_reason" field.
func (eeuo *EVMEventUpdateOne) SetFailedReason(s string) *EVMEventUpdateOne {
	eeuo.mutation.SetFailedReason(s)
	return eeuo
}

// SetNillableFailedReason sets the "failed_reason" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableFailedReason(s *string) *EVMEventUpdateOne {
	if s != nil {
		eeuo.SetFailedReason(*s)
	}
	return eeuo
}

// ClearFailedReason clears the value of the "failed_reason" field.
func (eeuo *EVMEventUpdateOne) ClearFailedReason() *EVMEventUpdateOne {
	eeuo.mutation.ClearFailedReason()
	return eeuo
}

// SetTimestamp sets the "timestamp" field.
func (eeuo *EVMEventUpdateOne) SetTimestamp(t time.Time) *EVMEventUpdateOne {
	eeuo.mutation.SetTimestamp(t)
	return eeuo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (eeuo *EVMEventUpdateOne) SetNillableTimestamp(t *time.Time) *EVMEventUpdateOne {
	if t != nil {
		eeuo.SetTimestamp(*t)
	}
	return eeuo
}

// Mutation returns the EVMEventMutation object of the builder.
func (eeuo *EVMEventUpdateOne) Mutation() *EVMEventMutation {
	return eeuo.mutation
}

// Where appends a list predicates to the EVMEventUpdate builder.
func (eeuo *EVMEventUpdateOne) Where(ps ...predicate.EVMEvent) *EVMEventUpdateOne {
	eeuo.mutation.Where(ps...)
	return eeuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (eeuo *EVMEventUpdateOne) Select(field string, fields ...string) *EVMEventUpdateOne {
	eeuo.fields = append([]string{field}, fields...)
	return eeuo
}

// Save executes the query and returns the updated EVMEvent entity.
func (eeuo *EVMEventUpdateOne) Save(ctx context.Context) (*EVMEvent, error) {
	return withHooks(ctx, eeuo.sqlSave, eeuo.mutation, eeuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eeuo *EVMEventUpdateOne) SaveX(ctx context.Context) *EVMEvent {
	node, err := eeuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (eeuo *EVMEventUpdateOne) Exec(ctx context.Context) error {
	_, err := eeuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eeuo *EVMEventUpdateOne) ExecX(ctx context.Context) {
	if err := eeuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eeuo *EVMEventUpdateOne) check() error {
	if v, ok := eeuo.mutation.TransactionHash(); ok {
		if err := evmevent.TransactionHashValidator(v); err != nil {
			return &ValidationError{Name: "transaction_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.transaction_hash": %w`, err)}
		}
	}
	if v, ok := eeuo.mutation.BlockHash(); ok {
		if err := evmevent.BlockHashValidator(v); err != nil {
			return &ValidationError{Name: "block_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.block_hash": %w`, err)}
		}
	}
	if v, ok := eeuo.mutation.Address(); ok {
		if err := evmevent.AddressValidator(v); err != nil {
			return &ValidationError{Name: "address", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.address": %w`, err)}
		}
	}
	if v, ok := eeuo.mutation.Topic0(); ok {
		if err := evmevent.Topic0Validator(v); err != nil {
			return &ValidationError{Name: "topic0", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0": %w`, err)}
		}
	}
	if v, ok := eeuo.mutation.Topic0Hex(); ok {
		if err := evmevent.Topic0HexValidator(v); err != nil {
			return &ValidationError{Name: "topic0_hex", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0_hex": %w`, err)}
		}
	}
	if v, ok := eeuo.mutation.Status(); ok {
		if err := evmevent.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.status": %w`, err)}
		}
	}
	return nil
}

func (eeuo *EVMEventUpdateOne) sqlSave(ctx context.Context) (_node *EVMEvent, err error) {
	if err := eeuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(evmevent.Table, evmevent.Columns, sqlgraph.NewFieldSpec(evmevent.FieldID, field.TypeInt))
	id, ok := eeuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EVMEvent.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := eeuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, evmevent.FieldID)
		for _, f := range fields {
			if !evmevent.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != evmevent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := eeuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eeuo.mutation.TransactionHash(); ok {
		_spec.SetField(evmevent.FieldTransactionHash, field.TypeString, value)
	}
	if value, ok := eeuo.mutation.TransactionIndex(); ok {
		_spec.SetField(evmevent.FieldTransactionIndex, field.TypeUint, value)
	}
	if value, ok := eeuo.mutation.AddedTransactionIndex(); ok {
		_spec.AddField(evmevent.FieldTransactionIndex, field.TypeUint, value)
	}
	if value, ok := eeuo.mutation.BlockHash(); ok {
		_spec.SetField(evmevent.FieldBlockHash, field.TypeString, value)
	}
	if value, ok := eeuo.mutation.BlockNumber(); ok {
		_spec.SetField(evmevent.FieldBlockNumber, field.TypeUint64, value)
	}
	if value, ok := eeuo.mutation.AddedBlockNumber(); ok {
		_spec.AddField(evmevent.FieldBlockNumber, field.TypeUint64, value)
	}
	if value, ok := eeuo.mutation.LogIndex(); ok {
		_spec.SetField(evmevent.FieldLogIndex, field.TypeUint, value)
	}
	if value, ok := eeuo.mutation.AddedLogIndex(); ok {
		_spec.AddField(evmevent.FieldLogIndex, field.TypeUint, value)
	}
	if value, ok := eeuo.mutation.Address(); ok {
		_spec.SetField(evmevent.FieldAddress, field.TypeString, value)
	}
	if value, ok := eeuo.mutation.Topic0(); ok {
		_spec.SetField(evmevent.FieldTopic0, field.TypeString, value)
	}
	if value, ok := eeuo.mutation.Topic0Hex(); ok {
		_spec.SetField(evmevent.FieldTopic0Hex, field.TypeString, value)
	}
	if value, ok := eeuo.mutation.Topic1(); ok {
		_spec.SetField(evmevent.FieldTopic1, field.TypeString, value)
	}
	if eeuo.mutation.Topic1Cleared() {
		_spec.ClearField(evmevent.FieldTopic1, field.TypeString)
	}
	if value, ok := eeuo.mutation.Topic1Hex(); ok {
		_spec.SetField(evmevent.FieldTopic1Hex, field.TypeString, value)
	}
	if eeuo.mutation.Topic1HexCleared() {
		_spec.ClearField(evmevent.FieldTopic1Hex, field.TypeString)
	}
	if value, ok := eeuo.mutation.Topic2(); ok {
		_spec.SetField(evmevent.FieldTopic2, field.TypeString, value)
	}
	if eeuo.mutation.Topic2Cleared() {
		_spec.ClearField(evmevent.FieldTopic2, field.TypeString)
	}
	if value, ok := eeuo.mutation.Topic2Hex(); ok {
		_spec.SetField(evmevent.FieldTopic2Hex, field.TypeString, value)
	}
	if eeuo.mutation.Topic2HexCleared() {
		_spec.ClearField(evmevent.FieldTopic2Hex, field.TypeString)
	}
	if value, ok := eeuo.mutation.Topic3(); ok {
		_spec.SetField(evmevent.FieldTopic3, field.TypeString, value)
	}
	if eeuo.mutation.Topic3Cleared() {
		_spec.ClearField(evmevent.FieldTopic3, field.TypeString)
	}
	if value, ok := eeuo.mutation.Topic3Hex(); ok {
		_spec.SetField(evmevent.FieldTopic3Hex, field.TypeString, value)
	}
	if eeuo.mutation.Topic3HexCleared() {
		_spec.ClearField(evmevent.FieldTopic3Hex, field.TypeString)
	}
	if value, ok := eeuo.mutation.Data(); ok {
		_spec.SetField(evmevent.FieldData, field.TypeString, value)
	}
	if eeuo.mutation.DataCleared() {
		_spec.ClearField(evmevent.FieldData, field.TypeString)
	}
	if value, ok := eeuo.mutation.DataHex(); ok {
		_spec.SetField(evmevent.FieldDataHex, field.TypeString, value)
	}
	if eeuo.mutation.DataHexCleared() {
		_spec.ClearField(evmevent.FieldDataHex, field.TypeString)
	}
	if value, ok := eeuo.mutation.Removed(); ok {
		_spec.SetField(evmevent.FieldRemoved, field.TypeBool, value)
	}
	if value, ok := eeuo.mutation.Status(); ok {
		_spec.SetField(evmevent.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := eeuo.mutation.FailedReason(); ok {
		_spec.SetField(evmevent.FieldFailedReason, field.TypeString, value)
	}
	if eeuo.mutation.FailedReasonCleared() {
		_spec.ClearField(evmevent.FieldFailedReason, field.TypeString)
	}
	if value, ok := eeuo.mutation.Timestamp(); ok {
		_spec.SetField(evmevent.FieldTimestamp, field.TypeTime, value)
	}
	_node = &EVMEvent{config: eeuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, eeuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{evmevent.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	eeuo.mutation.done = true
	return _node, nil
}
