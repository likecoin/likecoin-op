// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"likenft-indexer/ent/evmevent"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EVMEventCreate is the builder for creating a EVMEvent entity.
type EVMEventCreate struct {
	config
	mutation *EVMEventMutation
	hooks    []Hook
}

// SetTransactionHash sets the "transaction_hash" field.
func (eec *EVMEventCreate) SetTransactionHash(s string) *EVMEventCreate {
	eec.mutation.SetTransactionHash(s)
	return eec
}

// SetTransactionIndex sets the "transaction_index" field.
func (eec *EVMEventCreate) SetTransactionIndex(u uint) *EVMEventCreate {
	eec.mutation.SetTransactionIndex(u)
	return eec
}

// SetBlockHash sets the "block_hash" field.
func (eec *EVMEventCreate) SetBlockHash(s string) *EVMEventCreate {
	eec.mutation.SetBlockHash(s)
	return eec
}

// SetBlockNumber sets the "block_number" field.
func (eec *EVMEventCreate) SetBlockNumber(u uint64) *EVMEventCreate {
	eec.mutation.SetBlockNumber(u)
	return eec
}

// SetLogIndex sets the "log_index" field.
func (eec *EVMEventCreate) SetLogIndex(u uint) *EVMEventCreate {
	eec.mutation.SetLogIndex(u)
	return eec
}

// SetAddress sets the "address" field.
func (eec *EVMEventCreate) SetAddress(s string) *EVMEventCreate {
	eec.mutation.SetAddress(s)
	return eec
}

// SetTopic0 sets the "topic0" field.
func (eec *EVMEventCreate) SetTopic0(s string) *EVMEventCreate {
	eec.mutation.SetTopic0(s)
	return eec
}

// SetTopic0Hex sets the "topic0_hex" field.
func (eec *EVMEventCreate) SetTopic0Hex(s string) *EVMEventCreate {
	eec.mutation.SetTopic0Hex(s)
	return eec
}

// SetTopic1 sets the "topic1" field.
func (eec *EVMEventCreate) SetTopic1(s string) *EVMEventCreate {
	eec.mutation.SetTopic1(s)
	return eec
}

// SetNillableTopic1 sets the "topic1" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic1(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic1(*s)
	}
	return eec
}

// SetTopic1Hex sets the "topic1_hex" field.
func (eec *EVMEventCreate) SetTopic1Hex(s string) *EVMEventCreate {
	eec.mutation.SetTopic1Hex(s)
	return eec
}

// SetNillableTopic1Hex sets the "topic1_hex" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic1Hex(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic1Hex(*s)
	}
	return eec
}

// SetTopic2 sets the "topic2" field.
func (eec *EVMEventCreate) SetTopic2(s string) *EVMEventCreate {
	eec.mutation.SetTopic2(s)
	return eec
}

// SetNillableTopic2 sets the "topic2" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic2(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic2(*s)
	}
	return eec
}

// SetTopic2Hex sets the "topic2_hex" field.
func (eec *EVMEventCreate) SetTopic2Hex(s string) *EVMEventCreate {
	eec.mutation.SetTopic2Hex(s)
	return eec
}

// SetNillableTopic2Hex sets the "topic2_hex" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic2Hex(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic2Hex(*s)
	}
	return eec
}

// SetTopic3 sets the "topic3" field.
func (eec *EVMEventCreate) SetTopic3(s string) *EVMEventCreate {
	eec.mutation.SetTopic3(s)
	return eec
}

// SetNillableTopic3 sets the "topic3" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic3(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic3(*s)
	}
	return eec
}

// SetTopic3Hex sets the "topic3_hex" field.
func (eec *EVMEventCreate) SetTopic3Hex(s string) *EVMEventCreate {
	eec.mutation.SetTopic3Hex(s)
	return eec
}

// SetNillableTopic3Hex sets the "topic3_hex" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableTopic3Hex(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetTopic3Hex(*s)
	}
	return eec
}

// SetData sets the "data" field.
func (eec *EVMEventCreate) SetData(s string) *EVMEventCreate {
	eec.mutation.SetData(s)
	return eec
}

// SetNillableData sets the "data" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableData(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetData(*s)
	}
	return eec
}

// SetDataHex sets the "data_hex" field.
func (eec *EVMEventCreate) SetDataHex(s string) *EVMEventCreate {
	eec.mutation.SetDataHex(s)
	return eec
}

// SetNillableDataHex sets the "data_hex" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableDataHex(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetDataHex(*s)
	}
	return eec
}

// SetRemoved sets the "removed" field.
func (eec *EVMEventCreate) SetRemoved(b bool) *EVMEventCreate {
	eec.mutation.SetRemoved(b)
	return eec
}

// SetStatus sets the "status" field.
func (eec *EVMEventCreate) SetStatus(e evmevent.Status) *EVMEventCreate {
	eec.mutation.SetStatus(e)
	return eec
}

