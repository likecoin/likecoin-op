// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"likenft-indexer/ent/account"
	"likenft-indexer/ent/nft"
	"likenft-indexer/ent/nftclass"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/evm/model"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// NFTCreate is the builder for creating a NFT entity.
type NFTCreate struct {
	config
	mutation *NFTMutation
	hooks    []Hook
}

// SetContractAddress sets the "contract_address" field.
func (nc *NFTCreate) SetContractAddress(s string) *NFTCreate {
	nc.mutation.SetContractAddress(s)
	return nc
}

// SetTokenID sets the "token_id" field.
func (nc *NFTCreate) SetTokenID(t typeutil.Uint64) *NFTCreate {
	nc.mutation.SetTokenID(t)
	return nc
}

// SetTokenURI sets the "token_uri" field.
func (nc *NFTCreate) SetTokenURI(s string) *NFTCreate {
	nc.mutation.SetTokenURI(s)
	return nc
}

// SetNillableTokenURI sets the "token_uri" field if the given value is not nil.
func (nc *NFTCreate) SetNillableTokenURI(s *string) *NFTCreate {
	if s != nil {
		nc.SetTokenURI(*s)
	}
	return nc
}

// SetImage sets the "image" field.
func (nc *NFTCreate) SetImage(s string) *NFTCreate {
	nc.mutation.SetImage(s)
	return nc
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (nc *NFTCreate) SetNillableImage(s *string) *NFTCreate {
	if s != nil {
		nc.SetImage(*s)
	}
	return nc
}

// SetImageData sets the "image_data" field.
func (nc *NFTCreate) SetImageData(s string) *NFTCreate {
	nc.mutation.SetImageData(s)
	return nc
}

// SetNillableImageData sets the "image_data" field if the given value is not nil.
func (nc *NFTCreate) SetNillableImageData(s *string) *NFTCreate {
	if s != nil {
		nc.SetImageData(*s)
	}
	return nc
}

// SetExternalURL sets the "external_url" field.
func (nc *NFTCreate) SetExternalURL(s string) *NFTCreate {
	nc.mutation.SetExternalURL(s)
	return nc
}

// SetNillableExternalURL sets the "external_url" field if the given value is not nil.
func (nc *NFTCreate) SetNillableExternalURL(s *string) *NFTCreate {
	if s != nil {
		nc.SetExternalURL(*s)
	}
	return nc
}

// SetDescription sets the "description" field.
func (nc *NFTCreate) SetDescription(s string) *NFTCreate {
	nc.mutation.SetDescription(s)
	return nc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (nc *NFTCreate) SetNillableDescription(s *string) *NFTCreate {
	if s != nil {
		nc.SetDescription(*s)
	}
	return nc
}

// SetName sets the "name" field.
func (nc *NFTCreate) SetName(s string) *NFTCreate {
	nc.mutation.SetName(s)
	return nc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (nc *NFTCreate) SetNillableName(s *string) *NFTCreate {
	if s != nil {
		nc.SetName(*s)
	}
	return nc
}

// SetAttributes sets the "attributes" field.
func (nc *NFTCreate) SetAttributes(ma []model.ERC721MetadataAttribute) *NFTCreate {
	nc.mutation.SetAttributes(ma)
	return nc
}

// SetBackgroundColor sets the "background_color" field.
func (nc *NFTCreate) SetBackgroundColor(s string) *NFTCreate {
	nc.mutation.SetBackgroundColor(s)
	return nc
}

// SetNillableBackgroundColor sets the "background_color" field if the given value is not nil.
func (nc *NFTCreate) SetNillableBackgroundColor(s *string) *NFTCreate {
	if s != nil {
		nc.SetBackgroundColor(*s)
	}
	return nc
}

// SetAnimationURL sets the "animation_url" field.
func (nc *NFTCreate) SetAnimationURL(s string) *NFTCreate {
	nc.mutation.SetAnimationURL(s)
	return nc
}

// SetNillableAnimationURL sets the "animation_url" field if the given value is not nil.
func (nc *NFTCreate) SetNillableAnimationURL(s *string) *NFTCreate {
	if s != nil {
		nc.SetAnimationURL(*s)
	}
	return nc
}

// SetYoutubeURL sets the "youtube_url" field.
func (nc *NFTCreate) SetYoutubeURL(s string) *NFTCreate {
	nc.mutation.SetYoutubeURL(s)
	return nc
}

// SetNillableYoutubeURL sets the "youtube_url" field if the given value is not nil.
func (nc *NFTCreate) SetNillableYoutubeURL(s *string) *NFTCreate {
	if s != nil {
		nc.SetYoutubeURL(*s)
	}
	return nc
}

// SetOwnerAddress sets the "owner_address" field.
func (nc *NFTCreate) SetOwnerAddress(s string) *NFTCreate {
	nc.mutation.SetOwnerAddress(s)
	return nc
}

// SetMintedAt sets the "minted_at" field.
func (nc *NFTCreate) SetMintedAt(t time.Time) *NFTCreate {
	nc.mutation.SetMintedAt(t)
	return nc
}

// SetUpdatedAt sets the "updated_at" field.
func (nc *NFTCreate) SetUpdatedAt(t time.Time) *NFTCreate {
	nc.mutation.SetUpdatedAt(t)
	return nc
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (nc *NFTCreate) SetOwnerID(id int) *NFTCreate {
	nc.mutation.SetOwnerID(id)
	return nc
}

// SetNillableOwnerID sets the "owner" edge to the Account entity by ID if the given value is not nil.
func (nc *NFTCreate) SetNillableOwnerID(id *int) *NFTCreate {
	if id != nil {
		nc = nc.SetOwnerID(*id)
	}
	return nc
}

// SetOwner sets the "owner" edge to the Account entity.
func (nc *NFTCreate) SetOwner(a *Account) *NFTCreate {
	return nc.SetOwnerID(a.ID)
}

// SetClassID sets the "class" edge to the NFTClass entity by ID.
func (nc *NFTCreate) SetClassID(id int) *NFTCreate {
	nc.mutation.SetClassID(id)
	return nc
}

// SetNillableClassID sets the "class" edge to the NFTClass entity by ID if the given value is not nil.
func (nc *NFTCreate) SetNillableClassID(id *int) *NFTCreate {
	if id != nil {
		nc = nc.SetClassID(*id)
	}
	return nc
}

// SetClass sets the "class" edge to the NFTClass entity.
func (nc *NFTCreate) SetClass(n *NFTClass) *NFTCreate {
	return nc.SetClassID(n.ID)
}

// Mutation returns the NFTMutation object of the builder.
func (nc *NFTCreate) Mutation() *NFTMutation {
	return nc.mutation
}

// Save creates the NFT in the database.
func (nc *NFTCreate) Save(ctx context.Context) (*NFT, error) {
	return withHooks(ctx, nc.sqlSave, nc.mutation, nc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (nc *NFTCreate) SaveX(ctx context.Context) *NFT {
	v, err := nc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nc *NFTCreate) Exec(ctx context.Context) error {
	_, err := nc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nc *NFTCreate) ExecX(ctx context.Context) {
	if err := nc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nc *NFTCreate) check() error {
	if _, ok := nc.mutation.ContractAddress(); !ok {
		return &ValidationError{Name: "contract_address", err: errors.New(`ent: missing required field "NFT.contract_address"`)}
	}
	if v, ok := nc.mutation.ContractAddress(); ok {
		if err := nft.ContractAddressValidator(v); err != nil {
			return &ValidationError{Name: "contract_address", err: fmt.Errorf(`ent: validator failed for field "NFT.contract_address": %w`, err)}
		}
	}
	if _, ok := nc.mutation.TokenID(); !ok {
		return &ValidationError{Name: "token_id", err: errors.New(`ent: missing required field "NFT.token_id"`)}
	}
	if _, ok := nc.mutation.OwnerAddress(); !ok {
		return &ValidationError{Name: "owner_address", err: errors.New(`ent: missing required field "NFT.owner_address"`)}
	}
	if v, ok := nc.mutation.OwnerAddress(); ok {
		if err := nft.OwnerAddressValidator(v); err != nil {
			return &ValidationError{Name: "owner_address", err: fmt.Errorf(`ent: validator failed for field "NFT.owner_address": %w`, err)}
		}
	}
	if _, ok := nc.mutation.MintedAt(); !ok {
		return &ValidationError{Name: "minted_at", err: errors.New(`ent: missing required field "NFT.minted_at"`)}
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "NFT.updated_at"`)}
	}
	return nil
}

func (nc *NFTCreate) sqlSave(ctx context.Context) (*NFT, error) {
	if err := nc.check(); err != nil {
		return nil, err
	}
	_node, _spec, err := nc.createSpec()
	if err != nil {
		return nil, err
	}
	if err := sqlgraph.CreateNode(ctx, nc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	nc.mutation.id = &_node.ID
	nc.mutation.done = true
	return _node, nil
}

func (nc *NFTCreate) createSpec() (*NFT, *sqlgraph.CreateSpec, error) {
	var (
		_node = &NFT{config: nc.config}
		_spec = sqlgraph.NewCreateSpec(nft.Table, sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt))
	)
	if value, ok := nc.mutation.ContractAddress(); ok {
		_spec.SetField(nft.FieldContractAddress, field.TypeString, value)
		_node.ContractAddress = value
	}
	if value, ok := nc.mutation.TokenID(); ok {
		vv, err := nft.ValueScanner.TokenID.Value(value)
		if err != nil {
			return nil, nil, err
		}
		_spec.SetField(nft.FieldTokenID, field.TypeUint64, vv)
		_node.TokenID = value
	}
	if value, ok := nc.mutation.TokenURI(); ok {
		_spec.SetField(nft.FieldTokenURI, field.TypeString, value)
		_node.TokenURI = &value
	}
	if value, ok := nc.mutation.Image(); ok {
		_spec.SetField(nft.FieldImage, field.TypeString, value)
		_node.Image = &value
	}
	if value, ok := nc.mutation.ImageData(); ok {
		_spec.SetField(nft.FieldImageData, field.TypeString, value)
		_node.ImageData = &value
	}
	if value, ok := nc.mutation.ExternalURL(); ok {
		_spec.SetField(nft.FieldExternalURL, field.TypeString, value)
		_node.ExternalURL = &value
	}
	if value, ok := nc.mutation.Description(); ok {
		_spec.SetField(nft.FieldDescription, field.TypeString, value)
		_node.Description = &value
	}
	if value, ok := nc.mutation.Name(); ok {
		_spec.SetField(nft.FieldName, field.TypeString, value)
		_node.Name = &value
	}
	if value, ok := nc.mutation.Attributes(); ok {
		_spec.SetField(nft.FieldAttributes, field.TypeJSON, value)
		_node.Attributes = value
	}
	if value, ok := nc.mutation.BackgroundColor(); ok {
		_spec.SetField(nft.FieldBackgroundColor, field.TypeString, value)
		_node.BackgroundColor = &value
	}
	if value, ok := nc.mutation.AnimationURL(); ok {
		_spec.SetField(nft.FieldAnimationURL, field.TypeString, value)
		_node.AnimationURL = &value
	}
	if value, ok := nc.mutation.YoutubeURL(); ok {
		_spec.SetField(nft.FieldYoutubeURL, field.TypeString, value)
		_node.YoutubeURL = &value
	}
	if value, ok := nc.mutation.OwnerAddress(); ok {
		_spec.SetField(nft.FieldOwnerAddress, field.TypeString, value)
		_node.OwnerAddress = value
	}
	if value, ok := nc.mutation.MintedAt(); ok {
		_spec.SetField(nft.FieldMintedAt, field.TypeTime, value)
		_node.MintedAt = value
	}
	if value, ok := nc.mutation.UpdatedAt(); ok {
		_spec.SetField(nft.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := nc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nft.OwnerTable,
			Columns: []string{nft.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.account_nfts = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.ClassIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nft.ClassTable,
			Columns: []string{nft.ClassColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nftclass.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.nft_class_nfts = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec, nil
}

// NFTCreateBulk is the builder for creating many NFT entities in bulk.
type NFTCreateBulk struct {
	config
	err      error
	builders []*NFTCreate
}

// Save creates the NFT entities in the database.
func (ncb *NFTCreateBulk) Save(ctx context.Context) ([]*NFT, error) {
	if ncb.err != nil {
		return nil, ncb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ncb.builders))
	nodes := make([]*NFT, len(ncb.builders))
	mutators := make([]Mutator, len(ncb.builders))
	for i := range ncb.builders {
		func(i int, root context.Context) {
			builder := ncb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NFTMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i], err = builder.createSpec()
				if err != nil {
					return nil, err
				}
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ncb *NFTCreateBulk) SaveX(ctx context.Context) []*NFT {
	v, err := ncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ncb *NFTCreateBulk) Exec(ctx context.Context) error {
	_, err := ncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncb *NFTCreateBulk) ExecX(ctx context.Context) {
	if err := ncb.Exec(ctx); err != nil {
		panic(err)
	}
}
