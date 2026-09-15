package schema

import (
	"likecollective-indexer/ent/schema/typeutil"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// StakingStateHead records the newest chain head accounts, nft_classes and
// stakings have been re-read at. It holds at most one row.
type StakingStateHead struct {
	ent.Schema
}

// Fields of the StakingStateHead.
func (StakingStateHead) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("block_number").
			GoType(typeutil.Uint64(0)).
			SchemaType(typeutil.Uint64SchemaType).
			ValueScanner(typeutil.Uint64ValueScanner).
			Annotations(typeutil.Uint64Annotations("block_number")...),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
