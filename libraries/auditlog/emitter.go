// Package auditlog emits structured audit events for compliance.
//
// Events are published to the "audit" AMQP exchange and forwarded to the
// warehouse for long-term storage. The emitter is intentionally fire-and-forget
// at the call site — auditing must never block a request.
//
// TODO: implement Emit, BatchEmit, and a buffered async emitter backed by
// kit/messaging.
package auditlog

import (
	"context"
	"time"

	"github.com/scape-labs/monorepo/libraries/idgen"
)

// Event is one row in the audit log.
type Event struct {
	ID        idgen.ID    `json:"id"`
	Tenant    string      `json:"tenant"`
	Actor     string      `json:"actor"`     // user id or "system"
	Action    string      `json:"action"`    // e.g. "transfer.sent"
	Target    string      `json:"target"`    // e.g. "transfer:abc123"
	Metadata  Metadata    `json:"metadata"`
	Timestamp time.Time   `json:"timestamp"`
}

// Metadata is free-form structured detail.
type Metadata map[string]any

// Emitter publishes events.
//
// TODO: pick the real interface signature once kit/messaging is stable.
type Emitter interface {
	Emit(ctx context.Context, e Event) error
}
