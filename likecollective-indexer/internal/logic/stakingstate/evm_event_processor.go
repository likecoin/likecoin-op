package stakingstate

import (
	"context"
	"log/slog"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/loader"
	"likecollective-indexer/internal/logic/stakingstate/persistor"

	"github.com/ethereum/go-ethereum/common"
)

type StakingEvmEventProcessor interface {
	Process(
		ctx context.Context,
		logger *slog.Logger,
		evmEvents []*ent.EVMEvent,
	) error
}

type stakingEvmEventProcessor struct {
	evmClient             evm.EVMClient
	stakingStateLoader    loader.StakingStateLoader
	stakingStatePersistor persistor.StakingStatePersistor

	likeCollectiveAddress    common.Address
	likeStakePositionAddress common.Address

	reconcileFromChain bool
}

func MakeStakingEvmEventProcessor(
	evmClient evm.EVMClient,
	stakingStateLoader loader.StakingStateLoader,
	stakingStatePersistor persistor.StakingStatePersistor,
	likeCollectiveAddress common.Address,
	likeStakePositionAddress common.Address,
) StakingEvmEventProcessor {
	return &stakingEvmEventProcessor{
		evmClient,
		stakingStateLoader,
		stakingStatePersistor,
		likeCollectiveAddress,
		likeStakePositionAddress,
		true,
	}
}

// MakeSimulationStakingEvmEventProcessor builds a processor that does not
// re-read the totals from the chain.
//
// Simulation exists to check the delta arithmetic against a known expected
// state. Reconciling would overwrite the numbers under test with the chain's
// own, and simulate would then be comparing the chain with itself -- which
// would pass whatever the arithmetic did.
func MakeSimulationStakingEvmEventProcessor(
	evmClient evm.EVMClient,
	stakingStateLoader loader.StakingStateLoader,
	stakingStatePersistor persistor.StakingStatePersistor,
	likeCollectiveAddress common.Address,
	likeStakePositionAddress common.Address,
) StakingEvmEventProcessor {
	return &stakingEvmEventProcessor{
		evmClient,
		stakingStateLoader,
		stakingStatePersistor,
		likeCollectiveAddress,
		likeStakePositionAddress,
		false,
	}
}

func (e *stakingEvmEventProcessor) Process(
	ctx context.Context,
	logger *slog.Logger,
	evmEvents []*ent.EVMEvent,
) error {
	stakingEvents := make([]*ent.StakingEvent, 0)

	for _, evmEvent := range evmEvents {
		var (
			stakingEvent []*ent.StakingEvent
			err          error
		)
		if common.HexToAddress(evmEvent.Address) == e.likeCollectiveAddress {
			stakingEvent, err = GetStakingEventsFromLikeCollectiveEvent(ctx, e.evmClient, evmEvent)
		} else if common.HexToAddress(evmEvent.Address) == e.likeStakePositionAddress {
			stakingEvent, err = GetStakingEventsFromLikeStakePositionEvent(ctx, e.evmClient, evmEvent)
		}
		if err != nil {
			return err
		}
		stakingEvents = append(stakingEvents, stakingEvent...)
	}

	stakingState, err := loadStakingState(
		ctx, e.stakingStateLoader, stakingEvents, e.reconcileFromChain,
	)
	if err != nil {
		return err
	}

	processedState, processedStakingEvents, err := stakingState.Process(stakingEvents)
	if err != nil {
		return err
	}

	// The applications above have produced staking_events -- the history. On
	// chain-backed state they left the staked and pending totals as loaded,
	// neither moving them by the deltas nor checking the deltas against them:
	// the loaded totals may already be from a head past these events. Now
	// replace those totals with what the contract itself reports, so a delta
	// that was wrong, or an event that never arrived, does not leave a
	// permanent offset behind.
	if e.reconcileFromChain {
		if err := reconcileFromChain(
			ctx, logger, e.evmClient, processedState, newestBlockNumber(evmEvents),
		); err != nil {
			return err
		}
	}

	err = processedState.Persist(ctx, processedStakingEvents, e.stakingStatePersistor)
	if err != nil {
		return err
	}

	return nil
}

// newestBlockNumber is the block of the newest event in the batch: the reads
// that replace its totals must come from a head at least that new.
func newestBlockNumber(evmEvents []*ent.EVMEvent) *big.Int {
	var newest uint64
	for _, evmEvent := range evmEvents {
		if uint64(evmEvent.BlockNumber) > newest {
			newest = uint64(evmEvent.BlockNumber)
		}
	}
	return big.NewInt(0).SetUint64(newest)
}
