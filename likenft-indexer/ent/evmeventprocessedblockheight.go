// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"likenft-indexer/ent/evmeventprocessedblockheight"
	"likenft-indexer/ent/schema/typeutil"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// EVMEventProcessedBlockHeight is the model entity for the EVMEventProcessedBlockHeight schema.
type EVMEventProcessedBlockHeight struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// For reference only
	ContractType evmeventprocessedblockheight.ContractType `json:"contract_type,omitempty"`
	// ContractAddress holds the value of the "contract_address" field.
	ContractAddress string `json:"contract_address,omitempty"`
	// Event holds the value of the "event" field.
	Event evmeventprocessedblockheight.Event `json:"event,omitempty"`
	// BlockHeight holds the value of the "block_height" field.
	BlockHeight  typeutil.Uint64 `json:"block_height,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EVMEventProcessedBlockHeight) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case evmeventprocessedblockheight.FieldID:
			values[i] = new(sql.NullInt64)
		case evmeventprocessedblockheight.FieldContractType, evmeventprocessedblockheight.FieldContractAddress, evmeventprocessedblockheight.FieldEvent:
			values[i] = new(sql.NullString)
		case evmeventprocessedblockheight.FieldBlockHeight:
			values[i] = evmeventprocessedblockheight.ValueScanner.BlockHeight.ScanValue()
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EVMEventProcessedBlockHeight fields.
func (eepbh *EVMEventProcessedBlockHeight) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case evmeventprocessedblockheight.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			eepbh.ID = int(value.Int64)
		case evmeventprocessedblockheight.FieldContractType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field contract_type", values[i])
			} else if value.Valid {
				eepbh.ContractType = evmeventprocessedblockheight.ContractType(value.String)
			}
		case evmeventprocessedblockheight.FieldContractAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field contract_address", values[i])
			} else if value.Valid {
				eepbh.ContractAddress = value.String
			}
		case evmeventprocessedblockheight.FieldEvent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field event", values[i])
			} else if value.Valid {
				eepbh.Event = evmeventprocessedblockheight.Event(value.String)
			}
		case evmeventprocessedblockheight.FieldBlockHeight:
			if value, err := evmeventprocessedblockheight.ValueScanner.BlockHeight.FromValue(values[i]); err != nil {
				return err
			} else {
				eepbh.BlockHeight = value
			}
		default:
			eepbh.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the EVMEventProcessedBlockHeight.
// This includes values selected through modifiers, order, etc.
func (eepbh *EVMEventProcessedBlockHeight) Value(name string) (ent.Value, error) {
	return eepbh.selectValues.Get(name)
}

// Update returns a builder for updating this EVMEventProcessedBlockHeight.
// Note that you need to call EVMEventProcessedBlockHeight.Unwrap() before calling this method if this EVMEventProcessedBlockHeight
// was returned from a transaction, and the transaction was committed or rolled back.
func (eepbh *EVMEventProcessedBlockHeight) Update() *EVMEventProcessedBlockHeightUpdateOne {
	return NewEVMEventProcessedBlockHeightClient(eepbh.config).UpdateOne(eepbh)
}

// Unwrap unwraps the EVMEventProcessedBlockHeight entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (eepbh *EVMEventProcessedBlockHeight) Unwrap() *EVMEventProcessedBlockHeight {
	_tx, ok := eepbh.config.driver.(*txDriver)
	if !ok {
		panic("ent: EVMEventProcessedBlockHeight is not a transactional entity")
	}
	eepbh.config.driver = _tx.drv
	return eepbh
}

// String implements the fmt.Stringer.
func (eepbh *EVMEventProcessedBlockHeight) String() string {
	var builder strings.Builder
	builder.WriteString("EVMEventProcessedBlockHeight(")
	builder.WriteString(fmt.Sprintf("id=%v, ", eepbh.ID))
	builder.WriteString("contract_type=")
	builder.WriteString(fmt.Sprintf("%v", eepbh.ContractType))
	builder.WriteString(", ")
	builder.WriteString("contract_address=")
	builder.WriteString(eepbh.ContractAddress)
	builder.WriteString(", ")
	builder.WriteString("event=")
	builder.WriteString(fmt.Sprintf("%v", eepbh.Event))
	builder.WriteString(", ")
	builder.WriteString("block_height=")
	builder.WriteString(fmt.Sprintf("%v", eepbh.BlockHeight))
	builder.WriteByte(')')
	return builder.String()
}

// EVMEventProcessedBlockHeights is a parsable slice of EVMEventProcessedBlockHeight.
type EVMEventProcessedBlockHeights []*EVMEventProcessedBlockHeight
