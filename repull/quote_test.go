package repull

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func TestQuoteReservation_PostsStayAndDecodesAnswer(t *testing.T) {
	var gotPath string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.Method + " " + r.URL.Path
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"listingId":"4118","provider":"hostaway","available":true,"total":880,"currency":"USD","restrictions":[]}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL, WithBearer("sk_test_x"))
	if err != nil {
		t.Fatal(err)
	}
	adults := 2
	res, err := c.QuoteReservationWithResponse(context.Background(), nil, QuoteReservationJSONRequestBody{
		ListingId: 4118,
		CheckIn:   openapi_types.Date{Time: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		CheckOut:  openapi_types.Date{Time: time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)},
		Adults:    &adults,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "POST /v1/reservations/quote" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotBody["checkIn"] != "2026-10-01" || gotBody["adults"] != float64(2) {
		t.Fatalf("body = %v", gotBody)
	}
	if res.JSON200 == nil || res.JSON200.Available == nil || !*res.JSON200.Available || *res.JSON200.Total != 880 || *res.JSON200.ListingId != "4118" {
		t.Fatalf("decoded = %+v", res.JSON200)
	}
}

func TestCreateReservation_SendsPmsFields(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"91","listingId":"4118"}`))
	}))
	defer srv.Close()

	c, _ := NewClientWithResponses(srv.URL)
	total := float32(880)
	notes, unit := "Late arrival", "u-1"
	send := false
	status := ReservationCreateRequestStatus("tentative")
	first := "Ada"
	res, err := c.CreateReservationWithResponse(context.Background(), nil, CreateReservationJSONRequestBody{
		ListingId:             4118,
		CheckIn:               openapi_types.Date{Time: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		CheckOut:              openapi_types.Date{Time: time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)},
		Guest:                 ReservationGuestInput{FirstName: first},
		TotalPrice:            &total,
		Notes:                 &notes,
		UnitId:                &unit,
		SendConfirmationEmail: &send,
		Status:                &status,
	})
	if err != nil {
		t.Fatal(err)
	}
	for k, want := range map[string]any{"totalPrice": float64(880), "notes": "Late arrival", "unitId": "u-1", "sendConfirmationEmail": false, "status": "tentative"} {
		if gotBody[k] != want {
			t.Fatalf("%s = %v, want %v", k, gotBody[k], want)
		}
	}
	if res.JSON201 == nil || *res.JSON201.Id != "91" {
		t.Fatalf("decoded = %+v", res.JSON201)
	}
}
