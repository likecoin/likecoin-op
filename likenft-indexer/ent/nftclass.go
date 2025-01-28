// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"likenft-indexer/ent/account"
	"likenft-indexer/ent/nftclass"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// NFTClass is the model entity for the NFTClass schema.
type NFTClass struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Address holds the value of the "address" field.
	Address string `json:"address,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Symbol holds the value of the "symbol" field.
	Symbol string `json:"symbol,omitempty"`
	// OwnerAddress holds the value of the "owner_address" field.
	OwnerAddress *string `json:"owner_address,omitempty"`
	// MinterAddresses holds the value of the "minter_addresses" field.
	MinterAddresses []string `json:"minter_addresses,omitempty"`
	// TotalSupply holds the value of the "total_supply" field.
	TotalSupply int `json:"total_supply,omitempty"`
	// Metadata holds the value of the "metadata" field.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// BannerImage holds the value of the "banner_image" field.
	BannerImage string `json:"banner_image,omitempty"`
	// FeaturedImage holds the value of the "featured_image" field.
	FeaturedImage string `json:"featured_image,omitempty"`
	// DeployerAddress holds the value of the "deployer_address" field.
	DeployerAddress string `json:"deployer_address,omitempty"`
	// DeployedBlockNumber holds the value of the "deployed_block_number" field.
	DeployedBlockNumber string `json:"deployed_block_number,omitempty"`
	// MintedAt holds the value of the "minted_at" field.
	MintedAt time.Time `json:"minted_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NFTClassQuery when eager-loading is set.
	Edges               NFTClassEdges `json:"edges"`
	account_nft_classes *int
	selectValues        sql.SelectValues
}

