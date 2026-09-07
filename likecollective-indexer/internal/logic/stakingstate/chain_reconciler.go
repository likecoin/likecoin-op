package stakingstate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/util/parallel"

	"github.com/holiman/uint256"
)

// reconcileConcurrency bounds the reads one event can fan out into. A
// RewardDeposited loads every staker the pool has ever had, and each read is a
// contract call that itself loops over that user's positions.
const reconcileConcurrency = 8

// ErrChainBehindEvent is returned when the node's head is older than the block
// of the event being applied, so a read from it would predate the event.
var ErrChainBehindEvent = errors.New("chain head is behind the event")

// reconcileFromChain replaces the amounts in state with what LikeCollective
// reports.
//
// The totals used to be built by accumulating event deltas, which makes them
// unrecoverable: one missing or wrong event does not delay a number, it offsets
// it, and every later event builds on the wrong base. Worse, the deltas were
// not always the contract's own arithmetic -- RewardDeposited re-derived the
// per-staker split with one integer floor per account, where the contract takes
// one per position against a reward index. Those two disagree by dust on every
// deposit, so the totals drifted by construction rather than by accident.
//
// So the event stops being the source of the amount and becomes the trigger to
// re-read it, and the delta arithmetic is left to write staking_events, which
// is history and has to stay accumulated.
//
// The reads are deliberately NOT pinned to the event's block. Initial delivery
// is enqueued oldest-first, but asynq retries land late and
// retry-failed-evm-events re-drives old events on purpose, so a pinned read
// would let an older event commit an older truth on top of a newer one and
// leave it there. Reading the head, every writer converges on the same answer
// whatever order they run in, which is the property that makes the totals
// recoverable. It also means no archive state is required, where a pinned read
// would need it for every event once retries push one past a non-archive
// node's window.
//
// The head is resolved once and every read in the pass is pinned to it. A
// RewardDeposited fans out into a read per staker, and with each read taking
// whatever block the node had at that moment, a block landing mid-pass would
// leave some rows from before it and some from after -- a set of totals that
// no block ever held, and that nothing corrects if no later event touches the
// rows read early. The head is also required to be at least as new as the
// event being applied: a node that has not caught up to the event's own block
// would otherwise persist a truth older than the event, with nothing left to
// re-drive it.
//
// claimed_reward_amount is left alone. It is lifetime cumulative and the
// contract keeps no such counter, so there is nothing to read it from.
func reconcileFromChain(
	ctx context.Context,
	logger *slog.Logger,
	evmClient evm.EVMClient,
	state *stakingState,
	eventBlockNumber *big.Int,
) error {
	mylogger := logger.WithGroup("reconcileFromChain")

	head, err := evmClient.LatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}
	if head.Cmp(eventBlockNumber) < 0 {
		return fmt.Errorf(
			"%w: head %s is behind event block %s",
			ErrChainBehindEvent, head, eventBlockNumber,
		)
	}

	// A row holding nothing on either side contributes nothing to any total,
	// and a pool keeps its fully unstaked rows forever, so re-reading them
	// would make a deposit's cost grow with the pool's whole history rather
	// than with its current stakers. Drift banked on such a row is left for
	// `cli resync`, which reads every row rather than the loaded ones. A row an
	// event named is read regardless: the event is the evidence it may hold
	// something now, and on chain-backed state its delta was not applied.
	targets := make([]*model.Staking, 0, len(state.stakings))
	for _, staking := range state.stakings {
		if staking.StakedAmount.IsZero() && staking.PendingRewardAmount.IsZero() &&
			!state.isTouched(staking) {
			continue
		}
		targets = append(targets, staking)
	}

	// Read concurrently, apply sequentially: the corrections accumulate onto
	// account rows shared between stakings.
	amounts, err := parallel.MapWithLimit(ctx, reconcileConcurrency, targets, func(
		ctx context.Context,
		staking *model.Staking,
	) (*stakingAmounts, error) {
		return readStakingAmounts(ctx, evmClient, head, staking)
	})
	if err != nil {
		return err
	}

	for _, amount := range amounts {
		staking := amount.staking

		if !amount.staked.Eq(staking.StakedAmount) || !amount.pending.Eq(staking.PendingRewardAmount) {
			// On chain-backed state a row an event named still holds its
			// loaded amount, so it differs from the chain by that event's own
			// delta as a matter of course; only an untouched row moving is
			// drift worth a warning.
			level := slog.LevelWarn
			if state.chainBacked && state.isTouched(staking) {
				level = slog.LevelDebug
			}
			mylogger.Log(ctx, level, "corrected a staking against the chain",
				"account", staking.AccountEVMAddress,
				"book_nft", staking.BookNFTEvmAddress,
				"staked_was", staking.StakedAmount,
				"staked_now", amount.staked,
				"pending_was", staking.PendingRewardAmount,
				"pending_now", amount.pending,
			)
		}

		// An account row totals that account's stakings, but only the stakings
		// this event touched are loaded, so it cannot be recomputed as a sum
		// here. Moving it by exactly the correction applied to the staking
		// keeps the two consistent with each other.
		if account, ok := state.GetAccountByAddress(staking.AccountEVMAddress); ok {
			account.StakedAmount = applyCorrection(
				mylogger, account.StakedAmount, staking.StakedAmount, amount.staked,
			)
			account.PendingRewardAmount = applyCorrection(
				mylogger, account.PendingRewardAmount, staking.PendingRewardAmount, amount.pending,
			)
		}

		staking.StakedAmount = amount.staked
		staking.PendingRewardAmount = amount.pending
	}

	for _, nftClass := range state.nftClasses {
		totalStake, err := evmClient.GetTotalStake(ctx, head, nftClass.EVMAddress)
		if err != nil {
			return fmt.Errorf("failed to get total stake of %s: %w", nftClass.EVMAddress, err)
		}
		staked, err := toUint256(totalStake)
		if err != nil {
			return fmt.Errorf("total stake of %s: %w", nftClass.EVMAddress, err)
		}
		nftClass.StakedAmount = staked
	}

	return nil
}