// SetFailedReason sets the "failed_reason" field.
func (eec *EVMEventCreate) SetFailedReason(s string) *EVMEventCreate {
	eec.mutation.SetFailedReason(s)
	return eec
}

// SetNillableFailedReason sets the "failed_reason" field if the given value is not nil.
func (eec *EVMEventCreate) SetNillableFailedReason(s *string) *EVMEventCreate {
	if s != nil {
		eec.SetFailedReason(*s)
	}
	return eec
}

// SetTimestamp sets the "timestamp" field.
func (eec *EVMEventCreate) SetTimestamp(t time.Time) *EVMEventCreate {
	eec.mutation.SetTimestamp(t)
	return eec
}

// Mutation returns the EVMEventMutation object of the builder.
func (eec *EVMEventCreate) Mutation() *EVMEventMutation {
	return eec.mutation
}

// Save creates the EVMEvent in the database.
func (eec *EVMEventCreate) Save(ctx context.Context) (*EVMEvent, error) {
	return withHooks(ctx, eec.sqlSave, eec.mutation, eec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (eec *EVMEventCreate) SaveX(ctx context.Context) *EVMEvent {
	v, err := eec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (eec *EVMEventCreate) Exec(ctx context.Context) error {
	_, err := eec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eec *EVMEventCreate) ExecX(ctx context.Context) {
	if err := eec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eec *EVMEventCreate) check() error {
	if _, ok := eec.mutation.TransactionHash(); !ok {
		return &ValidationError{Name: "transaction_hash", err: errors.New(`ent: missing required field "EVMEvent.transaction_hash"`)}
	}
	if v, ok := eec.mutation.TransactionHash(); ok {
		if err := evmevent.TransactionHashValidator(v); err != nil {
			return &ValidationError{Name: "transaction_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.transaction_hash": %w`, err)}
		}
	}
	if _, ok := eec.mutation.TransactionIndex(); !ok {
		return &ValidationError{Name: "transaction_index", err: errors.New(`ent: missing required field "EVMEvent.transaction_index"`)}
	}
	if _, ok := eec.mutation.BlockHash(); !ok {
		return &ValidationError{Name: "block_hash", err: errors.New(`ent: missing required field "EVMEvent.block_hash"`)}
	}
	if v, ok := eec.mutation.BlockHash(); ok {
		if err := evmevent.BlockHashValidator(v); err != nil {
			return &ValidationError{Name: "block_hash", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.block_hash": %w`, err)}
		}
	}
	if _, ok := eec.mutation.BlockNumber(); !ok {
		return &ValidationError{Name: "block_number", err: errors.New(`ent: missing required field "EVMEvent.block_number"`)}
	}
	if _, ok := eec.mutation.LogIndex(); !ok {
		return &ValidationError{Name: "log_index", err: errors.New(`ent: missing required field "EVMEvent.log_index"`)}
	}
	if _, ok := eec.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "EVMEvent.address"`)}
	}
	if v, ok := eec.mutation.Address(); ok {
		if err := evmevent.AddressValidator(v); err != nil {
			return &ValidationError{Name: "address", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.address": %w`, err)}
		}
	}
	if _, ok := eec.mutation.Topic0(); !ok {
		return &ValidationError{Name: "topic0", err: errors.New(`ent: missing required field "EVMEvent.topic0"`)}
	}
	if v, ok := eec.mutation.Topic0(); ok {
		if err := evmevent.Topic0Validator(v); err != nil {
			return &ValidationError{Name: "topic0", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0": %w`, err)}
		}
	}
	if _, ok := eec.mutation.Topic0Hex(); !ok {
		return &ValidationError{Name: "topic0_hex", err: errors.New(`ent: missing required field "EVMEvent.topic0_hex"`)}
	}
	if v, ok := eec.mutation.Topic0Hex(); ok {
		if err := evmevent.Topic0HexValidator(v); err != nil {
			return &ValidationError{Name: "topic0_hex", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.topic0_hex": %w`, err)}
		}
	}
	if _, ok := eec.mutation.Removed(); !ok {
		return &ValidationError{Name: "removed", err: errors.New(`ent: missing required field "EVMEvent.removed"`)}
	}
	if _, ok := eec.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "EVMEvent.status"`)}
	}
	if v, ok := eec.mutation.Status(); ok {
		if err := evmevent.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "EVMEvent.status": %w`, err)}
		}
	}
	if _, ok := eec.mutation.Timestamp(); !ok {
		return &ValidationError{Name: "timestamp", err: errors.New(`ent: missing required field "EVMEvent.timestamp"`)}
	}
	return nil
}

