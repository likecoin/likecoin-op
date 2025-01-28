// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"likenft-indexer/ent/account"
	"likenft-indexer/ent/nft"
	"likenft-indexer/ent/nftclass"
	"likenft-indexer/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
)

// NFTClassUpdate is the builder for updating NFTClass entities.
type NFTClassUpdate struct {
	config
	hooks    []Hook
	mutation *NFTClassMutation
}

// Where appends a list predicates to the NFTClassUpdate builder.
func (ncu *NFTClassUpdate) Where(ps ...predicate.NFTClass) *NFTClassUpdate {
	ncu.mutation.Where(ps...)
	return ncu
}

// SetAddress sets the "address" field.
func (ncu *NFTClassUpdate) SetAddress(s string) *NFTClassUpdate {
	ncu.mutation.SetAddress(s)
	return ncu
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableAddress(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetAddress(*s)
	}
	return ncu
}

// SetName sets the "name" field.
func (ncu *NFTClassUpdate) SetName(s string) *NFTClassUpdate {
	ncu.mutation.SetName(s)
	return ncu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableName(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetName(*s)
	}
	return ncu
}

// SetSymbol sets the "symbol" field.
func (ncu *NFTClassUpdate) SetSymbol(s string) *NFTClassUpdate {
	ncu.mutation.SetSymbol(s)
	return ncu
}

// SetNillableSymbol sets the "symbol" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableSymbol(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetSymbol(*s)
	}
	return ncu
}

// SetOwnerAddress sets the "owner_address" field.
func (ncu *NFTClassUpdate) SetOwnerAddress(s string) *NFTClassUpdate {
	ncu.mutation.SetOwnerAddress(s)
	return ncu
}

// SetNillableOwnerAddress sets the "owner_address" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableOwnerAddress(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetOwnerAddress(*s)
	}
	return ncu
}

// ClearOwnerAddress clears the value of the "owner_address" field.
func (ncu *NFTClassUpdate) ClearOwnerAddress() *NFTClassUpdate {
	ncu.mutation.ClearOwnerAddress()
	return ncu
}

// SetMinterAddresses sets the "minter_addresses" field.
func (ncu *NFTClassUpdate) SetMinterAddresses(s []string) *NFTClassUpdate {
	ncu.mutation.SetMinterAddresses(s)
	return ncu
}

// AppendMinterAddresses appends s to the "minter_addresses" field.
func (ncu *NFTClassUpdate) AppendMinterAddresses(s []string) *NFTClassUpdate {
	ncu.mutation.AppendMinterAddresses(s)
	return ncu
}

// ClearMinterAddresses clears the value of the "minter_addresses" field.
func (ncu *NFTClassUpdate) ClearMinterAddresses() *NFTClassUpdate {
	ncu.mutation.ClearMinterAddresses()
	return ncu
}

// SetTotalSupply sets the "total_supply" field.
func (ncu *NFTClassUpdate) SetTotalSupply(i int) *NFTClassUpdate {
	ncu.mutation.ResetTotalSupply()
	ncu.mutation.SetTotalSupply(i)
	return ncu
}

// SetNillableTotalSupply sets the "total_supply" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableTotalSupply(i *int) *NFTClassUpdate {
	if i != nil {
		ncu.SetTotalSupply(*i)
	}
	return ncu
}

// AddTotalSupply adds i to the "total_supply" field.
func (ncu *NFTClassUpdate) AddTotalSupply(i int) *NFTClassUpdate {
	ncu.mutation.AddTotalSupply(i)
	return ncu
}

// SetMetadata sets the "metadata" field.
func (ncu *NFTClassUpdate) SetMetadata(m map[string]interface{}) *NFTClassUpdate {
	ncu.mutation.SetMetadata(m)
	return ncu
}

// ClearMetadata clears the value of the "metadata" field.
func (ncu *NFTClassUpdate) ClearMetadata() *NFTClassUpdate {
	ncu.mutation.ClearMetadata()
	return ncu
}

