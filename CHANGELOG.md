# Changelog

## v0.2.12 — 2026-09-15

### Additive

- **Regenerated against the live spec (174 → 175 operations).**
- `SetListingsStatus` / `SetListingsStatusWithResponse` — `POST /v1/listings/status`. Activate or deactivate up to 500 listings in one all-or-nothing call (`ListingStatusBatchRequest{ListingIds, Active}` → `ListingStatusBatchResponse{Active, Updated, Unchanged}`).
- `DeleteConnectionParams.AccountId` — optional `accountId` query param on `DELETE /v1/connect/{provider}`, required when a workspace has more than one account for the provider. The `200` response is now typed (`Disconnected`, `Provider`, `AccountId`, `ListingsDeactivated`); the account's listings are deactivated, not deleted.
- `ConnectStatus.Accounts` — every Airbnb account the workspace has connected.
- New `403 listing_inactive` error response (`ListingInactive`), declared on 83 operations.
- Airbnb calendar operations gain `BusySubtype`; `AirbnbPricingWriteRequest.ModelType` is now an enum.

### Changed

- Lists default to active listings: `GET /v1/listings` accepts `status=active|inactive|archived|all`, `GET /v1/properties` accepts `status=active|inactive|all`. Inactive rows carry identity fields only; reading or writing an inactive listing returns `403 listing_inactive`.
- Airbnb calendar writes (`PUT .../pricing`, `PUT .../availability`) validate more strictly (unknown fields such as `price` are refused with `422 invalid_params`) and declare new errors: `422 airbnb_rejected`, `403 connection_reauth_required`, `429 airbnb_rate_limited`.
- Sending `accessType` to `POST /v1/connect/airbnb` now locks the consent screen to that tier; omit it to let the host choose.
- **`DeleteConnection` / `DeleteConnectionWithResponse` take a `*DeleteConnectionParams` argument** after `provider` (pass `nil` to omit `accountId`). This is a signature change for existing callers.
- **Enum constants renamed by the generator.** A new enum with an overlapping value made oapi-codegen switch six constants to their type-prefixed names: `FullAccess`/`Messaging`/`ReadOnly` → `CreateConnectionJSONBodyAccessType{FullAccess,Messaging,ReadOnly}`, and `Availability`/`DerivedPricing`/`Rates` → `BookingAvailabilityUpdateRequestType{Availability,DerivedPricing,Rates}`. The old names still compile as deprecated aliases in `repull/compat.go`.

### Deprecated

- Booking.com webhooks endpoints (`GET`/`POST`/`DELETE /v1/channels/booking/webhooks`) are deprecated and always return `403`.
- The six short enum constant names listed above.

## v0.2.11 — 2026-09-11

### Fix

- **Regenerated against the live spec — 19 schema corrections, no path/operation change (still 124 paths / 174 operations).** Ten fields were `snake_case` in the SDK but the live API has always sent camelCase — the SDK now matches: `dataFreshness`, `lastSyncedAt`, `fixUrl`, `nextCursor`, `hasMore`, `monthlyRequests`, `dailyAiRequests`, `dailyAi`, `dynamicPricingListings`, `resetsAt`.
- **`BookingPropertyListResponse`, `BookingConversationListResponse`, `VrboListingListResponse` are now bare array type aliases** (`= []BookingProperty` etc.), matching what these three endpoints actually return. Previously generated as `{data, pagination}` wrapper structs, which never matched the wire response.
- **Four id fields are now `*string`, not `*int`:** `AirbnbAlteration.Id`, `AirbnbAlteration.ReservationId`, `AirbnbConnection.Id`, `AirbnbListing.ListingId`.
- **`Property.Latitude` / `Property.Longitude` are now `*string`** (decimal degrees as a string), not `*float64`.
- `scripts/check-spec-freshness.py` now diffs schema shapes, not just the path/operation inventory, so drift like this fails CI going forward instead of only catching added/removed endpoints.

### Note