type stakingAmounts struct {
	staking *model.Staking
	staked  *uint256.Int
	pending *uint256.Int
}

func readStakingAmounts(
	ctx context.Context,
	evmClient evm.EVMClient,
	blockNumber *big.Int,
	staking *model.Staking,
) (*stakingAmounts, error) {
	stakedAmount, err := evmClient.GetStakeForUser(
		ctx, blockNumber, staking.AccountEVMAddress, staking.BookNFTEvmAddress,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get stake for user %s on %s: %w",
			staking.AccountEVMAddress, staking.BookNFTEvmAddress, err,
		)
	}
	pendingRewardAmount, err := evmClient.GetPendingRewardsForUser(
		ctx, blockNumber, staking.AccountEVMAddress, staking.BookNFTEvmAddress,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get pending rewards for user %s on %s: %w",
			staking.AccountEVMAddress, staking.BookNFTEvmAddress, err,
		)
	}

	staked, err := toUint256(stakedAmount)
	if err != nil {
		return nil, fmt.Errorf("staked amount of %s: %w", staking.AccountEVMAddress, err)
	}
	pending, err := toUint256(pendingRewardAmount)
	if err != nil {
		return nil, fmt.Errorf("pending reward amount of %s: %w", staking.AccountEVMAddress, err)
	}

	return &stakingAmounts{staking: staking, staked: staked, pending: pending}, nil
}

// applyCorrection moves total by (after - before), clamped at zero.
//
// The clamp only fires when total is already lower than one of the stakings it
// is supposed to contain, which means the row was corrupt before this ran --
// and zeroing it discards whatever that account holds in pools this event did
// not load. That is worth hearing about, so it is reported rather than
// swallowed; `cli resync` rebuilds the row from every staking it has.
func applyCorrection(
	logger *slog.Logger,
	total *uint256.Int,
	before *uint256.Int,
	after *uint256.Int,
) *uint256.Int {
	if after.Cmp(before) >= 0 {
		return new(uint256.Int).Add(total, new(uint256.Int).Sub(after, before))
	}
	decrease := new(uint256.Int).Sub(before, after)
	if total.Lt(decrease) {
		logger.Error("account total is lower than the staking it contains; clamping to zero",
			"account_total", total,
			"staking_was", before,
			"staking_now", after,
		)
		return uint256.NewInt(0)
	}
	return new(uint256.Int).Sub(total, decrease)
}

func toUint256(value *big.Int) (*uint256.Int, error) {
	converted, overflow := uint256.FromBig(value)
	if overflow {
		return nil, fmt.Errorf("value %s overflows uint256", value)
	}
	return converted, nil
}
