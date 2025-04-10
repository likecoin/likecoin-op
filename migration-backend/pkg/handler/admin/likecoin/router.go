package likecoin

import (
	"database/sql"
	"net/http"

	"github.com/likecoin/like-migration-backend/pkg/handler/admin/likecoin/migration"
)

type LikeCoinRouter struct {
	Db          *sql.DB
}

func (h *LikeCoinRouter) Router() *http.ServeMux {
	router := http.NewServeMux()

	migrationRouter := migration.MigrationRouter{
		Db:          h.Db,
	}

	// FIXME: Find a way to handle CRUD paths
	// This is for paths without trailing /. e.g. GET / POST
	router.Handle("/migration", migrationRouter.Router())
	// This is for paths with trailing/intermediate /, e.g. GET / PUT
	router.Handle("/migration/", migrationRouter.Router())

	return router
}
