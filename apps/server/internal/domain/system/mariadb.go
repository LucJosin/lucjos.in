package system

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/database/tx"
)

type MariaDBRepository struct {
	db *mariadb.Database
}

func NewMariaDBRepository(db *mariadb.Database) Repository {
	return &MariaDBRepository{
		db: db,
	}
}

func (r *MariaDBRepository) FindOne(ctx context.Context) (System, error) {
	query := `
	SELECT
		s.owner_user_id,
		s.created_at,
		s.updated_at
	FROM 
		system s`

	row := r.db.QueryRowContext(ctx, query)
	entity, err := scanRow(row)
	if err != nil {
		return System{}, fmt.Errorf("scanning system: %w", err)
	}

	return entity, nil
}

func (r *MariaDBRepository) Create(ctx context.Context, ownerUserID uuid.UUID) (System, error) {
	entity := System{}
	err := tx.BeginFunc(ctx, r.db, func(tx *sql.Tx) error {
		now := time.Now()

		query := `
		INSERT INTO system (
			owner_user_id,
			created_at,
			updated_at
		) VALUES (?, ?, ?)`

		args := []any{
			ownerUserID,
			now,
			now,
		}
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return mariadb.ParseError(err, "saving system configuration")
		}

		entity.OwnerUserID = ownerUserID
		entity.CreatedAt = now
		entity.UpdatedAt = now
		return nil
	})
	return entity, err
}
