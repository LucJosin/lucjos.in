package collection

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lucjosin/qorv.in/internal/errs"
)

type Repository interface {
	List(ctx context.Context) ([]Collection, error)
}

func scanRow(row interface{ Scan(...any) error }) (Collection, error) {
	var e Collection
	err := row.Scan(
		&e.ID,
		&e.PublicID,
		&e.WorkspaceID,
		&e.Name,
		&e.Description,
		&e.Category,
		&e.IsPublic,
		&e.Color,
		&e.Icon,
		&e.PublicStats,
		&e.Pinned,
		&e.DeletedAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, errs.ErrNotFound
		}
		return e, err
	}
	return e, nil
}