// SetBannerImage sets the "banner_image" field.
func (ncu *NFTClassUpdate) SetBannerImage(s string) *NFTClassUpdate {
	ncu.mutation.SetBannerImage(s)
	return ncu
}

// SetNillableBannerImage sets the "banner_image" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableBannerImage(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetBannerImage(*s)
	}
	return ncu
}

// SetFeaturedImage sets the "featured_image" field.
func (ncu *NFTClassUpdate) SetFeaturedImage(s string) *NFTClassUpdate {
	ncu.mutation.SetFeaturedImage(s)
	return ncu
}

// SetNillableFeaturedImage sets the "featured_image" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableFeaturedImage(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetFeaturedImage(*s)
	}
	return ncu
}

// SetDeployerAddress sets the "deployer_address" field.
func (ncu *NFTClassUpdate) SetDeployerAddress(s string) *NFTClassUpdate {
	ncu.mutation.SetDeployerAddress(s)
	return ncu
}

// SetNillableDeployerAddress sets the "deployer_address" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableDeployerAddress(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetDeployerAddress(*s)
	}
	return ncu
}

// SetDeployedBlockNumber sets the "deployed_block_number" field.
func (ncu *NFTClassUpdate) SetDeployedBlockNumber(s string) *NFTClassUpdate {
	ncu.mutation.SetDeployedBlockNumber(s)
	return ncu
}

// SetNillableDeployedBlockNumber sets the "deployed_block_number" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableDeployedBlockNumber(s *string) *NFTClassUpdate {
	if s != nil {
		ncu.SetDeployedBlockNumber(*s)
	}
	return ncu
}

// SetMintedAt sets the "minted_at" field.
func (ncu *NFTClassUpdate) SetMintedAt(t time.Time) *NFTClassUpdate {
	ncu.mutation.SetMintedAt(t)
	return ncu
}

// SetNillableMintedAt sets the "minted_at" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableMintedAt(t *time.Time) *NFTClassUpdate {
	if t != nil {
		ncu.SetMintedAt(*t)
	}
	return ncu
}

