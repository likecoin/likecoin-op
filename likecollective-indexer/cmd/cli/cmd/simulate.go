package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"likecollective-indexer/internal/evm"
	"likecollective-indexer/internal/logic/simulate"
	"likecollective-indexer/internal/logic/stakingstate"
	"likecollective-indexer/internal/logic/stakingstate/loader"

	clicontext "likecollective-indexer/internal/cli/context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate <input>",
	Short: "Simulate",
	Long:  `Simulate`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]

		envCfg := clicontext.ConfigFromContext(cmd.Context())
		logger := slog.New(slog.Default().Handler())

		jsonBytes, err := os.ReadFile(input)
		if err != nil {
			panic(fmt.Errorf("failed to read file %s: %w", input, err))
		}

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
		verify, err := cmd.Flags().GetBool("verify")
		if err != nil {
			panic(fmt.Errorf("failed to get verify: %w", err))
		}

		var simulation simulate.Simulation
		err = json.Unmarshal(jsonBytes, &simulation)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal JSON: %w", err))
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

		stakingStateLoader := loader.MakeEmptyStateLoader()

		echoResult := func(
			state *simulate.State,
		) error {
			jsonBytes, err := json.Marshal(state)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonBytes))
			return nil
		}

		verifyResult := func(
			simulatedState *simulate.State,
		) error {
			return simulatedState.Verify(&simulation.State)
		}

		handleResult := func(
			state *simulate.State,
		) error {
			if verify {
				return verifyResult(state)
			}
			return echoResult(state)
		}

		persistor := simulate.MakeSimulationPersistor(handleResult)

		stakingEvmEventProcessor := stakingstate.MakeSimulationStakingEvmEventProcessor(
			evmClient,
			stakingStateLoader,
			persistor,
			common.HexToAddress(likeCollectiveAddress),
			common.HexToAddress(likeStakePositionAddress),
		)

		simulateLogProcessor := simulate.MakeSimulateLogProcessor(
			stakingEvmEventProcessor,
			common.HexToAddress(likeCollectiveAddress),
			common.HexToAddress(likeStakePositionAddress),
		)

		err = simulateLogProcessor.Process(
			cmd.Context(),
			logger,
			simulation.Logs,
		)
		if err != nil {
			panic(fmt.Errorf("failed to process logs: %w", err))
		}
	},
}

func init() {
	simulateCmd.Flags().String("rpc", "", "EVM RPC URL")
	simulateCmd.Flags().String("like-collective-address", "", "LikeCollective contract address")
	simulateCmd.Flags().String("like-stake-position-address", "", "LikeStakePosition contract address")
	simulateCmd.Flags().Bool("verify", false, "Verify the simulation result")
	rootCmd.AddCommand(simulateCmd)
}
