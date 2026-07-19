package collection

import (
	"context"
	"fmt"

	"github.com/lucjosin/qorv.in/internal/database"
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
)

const selectCollectionsStmt = `
	SELECT
		c.id,
		c.public_id,
		c.workspace_id,
		c.name,
		c.description,
		c.category,
		c.is_public,
		c.color,
		c.icon,
		c.public_stats,
		c.pinned,
		c.deleted_at,
		c.created_at,
		c.updated_at
	FROM 
		collections c`

type MariaDBRepository struct {
	db *mariadb.Database
}

func NewMariaDBRepository(db *mariadb.Database) Repository {
	return &MariaDBRepository{
		db: db,
	}
}

func (r *MariaDBRepository) List(ctx context.Context) ([]Collection, error) {
	query := selectCollectionsStmt + ` WHERE c.deleted_at IS NULL ORDER BY c.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer database.CloseRows(ctx, rows)

	var entities []Collection
	for rows.Next() {
		entity, err := scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning collection: %w", err)
		}
		entities = append(entities, entity)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return entities, nil
}
