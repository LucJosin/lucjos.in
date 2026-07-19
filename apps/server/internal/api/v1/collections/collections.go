package collections

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/domain/collection"
)

type Collection struct {
	ID          uuid.UUID  `json:"id"`
	PublicID    uuid.UUID  `json:"public_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitzero"`
	Category    *string    `json:"category,omitzero"`
	IsPublic    bool       `json:"is_public"`
	Color       *string    `json:"color,omitzero"`
	Icon        *string    `json:"icon,omitzero"`
	PublicStats bool       `json:"public_stats"`
	Pinned      bool       `json:"pinned"`
	DeletedAt   *time.Time `json:"deleted_at,omitzero"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (h *Handler) toResponse(e collection.Collection) (Collection, error) {
	if e.PublicID == uuid.Nil {
		return Collection{}, errors.New("workspace has empty public ID")
	}

	res := Collection{
		ID:          e.PublicID,
		WorkspaceID: e.WorkspaceID,
		Name:        e.Name,
		Category:    e.Category,
		IsPublic:    e.IsPublic,
		Color:       e.Color,
		Icon:        e.Icon,
		PublicStats: e.PublicStats,
		Pinned:      e.Pinned,
		Description: e.Description,
		DeletedAt:   e.DeletedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
	return res, nil
}
