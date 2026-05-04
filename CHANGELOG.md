# Changelog

## v0.2.1 — 2026-05-04

### Added — Studio routes (16 ops)

Repull Studio is now reachable from the Go SDK. The generated `repull` package
gains 16 new client methods spanning 10 paths under `/api/studio/*`:

- `ListStudioProjects`, `CreateStudioProject`, `GetStudioProject`,
  `UpdateStudioProject`, `DeleteStudioProject`
- `ListStudioProjectFiles`, `UpsertStudioProjectFile`,
  `DeleteStudioProjectFile`
- `CreateStudioProjectGeneration`, `GenerateStudioCompletion`
- `ListStudioDeployments`, `CreateStudioDeployment`,
  `GetStudioDeployment`, `DeleteStudioDeployment`,
  `SuspendStudioDeployment`, `WakeStudioDeployment`

New types: `StudioProject`, `StudioFile`, `StudioGeneration`,
`StudioDeployment`, `StudioError`.

## v0.2.0 — 2026-05-02

**MAJOR — coordinated breaking-change release across the Repull SDK fleet (TS, Python, PHP, Ruby, .NET, Go).** The Repull API converged on a single canonical envelope shape, camelCase field names, and string-typed IDs. This version of the Go SDK regenerates against that converged spec.

### Breaking

- **Canonical pagination envelope.** All list responses now look like `{ data: [...], pagination: { nextCursor, hasMore, total? } }`. `nextCursor` and `hasMore` are required; `total` is omitted when `?include_total=false`. Code that walked old per-endpoint shapes (e.g. `MarketsListResponse.Markets`, `MarketsListResponse.TotalInFilter`) must be updated to read `Data` / `Pagination.Total`.
- **camelCase across the board.** All response field names are camelCase (`checkIn`, `checkOut`, `confirmationCode`, `guestId`, `listingId`, `createdAt`, `nextCursor`, `hasMore`, `totalPrice`, `guestDetails`). Generated Go struct field names are `CamelCase` (oapi-codegen convention) with `json:"camelCase"` tags. Snake-case JSON tags are gone.
- **All IDs are strings.** `Reservation.Id`, `Reservation.GuestId`, `Reservation.ListingId`, `Review.Id`, etc. are now `string` (not `int`/`*int`). Any consumer that did integer arithmetic, formatting with `%d`, or numeric comparison on IDs must switch to string handling.
- **`POST /v1/connect/{provider}` response field rename.** `oauthUrl` → `url`. The body now contains `{ url, sessionId, ... }`. Code reading the OAuth consent URL field by name must update.
- **`GET /v1/markets` envelope change.** `markets` → `data`; `total_in_filter` → `pagination.total`. The list shape now matches every other list endpoint.
- **`GET /v1/reviews/{id}` returns a bare `Review` object.** Previously wrapped — now it is just the `Review` schema directly. Consumers that did `resp.JSON200.Review` (or similar wrapper field) should switch to using the `Review` fields directly off `JSON200`.
- **`/v1/channels/airbnb/*` list endpoints adopt the canonical envelope.** Listings, reservations, threads, messages, photos, and reviews now all return `{ data, pagination }` instead of bespoke per-endpoint top-level keys.
- **`Error` envelope shape.** The generated `Error` type's inner `Error` struct is no longer a pointer, and `Message` / `Code` / `Fix` / `DocsUrl` / `RequestId` are now non-pointer `string`. Existing call-sites that did `*err.Detail.Error.Message` or `err.Detail.Error != nil` must drop the deref / nil check. `helpers.go` and `helpers_test.go` updated accordingly.
- **`Reservation.Id`, `ConfirmationCode`, `CheckIn`, `CheckOut`, `Currency`, `GuestId`, `ListingId`, `Status`, `CreatedAt`, `GuestDetails`, `TotalPrice`** are non-pointer (already moved in v0.1.2; reaffirmed and now string-typed where applicable).

### Additive

