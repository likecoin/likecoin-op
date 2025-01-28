// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "cosmos_address", Type: field.TypeString, Nullable: true},
		{Name: "evm_address", Type: field.TypeString, Unique: true},
		{Name: "likeid", Type: field.TypeString, Unique: true, Nullable: true},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
	}
	// EvmEventsColumns holds the columns for the "evm_events" table.
	EvmEventsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "transaction_hash", Type: field.TypeString, Unique: true},
		{Name: "block_hash", Type: field.TypeString},
		{Name: "block_number", Type: field.TypeUint64},
		{Name: "log_index", Type: field.TypeUint64},
		{Name: "address", Type: field.TypeString},
		{Name: "topic0", Type: field.TypeString},
		{Name: "topic1", Type: field.TypeString},
		{Name: "topic2", Type: field.TypeString},
		{Name: "topic3", Type: field.TypeString},
		{Name: "data", Type: field.TypeString},
		{Name: "timestamp", Type: field.TypeTime},
	}
	// EvmEventsTable holds the schema information for the "evm_events" table.
	EvmEventsTable = &schema.Table{
		Name:       "evm_events",
		Columns:    EvmEventsColumns,
		PrimaryKey: []*schema.Column{EvmEventsColumns[0]},
	}
	// NftsColumns holds the columns for the "nfts" table.
	NftsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "token_id", Type: field.TypeUint64, Unique: true},
		{Name: "token_url", Type: field.TypeString, Nullable: true},
		{Name: "raw", Type: field.TypeJSON, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "image", Type: field.TypeJSON, Nullable: true},
		{Name: "attributes", Type: field.TypeJSON, Nullable: true},
		{Name: "owner_address", Type: field.TypeString},
		{Name: "minted_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "account_nfts", Type: field.TypeInt, Nullable: true},
		{Name: "nft_class_nfts", Type: field.TypeInt, Nullable: true},
	}
	// NftsTable holds the schema information for the "nfts" table.
	NftsTable = &schema.Table{
		Name:       "nfts",
		Columns:    NftsColumns,
		PrimaryKey: []*schema.Column{NftsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "nfts_accounts_nfts",
				Columns:    []*schema.Column{NftsColumns[11]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "nfts_nft_classes_nfts",
				Columns:    []*schema.Column{NftsColumns[12]},
				RefColumns: []*schema.Column{NftClassesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// NftClassesColumns holds the columns for the "nft_classes" table.
	NftClassesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "address", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString},
		{Name: "symbol", Type: field.TypeString},
		{Name: "owner_address", Type: field.TypeString, Nullable: true},
		{Name: "minter_addresses", Type: field.TypeJSON, Nullable: true},
		{Name: "total_supply", Type: field.TypeInt},
		{Name: "metadata", Type: field.TypeJSON, Nullable: true},
		{Name: "banner_image", Type: field.TypeString},
		{Name: "featured_image", Type: field.TypeString},
		{Name: "deployer_address", Type: field.TypeString},
		{Name: "deployed_block_number", Type: field.TypeString},
		{Name: "minted_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "account_nft_classes", Type: field.TypeInt, Nullable: true},
	}
	// NftClassesTable holds the schema information for the "nft_classes" table.
	NftClassesTable = &schema.Table{
		Name:       "nft_classes",
		Columns:    NftClassesColumns,
		PrimaryKey: []*schema.Column{NftClassesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "nft_classes_accounts_nft_classes",
				Columns:    []*schema.Column{NftClassesColumns[14]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		EvmEventsTable,
		NftsTable,
		NftClassesTable,
	}
)

func init() {
	NftsTable.ForeignKeys[0].RefTable = AccountsTable
	NftsTable.ForeignKeys[1].RefTable = NftClassesTable
	NftsTable.Annotation = &entsql.Annotation{
		Table: "nfts",
	}
	NftClassesTable.ForeignKeys[0].RefTable = AccountsTable
}
