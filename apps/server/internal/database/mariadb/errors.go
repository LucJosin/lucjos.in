package mariadb

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lucjosin/qorv.in/internal/errs"
)

const (
	// UniqueViolation represents a duplicate entry for UNIQUE or PRIMARY KEY
	UniqueViolation uint16 = 1062

	// ForeignKeyViolation represents a foreign key constraint violation when trying to add/update a child row
	ForeignKeyViolation uint16 = 1452

	// ForeignKeyReferenceViolation represents a foreign key constraint violation when trying to delete/update a parent row
	ForeignKeyReferenceViolation uint16 = 1451

	// NotNullViolation represents a NOT NULL constraint violation
	NotNullViolation uint16 = 1048

	// CheckViolation represents a CHECK constraint violation
	CheckViolation uint16 = 3819
)

// ParseError converts a MySQL error into a more specific application error based on the error code.
func ParseError(err error, fallback string) error {
	if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok {
		switch mysqlErr.Number {
		case UniqueViolation:
			return fmt.Errorf("%w: record already exists", errs.ErrConflict)
		case ForeignKeyViolation, ForeignKeyReferenceViolation:
			return fmt.Errorf("%w: invalid reference", errs.ErrInvalid)
		case NotNullViolation:
			return fmt.Errorf("%w: required field missing", errs.ErrInvalid)
		case CheckViolation:
			return fmt.Errorf("%w: check failed", errs.ErrInvalid)
		default:
			return fmt.Errorf("%s: %w", fallback, err)
		}
	}
	return fmt.Errorf("%s: %w", fallback, err)
}
