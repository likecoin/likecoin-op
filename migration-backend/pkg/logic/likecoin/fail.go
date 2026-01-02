package likecoin

import (
	"database/sql"
	"fmt"
	"time"

	appdb "github.com/likecoin/like-migration-backend/pkg/db"
	"github.com/likecoin/like-migration-backend/pkg/model"
)

var ErrMigrationCannotBeFailed = fmt.Errorf("migration cannot be failed")
var ErrMigrationTooRecent = fmt.Errorf("migration is too recent to be cancelled")

// FailPendingMigration fails a migration that is still in pending_cosmos_tx_hash state
// minAge parameter enforces a minimum age before allowing cancellation (0 = no restriction)
func FailPendingMigration(
	db *sql.DB,
	migration *model.LikeCoinMigration,
	reason string,
	minAge time.Duration,
) (*model.LikeCoinMigration, error) {
	// Only allow failing migrations that haven't submitted cosmos tx yet
	if migration.Status != model.LikeCoinMigrationStatusPendingCosmosTxHash {
		return nil, ErrMigrationCannotBeFailed
	}

	// Enforce minimum age if specified (safety check for admin operations)
	if minAge > 0 {
		age := time.Since(migration.CreatedAt)
		if age < minAge {
			return nil, ErrMigrationTooRecent
		}
	}

	migration.Status = model.LikeCoinMigrationStatusFailed
	migration.FailedReason = &reason

	err := appdb.UpdateLikeCoinMigration(db, migration)
	if err != nil {
		return nil, err
	}

	return migration, nil
}

// FailPendingMigrationByCosmosAddress fails a user's latest pending migration
func FailPendingMigrationByCosmosAddress(
	db *sql.DB,
	cosmosAddress string,
	reason string,
) (*model.LikeCoinMigration, error) {
	m, err := appdb.QueryLatestLikeCoinMigration(db, cosmosAddress)
	if err != nil {
		return nil, err
	}

	return FailPendingMigration(db, m, reason, 0) // No age restriction for user
}

// FailPendingMigrationById fails a migration by ID (admin operation)
func FailPendingMigrationById(
	db *sql.DB,
	migrationId uint64,
	reason string,
	minAge time.Duration,
) (*model.LikeCoinMigration, error) {
	m, err := appdb.QueryLikeCoinMigrationById(db, migrationId)
	if err != nil {
		return nil, err
	}

	return FailPendingMigration(db, m, reason, minAge)
}