- **Full error envelope.** `Error.Error` now ships `code`, `message`, `fix`, `docs_url`, `request_id`, plus optional `did_you_mean`, `endpoint`, `field`, `retry_after`, `support`, `valid_params`, `valid_values`, `value_received` for high-quality machine-and-human-friendly errors.
- **Rate-limit headers (documented in spec).** Responses carry standard `x-ratelimit-*` and `Retry-After` headers; surface them via the `*HTTPResponse.Header` getters on any `*ClientResponse`.
- **`X-Schema` header on read endpoints** (carried over from v0.1.2). All list/get reads still accept `XSchema *XSchemaHeader` to remap the response into your own field names. Built-ins: `native`, `calry`, `calry-v1`. Custom schemas are managed via the `*CustomSchema` CRUD ops.
- **Custom schema CRUD** (carried over from v0.1.2): `CreateCustomSchema`, `ListCustomSchemas`, `GetCustomSchema`, `UpdateCustomSchema`, `DeleteCustomSchema` on `/v1/schema/custom`.
- **Connect detail / sessions / providers.** `ListConnectProviders`, `SelectConnectProvider`, `VerifyBookingHotel`, `ListConnectBookingRooms`, `MapConnectBookingRooms`. Booking.com hotel verification + room-mapping flow exposed.
- **Markets browse + calendar.** `ListMarketBrowse`, `GetMarket`, `GetMarketCalendar` — paginated browse, per-city detail, daily occupancy/ADR calendar.
- **Airbnb channel detail endpoints.** Per-listing pricing / availability GET+PUT, photo list / upload, threaded messaging (list threads, list messages, send), reservation detail + actions, review respond, sync trigger.
- **Spec source path moved** to `https://api.repull.dev/api/repull/openapi.json`. `scripts/regen.sh` updated; the legacy `https://api.repull.dev/openapi.json` URL is still served (alias) but no longer canonical.

### Migration

```diff
- resp, _ := client.GetV1ReservationsWithResponse(ctx, &repull.GetV1ReservationsParams{Limit: &limit})
+ resp, _ := client.ListReservationsWithResponse(ctx, &repull.ListReservationsParams{Limit: &limit})

- fmt.Printf("%d  %s\n", *r.Id, *r.ConfirmationCode)
+ fmt.Printf("%s  %s\n", r.Id, r.ConfirmationCode)

- if total := resp.JSON200.Pagination.Total; total > 0 { ... }       // now *int
+ if p := resp.JSON200.Pagination; p != nil && p.Total != nil && *p.Total > 0 { ... }

- if e.Detail != nil && e.Detail.Error != nil && e.Detail.Error.Message != nil {
-     msg := *e.Detail.Error.Message
+ if e.Detail != nil && e.Detail.Error.Message != "" {
+     msg := e.Detail.Error.Message
```

## v0.1.2 — 2026-05-02

- **Custom schemas (additive).** New CRUD ops on `/v1/schema/custom`: `CreateCustomSchema`, `ListCustomSchemas`, `GetCustomSchema`, `UpdateCustomSchema`, `DeleteCustomSchema`. Lets you create workspace-scoped field-mapping schemas that reshape any read response into your own field names.
- **`X-Schema` header on read endpoints.** All list/get reads gain an optional `XSchema *XSchemaHeader` param. Built-in values: `native` (default), `calry`, `calry-v1`. Custom schemas created via the new CRUD endpoints can be passed by name. Unknown / inactive names fall back to `native`.
- **8 new types:** `CustomSchema`, `CustomSchemaCreate`, `CustomSchemaCreateResponse`, `CustomSchemaDeleteResponse`, `CustomSchemaListResponse`, `CustomSchemaMappings`, `CustomSchemaSummary`, `CustomSchemaUpdate`.
- **BREAKING — Reservation shape drift fix.** Several `Reservation` fields are now non-pointer, matching what the live API actually returns: `Id`, `ConfirmationCode`, `CheckIn`, `CheckOut`, `Currency`, `GuestId`, `ListingId`, `Status`, `CreatedAt`, `GuestDetails`. Code that did `*r.Id` / `*r.ConfirmationCode` must drop the deref. `examples/quickstart` updated accordingly.

## v0.1.1 — 2026-05-01

- Added conversations, guests, and reviews endpoints.
- Cursor-paginated reservations (`?cursor=`).

## v0.1.0 — 2026-05-01

- Initial release. Bootstrapped Go SDK for Repull from `api.repull.dev/openapi.json`.
