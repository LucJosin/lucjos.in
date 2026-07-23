package collections

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucjosin/qorv.in/internal/domain/collection"
	"github.com/lucjosin/qorv.in/internal/slogx"
)

type Handler struct {
	service collection.Service
}

func NewHandler(service collection.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/collections", func(r chi.Router) {
		r.Get("/", h.list)
	})
}

func (h *Handler) list(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := slogx.FromCtx(ctx)

	entities, err := h.service.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]CollectionResponse, len(entities))
	for i, entity := range entities {
		entity, err := h.toResponse(entity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response[i] = entity
	}

	// TODO: add specific httpx package to handle writing and errors
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("couldn't write response for collections", "error", err)
	}
}
