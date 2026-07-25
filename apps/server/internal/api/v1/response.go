package v1

import (
	"context"
	"net/http"

	"github.com/lucjosin/qorv.in/internal/httpx"
)

// APIResponse represents a standardized API response format,
// containing a data payload of any type and additional metadata.
type APIResponse[T any] struct {
	Data  T              `json:"data,omitzero"`
	Error *ErrorResponse `json:"error,omitzero"`
}

// internalError sends an [http.StatusInternalServerError] response and logs the error.
//
// No body is sent to the client.
func internalError(ctx context.Context, w http.ResponseWriter, err error) {
	httpx.WriteErrorf(ctx, w, err, http.StatusInternalServerError, nil)
}

// writeResponse sends a JSON response with the specified status code and body.
//
// It wraps the body in an APIResponse struct and includes pagination metadata if provided.
func writeResponse[T any](ctx context.Context, w http.ResponseWriter, status int, body T) {
	httpx.WriteJSON(ctx, w, status, APIResponse[T]{
		Data: body,
	})
}
