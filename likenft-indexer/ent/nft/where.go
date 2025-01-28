// Code generated by ent, DO NOT EDIT.

package nft

import (
	"likenft-indexer/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldID, id))
}

// TokenID applies equality check predicate on the "token_id" field. It's identical to TokenIDEQ.
func TokenID(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldTokenID, v))
}

// TokenURL applies equality check predicate on the "token_url" field. It's identical to TokenURLEQ.
func TokenURL(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldTokenURL, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldDescription, v))
}

// OwnerAddress applies equality check predicate on the "owner_address" field. It's identical to OwnerAddressEQ.
func OwnerAddress(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldOwnerAddress, v))
}

// MintedAt applies equality check predicate on the "minted_at" field. It's identical to MintedAtEQ.
func MintedAt(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldMintedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldUpdatedAt, v))
}

// TokenIDEQ applies the EQ predicate on the "token_id" field.
func TokenIDEQ(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldTokenID, v))
}

// TokenIDNEQ applies the NEQ predicate on the "token_id" field.
func TokenIDNEQ(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldTokenID, v))
}

// TokenIDIn applies the In predicate on the "token_id" field.
func TokenIDIn(vs ...uint64) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldTokenID, vs...))
}

// TokenIDNotIn applies the NotIn predicate on the "token_id" field.
func TokenIDNotIn(vs ...uint64) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldTokenID, vs...))
}

// TokenIDGT applies the GT predicate on the "token_id" field.
func TokenIDGT(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldTokenID, v))
}

// TokenIDGTE applies the GTE predicate on the "token_id" field.
func TokenIDGTE(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldTokenID, v))
}

// TokenIDLT applies the LT predicate on the "token_id" field.
func TokenIDLT(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldTokenID, v))
}

// TokenIDLTE applies the LTE predicate on the "token_id" field.
func TokenIDLTE(v uint64) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldTokenID, v))
}

// TokenURLEQ applies the EQ predicate on the "token_url" field.
func TokenURLEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldTokenURL, v))
}

// TokenURLNEQ applies the NEQ predicate on the "token_url" field.
func TokenURLNEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldTokenURL, v))
}

// TokenURLIn applies the In predicate on the "token_url" field.
func TokenURLIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldTokenURL, vs...))
}

// TokenURLNotIn applies the NotIn predicate on the "token_url" field.
func TokenURLNotIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldTokenURL, vs...))
}

// TokenURLGT applies the GT predicate on the "token_url" field.
func TokenURLGT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldTokenURL, v))
}

// TokenURLGTE applies the GTE predicate on the "token_url" field.
func TokenURLGTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldTokenURL, v))
}

// TokenURLLT applies the LT predicate on the "token_url" field.
func TokenURLLT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldTokenURL, v))
}

// TokenURLLTE applies the LTE predicate on the "token_url" field.
func TokenURLLTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldTokenURL, v))
}

// TokenURLContains applies the Contains predicate on the "token_url" field.
func TokenURLContains(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContains(FieldTokenURL, v))
}

// TokenURLHasPrefix applies the HasPrefix predicate on the "token_url" field.
func TokenURLHasPrefix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasPrefix(FieldTokenURL, v))
}

// TokenURLHasSuffix applies the HasSuffix predicate on the "token_url" field.
func TokenURLHasSuffix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasSuffix(FieldTokenURL, v))
}

// TokenURLIsNil applies the IsNil predicate on the "token_url" field.
func TokenURLIsNil() predicate.NFT {
	return predicate.NFT(sql.FieldIsNull(FieldTokenURL))
}

// TokenURLNotNil applies the NotNil predicate on the "token_url" field.
func TokenURLNotNil() predicate.NFT {
	return predicate.NFT(sql.FieldNotNull(FieldTokenURL))
}

// TokenURLEqualFold applies the EqualFold predicate on the "token_url" field.
func TokenURLEqualFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEqualFold(FieldTokenURL, v))
}

// TokenURLContainsFold applies the ContainsFold predicate on the "token_url" field.
func TokenURLContainsFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContainsFold(FieldTokenURL, v))
}

// RawIsNil applies the IsNil predicate on the "raw" field.
func RawIsNil() predicate.NFT {
	return predicate.NFT(sql.FieldIsNull(FieldRaw))
}

