package account

import (
	"embed"

	"github.com/vanti-dev/sc-bos/internal/sqlite"
)

// Migrations contains the SQL migrations for the queries schema.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

var schema = sqlite.MustLoadVersionedSchema(migrationsFS, "migrations")
