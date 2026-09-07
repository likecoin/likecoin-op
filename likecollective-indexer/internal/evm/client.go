package evm

import (
	"context"
	"math/big"

	"likecollective-indexer/internal/evm/like_stake_position"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMClient interface {
	// FilterLogs reads logs straight from the chain, bypassing the webhook.
	FilterLogs(
		ctx context.Context,
		fromBlock *big.Int,
		toBlock *big.Int,
	) ([]types.Log, error)

	LatestBlockNumber(ctx context.Context) (*big.Int, error)

	ListStakePositions(
		ctx context.Context,
		blockNumber *big.Int,
		concurrency int,
	) ([]*StakePosition, error)

	GetStakeForUser(
		ctx context.Context,
		blockNumber *big.Int,
		user common.Address,
		bookNFT common.Address,
	) (*big.Int, error)

	GetPendingRewardsForUser(
		ctx context.Context,
		blockNumber *big.Int,
		user common.Address,
		bookNFT common.Address,
	) (*big.Int, error)

	GetTotalStake(
		ctx context.Context,
		blockNumber *big.Int,
		bookNFT common.Address,
	) (*big.Int, error)

	GetRewardsOfPosition(
		ctx context.Context,
		blockNumber *big.Int,
		tokenId *big.Int,
	) (*big.Int, error)

	GetStakePosition(
		ctx context.Context,
		blockNumber *big.Int,
		tokenId *big.Int,
	) (*like_stake_position.LikeStakePositionPosition, error)
}

type evmClient struct {
	likeCollectiveAddress    common.Address
	likeStakePositionAddress common.Address
	client                   *ethclient.Client
}

func NewEVMClient(
	likeCollectiveAddress common.Address,
	likeStakePositionAddress common.Address,
	client *ethclient.Client,
) EVMClient {
	return &evmClient{
		likeCollectiveAddress,
		likeStakePositionAddress,
		client,
	}
}
