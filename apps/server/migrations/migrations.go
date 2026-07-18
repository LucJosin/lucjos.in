package migrations

import "embed"

// FS contains all SQL files used in mariadb migration.
//
//go:embed *.sql
var FS embed.FS
