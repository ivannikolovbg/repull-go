// Quickstart for the Repull Go SDK.
//
// Set REPULL_API_KEY in your environment, then:
//
//	go run ./examples/quickstart
//
// Prints up to 10 reservations across all your connected PMS platforms.
//
// NOTE: As of v0.1.0-alpha the live api.repull.dev response shape drifts from
// the published OpenAPI spec for a few fields (e.g. Reservation.id arrives as
// a string, totalPrice as a string). The example falls back to a raw-JSON
// parse when the strict type binding fails so the demo still works end-to-end.
// Tracked upstream — once the spec is corrected, regenerate and the typed
// access path will light up.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ivannikolovbg/repull-go/repull"
)

type rawReservation struct {
	ID               any    `json:"id"`
	CheckIn          string `json:"checkIn"`
	CheckOut         string `json:"checkOut"`
	Platform         string `json:"platform"`
	Status           string `json:"status"`
	ConfirmationCode string `json:"confirmationCode"`
}

type rawList struct {
	Data       []rawReservation `json:"data"`
	Pagination struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"pagination"`
}

func main() {
	apiKey := os.Getenv("REPULL_API_KEY")
	if apiKey == "" {
		log.Fatal("set REPULL_API_KEY in your environment (get one at https://repull.dev/dashboard)")
	}

	client, err := repull.NewClientWithResponses(repull.DefaultBaseURL, repull.WithBearer(apiKey))
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
	limit := 10
	resp, err := client.GetV1ReservationsWithResponse(ctx, &repull.GetV1ReservationsParams{Limit: &limit})
	if err != nil {
		// Type binding failed — fall back to a manual call so we can still demo
		// the HTTP path against the live API while the spec is being corrected.
		raw, rawErr := fetchRaw(ctx, apiKey, limit)
		if rawErr != nil {
			log.Fatalf("list reservations: %v (raw fallback: %v)", err, rawErr)
		}
		printList(raw)
		return
	}
	if resp.HTTPResponse.StatusCode >= 300 {
		log.Fatal(repull.NewAPIError(resp.HTTPResponse.StatusCode, resp.Body))
	}

	if resp.JSON200 != nil && resp.JSON200.Data != nil {
		// Spec-true path (currently unreachable until upstream fixes the drift).
		for _, r := range *resp.JSON200.Data {
			fmt.Printf("%-6d  %s  %s\n", deref(r.Id), platStr(r.Platform), confStr(r.ConfirmationCode))
		}
		return
	}

	// Loose parse from the raw body that the SDK already fetched.
	var list rawList
	if err := json.Unmarshal(resp.Body, &list); err != nil {
		log.Fatalf("decode body: %v\nbody: %s", err, string(resp.Body))
	}
	printList(&list)
}

func fetchRaw(ctx context.Context, apiKey string, limit int) (*rawList, error) {
	c, err := repull.NewClient(repull.DefaultBaseURL, repull.WithBearer(apiKey))
	if err != nil {
		return nil, err
	}
	httpResp, err := c.GetV1Reservations(ctx, &repull.GetV1ReservationsParams{Limit: &limit})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	var list rawList
	if err := json.NewDecoder(httpResp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return &list, nil
}

func printList(list *rawList) {
	if len(list.Data) == 0 {
		fmt.Println("no reservations yet — connect a PMS at https://repull.dev/dashboard")
		return
	}
	fmt.Printf("Total: %d   showing: %d\n", list.Pagination.Total, len(list.Data))
	for _, r := range list.Data {
		fmt.Printf("  %-10v  %s -> %s   %-12s  %s\n", r.ID, r.CheckIn, r.CheckOut, r.Platform, r.ConfirmationCode)
	}
}

func deref(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func platStr(p *repull.ReservationPlatform) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

func confStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
