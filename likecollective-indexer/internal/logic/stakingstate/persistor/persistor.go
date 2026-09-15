package persistor

import (
	"context"
	"fmt"
	"math/big"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/logic/stakingstate/model"
)

type StakingStatePersistor interface {
	// WithLock runs fn as the only writer of the staking state. A caller that
	// loads state, re-reads it from the chain and persists it holds the lock
	// across all three.
	WithLock(ctx context.Context, fn func(ctx context.Context) error) error

	// Persist writes the state. headBlockNumber is the block its amounts were
	// read from the chain at, and fails the write with
	// database.ErrStakingStateHeadBehind if the state has already been
	// committed at a newer one; nil writes without that check.
	Persist(
		ctx context.Context,
		headBlockNumber *big.Int,
		stakingEvents []*ent.StakingEvent,
		accounts []*model.Account,
		nftClasses []*model.NFTClass,
		stakings []*model.Staking,
	) error

	// AlreadyApplied reports whether the log's delta has been persisted. The
	// staking events are committed in the same transaction as the totals, so
	// one stored means the totals were written too, even if whatever the
	// caller does after Persist returned did not complete.
	AlreadyApplied(
		ctx context.Context,
		transactionHash string,
		transactionIndex uint,
		logIndex uint,
	) (bool, error)
}

type latestStakingStatePersistor struct {
	accountRepository      database.AccountRepository
	nftClassRepository     database.NFTClassRepository
	stakingRepository      database.StakingRepository
	stakingEventRepository database.StakingEventRepository
	headRepository         database.StakingStateHeadRepository
	dbService              database.Service
}

func MakeStakingStatePersistor(
	dbService database.Service,
) StakingStatePersistor {
	accountRepository := database.MakeAccountRepository(dbService)
	nftClassRepository := database.MakeNFTClassRepository(dbService)
	stakingRepository := database.MakeStakingRepository(dbService)
	stakingEventRepository := database.MakeStakingEventRepository(dbService)
	headRepository := database.MakeStakingStateHeadRepository(dbService)
	return &latestStakingStatePersistor{
		accountRepository,
		nftClassRepository,
		stakingRepository,
		stakingEventRepository,
		headRepository,
		dbService,
	}
}

func (p *latestStakingStatePersistor) WithLock(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return database.WithStakingStateLock(ctx, p.dbService, fn)
}

func (p *latestStakingStatePersistor) AlreadyApplied(
	ctx context.Context,
	transactionHash string,
	transactionIndex uint,
	logIndex uint,
) (bool, error) {
	return p.stakingEventRepository.HasStakingEventsForLog(ctx, transactionHash, transactionIndex, logIndex)
}

func (p *latestStakingStatePersistor) Persist(
	ctx context.Context,
	headBlockNumber *big.Int,
	stakingEvents []*ent.StakingEvent,
	accounts []*model.Account,
	nftClasses []*model.NFTClass,
	stakings []*model.Staking,
) error {
	err := database.WithStakingStateTx(ctx, p.dbService, func(tx *ent.Tx) error {
		allAlreadyStored, err := p.stakingEventRepository.InsertStakingEventsIfNeeded(
			ctx,
			tx,
			stakingEvents,
		)
		if err != nil {
			return fmt.Errorf("failed to insert staking events if needed: %w", err)
		}

		// A replay: an earlier attempt committed this transaction and failed
		// afterward. Callers are expected to ask AlreadyApplied first; this is
		// the atomic backstop. accounts, nftClasses and stakings are
		// loaded-state-plus-delta, and the loaded state already holds the
		// delta, so writing them again would apply the event twice.
		if allAlreadyStored {
			return nil
		}

		// The amounts are absolute values read at a head, so committing an
		// older head over a newer one would roll back every row this writes,
		// and nothing re-reads a row until some later event touches it. That
		// happens when a load-balanced RPC node lags the one the previous
		// writer asked; failing hands the event back to be retried.
		if headBlockNumber != nil {
			if !headBlockNumber.IsUint64() {
				return fmt.Errorf("head block number %s overflows uint64", headBlockNumber)
			}
			err := p.headRepository.AdvanceBlockNumber(ctx, tx, headBlockNumber.Uint64())
			if err != nil {
				return fmt.Errorf("failed to advance staking state head: %w", err)
			}
		}

		for _, account := range accounts {
			_, err := p.accountRepository.CreateOrUpdateAccount(
				ctx,
				tx,
				account.EVMAddress.String(),
				account.StakedAmount,
				account.PendingRewardAmount,
				account.ClaimedRewardAmount,
			)
			if err != nil {
				return fmt.Errorf("failed to create or update account: %w", err)
			}
		}

		for _, nftClass := range nftClasses {
			_, err := p.nftClassRepository.CreateOrUpdateNFTClass(
				ctx,
				tx,
				nftClass.EVMAddress.String(),
				nftClass.StakedAmount,
			)
			if err != nil {
				return fmt.Errorf("failed to create or update nft class: %w", err)
			}
		}

		for _, staking := range stakings {
			_, err := p.stakingRepository.CreateOrUpdateStaking(
				ctx,
				tx,
				staking.BookNFTEvmAddress.String(),
				staking.AccountEVMAddress.String(),
				staking.StakedAmount,
				staking.PendingRewardAmount,
				staking.ClaimedRewardAmount,
			)
			if err != nil {
				return fmt.Errorf("failed to create or update staking: %w", err)
			}
		}

		for _, nftClass := range nftClasses {
			_, err := p.nftClassRepository.RecomputeNumberOfStakersByNFTClassAddress(
				ctx,
				tx,
				nftClass.EVMAddress.String(),
			)
			if err != nil {
				return fmt.Errorf("failed to recompute number of stakers by nft class address: %w", err)
			}

			err = p.stakingRepository.RecomputePoolSharesByNFTClassAddress(
				ctx,
				tx,
				nftClass.EVMAddress.String(),
			)
			if err != nil {
				return fmt.Errorf("failed to recompute pool shares by nft class address: %w", err)
			}
		}

		return nil
	})

	return err
}
