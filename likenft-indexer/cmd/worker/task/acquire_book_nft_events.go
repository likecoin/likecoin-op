package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	appcontext "likenft-indexer/cmd/worker/context"
	"likenft-indexer/cmd/worker/config"
	"likenft-indexer/ent"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/database"
	"likenft-indexer/internal/logic/contractevmeventacquirer"
	"likenft-indexer/internal/worker/task"

	"github.com/hibiken/asynq"
)

const TypeAcquireBookNFTEventsTaskPayload = "acquire-book-nft-events"

type AcquireBookNFTEventsTaskPayload struct {
	ContractAddresses []string
}

func NewAcquireBookNFTEventsTask(contractAddresses []string) (*asynq.Task, error) {
	payload, err := json.Marshal(AcquireBookNFTEventsTaskPayload{
		ContractAddresses: contractAddresses,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeAcquireBookNFTEventsTaskPayload,
		payload,
		asynq.Queue(TypeAcquireBookNFTEventsTaskPayload),
	), nil
}

func groupByProperty[T any, K comparable](items []T, getProperty func(T) K) map[K][]T {
	grouped := make(map[K][]T)

	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}

func HandleAcquireBookNFTEventsTask(ctx context.Context, t *asynq.Task) error {
	logger := appcontext.LoggerFromContext(ctx)
	cfg := appcontext.ConfigFromContext(ctx)
	evmEventQueryClient := appcontext.EvmQueryClientFromContext(ctx)
	evmClient := appcontext.EvmClientFromContext(ctx)

	mylogger := logger.WithGroup("HandleAcquireBookNFTEventsTask")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var p AcquireBookNFTEventsTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		mylogger.Error("json.Unmarshal AcquireNewBookNFTTaskPayload", "err", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	dbService := database.New()

	evmEventRepository := database.MakeEVMEventRepository(dbService)
	nftClassRepository := database.MakeNFTClassRepository(dbService)

	nftClasses, err := nftClassRepository.QueryNFTClassesByAddressesExact(
		ctx, p.ContractAddresses,
	)

	if err != nil {
		mylogger.Error("nftClassRepository.QueryNFTClassesByAddressesExact", "err", err)
	}

	groupedByBlockHeight := groupByProperty(nftClasses, func(nftClass *ent.NFTClass) typeutil.Uint64 {
		return nftClass.LatestEventBlockNumber
	})

	for latestEventsBlockHeight, nftClasses := range groupedByBlockHeight {
		addresses := make([]string, len(nftClasses))
		for i, nftClass := range nftClasses {
			addresses[i] = nftClass.Address
		}

		acquirer := contractevmeventacquirer.NewContractEvmEventsAcquirer(
			evmEventQueryClient,
			evmEventRepository,
			evmEventQueryClient,
			evmClient,
			cfg.EvmEventQueryToBlockPadding,
			contractevmeventacquirer.ContractEvmEventsAcquirerContractTypeBookNFT,
			addresses,
		)

		newBlockHeight, _, err := acquirer.Acquire(
			ctx,
			logger,
			uint64(latestEventsBlockHeight),
			cfg.EvmEventQueryNumberOfBlocksLimit,
		)

		if err != nil && len(addresses) > 1 && isResponseTooLargeError(err) {
			mylogger.Warn("batched eth_getLogs response too large, falling back to per-address queries",
				"addressCount", len(addresses), "err", err)
			if fallbackErr := acquirePerAddress(
				ctx, logger, cfg, evmEventQueryClient, evmEventRepository,
				evmEventQueryClient, evmClient, nftClassRepository,
				addresses, uint64(latestEventsBlockHeight),
			); fallbackErr != nil {
				return fallbackErr
			}
			continue
		}

		if err != nil {
			mylogger.Error("acquirer.Acquire", "err", err)
			var errCannotConvertLog *contractevmeventacquirer.ErrCannotConvertLog
			if errors.As(err, &errCannotConvertLog) {
				return errors.Join(
					err,
					nftClassRepository.DisableForIndexing(ctx, errCannotConvertLog.Log.Address.Hex(), err.Error()),
				)
			}
			return err
		}

		err = nftClassRepository.UpdateNFTClassesLatestEventBlockNumber(
			ctx,
			addresses,
			typeutil.Uint64(newBlockHeight),
		)

		if err != nil {
			mylogger.Error("nftClassRepository.UpdateNFTClassesLatestEventsBlockNumber", "err", err)
		}
	}

	return nil
}

// isResponseTooLargeError checks if an RPC error indicates the eth_getLogs
// response exceeded the provider's size limit.
// Known error messages:
//   - Alchemy: "Log response size exceeded."
//   - Geth/others: "query returned more than 10000 results"
func isResponseTooLargeError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "log response size exceeded") ||
		strings.Contains(msg, "query returned more than")
}

// acquirePerAddress retries event acquisition one address at a time.
// Used as a fallback when a batched FilterLogs call exceeds response limits.
func acquirePerAddress(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.EnvConfig,
	abiManager contractevmeventacquirer.ABIManager,
	evmEventRepository database.EVMEventRepository,
	evmEventQueryClient contractevmeventacquirer.EvmEventQueryClient,
	evmClient contractevmeventacquirer.EvmClient,
	nftClassRepository database.NFTClassRepository,
	addresses []string,
	fromBlock uint64,
) error {
	mylogger := logger.WithGroup("acquirePerAddress")

	for _, addr := range addresses {
		singleAcquirer := contractevmeventacquirer.NewContractEvmEventsAcquirer(
			abiManager,
			evmEventRepository,
			evmEventQueryClient,
			evmClient,
			cfg.EvmEventQueryToBlockPadding,
			contractevmeventacquirer.ContractEvmEventsAcquirerContractTypeBookNFT,
			[]string{addr},
		)

		newBlockHeight, _, err := singleAcquirer.Acquire(
			ctx,
			logger,
			fromBlock,
			cfg.EvmEventQueryNumberOfBlocksLimit,
		)
		if err != nil {
			mylogger.Error("acquirer.Acquire", "addr", addr, "err", err)
			var errCannotConvertLog *contractevmeventacquirer.ErrCannotConvertLog
			if errors.As(err, &errCannotConvertLog) {
				if disableErr := nftClassRepository.DisableForIndexing(
					ctx, errCannotConvertLog.Log.Address.Hex(), err.Error(),
				); disableErr != nil {
					mylogger.Error("DisableForIndexing", "addr", addr, "err", disableErr)
				}
				continue
			}
			return err
		}

		if updateErr := nftClassRepository.UpdateNFTClassesLatestEventBlockNumber(
			ctx, []string{addr}, typeutil.Uint64(newBlockHeight),
		); updateErr != nil {
			mylogger.Error("UpdateNFTClassesLatestEventBlockNumber", "addr", addr, "err", updateErr)
		}
	}

	return nil
}

func init() {
	Tasks.Register(task.DefineTask(
		TypeAcquireBookNFTEventsTaskPayload,
		HandleAcquireBookNFTEventsTask,
	))
}
