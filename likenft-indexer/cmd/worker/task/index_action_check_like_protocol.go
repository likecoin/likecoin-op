package task

import (
	"context"
	"encoding/json"
	"fmt"

	appcontext "likenft-indexer/cmd/worker/context"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/database"
	"likenft-indexer/internal/logic/contractevmeventacquirer"
	"likenft-indexer/internal/worker/task"

	"github.com/hibiken/asynq"
)

const TypeIndexActionCheckLikeProtocolPayload = "index-action-check-like-protocol"

type IndexActionCheckLikeProtocolPayload struct {
	ContractAddress string
}

func NewIndexActionCheckLikeProtocolTask(contractAddress string) (*asynq.Task, error) {
	payload, err := json.Marshal(IndexActionCheckLikeProtocolPayload{
		ContractAddress: contractAddress,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeIndexActionCheckLikeProtocolPayload,
		payload,
		asynq.Queue(TypeIndexActionCheckLikeProtocolPayload),
	), nil
}

func HandleIndexActionCheckLikeProtocol(ctx context.Context, t *asynq.Task) error {
	logger := appcontext.LoggerFromContext(ctx)
	envCfg := appcontext.ConfigFromContext(ctx)
	asynqClient := appcontext.AsynqClientFromContext(ctx)
	evmQueryClient := appcontext.EvmQueryClientFromContext(ctx)
	evmClient := appcontext.EvmClientFromContext(ctx)

	evmClient.InvalidateBlockNumberCache()

	dbService := database.New()
	likeProtocolRepository := database.MakeLikeProtocolRepository(dbService)
	evmEventRepository := database.MakeEVMEventRepository(dbService)

	mylogger := logger.WithGroup("HandleIndexActionCheckLikeProtocol")

	var p IndexActionCheckLikeProtocolPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		mylogger.Error("json.Unmarshal IndexActionCheckLikeProtocolPayload", "err", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	likeProtocol, err := likeProtocolRepository.GetLikeProtocol(ctx, p.ContractAddress)
	if err != nil {
		return err
	}

	latestBlockNumber, err := evmClient.BlockNumber(ctx)
	if err != nil {
		return err
	}

	mylogger = mylogger.With("latestBlockNumber", latestBlockNumber)

	blockStarts := make([]uint64, 0)

	for i := uint64(likeProtocol.LatestEventBlockNumber); i < latestBlockNumber; i = i + envCfg.EvmEventQueryNumberOfBlocksLimit {
		blockStarts = append(blockStarts, i)
	}

	acquirer := contractevmeventacquirer.NewContractEvmEventsAcquirer(
		evmQueryClient,
		evmEventRepository,
		evmQueryClient,
		evmClient,
		contractevmeventacquirer.ContractEvmEventsAcquirerContractTypeLikeProtocol,
		[]string{p.ContractAddress},
	)

	for i, blockStart := range blockStarts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		mylogger := mylogger.With(
			"partition",
			fmt.Sprintf("%d/%d", i, len(blockStarts)),
		)
		newBlockNumber, evmEvents, err := acquirer.Acquire(
			ctx,
			mylogger,
			blockStart,
			envCfg.EvmEventQueryNumberOfBlocksLimit,
		)
		if err != nil {
			return err
		}

		if len(evmEvents) > 0 {
			task, err := NewIndexActionCheckReceivedEventsTask(p.ContractAddress)
			if err != nil {
				mylogger.
					WarnContext(ctx, "cannot create task. should eventually picked up by periodic worker. skip.", "err", err)
			}
			_, err = asynqClient.EnqueueContext(ctx, task, asynq.MaxRetry(0))
			if err != nil {
				mylogger.
					WarnContext(ctx, "cannot enqueue task. should eventually picked up by periodic worker. skip.", "err", err)
			}
		}

		err = likeProtocolRepository.CreateOrUpdateLatestEventBlockHeight(
			ctx,
			p.ContractAddress,
			typeutil.Uint64(newBlockNumber),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	Tasks.Register(task.DefineTask(
		TypeIndexActionCheckLikeProtocolPayload,
		HandleIndexActionCheckLikeProtocol,
	))
}
