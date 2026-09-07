package cmd

import (
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
the history that produced it.`,
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

		// Every read is pinned to one block. Left at "latest" the snapshot
		// would be stitched together from thousands of calls landing on
		// different blocks, and could not be internally consistent.
		blockNumber := big.NewInt(block)
		if block == 0 {
			header, err := ethClient.HeaderByNumber(ctx, nil)
			if err != nil {
				panic(fmt.Errorf("failed to get latest block header: %w", err))
			}
			blockNumber = new(big.Int).Sub(header.Number, big.NewInt(confirmations))
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

		snapshot, err := resync.Build(ctx, logger, evmClient, dbService, blockNumber, concurrency)
		if err != nil {
			panic(fmt.Errorf("failed to build snapshot: %w", err))
		}

		if report != "" {
			reportBytes, err := json.MarshalIndent(snapshot.Changes, "", "  ")
			if err != nil {
				panic(fmt.Errorf("failed to marshal report: %w", err))
			}
			if err := os.WriteFile(report, reportBytes, 0o600); err != nil {
				panic(fmt.Errorf("failed to write report %s: %w", report, err))
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
			return
		}

		err = persistor.MakeStakingStatePersistor(dbService).Persist(
			ctx,
			nil,
			snapshot.Accounts,
			snapshot.NFTClasses,
			snapshot.Stakings,
		)
		if err != nil {
			panic(fmt.Errorf("failed to persist snapshot: %w", err))
		}

		logger.Info("snapshot applied", "block_number", blockNumber)
	},
}

func init() {
	resyncCmd.Flags().String("rpc", "", "EVM RPC URL")
	resyncCmd.Flags().String("like-collective-address", "", "LikeCollective contract address")
	resyncCmd.Flags().String("like-stake-position-address", "", "LikeStakePosition contract address")
	resyncCmd.Flags().Int64("block", 0, "Block to snapshot; 0 reads the head less --confirmations")
	resyncCmd.Flags().Int64("confirmations", 5, "Blocks to stay behind the head when --block is 0")
	resyncCmd.Flags().Int("concurrency", 8, "Concurrent RPC calls")
	resyncCmd.Flags().Bool("apply", false, "Write the snapshot; without it the diff is only reported")
	resyncCmd.Flags().String("report", "", "Write the full diff to this path as JSON")
	rootCmd.AddCommand(resyncCmd)
}
