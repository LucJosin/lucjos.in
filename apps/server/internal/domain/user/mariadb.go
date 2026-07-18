package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/database/tx"
)

const selectUsersStmt = `
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		u.email,
		u.password,
		u.avatar_url,
		u.email_verified,
		u.deleted_at,
		u.created_at,
		u.updated_at
	FROM 
		users u`

type MariaDBRepository struct {
	db *mariadb.Database
}

func NewMariaDBRepository(db *mariadb.Database) Repository {
	return &MariaDBRepository{
		db: db,
	}
}

func (r *MariaDBRepository) FindByID(ctx context.Context, id uuid.UUID) (User, error) {
	query := selectUsersStmt + ` WHERE u.id = ? AND u.deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, id)
	entity, err := scanRow(row)
	if err != nil {
		return User{}, fmt.Errorf("scanning user: %w", err)
	}

	return entity, nil
}

func (r *MariaDBRepository) FindByUsername(ctx context.Context, username string) (User, error) {
	query := selectUsersStmt + ` WHERE u.username = ? AND u.deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, username)
	entity, err := scanRow(row)
	if err != nil {
		return User{}, fmt.Errorf("scanning user: %w", err)
	}

	return entity, nil
}

func (r *MariaDBRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS (
		%s
		WHERE 
			u.id = ? AND u.deleted_at IS NULL
		LIMIT 1
	)`, selectUsersStmt)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fetching user: %w", err)
	}

	return exists, nil
}

func (r *MariaDBRepository) Create(ctx context.Context, entity User) (User, error) {
	err := tx.BeginFunc(ctx, r.db, func(tx *sql.Tx) error {
		now := time.Now()

		query := `
		INSERT INTO users (
			first_name,
			last_name,
			username,
			email,
			password,
			avatar_url,
			email_verified,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

		args := []any{
			entity.FirstName,
			entity.LastName,
			entity.Username,
			entity.Email,
			entity.Password,
			entity.AvatarURL,
			entity.EmailVerified,
			now,
			now,
		}
		row := tx.QueryRowContext(ctx, query, args...)

		err := row.Scan(&entity.ID)
		if err != nil {
			return mariadb.ParseError(err, "creating user")
		}

		entity.CreatedAt = now
		entity.UpdatedAt = now
		return nil
	})
	return entity, err
}
