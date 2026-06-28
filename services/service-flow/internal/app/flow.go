// Package app wires the flow service together.
//
// TODO: implement Setup() to register handlers, jobs, RPC clients, and the
// Postgres store against the kit server.
package app

import (
	"context"

	"github.com/scape-labs/platform/kit/server"
)

// Service is the composition root for flow.
//
// TODO: hold references to the Store, the AMQP publisher, the metrics
// recorder, etc.
type Service struct {
	// TODO: fields.
}

// New returns a configured Service.
//
// TODO: accept every dependency through this constructor; wire-DI will fill
// them in.
func New() *Service {
	return &Service{}
}

// Setup registers everything flow needs with the kit server.
//
// TODO: register HTTP handlers, cron jobs, RPC clients, and the leaderlock.
func (s *Service) Setup(ctx context.Context, srv *server.Server) error {
	_ = ctx
	_ = srv
	// TODO: implement.
	return nil
}
