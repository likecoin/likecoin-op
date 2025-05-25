package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	appcontext "likenft-indexer/cmd/worker/context"
	"likenft-indexer/ent"
	"likenft-indexer/ent/evmevent"
	"likenft-indexer/internal/database"
	"likenft-indexer/internal/model"

	"github.com/hibiken/asynq"
)

const TypeCheckReceivedEVMEventsPayload = "check-received-evm-events"

type CheckReceivedEVMEventsPayload struct {
}

func NewCheckReceivedEVMEventsTask() (*asynq.Task, error) {
	payload, err := json.Marshal(CheckReceivedEVMEventsPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCheckReceivedEVMEventsPayload, payload), nil
}

func HandleCheckReceivedEVMEvents(ctx context.Context, t *asynq.Task) error {
	logger := appcontext.LoggerFromContext(ctx)
	asynqClient := appcontext.AsynqClientFromContext(ctx)

	mylogger := logger.WithGroup("HandleCheckReceivedEVMEvents")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var p CheckReceivedEVMEventsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		mylogger.Error("json.Unmarshal EnqueueProcessEvmEventPayload", "err", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	dbService := database.New()

	evmEventRepository := database.MakeEVMEventRepository(dbService)

	receivedEvents, err := evmEventRepository.GetEVMEventsByStatus(ctx, evmevent.StatusReceived)

	if err != nil {
		return err
	}
	mylogger.Info(fmt.Sprintf("%d events found", len(receivedEvents)))

	receivedEventIdsInProcessingOrder := slices.Clone(receivedEvents)
	slices.SortFunc(receivedEventIdsInProcessingOrder, model.EvmEventsProcessingComparator)
	receivedEventIds := make([]int, 0, len(receivedEventIdsInProcessingOrder))

	for _, evmEvent := range receivedEventIdsInProcessingOrder {
		err := handleCheckReceivedEVMEvents_enqueueProcessEVMEvent(
			mylogger,
			asynqClient,
			evmEvent,
		)
		if err != nil {
			mylogger.Error("enqueueProcessEVMEvent", "err", err)
			continue
		}
		receivedEventIds = append(receivedEventIds, evmEvent.ID)
	}

	err = evmEventRepository.BatchUpdateEvmEventStatusByIds(ctx, receivedEventIds, evmevent.StatusEnqueued)

	if err != nil {
		mylogger.Error("evmEventRepository.BatchUpdateEvmEventStatusByIds", "err", err)
		return err
	}

	return nil
}

func handleCheckReceivedEVMEvents_enqueueProcessEVMEvent(
	logger *slog.Logger,
	asynqClient *asynq.Client,

	evmEvent *ent.EVMEvent,
) error {
	mylogger := logger.WithGroup("enqueueProcessEVMEvent").
		With("evmEventId", evmEvent.ID)
	mylogger.Info("Enqueueing NewProcessEVMEvent task...")
	t, err := NewProcessEVMEvent(evmEvent.ID)
	if err != nil {
		mylogger.Error("Cannot create task", "err", err)
		return err
	}
	taskInfo, err := asynqClient.Enqueue(t, asynq.MaxRetry(0))
	if err != nil {
		mylogger.Error("Cannot enqueue task", "err", err)
		return err
	}
	mylogger.Info("task enqueued", "taskId", taskInfo.ID)
	return nil
}
