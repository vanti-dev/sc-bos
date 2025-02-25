package account

import (
	"embed"

	"github.com/vanti-dev/sc-bos/internal/database"
)

// Migrations contains the SQL migrations for the queries schema.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

var schema = database.MustLoadVersionedSchema(migrationsFS, "migrations")