// RawNotNil applies the NotNil predicate on the "raw" field.
func RawNotNil() predicate.NFT {
	return predicate.NFT(sql.FieldNotNull(FieldRaw))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.NFT {
	return predicate.NFT(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.NFT {
	return predicate.NFT(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContainsFold(FieldDescription, v))
}

// ImageIsNil applies the IsNil predicate on the "image" field.
func ImageIsNil() predicate.NFT {
	return predicate.NFT(sql.FieldIsNull(FieldImage))
}

// ImageNotNil applies the NotNil predicate on the "image" field.
func ImageNotNil() predicate.NFT {
	return predicate.NFT(sql.FieldNotNull(FieldImage))
}

// AttributesIsNil applies the IsNil predicate on the "attributes" field.
func AttributesIsNil() predicate.NFT {
	return predicate.NFT(sql.FieldIsNull(FieldAttributes))
}

// AttributesNotNil applies the NotNil predicate on the "attributes" field.
func AttributesNotNil() predicate.NFT {
	return predicate.NFT(sql.FieldNotNull(FieldAttributes))
}

// OwnerAddressEQ applies the EQ predicate on the "owner_address" field.
func OwnerAddressEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldOwnerAddress, v))
}

// OwnerAddressNEQ applies the NEQ predicate on the "owner_address" field.
func OwnerAddressNEQ(v string) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldOwnerAddress, v))
}

// OwnerAddressIn applies the In predicate on the "owner_address" field.
func OwnerAddressIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldOwnerAddress, vs...))
}

// OwnerAddressNotIn applies the NotIn predicate on the "owner_address" field.
func OwnerAddressNotIn(vs ...string) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldOwnerAddress, vs...))
}

// OwnerAddressGT applies the GT predicate on the "owner_address" field.
func OwnerAddressGT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldOwnerAddress, v))
}

// OwnerAddressGTE applies the GTE predicate on the "owner_address" field.
func OwnerAddressGTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldOwnerAddress, v))
}

// OwnerAddressLT applies the LT predicate on the "owner_address" field.
func OwnerAddressLT(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldOwnerAddress, v))
}

// OwnerAddressLTE applies the LTE predicate on the "owner_address" field.
func OwnerAddressLTE(v string) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldOwnerAddress, v))
}

// OwnerAddressContains applies the Contains predicate on the "owner_address" field.
func OwnerAddressContains(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContains(FieldOwnerAddress, v))
}

// OwnerAddressHasPrefix applies the HasPrefix predicate on the "owner_address" field.
func OwnerAddressHasPrefix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasPrefix(FieldOwnerAddress, v))
}

// OwnerAddressHasSuffix applies the HasSuffix predicate on the "owner_address" field.
func OwnerAddressHasSuffix(v string) predicate.NFT {
	return predicate.NFT(sql.FieldHasSuffix(FieldOwnerAddress, v))
}

// OwnerAddressEqualFold applies the EqualFold predicate on the "owner_address" field.
func OwnerAddressEqualFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldEqualFold(FieldOwnerAddress, v))
}

// OwnerAddressContainsFold applies the ContainsFold predicate on the "owner_address" field.
func OwnerAddressContainsFold(v string) predicate.NFT {
	return predicate.NFT(sql.FieldContainsFold(FieldOwnerAddress, v))
}

// MintedAtEQ applies the EQ predicate on the "minted_at" field.
func MintedAtEQ(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldMintedAt, v))
}

// MintedAtNEQ applies the NEQ predicate on the "minted_at" field.
func MintedAtNEQ(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldMintedAt, v))
}

// MintedAtIn applies the In predicate on the "minted_at" field.
func MintedAtIn(vs ...time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldMintedAt, vs...))
}

// MintedAtNotIn applies the NotIn predicate on the "minted_at" field.
func MintedAtNotIn(vs ...time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldMintedAt, vs...))
}

// MintedAtGT applies the GT predicate on the "minted_at" field.
func MintedAtGT(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldMintedAt, v))
}

// MintedAtGTE applies the GTE predicate on the "minted_at" field.
func MintedAtGTE(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldMintedAt, v))
}

// MintedAtLT applies the LT predicate on the "minted_at" field.
func MintedAtLT(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldMintedAt, v))
}

// MintedAtLTE applies the LTE predicate on the "minted_at" field.
func MintedAtLTE(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldMintedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.NFT {
	return predicate.NFT(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.NFT {
	return predicate.NFT(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.Account) predicate.NFT {
	return predicate.NFT(func(s *sql.Selector) {
		step := newOwnerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasClass applies the HasEdge predicate on the "class" edge.
func HasClass() predicate.NFT {
	return predicate.NFT(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ClassTable, ClassColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasClassWith applies the HasEdge predicate on the "class" edge with a given conditions (other predicates).
func HasClassWith(preds ...predicate.NFTClass) predicate.NFT {
	return predicate.NFT(func(s *sql.Selector) {
		step := newClassStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.NFT) predicate.NFT {
	return predicate.NFT(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.NFT) predicate.NFT {
	return predicate.NFT(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.NFT) predicate.NFT {
	return predicate.NFT(sql.NotPredicates(p))
}
