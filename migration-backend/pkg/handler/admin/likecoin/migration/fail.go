package migration

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/likecoin/like-migration-backend/pkg/handler"
	api_model "github.com/likecoin/like-migration-backend/pkg/handler/model"
	"github.com/likecoin/like-migration-backend/pkg/logic/likecoin"
)

const MinMigrationAgeForAdminCancel = 5 * time.Minute

type FailLikeCoinMigrationResponseBody struct {
	Migration *api_model.LikeCoinMigration `json:"migration,omitempty"`
}

type FailLikeCoinMigrationHandler struct {
	Db *sql.DB
}

func (h *FailLikeCoinMigrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hub := sentry.GetHubFromContext(r.Context())

	migrationIdStr := r.PathValue("migrationId")
	migrationId, err := strconv.ParseUint(migrationIdStr, 10, 64)
	if err != nil {
		handler.SendJSON(w, http.StatusBadRequest, handler.MakeErrorResponseBody(err))
		return
	}

	migration, err := h.handle(migrationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.SendJSON(w, http.StatusNotFound,
				handler.MakeErrorResponseBody(err).AsError(handler.ErrNotFound))
			return
		}
		if errors.Is(err, likecoin.ErrMigrationCannotBeFailed) {
			handler.SendJSON(w, http.StatusBadRequest,
				handler.MakeErrorResponseBody(err))
			return
		}
		if errors.Is(err, likecoin.ErrMigrationTooRecent) {
			handler.SendJSON(w, http.StatusBadRequest,
				handler.MakeErrorResponseBody(err))
			return
		}
		handler.SendJSON(w, http.StatusInternalServerError,
			handler.MakeErrorResponseBody(err).
				WithSentryReported(hub.CaptureException(err)).
				AsError(handler.ErrSomethingWentWrong))
		return
	}

	handler.SendJSON(w, http.StatusOK, &FailLikeCoinMigrationResponseBody{
		Migration: migration,
	})
}

func (h *FailLikeCoinMigrationHandler) handle(migrationId uint64) (*api_model.LikeCoinMigration, error) {
	m, err := likecoin.FailPendingMigrationById(
		h.Db,
		migrationId,
		"cancelled by admin",
		MinMigrationAgeForAdminCancel,
	)
	if err != nil {
		return nil, err
	}

	return api_model.LikeCoinMigrationFromModel(m), nil
}