// SetUpdatedAt sets the "updated_at" field.
func (ncu *NFTClassUpdate) SetUpdatedAt(t time.Time) *NFTClassUpdate {
	ncu.mutation.SetUpdatedAt(t)
	return ncu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableUpdatedAt(t *time.Time) *NFTClassUpdate {
	if t != nil {
		ncu.SetUpdatedAt(*t)
	}
	return ncu
}

// AddNftIDs adds the "nfts" edge to the NFT entity by IDs.
func (ncu *NFTClassUpdate) AddNftIDs(ids ...int) *NFTClassUpdate {
	ncu.mutation.AddNftIDs(ids...)
	return ncu
}

// AddNfts adds the "nfts" edges to the NFT entity.
func (ncu *NFTClassUpdate) AddNfts(n ...*NFT) *NFTClassUpdate {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return ncu.AddNftIDs(ids...)
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (ncu *NFTClassUpdate) SetOwnerID(id int) *NFTClassUpdate {
	ncu.mutation.SetOwnerID(id)
	return ncu
}

// SetNillableOwnerID sets the "owner" edge to the Account entity by ID if the given value is not nil.
func (ncu *NFTClassUpdate) SetNillableOwnerID(id *int) *NFTClassUpdate {
	if id != nil {
		ncu = ncu.SetOwnerID(*id)
	}
	return ncu
}

// SetOwner sets the "owner" edge to the Account entity.
func (ncu *NFTClassUpdate) SetOwner(a *Account) *NFTClassUpdate {
	return ncu.SetOwnerID(a.ID)
}

// Mutation returns the NFTClassMutation object of the builder.
func (ncu *NFTClassUpdate) Mutation() *NFTClassMutation {
	return ncu.mutation
}

// ClearNfts clears all "nfts" edges to the NFT entity.
func (ncu *NFTClassUpdate) ClearNfts() *NFTClassUpdate {
	ncu.mutation.ClearNfts()
	return ncu
}

// RemoveNftIDs removes the "nfts" edge to NFT entities by IDs.
func (ncu *NFTClassUpdate) RemoveNftIDs(ids ...int) *NFTClassUpdate {
	ncu.mutation.RemoveNftIDs(ids...)
	return ncu
}

// RemoveNfts removes "nfts" edges to NFT entities.
func (ncu *NFTClassUpdate) RemoveNfts(n ...*NFT) *NFTClassUpdate {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return ncu.RemoveNftIDs(ids...)
}

// ClearOwner clears the "owner" edge to the Account entity.
func (ncu *NFTClassUpdate) ClearOwner() *NFTClassUpdate {
	ncu.mutation.ClearOwner()
	return ncu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ncu *NFTClassUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ncu.sqlSave, ncu.mutation, ncu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ncu *NFTClassUpdate) SaveX(ctx context.Context) int {
	affected, err := ncu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ncu *NFTClassUpdate) Exec(ctx context.Context) error {
	_, err := ncu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncu *NFTClassUpdate) ExecX(ctx context.Context) {
	if err := ncu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ncu *NFTClassUpdate) check() error {
	if v, ok := ncu.mutation.Name(); ok {
		if err := nftclass.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "NFTClass.name": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.Symbol(); ok {
		if err := nftclass.SymbolValidator(v); err != nil {
			return &ValidationError{Name: "symbol", err: fmt.Errorf(`ent: validator failed for field "NFTClass.symbol": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.TotalSupply(); ok {
		if err := nftclass.TotalSupplyValidator(v); err != nil {
			return &ValidationError{Name: "total_supply", err: fmt.Errorf(`ent: validator failed for field "NFTClass.total_supply": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.BannerImage(); ok {
		if err := nftclass.BannerImageValidator(v); err != nil {
			return &ValidationError{Name: "banner_image", err: fmt.Errorf(`ent: validator failed for field "NFTClass.banner_image": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.FeaturedImage(); ok {
		if err := nftclass.FeaturedImageValidator(v); err != nil {
			return &ValidationError{Name: "featured_image", err: fmt.Errorf(`ent: validator failed for field "NFTClass.featured_image": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.DeployerAddress(); ok {
		if err := nftclass.DeployerAddressValidator(v); err != nil {
			return &ValidationError{Name: "deployer_address", err: fmt.Errorf(`ent: validator failed for field "NFTClass.deployer_address": %w`, err)}
		}
	}
	if v, ok := ncu.mutation.DeployedBlockNumber(); ok {
		if err := nftclass.DeployedBlockNumberValidator(v); err != nil {
			return &ValidationError{Name: "deployed_block_number", err: fmt.Errorf(`ent: validator failed for field "NFTClass.deployed_block_number": %w`, err)}
		}
	}
	return nil
}

func (ncu *NFTClassUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ncu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(nftclass.Table, nftclass.Columns, sqlgraph.NewFieldSpec(nftclass.FieldID, field.TypeInt))
	if ps := ncu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ncu.mutation.Address(); ok {
		_spec.SetField(nftclass.FieldAddress, field.TypeString, value)
	}
	if value, ok := ncu.mutation.Name(); ok {
		_spec.SetField(nftclass.FieldName, field.TypeString, value)
	}
	if value, ok := ncu.mutation.Symbol(); ok {
		_spec.SetField(nftclass.FieldSymbol, field.TypeString, value)
	}
	if value, ok := ncu.mutation.OwnerAddress(); ok {
		_spec.SetField(nftclass.FieldOwnerAddress, field.TypeString, value)
	}
	if ncu.mutation.OwnerAddressCleared() {
		_spec.ClearField(nftclass.FieldOwnerAddress, field.TypeString)
	}
	if value, ok := ncu.mutation.MinterAddresses(); ok {
		_spec.SetField(nftclass.FieldMinterAddresses, field.TypeJSON, value)
	}
	if value, ok := ncu.mutation.AppendedMinterAddresses(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, nftclass.FieldMinterAddresses, value)
		})
	}
	if ncu.mutation.MinterAddressesCleared() {
		_spec.ClearField(nftclass.FieldMinterAddresses, field.TypeJSON)
	}
	if value, ok := ncu.mutation.TotalSupply(); ok {
		_spec.SetField(nftclass.FieldTotalSupply, field.TypeInt, value)
	}
	if value, ok := ncu.mutation.AddedTotalSupply(); ok {
		_spec.AddField(nftclass.FieldTotalSupply, field.TypeInt, value)
	}
	if value, ok := ncu.mutation.Metadata(); ok {
		_spec.SetField(nftclass.FieldMetadata, field.TypeJSON, value)
	}
	if ncu.mutation.MetadataCleared() {
		_spec.ClearField(nftclass.FieldMetadata, field.TypeJSON)
	}
	if value, ok := ncu.mutation.BannerImage(); ok {
		_spec.SetField(nftclass.FieldBannerImage, field.TypeString, value)
	}
	if value, ok := ncu.mutation.FeaturedImage(); ok {
		_spec.SetField(nftclass.FieldFeaturedImage, field.TypeString, value)
	}
	if value, ok := ncu.mutation.DeployerAddress(); ok {
		_spec.SetField(nftclass.FieldDeployerAddress, field.TypeString, value)
	}
	if value, ok := ncu.mutation.DeployedBlockNumber(); ok {
		_spec.SetField(nftclass.FieldDeployedBlockNumber, field.TypeString, value)
	}
	if value, ok := ncu.mutation.MintedAt(); ok {
		_spec.SetField(nftclass.FieldMintedAt, field.TypeTime, value)
	}
	if value, ok := ncu.mutation.UpdatedAt(); ok {
		_spec.SetField(nftclass.FieldUpdatedAt, field.TypeTime, value)
	}
	if ncu.mutation.NftsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncu.mutation.RemovedNftsIDs(); len(nodes) > 0 && !ncu.mutation.NftsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncu.mutation.NftsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ncu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nftclass.OwnerTable,
			Columns: []string{nftclass.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nftclass.OwnerTable,
			Columns: []string{nftclass.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ncu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{nftclass.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ncu.mutation.done = true
	return n, nil
}

// NFTClassUpdateOne is the builder for updating a single NFTClass entity.
type NFTClassUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *NFTClassMutation
}

// SetAddress sets the "address" field.
func (ncuo *NFTClassUpdateOne) SetAddress(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetAddress(s)
	return ncuo
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableAddress(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetAddress(*s)
	}
	return ncuo
}

// SetName sets the "name" field.
func (ncuo *NFTClassUpdateOne) SetName(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetName(s)
	return ncuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableName(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetName(*s)
	}
	return ncuo
}

// SetSymbol sets the "symbol" field.
func (ncuo *NFTClassUpdateOne) SetSymbol(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetSymbol(s)
	return ncuo
}

// SetNillableSymbol sets the "symbol" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableSymbol(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetSymbol(*s)
	}
	return ncuo
}

// SetOwnerAddress sets the "owner_address" field.
func (ncuo *NFTClassUpdateOne) SetOwnerAddress(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetOwnerAddress(s)
	return ncuo
}

// SetNillableOwnerAddress sets the "owner_address" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableOwnerAddress(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetOwnerAddress(*s)
	}
	return ncuo
}

// ClearOwnerAddress clears the value of the "owner_address" field.
func (ncuo *NFTClassUpdateOne) ClearOwnerAddress() *NFTClassUpdateOne {
	ncuo.mutation.ClearOwnerAddress()
	return ncuo
}

// SetMinterAddresses sets the "minter_addresses" field.
func (ncuo *NFTClassUpdateOne) SetMinterAddresses(s []string) *NFTClassUpdateOne {
	ncuo.mutation.SetMinterAddresses(s)
	return ncuo
}

// AppendMinterAddresses appends s to the "minter_addresses" field.
func (ncuo *NFTClassUpdateOne) AppendMinterAddresses(s []string) *NFTClassUpdateOne {
	ncuo.mutation.AppendMinterAddresses(s)
	return ncuo
}

// ClearMinterAddresses clears the value of the "minter_addresses" field.
func (ncuo *NFTClassUpdateOne) ClearMinterAddresses() *NFTClassUpdateOne {
	ncuo.mutation.ClearMinterAddresses()
	return ncuo
}

// SetTotalSupply sets the "total_supply" field.
func (ncuo *NFTClassUpdateOne) SetTotalSupply(i int) *NFTClassUpdateOne {
	ncuo.mutation.ResetTotalSupply()
	ncuo.mutation.SetTotalSupply(i)
	return ncuo
}

// SetNillableTotalSupply sets the "total_supply" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableTotalSupply(i *int) *NFTClassUpdateOne {
	if i != nil {
		ncuo.SetTotalSupply(*i)
	}
	return ncuo
}

// AddTotalSupply adds i to the "total_supply" field.
func (ncuo *NFTClassUpdateOne) AddTotalSupply(i int) *NFTClassUpdateOne {
	ncuo.mutation.AddTotalSupply(i)
	return ncuo
}

// SetMetadata sets the "metadata" field.
func (ncuo *NFTClassUpdateOne) SetMetadata(m map[string]interface{}) *NFTClassUpdateOne {
	ncuo.mutation.SetMetadata(m)
	return ncuo
}

// ClearMetadata clears the value of the "metadata" field.
func (ncuo *NFTClassUpdateOne) ClearMetadata() *NFTClassUpdateOne {
	ncuo.mutation.ClearMetadata()
	return ncuo
}

// SetBannerImage sets the "banner_image" field.
func (ncuo *NFTClassUpdateOne) SetBannerImage(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetBannerImage(s)
	return ncuo
}

// SetNillableBannerImage sets the "banner_image" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableBannerImage(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetBannerImage(*s)
	}
	return ncuo
}

// SetFeaturedImage sets the "featured_image" field.
func (ncuo *NFTClassUpdateOne) SetFeaturedImage(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetFeaturedImage(s)
	return ncuo
}

// SetNillableFeaturedImage sets the "featured_image" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableFeaturedImage(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetFeaturedImage(*s)
	}
	return ncuo
}

// SetDeployerAddress sets the "deployer_address" field.
func (ncuo *NFTClassUpdateOne) SetDeployerAddress(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetDeployerAddress(s)
	return ncuo
}

// SetNillableDeployerAddress sets the "deployer_address" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableDeployerAddress(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetDeployerAddress(*s)
	}
	return ncuo
}

