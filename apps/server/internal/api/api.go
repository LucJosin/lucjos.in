package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes register all v1 HTTP API routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", h.ping)
}

func (h *Handler) ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		slog.Error("writing ping response", "error", err)
	}
}
