// Connect Airbnb example for the Repull Go SDK.
//
// Mints a white-label OAuth Connect session for Airbnb, prints the consent
// URL the user should visit, then polls connection status until it flips to
// connected (or you Ctrl-C).
//
//	REPULL_API_KEY=sk_test_... go run ./examples/connect_airbnb
//
// The redirect URL below is for demo purposes — register your real one with
// Repull before going to production.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ivannikolovbg/repull-go/repull"
)

func main() {
	apiKey := os.Getenv("REPULL_API_KEY")
	if apiKey == "" {
		log.Fatal("set REPULL_API_KEY in your environment")
	}

	client, err := repull.NewClientWithResponses(repull.DefaultBaseURL, repull.WithBearer(apiKey))
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
	provider := repull.Provider("airbnb")
	access := repull.PostV1ConnectProviderJSONBodyAccessType("full_access")
	redirect := "https://example.com/airbnb/return"

	mint, err := client.PostV1ConnectProviderWithResponse(ctx, provider, repull.PostV1ConnectProviderJSONRequestBody{
		AccessType:  &access,
		RedirectUrl: &redirect,
	})
	if err != nil {
		log.Fatalf("mint connect session: %v", err)
	}
	if mint.HTTPResponse.StatusCode >= 300 {
		log.Fatal(repull.NewAPIError(mint.HTTPResponse.StatusCode, mint.Body))
	}

	fmt.Println("Connect session minted. Open this URL in a browser, finish the consent flow:")
	fmt.Println()
	fmt.Printf("  status: %d\n", mint.HTTPResponse.StatusCode)
	fmt.Printf("  body:   %s\n", string(mint.Body))
	fmt.Println()
	fmt.Println("Polling connection status every 5s. Ctrl-C to exit.")

	for {
		st, err := client.GetV1ConnectProviderWithResponse(ctx, provider)
		if err != nil {
			log.Printf("status check: %v", err)
		} else {
			fmt.Printf("[%s] HTTP %d  body=%s\n", time.Now().Format(time.Kitchen), st.HTTPResponse.StatusCode, string(st.Body))
		}
		time.Sleep(5 * time.Second)
	}
}
