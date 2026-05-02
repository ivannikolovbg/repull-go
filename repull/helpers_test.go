package repull

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithBearer_AddsAuthorizationHeader(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := NewClient(srv.URL, WithBearer("sk_test_abc"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if _, err := c.GetHealth(context.Background()); err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if want := "Bearer sk_test_abc"; got != want {
		t.Fatalf("authorization header = %q, want %q", got, want)
	}
}

func TestNewAPIError_DecodesEnvelope(t *testing.T) {
	body := []byte(`{"error":{"code":"unauthorized","message":"bad key"}}`)
	e := NewAPIError(401, body)
	if e.StatusCode != 401 {
		t.Fatalf("StatusCode = %d, want 401", e.StatusCode)
	}
	if e.Detail == nil || e.Detail.Error == nil || e.Detail.Error.Message == nil || *e.Detail.Error.Message != "bad key" {
		t.Fatalf("envelope not decoded: %+v", e)
	}
	if got, want := e.Error(), "repull: 401 bad key"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestNewAPIError_FallsBackToBody(t *testing.T) {
	e := NewAPIError(500, []byte("internal boom"))
	if got, want := e.Error(), "repull: 500 internal boom"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}
