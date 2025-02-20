package accountdb

import (
	"embed"

	"github.com/vanti-dev/sc-bos/internal/database"
)

// Migrations contains the SQL migrations for the accountdb schema.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

var Schema = database.MustLoadVersionedSchema(migrationsFS, "migrations")
