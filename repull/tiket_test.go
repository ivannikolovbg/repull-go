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

func TestTiket_ListReservationsDecodesRows(t *testing.T) {
	var gotURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.Method + " " + r.URL.Path + "?" + r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"reservationId":"42","tiketReservationId":"T1","confirmationCode":"T1","checkIn":"2026-12-01","paymentMethod":"Bank Transfer","status":"confirmed","adults":2}],"pagination":{"limit":50,"offset":0,"hasMore":false}}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL, WithBearer("sk_test_x"))
	if err != nil {
		t.Fatal(err)
	}
	pid := "123"
	st := ListTiketReservationsParamsStatus("confirmed")
	res, err := c.ListTiketReservationsWithResponse(context.Background(), &ListTiketReservationsParams{PropertyId: &pid, Status: &st})
	if err != nil {
		t.Fatal(err)
	}
	if gotURL != "GET /v1/channels/tiket/reservations?property_id=123&status=confirmed" {
		t.Fatalf("url = %q", gotURL)
	}
	if res.JSON200 == nil || len(res.JSON200.Data) != 1 {
		t.Fatalf("decoded = %+v", res.JSON200)
	}
	row := res.JSON200.Data[0]
	if *row.ReservationId != "42" || *row.PaymentMethod != "Bank Transfer" || row.CheckIn.Time.Day() != 1 {
		t.Fatalf("row = %+v", row)
	}
}

func TestTiket_UpdateAvailabilitySendsRatesWrite(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/v1/channels/tiket/availability" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"propertyId":"123","type":"rates","updates":1,"requests":1,"echoTokens":["e1"],"warnings":[],"applied":"verified"}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL, WithBearer("sk_test_x"))
	if err != nil {
		t.Fatal(err)
	}
	room, plan, price := "R1", "P1", float64(500000)
	day := openapi_types.Date{Time: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)}
	res, err := c.UpdateTiketAvailabilityWithResponse(context.Background(), UpdateTiketAvailabilityJSONRequestBody{
		PropertyId: "123",
		Type:       TiketAvailabilityUpdateRequestType("rates"),
		Updates:    []TiketAriUpdate{{RoomTypeCode: &room, RatePlanCode: &plan, Start: day, End: day, Price: &price}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotBody["property_id"] != "123" || gotBody["type"] != "rates" {
		t.Fatalf("body = %v", gotBody)
	}
	u := gotBody["updates"].([]any)[0].(map[string]any)
	if u["start"] != "2026-12-01" || u["price"] != float64(500000) {
		t.Fatalf("update = %v", u)
	}
	if res.JSON200 == nil || string(res.JSON200.Applied) != "verified" {
		t.Fatalf("decoded = %+v", res.JSON200)
	}
}
