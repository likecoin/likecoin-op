package evm

import (
	"context"
	"fmt"
	"math/big"

	"likecollective-indexer/internal/evm/like_stake_position"
	"likecollective-indexer/internal/util/parallel"
	"likecollective-indexer/internal/util/retry"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

// maxStakePositions bounds the total supply ListStakePositions will enumerate.
// Far above any real pool, and low enough that a wrong address or an abnormal
// contract fails here instead of allocating an index and a result per token.
const maxStakePositions = 1_000_000

// StakePosition is one live LikeStakePosition token, resolved to its current
// owner and the book NFT it is staked against.
type StakePosition struct {
	TokenID      *big.Int
	Owner        common.Address
	BookNFT      common.Address
	StakedAmount *big.Int
}

// ListStakePositions enumerates every live position through ERC721Enumerable.
// blockNumber must be pinned: totalSupply and the per-token reads are separate
// calls, and a burn in between would make tokenByIndex skip a token.
func (e *evmClient) ListStakePositions(
	ctx context.Context,
	blockNumber *big.Int,
	concurrency int,
) ([]*StakePosition, error) {
	likeStakePositionClient, err := like_stake_position.NewLikeStakePosition(e.likeStakePositionAddress, e.client)
	if err != nil {
		return nil, err
	}

	totalSupply, err := likeStakePositionClient.TotalSupply(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %w", err)
	}

	// totalSupply is an unvalidated on-chain counter, and a slice is
	// allocated for it before any read, so it is bounded first.
	if !totalSupply.IsInt64() || totalSupply.Int64() > maxStakePositions {
		return nil, fmt.Errorf(
			"implausible total supply %s, above %d", totalSupply, maxStakePositions,
		)
	}
	indexes := make([]int64, totalSupply.Int64())
	for i := range indexes {
		indexes[i] = int64(i)
	}

	return parallel.MapWithLimit(ctx, concurrency, indexes, func(
		ctx context.Context,
		index int64,
	) (*StakePosition, error) {
		return retry.Value(ctx, retry.DefaultAttempts, retry.DefaultBase, func(ctx context.Context) (*StakePosition, error) {
			return e.stakePositionAtIndex(ctx, likeStakePositionClient, blockNumber, big.NewInt(index))
		})
	})
}

func (e *evmClient) stakePositionAtIndex(
	ctx context.Context,
	likeStakePositionClient *like_stake_position.LikeStakePosition,
	blockNumber *big.Int,
	index *big.Int,
) (*StakePosition, error) {
	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}

	tokenId, err := likeStakePositionClient.TokenByIndex(callOpts, index)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by index %s: %w", index, err)
	}

	owner, err := likeStakePositionClient.OwnerOf(callOpts, tokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner of token %s: %w", tokenId, err)
	}

	position, err := likeStakePositionClient.GetPosition(callOpts, tokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get position of token %s: %w", tokenId, err)
	}

	return &StakePosition{
		TokenID:      tokenId,
		Owner:        owner,
		BookNFT:      position.BookNFT,
		StakedAmount: position.StakedAmount,
	}, nil
}
