package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/domain/collection"
)

type CollectionResponse struct {
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

type CollectionHandler struct {
	service collection.Service
}

// NewCollectionHandler registers the collection HTTP handlers.
func NewCollectionHandler(r chi.Router, service collection.Service) {
	h := &CollectionHandler{service: service}
	r.Route("/collections", func(r chi.Router) {
		r.Get("/", h.list)
	})
}

func (h *CollectionHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entities, err := h.service.List(ctx)
	if err != nil {
		internalError(ctx, w, err)
		return
	}

	response := make([]CollectionResponse, len(entities))
	for _, entity := range entities {
		entity, err := h.toResponse(entity)
		if err != nil {
			internalError(ctx, w, err)
			return
		}
		response = append(response, entity)
	}

	// TODO: add paginated response
	writeResponse(ctx, w, http.StatusOK, response)
}

func (h *CollectionHandler) toResponse(e collection.Collection) (CollectionResponse, error) {
	if e.PublicID == uuid.Nil {
		return CollectionResponse{}, errors.New("workspace has empty public ID")
	}

	res := CollectionResponse{
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