- Jumping v0.2.8 → v0.2.11: `v0.2.9` and `v0.2.10` are already published git tags (cut 2026-07-30, before this file's `v0.2.7`/`v0.2.8` entries were written) pointing at unrelated commits — this changelog's numbering had drifted ahead of the actual tag history. Picking `v0.2.11` avoids colliding with an existing tag; it does not itself re-tag anything.

## v0.2.8 — 2026-09-11

### Additive

- **Regenerated against the live spec (170 → 174 operations).** Four write operations were added to the API as new methods on paths that already existed for `GET`, so the path count is unchanged at 124 and only the operation count moved:
  - `CreateGuest` / `CreateGuestWithResponse` — `POST /v1/guests`
  - `CreateReservation` / `CreateReservationWithResponse` — `POST /v1/reservations`
  - `UpdateReservation` / `UpdateReservationWithResponse` — `PATCH /v1/reservations/{id}`
  - `SendConversationMessage` / `SendConversationMessageWithResponse` — `POST /v1/conversations/{id}/messages`
- All three creating calls accept `IdempotencyKey` on their `*Params` struct. Send a unique string per distinct request: a repeat with the same key replays the stored response for 24 hours, a reuse with a changed payload returns `422 idempotency_key_reused`, and a reuse while the first request is in flight returns `409 idempotency_key_in_use`.
- No previously generated method was removed (229 → 237 `ClientWithResponses` methods).

## v0.2.7 — 2026-09-11

### Fix

- **Regenerated against the live spec (102 → 124 paths).** Removed the two dead `/v1/sandbox/*` paths (the sandbox was deleted from the API; `sk_test_*` keys now return 401). Added the 24 paths that were missing from the SDK: `PATCH /v1/availability/batch`, Airbnb alteration accept/decline, `GET /v1/channels/booking/properties/{id}/rooms`, credentials endpoints for the ten `/v1/connect/{provider}/credentials` PMS integrations, `POST /v1/connect/booking/callback`, `GET /v1/health/{atlas,auth,mcp,webhooks}` and `GET /v1/health/channels/{channel}`, listing photo upload endpoints, `GET /v1/quotes`, and `POST /v1/reviews/{id}/reply`.
- **`scripts/regen.sh` now points at the canonical spec URL** (`https://api.repull.dev/openapi.json`) instead of a second, previously-identical mirror path — this was a drift risk, not a behavior change. `REPULL_OPENAPI_URL` still overrides it.
- Dropped `sk_test_*` / sandbox mentions from doc comments, the test suite, the `connect_airbnb` example, and the README — the sandbox no longer exists; use `sk_live_*` everywhere.

## v0.2.6 — 2026-06-25

### Additive

- **Booking.com hosted Connect.** `POST /v1/connect/{provider}` (`CreateConnection`) now supports `provider = booking` for the hosted connect session, alongside the existing room-mapping helpers (`MapConnectBookingRooms`, `ListConnectBookingRooms`, plus the `/verify` and `/rooms` flows). `CreateConnectionJSONBody.RedirectUrl` is documented as Airbnb + Booking.com — where to redirect the user after the hosted connect flow completes.
- **`channel` filter on `GET /v1/properties`.** New `ListPropertiesParams.Channel` query param (`ListPropertiesParamsChannel`: `airbnb` / `booking` / `vrbo`) restricts results to properties with an active link on the given OTA. Omit to include every channel.
- **`channels` on `Property`.** New `Channels *[]string` field listing the OTAs each property is actively published on (e.g. `airbnb`, `booking`, `vrbo`); empty when the property has no active channel links.

### Note

- Regenerating against the new channel enum caused oapi-codegen to prefix the `ListReviewsParamsPlatform` constants. If you referenced the bare `Airbnb` / `Booking` / `Vrbo` constants for `ListReviews`, switch to `ListReviewsParamsPlatformAirbnb` / `...Booking` / `...Vrbo`. The enum values (`"airbnb"`, etc.) are unchanged.

## v0.2.5 — 2026-06-24

### Additive

- Add `messaging` Airbnb Connect access scope (read + send guest messages, no property management). Regenerated against the live spec; `CreateConnectionJSONBodyAccessType` now includes the `Messaging` value alongside `ReadOnly` and `FullAccess`. Unlike `full_access`, the `messaging` scope does not hold exclusive property management, so it can coexist with another app (e.g. an existing PMS) on the same Airbnb account.

## v0.2.4 — 2026-05-15

### Additive

- **`PaymentRequired` type added.** New 402 error envelope type alias (`type PaymentRequired = Error`) generated alongside the existing `BadRequest` / `Unauthorized` / `NotFound` / `TooManyRequests` aliases. Surfaces when the API returns `402 Payment Required` with `error.code = "listings_limit_exceeded"` for over-cap customers. Unlike 429, 402 is NOT a "wait and retry" condition — `Retry-After` is not set. Recovery: `DELETE` listings under the cap or upgrade at `repull.dev/dashboard/billing`. `/v1/health`, `/v1/usage/*`, and any `DELETE` are exempt. The 402 envelope mirrors `rate_limit_exceeded` and adds `tier`, `limit`, `active_listings`, `upgrade_url`. Tracks vanio-repull-api PR #66.

## v0.2.2 — 2026-05-07

### Additive

- **`?include=amenities` on `GET /v1/listings/{id}` and `GET /v1/properties/{id}`.** Both endpoints now expose an `Include` query param with the canonical `amenities` value. When passed, responses populate the new `Amenities *[]ListingAmenity` field on `Listing` and `Property` from the unified `listings_amenities` cache. Unknown include values return 422.
- **New `ListingAmenity` type.** Single amenity row from the unified `listings_amenities` table, with `amenityKey`, `isPresent`, and optional metadata. Surfaced on listing + property detail endpoints when `?include=amenities` is requested.

Refs: vanio-repull-api #59 (listings include), vanio-repull-api #61 (properties include).

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
