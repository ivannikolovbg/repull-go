# Changelog

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
