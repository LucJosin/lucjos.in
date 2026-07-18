package workspace

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/errs"
)

type Repository interface {
	// workspace

	FindByDomainID(ctx context.Context, domainID uuid.UUID) (Workspace, error)
	ExistsByID(ctx context.Context, ID uuid.UUID) (bool, error)
	Create(ctx context.Context, entity Workspace) (Workspace, error)

	// workspace user

	CreateWorkspaceUser(ctx context.Context, workspaceUser WorkspaceUser) error
}

func scanRow(row interface{ Scan(...any) error }) (Workspace, error) {
	var e Workspace
	err := row.Scan(
		&e.ID,
		&e.PublicID,
		&e.Name,
		&e.Description,
		&e.Color,
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
