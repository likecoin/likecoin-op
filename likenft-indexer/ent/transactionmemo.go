// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"likenft-indexer/ent/transactionmemo"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// TransactionMemo is the model entity for the TransactionMemo schema.
type TransactionMemo struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TransactionHash holds the value of the "transaction_hash" field.
	TransactionHash string `json:"transaction_hash,omitempty"`
	// contract address of book nft
	BookNftID string `json:"book_nft_id,omitempty"`
	// From holds the value of the "from" field.
	From string `json:"from,omitempty"`
	// To holds the value of the "to" field.
	To string `json:"to,omitempty"`
	// TokenID holds the value of the "token_id" field.
	TokenID uint64 `json:"token_id,omitempty"`
	// Memo holds the value of the "memo" field.
	Memo string `json:"memo,omitempty"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber  uint64 `json:"block_number,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TransactionMemo) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case transactionmemo.FieldID, transactionmemo.FieldTokenID, transactionmemo.FieldBlockNumber:
			values[i] = new(sql.NullInt64)
		case transactionmemo.FieldTransactionHash, transactionmemo.FieldBookNftID, transactionmemo.FieldFrom, transactionmemo.FieldTo, transactionmemo.FieldMemo:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TransactionMemo fields.
func (tm *TransactionMemo) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case transactionmemo.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tm.ID = int(value.Int64)
		case transactionmemo.FieldTransactionHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field transaction_hash", values[i])
			} else if value.Valid {
				tm.TransactionHash = value.String
			}
		case transactionmemo.FieldBookNftID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field book_nft_id", values[i])
			} else if value.Valid {
				tm.BookNftID = value.String
			}
		case transactionmemo.FieldFrom:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field from", values[i])
			} else if value.Valid {
				tm.From = value.String
			}
		case transactionmemo.FieldTo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field to", values[i])
			} else if value.Valid {
				tm.To = value.String
			}
		case transactionmemo.FieldTokenID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field token_id", values[i])
			} else if value.Valid {
				tm.TokenID = uint64(value.Int64)
			}
		case transactionmemo.FieldMemo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field memo", values[i])
			} else if value.Valid {
				tm.Memo = value.String
			}
		case transactionmemo.FieldBlockNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field block_number", values[i])
			} else if value.Valid {
				tm.BlockNumber = uint64(value.Int64)
			}
		default:
			tm.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TransactionMemo.
// This includes values selected through modifiers, order, etc.
func (tm *TransactionMemo) Value(name string) (ent.Value, error) {
	return tm.selectValues.Get(name)
}

// Update returns a builder for updating this TransactionMemo.
// Note that you need to call TransactionMemo.Unwrap() before calling this method if this TransactionMemo
// was returned from a transaction, and the transaction was committed or rolled back.
func (tm *TransactionMemo) Update() *TransactionMemoUpdateOne {
	return NewTransactionMemoClient(tm.config).UpdateOne(tm)
}

// Unwrap unwraps the TransactionMemo entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tm *TransactionMemo) Unwrap() *TransactionMemo {
	_tx, ok := tm.config.driver.(*txDriver)
	if !ok {
		panic("ent: TransactionMemo is not a transactional entity")
	}
	tm.config.driver = _tx.drv
	return tm
}

// String implements the fmt.Stringer.
func (tm *TransactionMemo) String() string {
	var builder strings.Builder
	builder.WriteString("TransactionMemo(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tm.ID))
	builder.WriteString("transaction_hash=")
	builder.WriteString(tm.TransactionHash)
	builder.WriteString(", ")
	builder.WriteString("book_nft_id=")
	builder.WriteString(tm.BookNftID)
	builder.WriteString(", ")
	builder.WriteString("from=")
	builder.WriteString(tm.From)
	builder.WriteString(", ")
	builder.WriteString("to=")
	builder.WriteString(tm.To)
	builder.WriteString(", ")
	builder.WriteString("token_id=")
	builder.WriteString(fmt.Sprintf("%v", tm.TokenID))
	builder.WriteString(", ")
	builder.WriteString("memo=")
	builder.WriteString(tm.Memo)
	builder.WriteString(", ")
	builder.WriteString("block_number=")
	builder.WriteString(fmt.Sprintf("%v", tm.BlockNumber))
	builder.WriteByte(')')
	return builder.String()
}

// TransactionMemos is a parsable slice of TransactionMemo.
type TransactionMemos []*TransactionMemo
