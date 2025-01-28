// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"likenft-indexer/ent/account"
	"likenft-indexer/ent/nft"
	"likenft-indexer/ent/nftclass"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// NFT is the model entity for the NFT schema.
type NFT struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TokenID holds the value of the "token_id" field.
	TokenID uint64 `json:"token_id,omitempty"`
	// TokenURL holds the value of the "token_url" field.
	TokenURL string `json:"token_url,omitempty"`
	// Raw holds the value of the "raw" field.
	Raw map[string]interface{} `json:"raw,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Image holds the value of the "image" field.
	Image map[string]interface{} `json:"image,omitempty"`
	// Attributes holds the value of the "attributes" field.
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	// OwnerAddress holds the value of the "owner_address" field.
	OwnerAddress string `json:"owner_address,omitempty"`
	// MintedAt holds the value of the "minted_at" field.
	MintedAt time.Time `json:"minted_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NFTQuery when eager-loading is set.
	Edges          NFTEdges `json:"edges"`
	account_nfts   *int
	nft_class_nfts *int
	selectValues   sql.SelectValues
}

// NFTEdges holds the relations/edges for other nodes in the graph.
type NFTEdges struct {
	// Owner holds the value of the owner edge.
	Owner *Account `json:"owner,omitempty"`
	// Class holds the value of the class edge.
	Class *NFTClass `json:"class,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e NFTEdges) OwnerOrErr() (*Account, error) {
	if e.Owner != nil {
		return e.Owner, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: account.Label}
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// ClassOrErr returns the Class value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e NFTEdges) ClassOrErr() (*NFTClass, error) {
	if e.Class != nil {
		return e.Class, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: nftclass.Label}
	}
	return nil, &NotLoadedError{edge: "class"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*NFT) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case nft.FieldRaw, nft.FieldImage, nft.FieldAttributes:
			values[i] = new([]byte)
		case nft.FieldID, nft.FieldTokenID:
			values[i] = new(sql.NullInt64)
		case nft.FieldTokenURL, nft.FieldName, nft.FieldDescription, nft.FieldOwnerAddress:
			values[i] = new(sql.NullString)
		case nft.FieldMintedAt, nft.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case nft.ForeignKeys[0]: // account_nfts
			values[i] = new(sql.NullInt64)
		case nft.ForeignKeys[1]: // nft_class_nfts
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the NFT fields.
func (n *NFT) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case nft.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			n.ID = int(value.Int64)
		case nft.FieldTokenID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field token_id", values[i])
			} else if value.Valid {
				n.TokenID = uint64(value.Int64)
			}
		case nft.FieldTokenURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field token_url", values[i])
			} else if value.Valid {
				n.TokenURL = value.String
			}
		case nft.FieldRaw:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field raw", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &n.Raw); err != nil {
					return fmt.Errorf("unmarshal field raw: %w", err)
				}
			}
		case nft.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				n.Name = value.String
			}
		case nft.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				n.Description = value.String
			}
		case nft.FieldImage:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field image", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &n.Image); err != nil {
					return fmt.Errorf("unmarshal field image: %w", err)
				}
			}
		case nft.FieldAttributes:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field attributes", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &n.Attributes); err != nil {
					return fmt.Errorf("unmarshal field attributes: %w", err)
				}
			}
		case nft.FieldOwnerAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field owner_address", values[i])
			} else if value.Valid {
				n.OwnerAddress = value.String
			}
		case nft.FieldMintedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field minted_at", values[i])
			} else if value.Valid {
				n.MintedAt = value.Time
			}
		case nft.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				n.UpdatedAt = value.Time
			}
		case nft.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field account_nfts", value)
			} else if value.Valid {
				n.account_nfts = new(int)
				*n.account_nfts = int(value.Int64)
			}
		case nft.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field nft_class_nfts", value)
			} else if value.Valid {
				n.nft_class_nfts = new(int)
				*n.nft_class_nfts = int(value.Int64)
			}
		default:
			n.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the NFT.
// This includes values selected through modifiers, order, etc.
func (n *NFT) Value(name string) (ent.Value, error) {
	return n.selectValues.Get(name)
}

// QueryOwner queries the "owner" edge of the NFT entity.
func (n *NFT) QueryOwner() *AccountQuery {
	return NewNFTClient(n.config).QueryOwner(n)
}

// QueryClass queries the "class" edge of the NFT entity.
func (n *NFT) QueryClass() *NFTClassQuery {
	return NewNFTClient(n.config).QueryClass(n)
}

// Update returns a builder for updating this NFT.
// Note that you need to call NFT.Unwrap() before calling this method if this NFT
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *NFT) Update() *NFTUpdateOne {
	return NewNFTClient(n.config).UpdateOne(n)
}

// Unwrap unwraps the NFT entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (n *NFT) Unwrap() *NFT {
	_tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: NFT is not a transactional entity")
	}
	n.config.driver = _tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *NFT) String() string {
	var builder strings.Builder
	builder.WriteString("NFT(")
	builder.WriteString(fmt.Sprintf("id=%v, ", n.ID))
	builder.WriteString("token_id=")
	builder.WriteString(fmt.Sprintf("%v", n.TokenID))
	builder.WriteString(", ")
	builder.WriteString("token_url=")
	builder.WriteString(n.TokenURL)
	builder.WriteString(", ")
	builder.WriteString("raw=")
	builder.WriteString(fmt.Sprintf("%v", n.Raw))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(n.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(n.Description)
	builder.WriteString(", ")
	builder.WriteString("image=")
	builder.WriteString(fmt.Sprintf("%v", n.Image))
	builder.WriteString(", ")
	builder.WriteString("attributes=")
	builder.WriteString(fmt.Sprintf("%v", n.Attributes))
	builder.WriteString(", ")
	builder.WriteString("owner_address=")
	builder.WriteString(n.OwnerAddress)
	builder.WriteString(", ")
	builder.WriteString("minted_at=")
	builder.WriteString(n.MintedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(n.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// NFTs is a parsable slice of NFT.
type NFTs []*NFT
