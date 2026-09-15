package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"

	"likecollective-indexer/internal/database"
	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/resync"
	"likecollective-indexer/internal/logic/stakingstate/persistor"

	clicontext "likecollective-indexer/internal/cli/context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var resyncCmd = &cobra.Command{
	Use:   "resync",
	Short: "Rebuild accounts, nft_classes and stakings from an on-chain snapshot",
	Long: `Rebuild the materialised staking state from a single on-chain snapshot.

The webhook is this indexer's only ingest path, so a delivery that never
arrives leaves the running totals in accounts, nft_classes and stakings wrong
indefinitely. This reads the same numbers straight off LikeCollective and
overwrites them.

Reports the diff and exits without writing unless --apply is given. Does not
touch staking_events or evm_events, so it repairs the current state only, not
the history that produced it.

With --apply it holds the lock the workers take around each event, so event
processing pauses until the snapshot is written. The snapshot block is never
older than the head the state was last committed at.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		envCfg := clicontext.ConfigFromContext(cmd.Context())
		logger := slog.New(slog.Default().Handler())
		ctx := cmd.Context()

		rpc, err := cmd.Flags().GetString("rpc")
		if err != nil {
			panic(fmt.Errorf("failed to get rpc: %w", err))
		}
		if rpc == "" {
			rpc = envCfg.EthNetworkPublicRPCURL
		}
		likeCollectiveAddress, err := cmd.Flags().GetString("like-collective-address")
		if err != nil {
			panic(fmt.Errorf("failed to get like collective address: %w", err))
		}
		if likeCollectiveAddress == "" {
			likeCollectiveAddress = envCfg.LikeCollectiveAddress
		}
		likeStakePositionAddress, err := cmd.Flags().GetString("like-stake-position-address")
		if err != nil {
			panic(fmt.Errorf("failed to get like stake position address: %w", err))
		}
		if likeStakePositionAddress == "" {
			likeStakePositionAddress = envCfg.LikeStakePositionAddress
		}
		block, err := cmd.Flags().GetInt64("block")
		if err != nil {
			panic(fmt.Errorf("failed to get block: %w", err))
		}
		confirmations, err := cmd.Flags().GetInt64("confirmations")
		if err != nil {
			panic(fmt.Errorf("failed to get confirmations: %w", err))
		}
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			panic(fmt.Errorf("failed to get concurrency: %w", err))
		}
		apply, err := cmd.Flags().GetBool("apply")
		if err != nil {
			panic(fmt.Errorf("failed to get apply: %w", err))
		}
		report, err := cmd.Flags().GetString("report")
		if err != nil {
			panic(fmt.Errorf("failed to get report: %w", err))
		}

		ethClient, err := ethclient.Dial(rpc)
		if err != nil {
			panic(fmt.Errorf("failed to dial eth client: %w", err))
		}

		evmClient := evm.NewEVMClient(
			common.HexToAddress(likeCollectiveAddress),
			common.HexToAddress(likeStakePositionAddress),
			ethClient,
		)

		dbService := database.New()
		defer func() {
			if err := dbService.Close(); err != nil {
				logger.Error("failed to close database", "error", err)
			}
		}()

		stakingStatePersistor := persistor.MakeStakingStatePersistor(dbService)

		run := func(ctx context.Context) error {
			// Every read is pinned to one block. Left at "latest" the snapshot
			// would be stitched together from thousands of calls landing on
			// different blocks, and could not be internally consistent.
			blockNumber := big.NewInt(block)
			if block == 0 {
				header, err := ethClient.HeaderByNumber(ctx, nil)
				if err != nil {
					return fmt.Errorf("failed to get latest block header: %w", err)
				}
				blockNumber = new(big.Int).Sub(header.Number, big.NewInt(confirmations))
			}

			// The workers re-read rows at the head, not the head less
			// confirmations, so a row may already hold a newer truth than this
			// snapshot would. Writing it would roll that row back, and the
			// event that wrote it is processed, so nothing would re-read it.
			committedHead, err := database.MakeStakingStateHeadRepository(dbService).GetBlockNumber(ctx)
			if err != nil {
				return fmt.Errorf("failed to get committed staking state head: %w", err)
			}
			committedBlockNumber := new(big.Int).SetUint64(committedHead)
			switch {
			case blockNumber.Cmp(committedBlockNumber) >= 0:
			case block == 0:
				logger.Info("raising snapshot block to the committed head",
					"block_number", blockNumber,
					"committed_block_number", committedBlockNumber,
				)
				blockNumber = committedBlockNumber
			case apply:
				return fmt.Errorf(
					"--block %s is older than block %s the staking state is committed at",
					blockNumber, committedBlockNumber,
				)
			default:
				logger.Warn("snapshot block is older than the committed head; the diff includes changes since",
					"block_number", blockNumber,
					"committed_block_number", committedBlockNumber,
				)
			}

			snapshot, err := resync.Build(ctx, logger, evmClient, dbService, blockNumber, concurrency)
			if err != nil {
				return fmt.Errorf("failed to build snapshot: %w", err)
			}

			if report != "" {
				reportBytes, err := json.MarshalIndent(snapshot.Changes, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal report: %w", err)
				}
				if err := os.WriteFile(report, reportBytes, 0o600); err != nil {
					return fmt.Errorf("failed to write report %s: %w", report, err)
				}
			}

			for _, change := range snapshot.Changes {
				fmt.Printf("%-11s %-45s %-22s %s -> %s\n",
					change.Table, change.Key, change.Field, change.Old, change.New)
			}

			logger.Info("snapshot built",
				"block_number", blockNumber,
				"positions", snapshot.Positions,
				"accounts", len(snapshot.Accounts),
				"nft_classes", len(snapshot.NFTClasses),
				"stakings", len(snapshot.Stakings),
				"changes", len(snapshot.Changes),
			)

			if !apply {
				logger.Info("dry run, nothing written; pass --apply to write")
				return nil
			}

			err = stakingStatePersistor.Persist(
				ctx,
				blockNumber,
				nil,
				snapshot.Accounts,
				snapshot.NFTClasses,
				snapshot.Stakings,
			)
			if err != nil {
				return fmt.Errorf("failed to persist snapshot: %w", err)
			}

			logger.Info("snapshot applied", "block_number", blockNumber)
			return nil
		}

		// Applying holds the staking state lock from before the block is
		// chosen until the write commits, so workers wait rather than commit
		// in between. Otherwise a claim processed mid-build would be rolled
		// back by the claimed_reward_amount the snapshot carried over, and a
		// head committed mid-build would be newer than the snapshot.
		if apply {
			err = stakingStatePersistor.WithLock(ctx, run)
		} else {
			err = run(ctx)
		}
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	resyncCmd.Flags().String("rpc", "", "EVM RPC URL")
	resyncCmd.Flags().String("like-collective-address", "", "LikeCollective contract address")
	resyncCmd.Flags().String("like-stake-position-address", "", "LikeStakePosition contract address")
	resyncCmd.Flags().Int64("block", 0, "Block to snapshot; 0 reads the head less --confirmations")
	resyncCmd.Flags().Int64("confirmations", 5, "Blocks to stay behind the head when --block is 0, never behind the committed head")
	resyncCmd.Flags().Int("concurrency", 8, "Concurrent RPC calls")
	resyncCmd.Flags().Bool("apply", false, "Write the snapshot; without it the diff is only reported")
	resyncCmd.Flags().String("report", "", "Write the full diff to this path as JSON")
	rootCmd.AddCommand(resyncCmd)
}
