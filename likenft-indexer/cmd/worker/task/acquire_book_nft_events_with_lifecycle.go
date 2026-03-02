package task

import (
	"context"
	"encoding/json"
	"fmt"

	appcontext "likenft-indexer/cmd/worker/context"
	"likenft-indexer/ent"
	"likenft-indexer/internal/database"
	"likenft-indexer/internal/logic/nftclassacquirebooknftevent"
	"likenft-indexer/internal/worker/task"

	"github.com/hibiken/asynq"
)

const TypeAcquireBookNFTEventsTaskPayloadWithLifecyclePayload = "acquire-book-nft-events-with-lifecycle"

type AcquireBookNFTEventsTaskPayloadWithLifecyclePayload struct {
	// Deprecated: use ContractAddresses instead. Kept for backward compatibility
	// with in-flight tasks during rolling deploy.
	ContractAddress   string   `json:"ContractAddress,omitempty"`
	ContractAddresses []string `json:"ContractAddresses,omitempty"`
}

func (p *AcquireBookNFTEventsTaskPayloadWithLifecyclePayload) GetAddresses() []string {
	if len(p.ContractAddresses) > 0 {
		return p.ContractAddresses
	}
	if p.ContractAddress != "" {
		return []string{p.ContractAddress}
	}
	return nil
}

func NewTypeAcquireBookNFTEventsTaskPayloadWithLifecycle(contractAddresses []string) (*asynq.Task, error) {
	payload, err := json.Marshal(AcquireBookNFTEventsTaskPayloadWithLifecyclePayload{
		ContractAddresses: contractAddresses,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeAcquireBookNFTEventsTaskPayloadWithLifecyclePayload,
		payload,
		asynq.Queue(TypeAcquireBookNFTEventsTaskPayloadWithLifecyclePayload),
	), nil
}

func handlerWithLifecycle(
	handler func(ctx context.Context, t *asynq.Task) error,
) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		config := appcontext.ConfigFromContext(ctx)
		logger := appcontext.LoggerFromContext(ctx)
		dbService := database.New()
		nftClassAcquireBookNFTEventsRepository := database.MakeNFTClassAcquireBookNFTEventsRepository(dbService)

		mylogger := logger.WithGroup("HandlerWithLifecycle")

		var p AcquireBookNFTEventsTaskPayloadWithLifecyclePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}

		addresses := p.GetAddresses()
		if len(addresses) == 0 {
			mylogger.Error("no contract addresses in payload")
			return nil
		}

		innerTask, err := NewAcquireBookNFTEventsTask(addresses)
		if err != nil {
			return fmt.Errorf("NewAcquireBookNFTEventsTask: %v", err)
		}

		calculateNextScoreFn := nftclassacquirebooknftevent.MakeCalculateNextProcessingScoreFn(
			config.TaskAcquireBookNFTNextProcessingBlockHeightWeight,
			config.TaskAcquireBookNFTNextProcessingTimeFloor,
			config.TaskAcquireBookNFTNextProcessingTimeCeiling,
			config.TaskAcquireBookNFTNextProcessingTimeWeight,
		)
		calculateTimeoutScoreFn := nftclassacquirebooknftevent.MakeCalculateTimeoutScoreFn(
			config.TaskAcquireBookNFTInProgressTimeoutSeconds,
		)
		calculateRetryScoreFn := nftclassacquirebooknftevent.MakeCalculateRetryScoreFn(
			config.TaskAcquireBookNFTRetryInitialTimeoutSeconds,
			config.TaskAcquireBookNFTRetryExponentialBackoffCoeff,
			config.TaskAcquireBookNFTRetryMaxTimeoutSeconds,
		)

		// Build lifecycle objects for all addresses in the batch
		var lifecycles []nftclassacquirebooknftevent.NFTClassAcquireBookNFTEventLifecycle
		for _, addr := range addresses {
			lifecycle, err := nftclassacquirebooknftevent.MakeNFTClassAcquireBookNFTEventLifecycleFromAddress(
				ctx,
				nftClassAcquireBookNFTEventsRepository,
				addr,
				calculateNextScoreFn,
				calculateTimeoutScoreFn,
				calculateRetryScoreFn,
			)
			if err != nil {
				mylogger.Error("MakeNFTClassAcquireBookNFTEventLifecycleFromAddress", "addr", addr, "err", err)
				continue
			}
			lifecycles = append(lifecycles, lifecycle)
		}

		if len(lifecycles) == 0 {
			mylogger.Error("no valid lifecycles for batch")
			return nil
		}

		// Run the inner handler once, then transition all lifecycles.
		// The first lifecycle executes the handler; the rest reuse its result.
		var handlerErr error
		handlerExecuted := false

		for _, lifecycle := range lifecycles {
			if err := lifecycle.WithEnqueued(ctx, func(nftClass *ent.NFTClass) error {
				if !handlerExecuted {
					handlerExecuted = true
					handlerErr = handler(ctx, innerTask)
				}
				return handlerErr
			}); err != nil {
				mylogger.Error("lifecycle.WithEnqueued", "err", err)
			}
		}

		return nil
	}
}

func init() {
	Tasks.Register(task.DefineTask(
		TypeAcquireBookNFTEventsTaskPayloadWithLifecyclePayload,
		handlerWithLifecycle(HandleAcquireBookNFTEventsTask),
	))
}
