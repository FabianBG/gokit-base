package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

var (
	// ErrBadRouting Error on the routing request
	ErrBadRouting = func(e error) error {
		return fmt.Errorf("bad routing: %s", e.Error())
	}
	// ErrBadRequest Error on the request
	ErrBadRequest = func(e error) error {
		return fmt.Errorf("the request was malformed: %s", e.Error())
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		encodeError(ctx, e.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	if strings.HasPrefix(err.Error(), "the request was malformed:") {
		return http.StatusBadRequest
	}
	if strings.HasPrefix(err.Error(), "bad routing:") {
		return http.StatusBadRequest
	}
	if strings.HasPrefix(err.Error(), "auth:") {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
