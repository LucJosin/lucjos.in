package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/errs"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, user User) (User, error)
}

func scanRow(row interface{ Scan(...any) error }) (User, error) {
	var e User
	err := row.Scan(
		&e.ID,
		&e.FirstName,
		&e.LastName,
		&e.Username,
		&e.Email,
		&e.Password,
		&e.AvatarURL,
		&e.EmailVerified,
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
