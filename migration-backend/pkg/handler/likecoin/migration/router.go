package migration

import (
	"database/sql"
	"net/http"

	"github.com/hibiken/asynq"

	"github.com/likecoin/like-migration-backend/pkg/likecoin/cosmos"
	"github.com/likecoin/like-migration-backend/pkg/likecoin/evm"
	"github.com/likecoin/like-migration-backend/pkg/signer"
)

type MigrationRouter struct {
	Db                           *sql.DB
	AsynqClient                  *asynq.Client
	Signer                       *signer.SignerClient
	EvmLikeCoinClient            *evm.LikeCoin
	CosmosLikcCoinClient         *cosmos.LikeCoin
	LikecoinBurningCosmosAddress string
}

func (h *MigrationRouter) Router() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /migration/evm-pool-balance", &GetEvmPoolBalanceHandler{
		Db:                h.Db,
		Signer:            h.Signer,
		EvmLikeCoinClient: h.EvmLikeCoinClient,
	})

	router.Handle("POST /migration/eth-signing-message", &CreateEthSigningMessageHandler{
		EvmLikeCoinClient:    h.EvmLikeCoinClient,
		CosmosLikcCoinClient: h.CosmosLikcCoinClient,
	})
	router.Handle("POST /migration/cosmos-memo-data", &CreateCosmosMemoDataHandler{})
	router.Handle("GET /migration/{cosmosWalletAddress}", &GetLatestLikeCoinMigrationHandler{
		Db: h.Db,
	})
	router.Handle("POST /migration", &CreateLikeCoinMigrationHandler{
		Db:                           h.Db,
		Signer:                       h.Signer,
		LikecoinBurningCosmosAddress: h.LikecoinBurningCosmosAddress,
	})
	router.Handle("PUT /migration/{cosmosWalletAddress}/cosmos-tx-hash", &UpdateLikeCoinMigrationCosmosHandler{
		Db:          h.Db,
		AsynqClient: h.AsynqClient,
	})
	router.Handle("DELETE /migration/{cosmosWalletAddress}", &CancelLikeCoinMigrationHandler{
		Db: h.Db,
	})

	return router
}