func (eec *EVMEventCreate) sqlSave(ctx context.Context) (*EVMEvent, error) {
	if err := eec.check(); err != nil {
		return nil, err
	}
	_node, _spec := eec.createSpec()
	if err := sqlgraph.CreateNode(ctx, eec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	eec.mutation.id = &_node.ID
	eec.mutation.done = true
	return _node, nil
}

func (eec *EVMEventCreate) createSpec() (*EVMEvent, *sqlgraph.CreateSpec) {
	var (
		_node = &EVMEvent{config: eec.config}
		_spec = sqlgraph.NewCreateSpec(evmevent.Table, sqlgraph.NewFieldSpec(evmevent.FieldID, field.TypeInt))
	)
	if value, ok := eec.mutation.TransactionHash(); ok {
		_spec.SetField(evmevent.FieldTransactionHash, field.TypeString, value)
		_node.TransactionHash = value
	}
	if value, ok := eec.mutation.TransactionIndex(); ok {
		_spec.SetField(evmevent.FieldTransactionIndex, field.TypeUint, value)
		_node.TransactionIndex = value
	}
	if value, ok := eec.mutation.BlockHash(); ok {
		_spec.SetField(evmevent.FieldBlockHash, field.TypeString, value)
		_node.BlockHash = value
	}
	if value, ok := eec.mutation.BlockNumber(); ok {
		_spec.SetField(evmevent.FieldBlockNumber, field.TypeUint64, value)
		_node.BlockNumber = value
	}
	if value, ok := eec.mutation.LogIndex(); ok {
		_spec.SetField(evmevent.FieldLogIndex, field.TypeUint, value)
		_node.LogIndex = value
	}
	if value, ok := eec.mutation.Address(); ok {
		_spec.SetField(evmevent.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if value, ok := eec.mutation.Topic0(); ok {
		_spec.SetField(evmevent.FieldTopic0, field.TypeString, value)
		_node.Topic0 = value
	}
	if value, ok := eec.mutation.Topic0Hex(); ok {
		_spec.SetField(evmevent.FieldTopic0Hex, field.TypeString, value)
		_node.Topic0Hex = value
	}
	if value, ok := eec.mutation.Topic1(); ok {
		_spec.SetField(evmevent.FieldTopic1, field.TypeString, value)
		_node.Topic1 = &value
	}
	if value, ok := eec.mutation.Topic1Hex(); ok {
		_spec.SetField(evmevent.FieldTopic1Hex, field.TypeString, value)
		_node.Topic1Hex = &value
	}
	if value, ok := eec.mutation.Topic2(); ok {
		_spec.SetField(evmevent.FieldTopic2, field.TypeString, value)
		_node.Topic2 = &value
	}
	if value, ok := eec.mutation.Topic2Hex(); ok {
		_spec.SetField(evmevent.FieldTopic2Hex, field.TypeString, value)
		_node.Topic2Hex = &value
	}
	if value, ok := eec.mutation.Topic3(); ok {
		_spec.SetField(evmevent.FieldTopic3, field.TypeString, value)
		_node.Topic3 = &value
	}
	if value, ok := eec.mutation.Topic3Hex(); ok {
		_spec.SetField(evmevent.FieldTopic3Hex, field.TypeString, value)
		_node.Topic3Hex = &value
	}
	if value, ok := eec.mutation.Data(); ok {
		_spec.SetField(evmevent.FieldData, field.TypeString, value)
		_node.Data = &value
	}
	if value, ok := eec.mutation.DataHex(); ok {
		_spec.SetField(evmevent.FieldDataHex, field.TypeString, value)
		_node.DataHex = &value
	}
	if value, ok := eec.mutation.Removed(); ok {
		_spec.SetField(evmevent.FieldRemoved, field.TypeBool, value)
		_node.Removed = value
	}
	if value, ok := eec.mutation.Status(); ok {
		_spec.SetField(evmevent.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if value, ok := eec.mutation.FailedReason(); ok {
		_spec.SetField(evmevent.FieldFailedReason, field.TypeString, value)
		_node.FailedReason = &value
	}
	if value, ok := eec.mutation.Timestamp(); ok {
		_spec.SetField(evmevent.FieldTimestamp, field.TypeTime, value)
		_node.Timestamp = value
	}
	return _node, _spec
}

// EVMEventCreateBulk is the builder for creating many EVMEvent entities in bulk.
type EVMEventCreateBulk struct {
	config
	err      error
	builders []*EVMEventCreate
}

// Save creates the EVMEvent entities in the database.
func (eecb *EVMEventCreateBulk) Save(ctx context.Context) ([]*EVMEvent, error) {
	if eecb.err != nil {
		return nil, eecb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(eecb.builders))
	nodes := make([]*EVMEvent, len(eecb.builders))
	mutators := make([]Mutator, len(eecb.builders))
	for i := range eecb.builders {
		func(i int, root context.Context) {
			builder := eecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EVMEventMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, eecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, eecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, eecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (eecb *EVMEventCreateBulk) SaveX(ctx context.Context) []*EVMEvent {
	v, err := eecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (eecb *EVMEventCreateBulk) Exec(ctx context.Context) error {
	_, err := eecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eecb *EVMEventCreateBulk) ExecX(ctx context.Context) {
	if err := eecb.Exec(ctx); err != nil {
		panic(err)
	}
}
