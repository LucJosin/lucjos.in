package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucjosin/qorv.in/internal/slogx"
)

type Handler struct {
}

// NewHandler registers the api HTTP handlers.
func NewHandler(r chi.Router) {
	h := &Handler{}
	r.Get("/ping", h.ping)
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		slogx.FromCtx(r.Context()).Error("writing ping response", "error", err)
	}
}
