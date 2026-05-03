// Package repull is the official Go SDK for Repull (https://api.repull.dev).
//
// The unified API for vacation-rental tech: 50+ PMS platforms, 4 OTA channels,
// AI guest communication, pricing, and listing optimization — through one REST API.
//
// The bulk of this package (types.gen.go, client.gen.go) is generated from the
// public OpenAPI spec at https://api.repull.dev/openapi.json by oapi-codegen.
// This file holds the small hand-written wrappers that make the SDK ergonomic:
// the bearer-token auth helper, a typed APIError, and a default base URL.
//
// Quick start:
//
//	client, err := repull.NewClientWithResponses(
//	    repull.DefaultBaseURL,
//	    repull.WithBearer(os.Getenv("REPULL_API_KEY")),
//	)
//	if err != nil { return err }
//
//	resp, err := client.ListReservationsWithResponse(ctx, &repull.ListReservationsParams{})
//	if err != nil { return err }
//	if resp.JSON200 == nil { return repull.NewAPIError(resp.HTTPResponse.StatusCode, resp.Body) }
//
//	for _, r := range *resp.JSON200.Data {
//	    fmt.Println(r.Id, r.ConfirmationCode)
//	}
package repull

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DefaultBaseURL is the production Repull API base URL.
const DefaultBaseURL = "https://api.repull.dev"

// WithBearer returns a ClientOption that adds an `Authorization: Bearer <token>`
// header to every request. Pass your Repull API key (sk_test_* for sandbox,
// sk_live_* for production).
//
// oapi-codegen does not generate auth handling out of the box — this is the
// canonical way to wire it in.
func WithBearer(token string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	})
}

// APIError is the typed error returned when an API call returns a non-2xx
// status. The raw response body is preserved on Body, and decoded into Detail
// when the server returns the standard `{ error, message, code }` envelope.
type APIError struct {
	StatusCode int
	Body       []byte
	Detail     *Error
}

// NewAPIError constructs an APIError from a status code and raw body, decoding
// the body into the Error envelope when possible.
func NewAPIError(statusCode int, body []byte) *APIError {
	e := &APIError{StatusCode: statusCode, Body: body}
	if len(body) > 0 {
		var detail Error
		if err := json.Unmarshal(body, &detail); err == nil {
			e.Detail = &detail
		}
	}
	return e
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Detail != nil && e.Detail.Error.Message != "" {
		return fmt.Sprintf("repull: %d %s", e.StatusCode, e.Detail.Error.Message)
	}
	if len(e.Body) > 0 {
		return fmt.Sprintf("repull: %d %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("repull: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}
