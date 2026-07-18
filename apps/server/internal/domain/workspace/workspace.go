package workspace

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Workspace struct {
	ID          uuid.UUID
	PublicID    uuid.UUID
	DomainID    uuid.UUID
	Slug        string
	Name        string
	Description string
	Color       string
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (w *Workspace) ValidateAndPrepare() error {
	w.Prepare()
	return w.Validate()
}

func (w *Workspace) Validate() error {
	return nil
}

func (w *Workspace) Prepare() {
}
