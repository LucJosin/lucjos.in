package system

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type System struct {
	OwnerUserID uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
