// Code generated by ent, DO NOT EDIT.

package transactionmemo

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the transactionmemo type in the database.
	Label = "transaction_memo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTransactionHash holds the string denoting the transaction_hash field in the database.
	FieldTransactionHash = "transaction_hash"
	// FieldBookNftID holds the string denoting the book_nft_id field in the database.
	FieldBookNftID = "book_nft_id"
	// FieldFrom holds the string denoting the from field in the database.
	FieldFrom = "from"
	// FieldTo holds the string denoting the to field in the database.
	FieldTo = "to"
	// FieldTokenID holds the string denoting the token_id field in the database.
	FieldTokenID = "token_id"
	// FieldMemo holds the string denoting the memo field in the database.
	FieldMemo = "memo"
	// FieldBlockNumber holds the string denoting the block_number field in the database.
	FieldBlockNumber = "block_number"
	// Table holds the table name of the transactionmemo in the database.
	Table = "transaction_memos"
)

// Columns holds all SQL columns for transactionmemo fields.
var Columns = []string{
	FieldID,
	FieldTransactionHash,
	FieldBookNftID,
	FieldFrom,
	FieldTo,
	FieldTokenID,
	FieldMemo,
	FieldBlockNumber,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// TransactionHashValidator is a validator for the "transaction_hash" field. It is called by the builders before save.
	TransactionHashValidator func(string) error
	// BookNftIDValidator is a validator for the "book_nft_id" field. It is called by the builders before save.
	BookNftIDValidator func(string) error
	// FromValidator is a validator for the "from" field. It is called by the builders before save.
	FromValidator func(string) error
	// ToValidator is a validator for the "to" field. It is called by the builders before save.
	ToValidator func(string) error
)

// OrderOption defines the ordering options for the TransactionMemo queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTransactionHash orders the results by the transaction_hash field.
func ByTransactionHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTransactionHash, opts...).ToFunc()
}

// ByBookNftID orders the results by the book_nft_id field.
func ByBookNftID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBookNftID, opts...).ToFunc()
}

// ByFrom orders the results by the from field.
func ByFrom(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFrom, opts...).ToFunc()
}

// ByTo orders the results by the to field.
func ByTo(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTo, opts...).ToFunc()
}

// ByTokenID orders the results by the token_id field.
func ByTokenID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTokenID, opts...).ToFunc()
}

// ByMemo orders the results by the memo field.
func ByMemo(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMemo, opts...).ToFunc()
}

// ByBlockNumber orders the results by the block_number field.
func ByBlockNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBlockNumber, opts...).ToFunc()
}
