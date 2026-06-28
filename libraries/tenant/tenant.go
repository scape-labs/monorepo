// Package tenant identifies the tenant that owns a request.
//
// Every inbound HTTP / AMQP / RPC request in the monorepo carries a tenant.
// Services must never operate across tenants without an explicit system-wide
// override.
//
// TODO: implement ResolveFromContext, ResolveFromHeader, and a Postgres
//-backed cache for tenant metadata.
package tenant

import "context"

// ID is a stable identifier for a tenant (UUID).
type ID string

// Tenant is the metadata we cache about a tenant.
type Tenant struct {
	ID      ID
	Name    string
	Plan    string // "free" | "growth" | "scale" | "enterprise"
	Active  bool
}

// ctxKey is unexported to prevent collisions with other packages.
type ctxKey struct{}

// With returns a new context that carries the given tenant.
//
// TODO: add a WithID helper for the common case where only the ID is known.
func With(ctx context.Context, t Tenant) context.Context {
	return context.WithValue(ctx, ctxKey{}, t)
}

// FromContext returns the tenant attached to ctx, or the zero value if none.
//
// TODO: distinguish "no tenant" from "anonymous tenant" via a sentinel.
func FromContext(ctx context.Context) Tenant {
	t, _ := ctx.Value(ctxKey{}).(Tenant)
	return t
}
