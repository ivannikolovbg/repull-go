// Quickstart for the Repull Go SDK.
//
// Set REPULL_API_KEY in your environment, then:
//
//	go run ./examples/quickstart
//
// Prints up to 10 reservations across all your connected PMS platforms.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ivannikolovbg/repull-go/repull"
)

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
	resp, err := client.ListReservationsWithResponse(ctx, &repull.ListReservationsParams{Limit: &limit})
	if err != nil {
		log.Fatalf("list reservations: %v", err)
	}
	if resp.HTTPResponse.StatusCode >= 300 {
		log.Fatal(repull.NewAPIError(resp.HTTPResponse.StatusCode, resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.Data == nil {
		fmt.Println("no reservations yet — connect a PMS at https://repull.dev/dashboard")
		return
	}

	data := *resp.JSON200.Data
	if len(data) == 0 {
		fmt.Println("no reservations yet — connect a PMS at https://repull.dev/dashboard")
		return
	}

	if p := resp.JSON200.Pagination; p != nil {
		total := -1
		if p.Total != nil {
			total = *p.Total
		}
		fmt.Printf("Total: %d   showing: %d   hasMore: %t\n", total, len(data), p.HasMore)
	}
	for _, r := range data {
		fmt.Printf("  %-12s  %s -> %s   %-12s  %s\n",
			r.Id,
			r.CheckIn.Format("2006-01-02"),
			r.CheckOut.Format("2006-01-02"),
			platStr(r.Platform),
			r.ConfirmationCode,
		)
	}
}

func platStr(p *repull.ReservationPlatform) string {
	if p == nil {
		return ""
	}
	return string(*p)
}
