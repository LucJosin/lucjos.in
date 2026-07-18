package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/database/tx"
)

const selectDomainsStmt = `
	SELECT
		d.id,
		d.domain,
		d.deleted_at,
		d.created_at,
		d.updated_at
	FROM 
		domains d`

type MariaDBRepository struct {
	db *mariadb.Database
}

func NewMariaDBRepository(db *mariadb.Database) Repository {
	return &MariaDBRepository{db: db}
}

func (r *MariaDBRepository) FindByDomain(ctx context.Context, domain string) (Domain, error) {
	query := selectDomainsStmt + ` WHERE d.domain = ? AND d.deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, domain)
	entity, err := scanRow(row)
	if err != nil {
		return Domain{}, err
	}

	return entity, nil
}

func (r *MariaDBRepository) Create(ctx context.Context, entity Domain) (Domain, error) {
	err := tx.BeginFunc(ctx, r.db, func(tx *sql.Tx) error {
		now := time.Now()

		query := `
		INSERT INTO domains (
			domain,
			created_at,
			updated_at
		) VALUES (?, ?, ?) RETURNING id`

		args := []any{
			entity.Domain,
			now,
			now,
		}
		row := tx.QueryRowContext(ctx, query, args...)

		err := row.Scan(&entity.ID)
		if err != nil {
			return mariadb.ParseError(err, "creating domain")
		}

		entity.CreatedAt = now
		entity.UpdatedAt = now
		return nil
	})
	return entity, err
}
