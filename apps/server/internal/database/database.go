package database

import (
	"context"
	"database/sql"

	"github.com/lucjosin/qorv.in/internal/slogx"
)

// Close closes the [*sql.Rows], logging a warning if an error occurs.
func Close(ctx context.Context, rows *sql.Rows) {
	if rows != nil {
		if err := rows.Close(); err != nil {
			slogx.FromCtx(ctx).Error("closing rows", "error", rows.Close())
		}
	}
}
