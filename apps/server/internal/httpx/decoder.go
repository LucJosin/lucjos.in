package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const unknownFieldPrefix = "json: unknown field "

// DecodeJSON decodes the JSON request body into v.
//
// The v param must be a non-nil pointer to a struct and unknown JSON fields are rejected.
func DecodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("empty body")
		}

		switch e := err.(type) {
		case *json.SyntaxError:
			return errors.New("invalid JSON syntax")
		case *json.UnmarshalTypeError:
			return fmt.Errorf("invalid JSON type for field %s", e.Field)
		}

		if after, ok := strings.CutPrefix(err.Error(), unknownFieldPrefix); ok {
			field := strings.Trim(after, `"`)
			return fmt.Errorf("unknown field %s", field)
		}
		return err
	}

	return nil
}