// SetDeployedBlockNumber sets the "deployed_block_number" field.
func (ncuo *NFTClassUpdateOne) SetDeployedBlockNumber(s string) *NFTClassUpdateOne {
	ncuo.mutation.SetDeployedBlockNumber(s)
	return ncuo
}

// SetNillableDeployedBlockNumber sets the "deployed_block_number" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableDeployedBlockNumber(s *string) *NFTClassUpdateOne {
	if s != nil {
		ncuo.SetDeployedBlockNumber(*s)
	}
	return ncuo
}

// SetMintedAt sets the "minted_at" field.
func (ncuo *NFTClassUpdateOne) SetMintedAt(t time.Time) *NFTClassUpdateOne {
	ncuo.mutation.SetMintedAt(t)
	return ncuo
}

// SetNillableMintedAt sets the "minted_at" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableMintedAt(t *time.Time) *NFTClassUpdateOne {
	if t != nil {
		ncuo.SetMintedAt(*t)
	}
	return ncuo
}

// SetUpdatedAt sets the "updated_at" field.
func (ncuo *NFTClassUpdateOne) SetUpdatedAt(t time.Time) *NFTClassUpdateOne {
	ncuo.mutation.SetUpdatedAt(t)
	return ncuo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableUpdatedAt(t *time.Time) *NFTClassUpdateOne {
	if t != nil {
		ncuo.SetUpdatedAt(*t)
	}
	return ncuo
}

// AddNftIDs adds the "nfts" edge to the NFT entity by IDs.
func (ncuo *NFTClassUpdateOne) AddNftIDs(ids ...int) *NFTClassUpdateOne {
	ncuo.mutation.AddNftIDs(ids...)
	return ncuo
}

// AddNfts adds the "nfts" edges to the NFT entity.
func (ncuo *NFTClassUpdateOne) AddNfts(n ...*NFT) *NFTClassUpdateOne {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return ncuo.AddNftIDs(ids...)
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (ncuo *NFTClassUpdateOne) SetOwnerID(id int) *NFTClassUpdateOne {
	ncuo.mutation.SetOwnerID(id)
	return ncuo
}

// SetNillableOwnerID sets the "owner" edge to the Account entity by ID if the given value is not nil.
func (ncuo *NFTClassUpdateOne) SetNillableOwnerID(id *int) *NFTClassUpdateOne {
	if id != nil {
		ncuo = ncuo.SetOwnerID(*id)
	}
	return ncuo
}

// SetOwner sets the "owner" edge to the Account entity.
func (ncuo *NFTClassUpdateOne) SetOwner(a *Account) *NFTClassUpdateOne {
	return ncuo.SetOwnerID(a.ID)
}

// Mutation returns the NFTClassMutation object of the builder.
func (ncuo *NFTClassUpdateOne) Mutation() *NFTClassMutation {
	return ncuo.mutation
}

// ClearNfts clears all "nfts" edges to the NFT entity.
func (ncuo *NFTClassUpdateOne) ClearNfts() *NFTClassUpdateOne {
	ncuo.mutation.ClearNfts()
	return ncuo
}

// RemoveNftIDs removes the "nfts" edge to NFT entities by IDs.
func (ncuo *NFTClassUpdateOne) RemoveNftIDs(ids ...int) *NFTClassUpdateOne {
	ncuo.mutation.RemoveNftIDs(ids...)
	return ncuo
}

// RemoveNfts removes "nfts" edges to NFT entities.
func (ncuo *NFTClassUpdateOne) RemoveNfts(n ...*NFT) *NFTClassUpdateOne {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return ncuo.RemoveNftIDs(ids...)
}

// ClearOwner clears the "owner" edge to the Account entity.
func (ncuo *NFTClassUpdateOne) ClearOwner() *NFTClassUpdateOne {
	ncuo.mutation.ClearOwner()
	return ncuo
}

// Where appends a list predicates to the NFTClassUpdate builder.
func (ncuo *NFTClassUpdateOne) Where(ps ...predicate.NFTClass) *NFTClassUpdateOne {
	ncuo.mutation.Where(ps...)
	return ncuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ncuo *NFTClassUpdateOne) Select(field string, fields ...string) *NFTClassUpdateOne {
	ncuo.fields = append([]string{field}, fields...)
	return ncuo
}

// Save executes the query and returns the updated NFTClass entity.
func (ncuo *NFTClassUpdateOne) Save(ctx context.Context) (*NFTClass, error) {
	return withHooks(ctx, ncuo.sqlSave, ncuo.mutation, ncuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ncuo *NFTClassUpdateOne) SaveX(ctx context.Context) *NFTClass {
	node, err := ncuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ncuo *NFTClassUpdateOne) Exec(ctx context.Context) error {
	_, err := ncuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncuo *NFTClassUpdateOne) ExecX(ctx context.Context) {
	if err := ncuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ncuo *NFTClassUpdateOne) check() error {
	if v, ok := ncuo.mutation.Name(); ok {
		if err := nftclass.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "NFTClass.name": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.Symbol(); ok {
		if err := nftclass.SymbolValidator(v); err != nil {
			return &ValidationError{Name: "symbol", err: fmt.Errorf(`ent: validator failed for field "NFTClass.symbol": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.TotalSupply(); ok {
		if err := nftclass.TotalSupplyValidator(v); err != nil {
			return &ValidationError{Name: "total_supply", err: fmt.Errorf(`ent: validator failed for field "NFTClass.total_supply": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.BannerImage(); ok {
		if err := nftclass.BannerImageValidator(v); err != nil {
			return &ValidationError{Name: "banner_image", err: fmt.Errorf(`ent: validator failed for field "NFTClass.banner_image": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.FeaturedImage(); ok {
		if err := nftclass.FeaturedImageValidator(v); err != nil {
			return &ValidationError{Name: "featured_image", err: fmt.Errorf(`ent: validator failed for field "NFTClass.featured_image": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.DeployerAddress(); ok {
		if err := nftclass.DeployerAddressValidator(v); err != nil {
			return &ValidationError{Name: "deployer_address", err: fmt.Errorf(`ent: validator failed for field "NFTClass.deployer_address": %w`, err)}
		}
	}
	if v, ok := ncuo.mutation.DeployedBlockNumber(); ok {
		if err := nftclass.DeployedBlockNumberValidator(v); err != nil {
			return &ValidationError{Name: "deployed_block_number", err: fmt.Errorf(`ent: validator failed for field "NFTClass.deployed_block_number": %w`, err)}
		}
	}
	return nil
}

func (ncuo *NFTClassUpdateOne) sqlSave(ctx context.Context) (_node *NFTClass, err error) {
	if err := ncuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(nftclass.Table, nftclass.Columns, sqlgraph.NewFieldSpec(nftclass.FieldID, field.TypeInt))
	id, ok := ncuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "NFTClass.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ncuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, nftclass.FieldID)
		for _, f := range fields {
			if !nftclass.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != nftclass.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ncuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ncuo.mutation.Address(); ok {
		_spec.SetField(nftclass.FieldAddress, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.Name(); ok {
		_spec.SetField(nftclass.FieldName, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.Symbol(); ok {
		_spec.SetField(nftclass.FieldSymbol, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.OwnerAddress(); ok {
		_spec.SetField(nftclass.FieldOwnerAddress, field.TypeString, value)
	}
	if ncuo.mutation.OwnerAddressCleared() {
		_spec.ClearField(nftclass.FieldOwnerAddress, field.TypeString)
	}
	if value, ok := ncuo.mutation.MinterAddresses(); ok {
		_spec.SetField(nftclass.FieldMinterAddresses, field.TypeJSON, value)
	}
	if value, ok := ncuo.mutation.AppendedMinterAddresses(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, nftclass.FieldMinterAddresses, value)
		})
	}
	if ncuo.mutation.MinterAddressesCleared() {
		_spec.ClearField(nftclass.FieldMinterAddresses, field.TypeJSON)
	}
	if value, ok := ncuo.mutation.TotalSupply(); ok {
		_spec.SetField(nftclass.FieldTotalSupply, field.TypeInt, value)
	}
	if value, ok := ncuo.mutation.AddedTotalSupply(); ok {
		_spec.AddField(nftclass.FieldTotalSupply, field.TypeInt, value)
	}
	if value, ok := ncuo.mutation.Metadata(); ok {
		_spec.SetField(nftclass.FieldMetadata, field.TypeJSON, value)
	}
	if ncuo.mutation.MetadataCleared() {
		_spec.ClearField(nftclass.FieldMetadata, field.TypeJSON)
	}
	if value, ok := ncuo.mutation.BannerImage(); ok {
		_spec.SetField(nftclass.FieldBannerImage, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.FeaturedImage(); ok {
		_spec.SetField(nftclass.FieldFeaturedImage, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.DeployerAddress(); ok {
		_spec.SetField(nftclass.FieldDeployerAddress, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.DeployedBlockNumber(); ok {
		_spec.SetField(nftclass.FieldDeployedBlockNumber, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.MintedAt(); ok {
		_spec.SetField(nftclass.FieldMintedAt, field.TypeTime, value)
	}
	if value, ok := ncuo.mutation.UpdatedAt(); ok {
		_spec.SetField(nftclass.FieldUpdatedAt, field.TypeTime, value)
	}
	if ncuo.mutation.NftsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncuo.mutation.RemovedNftsIDs(); len(nodes) > 0 && !ncuo.mutation.NftsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncuo.mutation.NftsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nftclass.NftsTable,
			Columns: []string{nftclass.NftsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nft.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ncuo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nftclass.OwnerTable,
			Columns: []string{nftclass.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ncuo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nftclass.OwnerTable,
			Columns: []string{nftclass.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &NFTClass{config: ncuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ncuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{nftclass.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ncuo.mutation.done = true
	return _node, nil
}
