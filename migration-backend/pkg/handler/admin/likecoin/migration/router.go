package migration

import (
	"database/sql"
	"net/http"

	"github.com/hibiken/asynq"
)

type MigrationRouter struct {
	Db          *sql.DB
	AsynqClient *asynq.Client
}

func (h *MigrationRouter) Router() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /migration", &ListLikeCoinMigrationHandler{
		Db: h.Db,
	})

	router.Handle("GET /migration/{migrationId}", &GetLikeCoinMigrationHandler{
		Db: h.Db,
	})

	router.Handle("PUT /migration/{migrationId}", &RetryLikeCoinMigrationHandler{
		Db:          h.Db,
		AsynqClient: h.AsynqClient,
	})

	router.Handle("DELETE /migration/{migrationId}", &FailLikeCoinMigrationHandler{
		Db: h.Db,
	})

	return router
}
