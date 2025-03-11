package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/likecoin/like-migration-backend/pkg/likenft/evm/book_nft"
)

func (l *BookNFT) TransferClass(
	evmClassId common.Address,
	newOwner common.Address,
) (*types.Transaction, error) {
	mylogger := l.Logger.WithGroup("TransferClass")

	opts, err := l.transactOpts()
	if err != nil {
		return nil, fmt.Errorf("err l.transactOpts: %v", err)
	}
	opts.NoSend = true

	instance, err := book_nft.NewBookNft(evmClassId, l.Client)
	if err != nil {
		return nil, fmt.Errorf("err book_nft.NewLikenftClass: %v", err)
	}
	tx, err := instance.TransferOwnership(opts, newOwner)
	if err != nil {
		mylogger.Error("instance.TransferOwnership", "err", err)
		return nil, fmt.Errorf("err book_nft.TransferOwnership: %v", err)
	}
	mylogger = mylogger.With("txHash", tx.Hash().Hex())

	err = l.Client.SendTransaction(opts.Context, tx)
	if err != nil {
		mylogger.Error("l.Client.SendTransaction", "err", err)
	}
	return tx, err
}
