package openapi

import (
	"log"
	"net/http"

	"likecollective-indexer/ent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/openapi/api"
)

type openAPIHandler struct {
	db *ent.Client

	stakingRepository      database.StakingRepository
	accountRepository      database.AccountRepository
	nftClassRepository     database.NFTClassRepository
	stakingEventRepository database.StakingEventRepository
	evmEventRepository     database.EVMEventRepository

	bookNFTDeltaTimeBucketRepository database.BookNFTDeltaTimeBucketRepository
}

var _ api.Handler = &openAPIHandler{}

func NewOpenAPIHandler(
	db database.Service,
	timescaleDbService database.TimescaleService,
) http.Handler {
	handler := &openAPIHandler{
		db: db.Client(),

		stakingRepository:      database.MakeStakingRepository(db),
		accountRepository:      database.MakeAccountRepository(db),
		nftClassRepository:     database.MakeNFTClassRepository(db),
		stakingEventRepository: database.MakeStakingEventRepository(db),
		evmEventRepository:     database.MakeEVMEventRepository(db),

		bookNFTDeltaTimeBucketRepository: database.MakeBookNFTDeltaTimeBucketRepository(timescaleDbService),
	}

	srv, err := api.NewServer(
		handler,
	)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}
