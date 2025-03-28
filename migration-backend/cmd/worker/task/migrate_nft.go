package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hibiken/asynq"

	appcontext "github.com/likecoin/like-migration-backend/cmd/worker/context"
	"github.com/likecoin/like-migration-backend/pkg/likenft/cosmos"
	"github.com/likecoin/like-migration-backend/pkg/likenft/evm"
	"github.com/likecoin/like-migration-backend/pkg/logic/likenft"
	apptask "github.com/likecoin/like-migration-backend/pkg/task"
)

func HandleMigrateNFTTask(ctx context.Context, t *asynq.Task) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	envCfg := appcontext.ConfigFromContext(ctx)
	db := appcontext.DBFromContext(ctx)
	logger := appcontext.LoggerFromContext(ctx)

	mylogger := logger.WithGroup("HandleMigrateNFTTask")

	var p apptask.MigrateNFTPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	mylogger = mylogger.With("LikenftAssetMigrationNFTId", p.LikenftAssetMigrationNFTId)

	likenftClient := &cosmos.LikeNFTCosmosClient{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		NodeURL: envCfg.CosmosNodeUrl,
	}

	ethClient, err := ethclient.Dial(envCfg.EthNetworkPublicRPCURL)
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(envCfg.EthWalletPrivateKey)
	if err != nil {
		return err
	}
	contractAddress := common.HexToAddress(envCfg.EthLikeNFTContractAddress)

	evmLikeProtocolClient := evm.NewLikeProtocol(
		logger,
		ethClient,
		privateKey,
		envCfg.EthChainId,
		contractAddress,
	)
	evmLikeNFTClient := evm.NewBookNFT(
		logger,
		ethClient,
		privateKey,
		envCfg.EthChainId,
	)

	mylogger.Info("running migrate nft")
	mn, err := likenft.MigrateNFTFromAssetMigration(
		ctx,
		logger,
		db,
		likenftClient,
		&evmLikeProtocolClient,
		&evmLikeNFTClient,
		envCfg.InitialNewClassOwner,
		envCfg.InitialNewClassMinter,
		envCfg.InitialNewClassUpdater,
		envCfg.InitialBatchMintNFTsOwner,
		p.LikenftAssetMigrationNFTId,
	)

	if err != nil {
		mylogger.Error("running migrate nft failed", "error", err)
		return err
	}

	mylogger.Info("migrate nft completed", "evmTxHash", *mn.EvmTxHash)
	return nil
}
