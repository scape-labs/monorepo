// Package money provides a currency-aware value type shared across services.
//
// TODO: implement Amount arithmetic with overflow checks, MarshalJSON for the
// API layer, and a SQL scanner/valuer for the repository layer.
package money

// Currency is an ISO-4217 currency code (e.g. "USD", "ZAR", "KES").
// TODO: validate against a known list at construction time.
type Currency string

// Amount is a value in the currency's minor unit (cents, ngwee, etc.).
// Stored as int64 so we never round through float.
//
// TODO: enforce non-negative at construction (or provide Signed/Unsigned variants).
type Amount struct {
	Currency Currency
	Minor    int64
}
