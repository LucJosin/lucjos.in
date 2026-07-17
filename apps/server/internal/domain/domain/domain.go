package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Domain struct {
	ID        uuid.UUID
	Domain    string
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
