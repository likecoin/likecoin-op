package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	appcontext "likecollective-indexer/cmd/worker/context"
	"likecollective-indexer/ent"
	"likecollective-indexer/ent/evmevent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/model"
	"likecollective-indexer/internal/worker/task"

	"github.com/hibiken/asynq"
)

const TypeCheckReceivedEVMEventsPayload = "check-received-evm-events"

// processEVMEventMaxRetry lets asynq retry a process-evm-event task that failed
// transiently -- an RPC blip, a lock wait. Without it a single failure parked
// the event at `failed`, where nothing picked it up again: the running totals
// are accumulated, so a skipped event does not delay a number, it offsets it
// permanently. retry-failed-evm-events is the slower backstop for whatever
// survives these attempts.
const processEVMEventMaxRetry = 5

type CheckReceivedEVMEventsPayload struct {
}

func NewCheckReceivedEVMEventsTask() (*asynq.Task, error) {
	payload, err := json.Marshal(CheckReceivedEVMEventsPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeCheckReceivedEVMEventsPayload,
		payload,
		asynq.Queue(TypeCheckReceivedEVMEventsPayload),
	), nil
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

	// Enqueue oldest first. GetEVMEventsByStatus does not order, so without
	// this the deltas that build staking_events are recorded in whatever order
	// the query happened to return.
	slices.SortFunc(receivedEvents, model.EvmEventsProcessingComparator)
	receivedEventIds := make([]int, 0, len(receivedEvents))

	for _, evmEvent := range receivedEvents {
		err := enqueueProcessEVMEvent(
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

func enqueueProcessEVMEvent(
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
	taskInfo, err := asynqClient.Enqueue(t, asynq.MaxRetry(processEVMEventMaxRetry))
	if err != nil {
		mylogger.Error("Cannot enqueue task", "err", err)
		return err
	}
	mylogger.Info("task enqueued", "taskId", taskInfo.ID)
	return nil
}

func init() {
	t := Tasks.Register(task.DefineTask(
		TypeCheckReceivedEVMEventsPayload,
		HandleCheckReceivedEVMEvents,
	))
	PeriodicTasks.Register(task.DefinePeriodicTask(
		t,
		NewCheckReceivedEVMEventsTask,
	))
}
