// Code generated by ent, DO NOT EDIT.

package nftclass

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the nftclass type in the database.
	Label = "nft_class"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldSymbol holds the string denoting the symbol field in the database.
	FieldSymbol = "symbol"
	// FieldOwnerAddress holds the string denoting the owner_address field in the database.
	FieldOwnerAddress = "owner_address"
	// FieldMinterAddresses holds the string denoting the minter_addresses field in the database.
	FieldMinterAddresses = "minter_addresses"
	// FieldTotalSupply holds the string denoting the total_supply field in the database.
	FieldTotalSupply = "total_supply"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// FieldBannerImage holds the string denoting the banner_image field in the database.
	FieldBannerImage = "banner_image"
	// FieldFeaturedImage holds the string denoting the featured_image field in the database.
	FieldFeaturedImage = "featured_image"
	// FieldDeployerAddress holds the string denoting the deployer_address field in the database.
	FieldDeployerAddress = "deployer_address"
	// FieldDeployedBlockNumber holds the string denoting the deployed_block_number field in the database.
	FieldDeployedBlockNumber = "deployed_block_number"
	// FieldMintedAt holds the string denoting the minted_at field in the database.
	FieldMintedAt = "minted_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeNfts holds the string denoting the nfts edge name in mutations.
	EdgeNfts = "nfts"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// Table holds the table name of the nftclass in the database.
	Table = "nft_classes"
	// NftsTable is the table that holds the nfts relation/edge.
	NftsTable = "nf_ts"
	// NftsInverseTable is the table name for the NFT entity.
	// It exists in this package in order to avoid circular dependency with the "nft" package.
	NftsInverseTable = "nf_ts"
	// NftsColumn is the table column denoting the nfts relation/edge.
	NftsColumn = "nft_class_nfts"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "nft_classes"
	// OwnerInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	OwnerInverseTable = "accounts"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "account_nft_classes"
)

// Columns holds all SQL columns for nftclass fields.
var Columns = []string{
	FieldID,
	FieldAddress,
	FieldName,
	FieldSymbol,
	FieldOwnerAddress,
	FieldMinterAddresses,
	FieldTotalSupply,
	FieldMetadata,
	FieldBannerImage,
	FieldFeaturedImage,
	FieldDeployerAddress,
	FieldDeployedBlockNumber,
	FieldMintedAt,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "nft_classes"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"account_nft_classes",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// SymbolValidator is a validator for the "symbol" field. It is called by the builders before save.
	SymbolValidator func(string) error
	// TotalSupplyValidator is a validator for the "total_supply" field. It is called by the builders before save.
	TotalSupplyValidator func(int) error
	// BannerImageValidator is a validator for the "banner_image" field. It is called by the builders before save.
	BannerImageValidator func(string) error
	// FeaturedImageValidator is a validator for the "featured_image" field. It is called by the builders before save.
	FeaturedImageValidator func(string) error
	// DeployerAddressValidator is a validator for the "deployer_address" field. It is called by the builders before save.
	DeployerAddressValidator func(string) error
	// DeployedBlockNumberValidator is a validator for the "deployed_block_number" field. It is called by the builders before save.
	DeployedBlockNumberValidator func(string) error
)

// OrderOption defines the ordering options for the NFTClass queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByAddress orders the results by the address field.
func ByAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAddress, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// BySymbol orders the results by the symbol field.
func BySymbol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSymbol, opts...).ToFunc()
}

// ByOwnerAddress orders the results by the owner_address field.
func ByOwnerAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOwnerAddress, opts...).ToFunc()
}

// ByTotalSupply orders the results by the total_supply field.
func ByTotalSupply(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTotalSupply, opts...).ToFunc()
}

// ByBannerImage orders the results by the banner_image field.
func ByBannerImage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBannerImage, opts...).ToFunc()
}

// ByFeaturedImage orders the results by the featured_image field.
func ByFeaturedImage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFeaturedImage, opts...).ToFunc()
}

// ByDeployerAddress orders the results by the deployer_address field.
func ByDeployerAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeployerAddress, opts...).ToFunc()
}

// ByDeployedBlockNumber orders the results by the deployed_block_number field.
func ByDeployedBlockNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeployedBlockNumber, opts...).ToFunc()
}

// ByMintedAt orders the results by the minted_at field.
func ByMintedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMintedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByNftsCount orders the results by nfts count.
func ByNftsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newNftsStep(), opts...)
	}
}

// ByNfts orders the results by nfts terms.
func ByNfts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNftsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByOwnerField orders the results by owner field.
func ByOwnerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), sql.OrderByField(field, opts...))
	}
}
func newNftsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NftsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, NftsTable, NftsColumn),
	)
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
	)
}
