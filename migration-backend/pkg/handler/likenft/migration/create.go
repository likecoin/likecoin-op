package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/likecoin/like-migration-backend/pkg/db"
	"github.com/likecoin/like-migration-backend/pkg/handler"
	api_model "github.com/likecoin/like-migration-backend/pkg/handler/model"
	likecoin_api "github.com/likecoin/like-migration-backend/pkg/likecoin/api"
	"github.com/likecoin/like-migration-backend/pkg/model"
	"github.com/likecoin/like-migration-backend/pkg/task"
)

type CreateMigrationRequestBody struct {
	AssetSnapshotId uint64 `json:"asset_snapshot_id"`
	CosmosAddress   string `json:"cosmos_address"`
	EthAddress      string `json:"eth_address"`
}

type CreateMigrationResponseBody struct {
	Migration *api_model.LikeNFTAssetMigration `json:"migration,omitempty"`
}

type CreateMigrationHandler struct {
	Db          *sql.DB
	LikecoinAPI *likecoin_api.LikecoinAPI
	AsynqClient *asynq.Client
}

var ErrMigrationExists = fmt.Errorf("error migration exists")
var ErrSignedEthAddressNotMatch = fmt.Errorf("error signed eth address not match")

func (h *CreateMigrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hub := sentry.GetHubFromContext(r.Context())

	decoder := json.NewDecoder(r.Body)
	var data CreateMigrationRequestBody
	err := decoder.Decode(&data)

	if err != nil {
		handler.SendJSON(w, http.StatusBadRequest, handler.MakeErrorResponseBody(err))
		return
	}

	migration, err := h.handle(data)

	if err != nil {
		handler.SendJSON(w, http.StatusInternalServerError,
			handler.MakeErrorResponseBody(err).
				WithSentryReported(hub.CaptureException(err)).
				AsError(handler.ErrSomethingWentWrong))
		return
	}

	handler.SendJSON(w, http.StatusCreated, &CreateMigrationResponseBody{
		Migration: migration,
	})

	go h.enqueueMigrationTasks(migration.Id)
}

func (h *CreateMigrationHandler) handle(req CreateMigrationRequestBody) (*api_model.LikeNFTAssetMigration, error) {
	userEVMMigrateResp, err := h.LikecoinAPI.GetUserEVMMigrate(req.CosmosAddress)

	if err != nil {
		return nil, err
	}

	if userEVMMigrateResp.EVMWallet == nil || !strings.EqualFold(*userEVMMigrateResp.EVMWallet, req.EthAddress) {
		return nil, ErrSignedEthAddressNotMatch
	}

	_, err = db.QueryLikeNFTAssetMigrationByCosmosAddress(h.Db, req.CosmosAddress)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			classes, err := db.QueryLikeNFTAssetSnapshotClassesByNFTSnapshotId(h.Db, req.AssetSnapshotId)

			if err != nil {
				return nil, err
			}

			nfts, err := db.QueryLikeNFTAssetSnapshotNFTsByNFTSnapshotId(h.Db, req.AssetSnapshotId)

			if err != nil {
				return nil, err
			}

			tx, err := h.Db.Begin()
			if err != nil {
				return nil, err
			}

			migration, err := db.InsertLikeNFTAssetMigration(tx, &model.LikeNFTAssetMigration{
				LikeNFTAssetSnapshotId: req.AssetSnapshotId,
				CosmosAddress:          req.CosmosAddress,
				EthAddress:             req.EthAddress,
				Status:                 model.NFTMigrationStatusInProgress,
			})
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			migrationClasses := make([]model.LikeNFTAssetMigrationClass, 0)

			for _, class := range classes {
				migrationClasses = append(migrationClasses, model.LikeNFTAssetMigrationClass{
					LikeNFTAssetMigrationId: migration.Id,
					CosmosClassId:           class.CosmosClassId,
					Name:                    class.Name,
					Image:                   class.Image,
					Status:                  model.LikeNFTAssetMigrationClassStatusInit,
				})
			}

			migrationNFTs := make([]model.LikeNFTAssetMigrationNFT, 0)

			for _, nft := range nfts {
				migrationNFTs = append(migrationNFTs, model.LikeNFTAssetMigrationNFT{
					LikeNFTAssetMigrationId: migration.Id,
					CosmosClassId:           nft.CosmosClassId,
					CosmosNFTId:             nft.CosmosNFTId,
					Name:                    nft.Name,
					Image:                   nft.Image,
					Status:                  model.LikeNFTAssetMigrationNFTStatusInit,
				})
			}

			err = db.InsertLikeNFTAssetMigrationClasses(tx, migrationClasses)

			if err != nil {
				tx.Rollback()
				return nil, err
			}

			err = db.InsertLikeNFTAssetMigrationNFTs(tx, migrationNFTs)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			err = tx.Commit()
			if err != nil {
				return nil, err
			}

			return api_model.LikeNFTAssetMigrationFromModel(migration, migrationClasses, migrationNFTs), nil
		} else {
			return nil, err
		}
	}

	return nil, ErrMigrationExists
}

func (h *CreateMigrationHandler) enqueueMigrationTasks(migrationId uint64) error {
	t, err := task.NewEnqueueLikeNFTAssetMigrationTask(migrationId)
	if err != nil {
		return err
	}
	_, err = h.AsynqClient.Enqueue(t, asynq.MaxRetry(0))
	if err != nil {
		return err
	}
	return nil
}
