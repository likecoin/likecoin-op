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

	batchSize := config.TaskAcquireBookNFTBatchSize
	if batchSize <= 0 {
		batchSize = 50
	}

	queueInfo, err := inspector.GetQueueInfo(TypeAcquireBookNFTEventsTaskPayloadWithLifecyclePayload)
	currentLength := 0
	if err != nil {
		mylogger.Info("get queue info error", "err", err)
	} else {
		currentLength = queueInfo.Pending + queueInfo.Active
	}
	// Each queued task now represents a batch of contracts
	estimatedContractsInFlight := currentLength * batchSize
	numberOfItemsToBeFetched := config.TaskAcquireBookNFTMaxQueueLength - estimatedContractsInFlight
	if numberOfItemsToBeFetched <= 0 {
		mylogger.Info("queue is full, skipping", "estimatedContractsInFlight", estimatedContractsInFlight)
		return nil
	}

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

	mylogger.Info("nft classes found", "count", len(nftClasses))

	calculateNextScoreFn := nftclassacquirebooknftevent.MakeCalculateNextProcessingScoreFn(
		config.TaskAcquireBookNFTNextProcessingBlockHeightWeight,
		config.TaskAcquireBookNFTNextProcessingTimeFloor,
		config.TaskAcquireBookNFTNextProcessingTimeCeiling,
		config.TaskAcquireBookNFTNextProcessingTimeWeight,
	)
	calculateRetryScoreFn := nftclassacquirebooknftevent.MakeCalculateRetryScoreFn(
		config.TaskAcquireBookNFTRetryInitialTimeoutSeconds,
		config.TaskAcquireBookNFTRetryExponentialBackoffCoeff,
		config.TaskAcquireBookNFTRetryMaxTimeoutSeconds,
	)

	lifecycleObjects := nftclassacquirebooknftevent.MakeNFTClassAcquireBookNFTEventLifecycles(
		nftClassAcquireBookNFTEventsRepository,
		nftClasses,
		calculateNextScoreFn,
		timeoutScoreFn,
		calculateRetryScoreFn,
	)

	// Split lifecycle objects into batches and enqueue one task per batch.
	// nftClasses and lifecycleObjects are parallel arrays (same length, same order).
	mylogger.Info("Enqueueing batched EnqueueAcquireBookNFTTask tasks...", "batchSize", batchSize)
	for i := 0; i < len(lifecycleObjects); i += batchSize {
		end := i + batchSize
		if end > len(lifecycleObjects) {
			end = len(lifecycleObjects)
		}

		// Get addresses from nftClasses directly (parallel to lifecycleObjects)
		batchAddresses := make([]string, 0, end-i)
		for j := i; j < end; j++ {
			batchAddresses = append(batchAddresses, nftClasses[j].Address)
		}

		// Transition each lifecycle via WithEnqueueing. The first successful
		// callback performs the real enqueue; subsequent callbacks reuse the
		// result. If the enqueue fails, all lifecycles get enqueueFailed
		// (no orphaned "enqueued" states).
		batch := lifecycleObjects[i:end]
		var taskInfo *asynq.TaskInfo
		var enqueueErr error
		taskEnqueued := false

		for _, lifecycleObject := range batch {
			if _, err := lifecycleObject.WithEnqueueing(ctx, func(nftClass *ent.NFTClass) (*asynq.TaskInfo, error) {
				if !taskEnqueued {
					taskEnqueued = true
					taskInfo, enqueueErr = EnqueueAcquireBookNFTTask(ctx, batchAddresses)
				}
				return taskInfo, enqueueErr
			}); err != nil {
				mylogger.ErrorContext(ctx, "lifecycleObject.WithEnqueueing", "err", err)
			}
		}
	}

	return nil
}

func EnqueueAcquireBookNFTTask(ctx context.Context, contractAddresses []string) (*asynq.TaskInfo, error) {
	asynqClient := appcontext.AsynqClientFromContext(ctx)

	t, err := NewTypeAcquireBookNFTEventsTaskPayloadWithLifecycle(contractAddresses)
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
