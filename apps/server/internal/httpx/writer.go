package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lucjosin/qorv.in/internal/slogx"
)

// WriteError writes an HTTP error response without a custom message.
func WriteError(ctx context.Context, w http.ResponseWriter, err error, statusCode int) {
	WriteErrorf(ctx, w, err, statusCode, nil)
}

// WriteErrorf writes an HTTP error response with a custom message.
func WriteErrorf(ctx context.Context, w http.ResponseWriter, err error, statusCode int, body any) {
	log := slogx.FromCtx(ctx)

	if statusCode == http.StatusInternalServerError {
		log.Error("server error", "error", err, "status", statusCode, "message", body)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	write(ctx, w, statusCode, body)
}

// WriteJSON sends a JSON response (JSONUTF8ContentType) with specified 'status' and 'body' (if not null)
func WriteJSON(ctx context.Context, w http.ResponseWriter, status int, body any) {
	write(ctx, w, status, body)
}

func write(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", JSONUTF8ContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slogx.FromCtx(ctx).Error("couldn't write error", "error", err, "status", statusCode)
	}
}
