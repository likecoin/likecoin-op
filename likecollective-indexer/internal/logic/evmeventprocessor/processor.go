package evmeventprocessor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/evmevent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/evm"
	stakingstateloader "likecollective-indexer/internal/logic/stakingstate/loader"
	stakingstatepersistor "likecollective-indexer/internal/logic/stakingstate/persistor"

	"github.com/ethereum/go-ethereum/common"
)

var ErrAlreadyProcessing = errors.New("evmevent already processing")

const failedStatusUpdateTimeout = 30 * time.Second

type EVMEventProcessor interface {
	Process(
		ctx context.Context,
		logger *slog.Logger,
		evmEvent *ent.EVMEvent,
	) error
}

type evmEventProcessor struct {
	evmClient             evm.EVMClient
	evmEventRepository    database.EVMEventRepository
	stakingStateLoader    stakingstateloader.StakingStateLoader
	stakingStatePersistor stakingstatepersistor.StakingStatePersistor

	likeCollectiveAddress    common.Address
	likeStakePositionAddress common.Address
}

func MakeEVMEventProcessor(
	evmClient evm.EVMClient,
	evmEventRepository database.EVMEventRepository,
	stakingStateLoader stakingstateloader.StakingStateLoader,
	stakingStatePersistor stakingstatepersistor.StakingStatePersistor,
	likeCollectiveAddress common.Address,
	likeStakePositionAddress common.Address,
) EVMEventProcessor {
	return &evmEventProcessor{
		evmClient,
		evmEventRepository,
		stakingStateLoader,
		stakingStatePersistor,
		likeCollectiveAddress,
		likeStakePositionAddress,
	}
}

func (e *evmEventProcessor) Process(
	ctx context.Context,
	logger *slog.Logger,

	evmEvent *ent.EVMEvent,
) (err error) {
	mylogger := logger.WithGroup("Process").
		With("evmEventId", evmEvent.ID)

	switch evmEvent.Status {
	case evmevent.StatusSkipped:
		mylogger.Warn("The event has already skipped")
		return nil
	case evmevent.StatusProcessing:
		mylogger.Warn("The event has already processing")
		return ErrAlreadyProcessing
	case evmevent.StatusProcessed:
		mylogger.Info("The event has already processed")
		return nil
	case evmevent.StatusFailed:
		mylogger.Info("record failed. retrying...")
	case evmevent.StatusReceived:
	case evmevent.StatusEnqueued:
		break
	default:
		return fmt.Errorf("unknown status %v", evmEvent.Status)
	}

	defer func() {
		if err != nil {
			errMsg := err.Error()
			mylogger.Error("something went wrong", "err", err)
			// Not on ctx: it is often the cancellation being reported, e.g. a
			// task timing out behind the staking state lock, and an event left
			// processing is never picked up again. Bounded so a stuck database
			// cannot hold up shutdown instead.
			updateCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), failedStatusUpdateTimeout)
			defer cancel()
			_, updateErr := e.evmEventRepository.UpdateEvmEventStatus(
				updateCtx, evmEvent, evmevent.StatusFailed, &errMsg,
			)
			if updateErr != nil {
				mylogger.Error("failed to mark evm event failed", "err", updateErr)
			}
		}
	}()

	// A failed event may already have been applied: the `processed` status
	// below is written after Persist committed, and if that write is what
	// failed, re-applying the delta against state that already holds it fails
	// validation before Persist's own replay guard runs. Ask first.
	if evmEvent.Status == evmevent.StatusFailed {
		var applied bool
		applied, err = e.stakingStatePersistor.AlreadyApplied(
			ctx,
			evmEvent.TransactionHash,
			evmEvent.TransactionIndex,
			evmEvent.LogIndex,
		)
		if err != nil {
			mylogger.Error("e.stakingStatePersistor.AlreadyApplied", "err", err)
			return err
		}
		if applied {
			mylogger.Info("staking events already persisted by an earlier attempt, marking processed")
			_, err = e.evmEventRepository.UpdateEvmEventStatus(ctx, evmEvent, evmevent.StatusProcessed, nil)
			return err
		}
	}

	evmEvent, err = e.evmEventRepository.UpdateEvmEventStatus(ctx, evmEvent, evmevent.StatusProcessing, nil)

	if err != nil {
		mylogger.Error("e.evmEventRepository.UpdateEvmEventStatus", "err", err)
		return err
	}

	processorCreator, err := e.getProcessorCreator(evmEvent)
	if err != nil {
		if errors.Is(err, UnknownEvent) {
			evmEvent, err = e.evmEventRepository.UpdateEvmEventStatus(ctx, evmEvent, evmevent.StatusSkipped, nil)
			return nil
		}
		mylogger.Error("e.getProcessorCreator", "err", err)
		return err
	}

	mylogger = mylogger.With(
		"processorCreator",
		runtime.FuncForPC(reflect.ValueOf(processorCreator).Pointer()).Name())

	processor := processorCreator(makeEventProcessorDeps(
		e.evmClient,
		e.stakingStateLoader,
		e.stakingStatePersistor,
		e.likeCollectiveAddress,
		e.likeStakePositionAddress,
	))

	err = processor.Process(ctx, logger, evmEvent)
	if err != nil {
		if errors.Is(err, UnknownEvent) {
			reason := err.Error()
			evmEvent, err = e.evmEventRepository.UpdateEvmEventStatus(
				ctx,
				evmEvent,
				evmevent.StatusSkipped,
				&reason,
			)
			return nil
		}
		mylogger.Error("processor.Process", "err", err)
		return err
	}

	evmEvent, err = e.evmEventRepository.UpdateEvmEventStatus(ctx, evmEvent, evmevent.StatusProcessed, nil)
	mylogger.Info("updated evm event status to processed")

	if err != nil {
		mylogger.Error("e.evmEventRepository.UpdateEvmEventStatus", "err", err)
		return err
	}

	return nil
}

func (e *evmEventProcessor) getProcessorCreator(
	evmEvent *ent.EVMEvent,
) (eventProcessorCreator, error) {
	processorCreator, err := getEventProcessor(evmEvent.Topic0)
	if err != nil {
		return nil, err
	}
	return processorCreator, nil
}
