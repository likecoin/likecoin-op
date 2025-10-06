package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appcontext "likenft-indexer/cmd/worker/context"
	"likenft-indexer/ent"
	"likenft-indexer/internal/database"
	"likenft-indexer/internal/logic/nftclassacquirebooknftevent"
	"likenft-indexer/internal/worker/task"

	"github.com/hibiken/asynq"
)

const TypeCheckBookNFTsPayload = "check-book-nfts"

type CheckBookNFTsPayload struct {
}

func NewCheckBookNFTsTask() (*asynq.Task, error) {
	payload, err := json.Marshal(CheckBookNFTsPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeCheckBookNFTsPayload,
		payload,
		asynq.Queue(TypeCheckBookNFTsPayload),
	), nil
}

func HandleCheckBookNFTs(ctx context.Context, t *asynq.Task) error {
	logger := appcontext.LoggerFromContext(ctx)
	config := appcontext.ConfigFromContext(ctx)
	inspector := appcontext.AsynqInspectorFromContext(ctx)

	mylogger := logger.WithGroup("HandleCheckBookNFTs")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var p CheckBookNFTsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	queueInfo, err := inspector.GetQueueInfo(TypeAcquireBookNFTEventsTaskPayload)
	currentLength := 0
	if err != nil {
		mylogger.Info("get queue info error", "err", err)
	} else {
		currentLength = queueInfo.Pending + queueInfo.Active
	}
	numberOfItemsToBeFetched := config.TaskAcquireBookNFTMaxQueueLength - currentLength

	dbService := database.New()
	nftClassAcquireBookNFTEventsRepository := database.MakeNFTClassAcquireBookNFTEventsRepository(dbService)

	timeoutScoreFn := nftclassacquirebooknftevent.MakeCalculateTimeoutScoreFn(
		config.TaskAcquireBookNFTInProgressTimeoutSeconds,
	)

	nftClasses, err := nftClassAcquireBookNFTEventsRepository.
		RequestForEnqueue(
			ctx,
			numberOfItemsToBeFetched,
			timeoutScoreFn(time.Now().UTC()),
		)

	if err != nil {
		return fmt.Errorf("nftClassAcquireBookNFTEventsRepository.RequestForEnqueue: %v", err)
	}

	myLogger := logger.With("batchSize", len(nftClasses))
	mylogger.Info("nft classes found")

	lifecycleObjects := nftclassacquirebooknftevent.MakeNFTClassAcquireBookNFTEventLifecycles(
		nftClassAcquireBookNFTEventsRepository,
		nftClasses,
		nftclassacquirebooknftevent.MakeCalculateNextProcessingScoreFn(
			config.TaskAcquireBookNFTNextProcessingBlockHeightWeight,
			config.TaskAcquireBookNFTNextProcessingTimeFloor,
			config.TaskAcquireBookNFTNextProcessingTimeCeiling,
			config.TaskAcquireBookNFTNextProcessingTimeWeight,
		),
		timeoutScoreFn,
		nftclassacquirebooknftevent.MakeCalculateRetryScoreFn(
			config.TaskAcquireBookNFTRetryInitialTimeoutSeconds,
			config.TaskAcquireBookNFTRetryExponentialBackoffCoeff,
			config.TaskAcquireBookNFTRetryMaxTimeoutSeconds,
		),
	)

	myLogger.Info("Enqueueing EnqueueAcquireBookNFTTask tasks...")
	for _, lifecycleObject := range lifecycleObjects {
		if _, err := lifecycleObject.WithEnqueueing(ctx, func(nftClass *ent.NFTClass) (*asynq.TaskInfo, error) {
			return EnqueueAcquireBookNFTTask(ctx, nftClass.Address)
		}); err != nil {
			myLogger.ErrorContext(ctx, "lifecycleObject.WithEnqueueing", "err", err)
		}
	}

	return nil
}

func EnqueueAcquireBookNFTTask(ctx context.Context, contractAddress string) (*asynq.TaskInfo, error) {
	asynqClient := appcontext.AsynqClientFromContext(ctx)

	t, err := NewTypeAcquireBookNFTEventsTaskPayloadWithLifecycle(contractAddress)
	if err != nil {
		return nil, err
	}

	taskInfo, err := asynqClient.Enqueue(t, asynq.MaxRetry(0))
	if err != nil {
		return nil, err
	}
	return taskInfo, nil
}

func init() {
	t := Tasks.Register(task.DefineTask(
		TypeCheckBookNFTsPayload,
		HandleCheckBookNFTs,
	))
	PeriodicTasks.Register(task.DefinePeriodicTask(
		t,
		NewCheckBookNFTsTask,
	))
}
