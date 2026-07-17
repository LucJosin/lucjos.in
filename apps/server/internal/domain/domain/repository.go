package domain

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lucjosin/qorv.in/internal/errs"
)

type Repository interface {
	FindByDomain(ctx context.Context, domain string) (Domain, error)
	Create(ctx context.Context, entity Domain) (Domain, error)
}

func scanRow(row interface{ Scan(...any) error }) (Domain, error) {
	var e Domain
	err := row.Scan(
		&e.ID,
		&e.Domain,
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
