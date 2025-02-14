package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hpetrov29/resttemplate/internal/validate"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Errors validate.FieldErrors `json:"errors"`
}

type PayloadResponse struct {
	Paylaod any `json:"payload"`
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, statusCode int, data any, args ...any) error {

	SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	var jsonData []byte
	var err error

	switch e := data.(type) {
	case validate.FieldErrors:
		jsonData, err = json.Marshal(ValidationErrorResponse{Errors: e})
	case error:
		jsonData, err = json.Marshal(ErrorResponse{Error: e.Error()})
	default:
		jsonData, err = json.Marshal(PayloadResponse{Paylaod: e})
	}
	
	if err != nil {
		w.Write([]byte(`{"error": "failed to encode error response"}`))
		return err
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
