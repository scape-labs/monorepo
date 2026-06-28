// Package apperror defines the typed errors flow-ussd returns to handlers and RPC
// clients. Every error must carry an HTTP / AMQP status code and a stable
// error code so callers can pattern-match without parsing strings.
//
// TODO: implement Code constants and the Error type, plus helper constructors.
package apperror

// Code is a stable, machine-readable error identifier.
//
// TODO: define the full set for flow-ussd.
type Code string

// TODO: define codes — e.g. CodeNotFound, CodeConflict, CodeUnauthorized.
const (
	CodeUnknown Code = "unknown"
)

// Error is the typed error flow-ussd returns.
//
// TODO: implement.
type Error struct {
	Code    Code
	Message string
	Cause   error
}
