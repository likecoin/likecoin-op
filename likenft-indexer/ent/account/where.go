// Code generated by ent, DO NOT EDIT.

package account

import (
	"likenft-indexer/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Account {
	return predicate.Account(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Account {
	return predicate.Account(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Account {
	return predicate.Account(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Account {
	return predicate.Account(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Account {
	return predicate.Account(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Account {
	return predicate.Account(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Account {
	return predicate.Account(sql.FieldLTE(FieldID, id))
}

// CosmosAddress applies equality check predicate on the "cosmos_address" field. It's identical to CosmosAddressEQ.
func CosmosAddress(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldCosmosAddress, v))
}

// EvmAddress applies equality check predicate on the "evm_address" field. It's identical to EvmAddressEQ.
func EvmAddress(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldEvmAddress, v))
}

// Likeid applies equality check predicate on the "likeid" field. It's identical to LikeidEQ.
func Likeid(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldLikeid, v))
}

// CosmosAddressEQ applies the EQ predicate on the "cosmos_address" field.
func CosmosAddressEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldCosmosAddress, v))
}

// CosmosAddressNEQ applies the NEQ predicate on the "cosmos_address" field.
func CosmosAddressNEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldNEQ(FieldCosmosAddress, v))
}

// CosmosAddressIn applies the In predicate on the "cosmos_address" field.
func CosmosAddressIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldIn(FieldCosmosAddress, vs...))
}

// CosmosAddressNotIn applies the NotIn predicate on the "cosmos_address" field.
func CosmosAddressNotIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldNotIn(FieldCosmosAddress, vs...))
}

// CosmosAddressGT applies the GT predicate on the "cosmos_address" field.
func CosmosAddressGT(v string) predicate.Account {
	return predicate.Account(sql.FieldGT(FieldCosmosAddress, v))
}

// CosmosAddressGTE applies the GTE predicate on the "cosmos_address" field.
func CosmosAddressGTE(v string) predicate.Account {
	return predicate.Account(sql.FieldGTE(FieldCosmosAddress, v))
}

// CosmosAddressLT applies the LT predicate on the "cosmos_address" field.
func CosmosAddressLT(v string) predicate.Account {
	return predicate.Account(sql.FieldLT(FieldCosmosAddress, v))
}

// CosmosAddressLTE applies the LTE predicate on the "cosmos_address" field.
func CosmosAddressLTE(v string) predicate.Account {
	return predicate.Account(sql.FieldLTE(FieldCosmosAddress, v))
}

// CosmosAddressContains applies the Contains predicate on the "cosmos_address" field.
func CosmosAddressContains(v string) predicate.Account {
	return predicate.Account(sql.FieldContains(FieldCosmosAddress, v))
}

// CosmosAddressHasPrefix applies the HasPrefix predicate on the "cosmos_address" field.
func CosmosAddressHasPrefix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasPrefix(FieldCosmosAddress, v))
}

// CosmosAddressHasSuffix applies the HasSuffix predicate on the "cosmos_address" field.
func CosmosAddressHasSuffix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasSuffix(FieldCosmosAddress, v))
}

// CosmosAddressIsNil applies the IsNil predicate on the "cosmos_address" field.
func CosmosAddressIsNil() predicate.Account {
	return predicate.Account(sql.FieldIsNull(FieldCosmosAddress))
}

// CosmosAddressNotNil applies the NotNil predicate on the "cosmos_address" field.
func CosmosAddressNotNil() predicate.Account {
	return predicate.Account(sql.FieldNotNull(FieldCosmosAddress))
}

// CosmosAddressEqualFold applies the EqualFold predicate on the "cosmos_address" field.
func CosmosAddressEqualFold(v string) predicate.Account {
	return predicate.Account(sql.FieldEqualFold(FieldCosmosAddress, v))
}

// CosmosAddressContainsFold applies the ContainsFold predicate on the "cosmos_address" field.
func CosmosAddressContainsFold(v string) predicate.Account {
	return predicate.Account(sql.FieldContainsFold(FieldCosmosAddress, v))
}

// EvmAddressEQ applies the EQ predicate on the "evm_address" field.
func EvmAddressEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldEvmAddress, v))
}

// EvmAddressNEQ applies the NEQ predicate on the "evm_address" field.
func EvmAddressNEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldNEQ(FieldEvmAddress, v))
}

// EvmAddressIn applies the In predicate on the "evm_address" field.
func EvmAddressIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldIn(FieldEvmAddress, vs...))
}

// EvmAddressNotIn applies the NotIn predicate on the "evm_address" field.
func EvmAddressNotIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldNotIn(FieldEvmAddress, vs...))
}

// EvmAddressGT applies the GT predicate on the "evm_address" field.
func EvmAddressGT(v string) predicate.Account {
	return predicate.Account(sql.FieldGT(FieldEvmAddress, v))
}

// EvmAddressGTE applies the GTE predicate on the "evm_address" field.
func EvmAddressGTE(v string) predicate.Account {
	return predicate.Account(sql.FieldGTE(FieldEvmAddress, v))
}

// EvmAddressLT applies the LT predicate on the "evm_address" field.
func EvmAddressLT(v string) predicate.Account {
	return predicate.Account(sql.FieldLT(FieldEvmAddress, v))
}

