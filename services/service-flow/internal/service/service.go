// Package service contains the main business logic for flow.
//
// TODO: implement. Service is the orchestrator that handlers call into.
// All persistence goes through Store; all side effects go through the
// Publisher / RPC clients injected by the app layer.
package service

// Service is the main business-logic struct.
//
// TODO: hold dependencies (Store, Publisher, RPC clients) and the per-request
// invariants.
type Service struct {
	store Store
	// TODO: other dependencies.
}

// New constructs a Service.
//
// TODO: validate inputs, return *Service.
func New(store Store) *Service {
	return &Service{store: store}
}
