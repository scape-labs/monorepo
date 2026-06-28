// Package ussd holds the text-based USSD protocol adapter for flow.
//
// USSD sessions are short-lived (typically < 30s) and stateless on the
// gateway side, so we own all session state in Postgres keyed by the
// gateway's session id.
package ussd

import (
	"context"
	"time"

	"github.com/scape-labs/monorepo/shared/idgen"
)

// Session is one in-progress USSD session.
//
// TODO: add Menu, Breadcrumbs, LastInput, ExpiresAt fields.
type Session struct {
	ID        idgen.ID
	GatewayID string
	Phone     string
	CreatedAt time.Time
}

// Store is the data-access boundary for sessions.
//
// TODO: implement on top of Postgres.
type Store interface {
	Create(ctx context.Context, s Session) error
	Get(ctx context.Context, gatewayID string) (Session, error)
	Delete(ctx context.Context, gatewayID string) error
}