// NFTClassEdges holds the relations/edges for other nodes in the graph.
type NFTClassEdges struct {
	// Nfts holds the value of the nfts edge.
	Nfts []*NFT `json:"nfts,omitempty"`
	// Owner holds the value of the owner edge.
	Owner *Account `json:"owner,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// NftsOrErr returns the Nfts value or an error if the edge
// was not loaded in eager-loading.
func (e NFTClassEdges) NftsOrErr() ([]*NFT, error) {
	if e.loadedTypes[0] {
		return e.Nfts, nil
	}
	return nil, &NotLoadedError{edge: "nfts"}
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e NFTClassEdges) OwnerOrErr() (*Account, error) {
	if e.Owner != nil {
		return e.Owner, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: account.Label}
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*NFTClass) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case nftclass.FieldMinterAddresses, nftclass.FieldMetadata:
			values[i] = new([]byte)
		case nftclass.FieldID, nftclass.FieldTotalSupply:
			values[i] = new(sql.NullInt64)
		case nftclass.FieldAddress, nftclass.FieldName, nftclass.FieldSymbol, nftclass.FieldOwnerAddress, nftclass.FieldBannerImage, nftclass.FieldFeaturedImage, nftclass.FieldDeployerAddress, nftclass.FieldDeployedBlockNumber:
			values[i] = new(sql.NullString)
		case nftclass.FieldMintedAt, nftclass.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case nftclass.ForeignKeys[0]: // account_nft_classes
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the NFTClass fields.
func (nc *NFTClass) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case nftclass.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			nc.ID = int(value.Int64)
		case nftclass.FieldAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address", values[i])
			} else if value.Valid {
				nc.Address = value.String
			}
		case nftclass.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				nc.Name = value.String
			}
		case nftclass.FieldSymbol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field symbol", values[i])
			} else if value.Valid {
				nc.Symbol = value.String
			}
		case nftclass.FieldOwnerAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field owner_address", values[i])
			} else if value.Valid {
				nc.OwnerAddress = new(string)
				*nc.OwnerAddress = value.String
			}
		case nftclass.FieldMinterAddresses:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field minter_addresses", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &nc.MinterAddresses); err != nil {
					return fmt.Errorf("unmarshal field minter_addresses: %w", err)
				}
			}
		case nftclass.FieldTotalSupply:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field total_supply", values[i])
			} else if value.Valid {
				nc.TotalSupply = int(value.Int64)
			}
		case nftclass.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &nc.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case nftclass.FieldBannerImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field banner_image", values[i])
			} else if value.Valid {
				nc.BannerImage = value.String
			}
		case nftclass.FieldFeaturedImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field featured_image", values[i])
			} else if value.Valid {
				nc.FeaturedImage = value.String
			}
		case nftclass.FieldDeployerAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field deployer_address", values[i])
			} else if value.Valid {
				nc.DeployerAddress = value.String
			}
		case nftclass.FieldDeployedBlockNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field deployed_block_number", values[i])
			} else if value.Valid {
				nc.DeployedBlockNumber = value.String
			}
		case nftclass.FieldMintedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field minted_at", values[i])
			} else if value.Valid {
				nc.MintedAt = value.Time
			}
		case nftclass.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				nc.UpdatedAt = value.Time
			}
		case nftclass.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field account_nft_classes", value)
			} else if value.Valid {
				nc.account_nft_classes = new(int)
				*nc.account_nft_classes = int(value.Int64)
			}
		default:
			nc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the NFTClass.
// This includes values selected through modifiers, order, etc.
func (nc *NFTClass) Value(name string) (ent.Value, error) {
	return nc.selectValues.Get(name)
}

// QueryNfts queries the "nfts" edge of the NFTClass entity.
func (nc *NFTClass) QueryNfts() *NFTQuery {
	return NewNFTClassClient(nc.config).QueryNfts(nc)
}

// QueryOwner queries the "owner" edge of the NFTClass entity.
func (nc *NFTClass) QueryOwner() *AccountQuery {
	return NewNFTClassClient(nc.config).QueryOwner(nc)
}

// Update returns a builder for updating this NFTClass.
// Note that you need to call NFTClass.Unwrap() before calling this method if this NFTClass
// was returned from a transaction, and the transaction was committed or rolled back.
func (nc *NFTClass) Update() *NFTClassUpdateOne {
	return NewNFTClassClient(nc.config).UpdateOne(nc)
}

// Unwrap unwraps the NFTClass entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (nc *NFTClass) Unwrap() *NFTClass {
	_tx, ok := nc.config.driver.(*txDriver)
	if !ok {
		panic("ent: NFTClass is not a transactional entity")
	}
	nc.config.driver = _tx.drv
	return nc
}

// String implements the fmt.Stringer.
func (nc *NFTClass) String() string {
	var builder strings.Builder
	builder.WriteString("NFTClass(")
	builder.WriteString(fmt.Sprintf("id=%v, ", nc.ID))
	builder.WriteString("address=")
	builder.WriteString(nc.Address)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(nc.Name)
	builder.WriteString(", ")
	builder.WriteString("symbol=")
	builder.WriteString(nc.Symbol)
	builder.WriteString(", ")
	if v := nc.OwnerAddress; v != nil {
		builder.WriteString("owner_address=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("minter_addresses=")
	builder.WriteString(fmt.Sprintf("%v", nc.MinterAddresses))
	builder.WriteString(", ")
	builder.WriteString("total_supply=")
	builder.WriteString(fmt.Sprintf("%v", nc.TotalSupply))
	builder.WriteString(", ")
	builder.WriteString("metadata=")
	builder.WriteString(fmt.Sprintf("%v", nc.Metadata))
	builder.WriteString(", ")
	builder.WriteString("banner_image=")
	builder.WriteString(nc.BannerImage)
	builder.WriteString(", ")
	builder.WriteString("featured_image=")
	builder.WriteString(nc.FeaturedImage)
	builder.WriteString(", ")
	builder.WriteString("deployer_address=")
	builder.WriteString(nc.DeployerAddress)
	builder.WriteString(", ")
	builder.WriteString("deployed_block_number=")
	builder.WriteString(nc.DeployedBlockNumber)
	builder.WriteString(", ")
	builder.WriteString("minted_at=")
	builder.WriteString(nc.MintedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(nc.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// NFTClasses is a parsable slice of NFTClass.
type NFTClasses []*NFTClass
