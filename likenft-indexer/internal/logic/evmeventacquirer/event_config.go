package evmeventacquirer

import (
	"context"
	"fmt"
	"log/slog"

	"likenft-indexer/ent/evmeventprocessedblockheight"
	"likenft-indexer/internal/evm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

type LogsRetriever func(
	ctx context.Context,
	logger *slog.Logger,
	contractAddress string,
	startBlock uint64,
) ([]types.Log, error)

type eventConfig struct {
	ContractType  evmeventprocessedblockheight.ContractType
	Abi           *abi.ABI
	LogsRetriever LogsRetriever
}

type eventConfigCreator func(evmClient *evm.EvmClient) eventConfig

var eventConfigMap = make(map[evmeventprocessedblockheight.Event]eventConfigCreator)

func registerEventConfig(event evmeventprocessedblockheight.Event, creator eventConfigCreator) {
	eventConfigMap[event] = creator
}

func getEventConfig(event evmeventprocessedblockheight.Event) (eventConfigCreator, error) {
	eventConfig, ok := eventConfigMap[event]
	if !ok {
		return nil, fmt.Errorf("unknown event %s", event)
	}
	return eventConfig, nil
}
