package user

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	Username      string
	Email         string
	Password      string
	EmailVerified bool
	AvatarURL     *string
	DeletedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