// EvmAddressLTE applies the LTE predicate on the "evm_address" field.
func EvmAddressLTE(v string) predicate.Account {
	return predicate.Account(sql.FieldLTE(FieldEvmAddress, v))
}

// EvmAddressContains applies the Contains predicate on the "evm_address" field.
func EvmAddressContains(v string) predicate.Account {
	return predicate.Account(sql.FieldContains(FieldEvmAddress, v))
}

// EvmAddressHasPrefix applies the HasPrefix predicate on the "evm_address" field.
func EvmAddressHasPrefix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasPrefix(FieldEvmAddress, v))
}

// EvmAddressHasSuffix applies the HasSuffix predicate on the "evm_address" field.
func EvmAddressHasSuffix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasSuffix(FieldEvmAddress, v))
}

// EvmAddressEqualFold applies the EqualFold predicate on the "evm_address" field.
func EvmAddressEqualFold(v string) predicate.Account {
	return predicate.Account(sql.FieldEqualFold(FieldEvmAddress, v))
}

// EvmAddressContainsFold applies the ContainsFold predicate on the "evm_address" field.
func EvmAddressContainsFold(v string) predicate.Account {
	return predicate.Account(sql.FieldContainsFold(FieldEvmAddress, v))
}

// LikeidEQ applies the EQ predicate on the "likeid" field.
func LikeidEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldEQ(FieldLikeid, v))
}

// LikeidNEQ applies the NEQ predicate on the "likeid" field.
func LikeidNEQ(v string) predicate.Account {
	return predicate.Account(sql.FieldNEQ(FieldLikeid, v))
}

// LikeidIn applies the In predicate on the "likeid" field.
func LikeidIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldIn(FieldLikeid, vs...))
}

// LikeidNotIn applies the NotIn predicate on the "likeid" field.
func LikeidNotIn(vs ...string) predicate.Account {
	return predicate.Account(sql.FieldNotIn(FieldLikeid, vs...))
}

// LikeidGT applies the GT predicate on the "likeid" field.
func LikeidGT(v string) predicate.Account {
	return predicate.Account(sql.FieldGT(FieldLikeid, v))
}

// LikeidGTE applies the GTE predicate on the "likeid" field.
func LikeidGTE(v string) predicate.Account {
	return predicate.Account(sql.FieldGTE(FieldLikeid, v))
}

// LikeidLT applies the LT predicate on the "likeid" field.
func LikeidLT(v string) predicate.Account {
	return predicate.Account(sql.FieldLT(FieldLikeid, v))
}

// LikeidLTE applies the LTE predicate on the "likeid" field.
func LikeidLTE(v string) predicate.Account {
	return predicate.Account(sql.FieldLTE(FieldLikeid, v))
}

// LikeidContains applies the Contains predicate on the "likeid" field.
func LikeidContains(v string) predicate.Account {
	return predicate.Account(sql.FieldContains(FieldLikeid, v))
}

// LikeidHasPrefix applies the HasPrefix predicate on the "likeid" field.
func LikeidHasPrefix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasPrefix(FieldLikeid, v))
}

// LikeidHasSuffix applies the HasSuffix predicate on the "likeid" field.
func LikeidHasSuffix(v string) predicate.Account {
	return predicate.Account(sql.FieldHasSuffix(FieldLikeid, v))
}

// LikeidIsNil applies the IsNil predicate on the "likeid" field.
func LikeidIsNil() predicate.Account {
	return predicate.Account(sql.FieldIsNull(FieldLikeid))
}

// LikeidNotNil applies the NotNil predicate on the "likeid" field.
func LikeidNotNil() predicate.Account {
	return predicate.Account(sql.FieldNotNull(FieldLikeid))
}

// LikeidEqualFold applies the EqualFold predicate on the "likeid" field.
func LikeidEqualFold(v string) predicate.Account {
	return predicate.Account(sql.FieldEqualFold(FieldLikeid, v))
}

// LikeidContainsFold applies the ContainsFold predicate on the "likeid" field.
func LikeidContainsFold(v string) predicate.Account {
	return predicate.Account(sql.FieldContainsFold(FieldLikeid, v))
}

// HasNftClasses applies the HasEdge predicate on the "nft_classes" edge.
func HasNftClasses() predicate.Account {
	return predicate.Account(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, NftClassesTable, NftClassesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasNftClassesWith applies the HasEdge predicate on the "nft_classes" edge with a given conditions (other predicates).
func HasNftClassesWith(preds ...predicate.NFTClass) predicate.Account {
	return predicate.Account(func(s *sql.Selector) {
		step := newNftClassesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasNfts applies the HasEdge predicate on the "nfts" edge.
func HasNfts() predicate.Account {
	return predicate.Account(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, NftsTable, NftsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasNftsWith applies the HasEdge predicate on the "nfts" edge with a given conditions (other predicates).
func HasNftsWith(preds ...predicate.NFT) predicate.Account {
	return predicate.Account(func(s *sql.Selector) {
		step := newNftsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Account) predicate.Account {
	return predicate.Account(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Account) predicate.Account {
	return predicate.Account(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Account) predicate.Account {
	return predicate.Account(sql.NotPredicates(p))
}
