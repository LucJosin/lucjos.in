package collection

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Collection struct {
	ID          uuid.UUID
	PublicID    uuid.UUID
	WorkspaceID uuid.UUID
	Name        string
	Description *string
	Category    *string
	IsPublic    bool
	Color       *string
	Icon        *string
	PublicStats bool
	Pinned      bool
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
