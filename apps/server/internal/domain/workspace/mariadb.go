package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/database/tx"
)

const selectWorkspacesStmt = `
	SELECT
		w.id,
		w.public_id,
		w.name,
		w.description,
		w.color,
		w.deleted_at,
		w.created_at,
		w.updated_at
	FROM 
		workspaces w`

type MariaDBRepository struct {
	db *mariadb.Database
}

func NewMariaDBRepository(db *mariadb.Database) Repository {
	return &MariaDBRepository{
		db: db,
	}
}

func (r *MariaDBRepository) FindByID(ctx context.Context, id uuid.UUID) (Workspace, error) {
	query := selectWorkspacesStmt + ` WHERE w.id = ? AND w.deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, id)
	entity, err := scanRow(row)
	if err != nil {
		return Workspace{}, fmt.Errorf("scanning workspace: %w", err)
	}

	return entity, nil
}

func (r *MariaDBRepository) FindByDomainID(ctx context.Context, domainID uuid.UUID) (Workspace, error) {
	query := selectWorkspacesStmt + ` WHERE w.domain_id = ? AND w.deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, domainID)
	entity, err := scanRow(row)
	if err != nil {
		return Workspace{}, fmt.Errorf("scanning workspace: %w", err)
	}

	return entity, nil
}

func (r *MariaDBRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS (
		%s
		WHERE 
			w.id = ? AND w.deleted_at IS NULL
		LIMIT 1
	)`, selectWorkspacesStmt)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fetching user: %w", err)
	}

	return exists, nil
}

func (r *MariaDBRepository) Create(ctx context.Context, entity Workspace) (Workspace, error) {
	err := tx.BeginFunc(ctx, r.db, func(tx *sql.Tx) error {
		now := time.Now()

		query := `
		INSERT INTO workspaces (
			domain_id,
			name,
			slug,
			description
		) VALUES (?, ?, ?, ?) RETURNING id, public_id`

		args := []any{
			entity.DomainID,
			entity.Name,
			entity.Slug,
			entity.Description,
		}
		row := tx.QueryRowContext(ctx, query, args...)

		err := row.Scan(&entity.ID, &entity.PublicID)
		if err != nil {
			return mariadb.ParseError(err, "creating workspace")
		}

		entity.CreatedAt = now
		entity.UpdatedAt = now
		return nil
	})
	return entity, err
}

func (r *MariaDBRepository) CreateWorkspaceUser(ctx context.Context, entity WorkspaceUser) error {
	return tx.BeginFunc(ctx, r.db, func(tx *sql.Tx) error {
		now := time.Now()

		query := `
		INSERT INTO workspace_users (
			workspace_id,
			user_id,
			invited_by,
			role,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)`

		args := []any{
			entity.WorkspaceID,
			entity.UserID,
			entity.InvitedBy,
			entity.Role,
			now,
			now,
		}
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return mariadb.ParseError(err, "couldn't add user to workspace")
		}

		entity.CreatedAt = now
		entity.UpdatedAt = now
		return nil
	})
}
