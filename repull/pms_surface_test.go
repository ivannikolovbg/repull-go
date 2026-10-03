package repull

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGuest_PatchesGuestAndDecodesPmsOutcome(t *testing.T) {
	var gotPath, gotKey string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.Method + " " + r.URL.Path
		gotKey = r.Header.Get("Idempotency-Key")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":42,"firstName":"Ana","lastName":"Lopez","contacts":[{"type":"email","value":"ana@example.com","isPrimary":true}],"pms":[{"provider":"guesty","applied":["guest"]}]}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL, WithBearer("sk_test_x"))
	if err != nil {
		t.Fatal(err)
	}
	last := "Lopez"
	key := IdempotencyKey("idem-1")
	res, err := c.UpdateGuestWithResponse(context.Background(), 42, &UpdateGuestParams{IdempotencyKey: &key}, UpdateGuestJSONRequestBody{LastName: &last})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "PATCH /v1/guests/42" || gotKey != "idem-1" {
		t.Fatalf("path = %q key = %q", gotPath, gotKey)
	}
	if gotBody["lastName"] != "Lopez" || len(gotBody) != 1 {
		t.Fatalf("body = %v", gotBody)
	}
	if res.JSON200 == nil || res.JSON200.Pms == nil || len(*res.JSON200.Pms) != 1 || *(*res.JSON200.Pms)[0].Provider != "guesty" {
		t.Fatalf("decoded = %+v", res.JSON200)
	}
}

func TestPmsCapabilities_DecodesFromListing(t *testing.T) {
	raw := []byte(`{"provider":"hostaway","connected":true,"reviews":{"read":true,"reply":false},"guests":{"create":false,"update":false},"conversations":{"send":true,"attachments":false,"channelSelect":true},"calendar":{"write":true}}`)
	var caps PmsCapabilities
	if err := json.Unmarshal(raw, &caps); err != nil {
		t.Fatal(err)
	}
	if caps.Reviews == nil || *caps.Reviews.Reply || !*caps.Conversations.ChannelSelect || *caps.Conversations.Attachments {
		t.Fatalf("decoded = %+v", caps)
	}
}
