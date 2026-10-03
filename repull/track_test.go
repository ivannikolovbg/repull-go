package repull

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitTrackCredentials_PostsCredentialsAndDecodesAnswer(t *testing.T) {
	var gotPath string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.Method + " " + r.URL.Path
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"provider":"track","connected":true,"pmsConnectionId":"pc_1","created":true,"accountInfo":{"domain":"acme.trackhs.com","keyType":"server","authMode":"hmac"},"firstSync":{"queued":true}}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL, WithBearer("sk_test_x"))
	if err != nil {
		t.Fatal(err)
	}
	var body SubmitTrackCredentialsJSONRequestBody
	body.Credentials.Domain = "acme.trackhs.com"
	body.Credentials.ApiKey = "trk_live_x"
	body.Credentials.ApiSecret = "c2VjcmV0"
	keyType := SubmitTrackCredentialsJSONBodyCredentialsKeyTypeServer
	body.Credentials.KeyType = &keyType
	moveReason := 3
	body.Credentials.MoveReasonId = &moveReason

	res, err := c.SubmitTrackCredentialsWithResponse(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "POST /v1/connect/track/credentials" {
		t.Fatalf("path = %q", gotPath)
	}
	creds, _ := gotBody["credentials"].(map[string]any)
	if creds["domain"] != "acme.trackhs.com" || creds["apiKey"] != "trk_live_x" || creds["keyType"] != "server" || creds["moveReasonId"] != float64(3) {
		t.Fatalf("body = %v", gotBody)
	}
	if res.JSON200 == nil || res.JSON200.AccountInfo == nil || *res.JSON200.AccountInfo.Domain != "acme.trackhs.com" {
		t.Fatalf("decoded = %+v", res.JSON200)
	}
}
