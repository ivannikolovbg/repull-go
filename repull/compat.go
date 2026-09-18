package repull

// Backward-compatible aliases for enum constants that oapi-codegen renamed.
//
// oapi-codegen emits short constant names (FullAccess, Rates, ...) only while
// they are unique across the whole spec. When a later spec adds another enum
// with an overlapping value, it switches every colliding constant to the
// type-prefixed form. These aliases keep code written against the short names
// compiling. Hand-written: regen.sh does not touch this file. If a future regen
// restores the short names, `go build` fails on a duplicate declaration --
// delete the matching alias here.

// Deprecated: use CreateConnectionJSONBodyAccessTypeFullAccess.
const FullAccess = CreateConnectionJSONBodyAccessTypeFullAccess

// Deprecated: use CreateConnectionJSONBodyAccessTypeMessaging.
const Messaging = CreateConnectionJSONBodyAccessTypeMessaging

// Deprecated: use CreateConnectionJSONBodyAccessTypeReadOnly.
const ReadOnly = CreateConnectionJSONBodyAccessTypeReadOnly

// Deprecated: use BookingAvailabilityUpdateRequestTypeAvailability.
const Availability = BookingAvailabilityUpdateRequestTypeAvailability

// Deprecated: use BookingAvailabilityUpdateRequestTypeDerivedPricing.
const DerivedPricing = BookingAvailabilityUpdateRequestTypeDerivedPricing

// Deprecated: use BookingAvailabilityUpdateRequestTypeRates.
const Rates = BookingAvailabilityUpdateRequestTypeRates

// Backward-compatible alias for a request body the spec moved from an inline
// schema to a named component (2026-09-18 regen, 175 -> 191 operations). The
// generator dropped the old struct name in favor of the new component name;
// code referencing the old type name still compiles against the same shape.
//
// Deprecated: use AirbnbAlterationCreateRequest.
type CreateAirbnbAlterationJSONBody = AirbnbAlterationCreateRequest
