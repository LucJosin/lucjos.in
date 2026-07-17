package system

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/errs"
)

type Repository interface {
	FindOne(ctx context.Context) (System, error)
	Create(ctx context.Context, ownerUserID uuid.UUID) (System, error)
}

func scanRow(row interface{ Scan(...any) error }) (System, error) {
	var e System
	err := row.Scan(
		&e.OwnerUserID,
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
