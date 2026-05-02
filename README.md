# repull-go

> Go SDK for Repull. Generated from OpenAPI. Context-first.

The official Go SDK for [api.repull.dev](https://api.repull.dev) — the unified
API for vacation-rental tech (50+ PMS platforms, Airbnb / Booking.com / VRBO /
Plumguide channels, AI ops, white-label OAuth).

> **Status:** v0.1.2 — alpha. `pkg.go.dev` listing pending (auto-publishes on
> first import).

## Install

```bash
go get github.com/ivannikolovbg/repull-go/repull
```

Requires Go 1.24+.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/ivannikolovbg/repull-go/repull"
)

func main() {
    client, err := repull.NewClientWithResponses(
        repull.DefaultBaseURL,
        repull.WithBearer(os.Getenv("REPULL_API_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    limit := 10
    resp, err := client.GetV1ReservationsWithResponse(ctx, &repull.GetV1ReservationsParams{Limit: &limit})
    if err != nil {
        log.Fatal(err)
    }
    if resp.HTTPResponse.StatusCode >= 300 {
        log.Fatal(repull.NewAPIError(resp.HTTPResponse.StatusCode, resp.Body))
    }

    for _, r := range *resp.JSON200.Data {
        fmt.Printf("%d  %s -> %s\n", *r.Id, r.CheckIn.Format("2006-01-02"), r.CheckOut.Format("2006-01-02"))
    }
}
```

Run the bundled quickstart:

```bash
REPULL_API_KEY=sk_live_... go run github.com/ivannikolovbg/repull-go/examples/quickstart@latest
```

## Authentication

All Repull API calls require a bearer token. Get an API key at
[repull.dev/dashboard](https://repull.dev/dashboard) — `sk_test_*` for sandbox,
`sk_live_*` for production.

```go
client, _ := repull.NewClientWithResponses(
    repull.DefaultBaseURL,
    repull.WithBearer(os.Getenv("REPULL_API_KEY")),
)
```

`WithBearer` is a thin `RequestEditorFn` wrapper that attaches
`Authorization: Bearer <token>` to every request. You can stack additional
editors (logging, tracing, retries) by passing more `repull.ClientOption`s to
`NewClientWithResponses`.

## Examples

| Path | What it shows |
|---|---|
| `examples/quickstart` | List your reservations across all connected PMS platforms. |
| `examples/connect_airbnb` | Mint a white-label Airbnb OAuth Connect session and poll for completion. |

```bash
REPULL_API_KEY=sk_test_... go run ./examples/quickstart
REPULL_API_KEY=sk_test_... go run ./examples/connect_airbnb
```

## Layout

```
repull-go/
  repull/
    client.gen.go     generated client (one method per OpenAPI op + WithResponse variants)
    types.gen.go      generated request/response types
    helpers.go        WithBearer auth helper, APIError, DefaultBaseURL — hand-written
    helpers_test.go   tiny smoke tests for the helpers
  examples/
    quickstart/       go run-able demo against api.repull.dev
    connect_airbnb/   white-label OAuth Connect example
  openapi/
    v1.json           snapshotted OpenAPI source
  scripts/
    regen.sh          re-fetch spec + regenerate the *.gen.go files
```

## Regenerating

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
./scripts/regen.sh
```

This re-snapshots `openapi/v1.json` from `https://api.repull.dev/openapi.json`,
regenerates `repull/types.gen.go` and `repull/client.gen.go`, and runs
`go mod tidy` + `go build` + `go vet`. Commit the snapshot and the regenerated
files together so the SDK and its source of truth stay in lockstep.

## Reference

Full API reference and guides: [repull.dev/docs](https://repull.dev/docs).

## License

MIT — see [LICENSE](./LICENSE).

## Custom schemas

`X-Schema` lets you reshape any read response into your own field names. Built-in
schemas: `native` (default), `calry`, `calry-v1`. Manage workspace-scoped custom
schemas via `CreateCustomSchema` / `ListCustomSchemas` / `GetCustomSchema` /
`UpdateCustomSchema` / `DeleteCustomSchema`, then pass the name on any read:

```go
schema := repull.XSchemaHeader("calry")
resp, err := client.ListReservationsWithResponse(ctx, &repull.ListReservationsParams{
    XSchema: &schema,
    Limit:   &limit,
})
```

## Status

v0.1.2 — alpha. The API surface tracks `https://api.repull.dev/openapi.json`
1:1 and may break before 1.0. Open an issue if you hit drift between the
generated client and the live API — both are still settling.

---

Powered by Repull. AI features powered by Vanio AI.
