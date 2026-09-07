package evm

import (
	"context"
	"math/big"

	"likecollective-indexer/internal/evm/like_collective"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

func (e *evmClient) GetStakeForUser(
	ctx context.Context,
	blockNumber *big.Int,
	user common.Address,
	bookNFT common.Address,
) (*big.Int, error) {
	likeCollectiveClient, err := like_collective.NewLikeCollective(e.likeCollectiveAddress, e.client)
	if err != nil {
		return nil, err
	}

	return likeCollectiveClient.GetStakeForUser(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}, user, bookNFT)
}
