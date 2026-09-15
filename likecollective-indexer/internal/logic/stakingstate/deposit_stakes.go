package stakingstate

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/stakingstate/model"
	"likecollective-indexer/internal/util/parallel"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// ErrHolderNotIndexed is returned when the loaded stakings of a pool add up to
// less than its total stake at a deposit, so some holder has no staking row.
var ErrHolderNotIndexed = errors.New("a holder of the pool is not indexed")

// depositStakes is how a pool's stake stood at the block a RewardDeposited was
// emitted in: what the deposit is split by.
type depositStakes struct {
	total  *uint256.Int
	staked map[common.Address]*uint256.Int
}

// readDepositStakes reads, for every RewardDeposited among stakingEvents, the
// pool's stake at the deposit's own block, and records it on state for the
// application to split the reward by.
//
// On chain-backed state the loaded amounts cannot be used for the split. They
// are the ones reconcileFromChain last wrote, read at whatever head it saw --
// when the queue is behind, an earlier event can have read them at a head past
// this deposit, so splitting by them hands the reward to whoever held the stake
// later. And the applications do not move them, so a stake just before the
// deposit would not count either. Those amounts are materialized totals; the
// reward_deposit_distributed rows are history, and have to be split by the
// stake the deposit actually met.
//
// Unlike the reconcile, these reads are pinned to the event's block, and so
// need archive state once the event is older than the node's window. That is
// safe here where it is not for the totals: what a pool held at a past block
// never changes, so a late or re-driven deposit reads the same answer as a
// prompt one.
//
// The read is of the state after the whole block, so a stake or unstake landing
// in the same block after the deposit is counted as if it came before. The
// contract offers no read between two logs of one block.
//
// Every loaded staking of the pool is read, empty or not: holding nothing at
// the head says nothing about what a row held at the deposit.
func readDepositStakes(
	ctx context.Context,
	evmClient evm.EVMClient,
	state *stakingState,
	stakingEvents []*ent.StakingEvent,
) error {
	for _, stakingEvent := range stakingEvents {
		if stakingEvent.EventType != stakingevent.EventTypeRewardDeposited {
			continue
		}

		blockNumber := big.NewInt(0).SetUint64(uint64(stakingEvent.BlockNumber))
		nftClassAddress := common.HexToAddress(stakingEvent.NftClassAddress)

		// The total is one more read beside the stakers', so it goes in the same
		// fan-out as a nil entry rather than a round trip ahead of it.
		stakings := state.GetStakingsByNFTClassAddress(nftClassAddress)
		amounts, err := parallel.MapWithLimit(ctx, reconcileConcurrency, append(
			[]*model.Staking{nil}, stakings...,
		), func(
			ctx context.Context,
			staking *model.Staking,
		) (*uint256.Int, error) {
			if staking == nil {
				return readTotalStake(ctx, evmClient, blockNumber, nftClassAddress)
			}
			return readStake(ctx, evmClient, blockNumber, staking)
		})
		if err != nil {
			return err
		}
		total, amounts := amounts[0], amounts[1:]

		stakes := &depositStakes{
			total:  total,
			staked: make(map[common.Address]*uint256.Int, len(stakings)),
		}
		sum := uint256.NewInt(0)
		for i, staking := range stakings {
			stakes.staked[staking.AccountEVMAddress] = amounts[i]
			sum.Add(sum, amounts[i])
		}
		// Only loaded stakings are read, so a holder whose row is missing -- a
		// Staked not yet processed, or never delivered -- would be left out of
		// the split for good. Fail instead, and retry once the row exists.
		//
		// Only a shortfall is refused. The contract's total has drifted below
		// the sum of its positions before (restake and increaseStakeToPosition
		// left pending rewards out of it), and a past block's reads never
		// change, so refusing a surplus would fail such a deposit forever.
		if sum.Lt(total) {
			return fmt.Errorf(
				"%w: stakings of %s sum to %s at block %s, below the total stake %s",
				ErrHolderNotIndexed, nftClassAddress, sum, blockNumber, total,
			)
		}
		state.setDepositStakes(stakingEvent, stakes)
	}
	return nil
}
