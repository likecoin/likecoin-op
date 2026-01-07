package migration

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/likecoin/like-migration-backend/pkg/handler"
	api_model "github.com/likecoin/like-migration-backend/pkg/handler/model"
	"github.com/likecoin/like-migration-backend/pkg/logic/likecoin"
	"github.com/likecoin/like-migration-backend/pkg/model"
)

type CancelLikeCoinMigrationHandler struct {
	Db *sql.DB
}

func (h *CancelLikeCoinMigrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hub := sentry.GetHubFromContext(r.Context())

	cosmosAddress := r.PathValue("cosmosWalletAddress")

	m, err := h.handle(cosmosAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.SendJSON(w, http.StatusNotFound,
				handler.MakeErrorResponseBody(err).
					AsError(handler.ErrNotFound))
			return
		}
		handler.SendJSON(w, http.StatusInternalServerError,
			handler.MakeErrorResponseBody(err).
				WithSentryReported(hub.CaptureException(err)).
				AsError(handler.ErrSomethingWentWrong))
		return
	}

	handler.SendJSON(w, http.StatusOK, api_model.LikeCoinMigrationFromModel(m))
}

func (h *CancelLikeCoinMigrationHandler) handle(cosmosAddress string) (*model.LikeCoinMigration, error) {
	return likecoin.FailPendingMigrationByCosmosAddress(
		h.Db,
		cosmosAddress,
		"cancelled by user",
	)
}
