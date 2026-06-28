package ussd

import (
	"context"
	"net/http"
)

// Handler is the HTTP entrypoint the USSD gateway POSTs to on every input.
//
// TODO: implement — load/create Session, resolve Menu via flow, render
// response, persist next Session state.
type Handler struct {
	store Store
}

// NewHandler constructs a Handler.
//
// TODO: validate the store is non-nil.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// ServeHTTP handles a single USSD callback.
//
// TODO: implement.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r
	// TODO: parse gateway callback, call into the menu graph, write response.
	w.WriteHeader(http.StatusOK)
}

// CallbackInput is the parsed USSD gateway callback.
//
// TODO: replace with the real schema once the gateway vendor is picked.
type CallbackInput struct {
	SessionID string
	Phone     string
	Text      string
}

// CallbackOutput is the response we send back to the gateway.
//
// TODO: replace with the real schema.
type CallbackOutput struct {
	Action string // "continue" | "end"
	Text   string
}

// Handle is the testable core of the handler (no *http.Request).
//
// TODO: implement — load Session, run the menu graph, persist next state.
func (h *Handler) Handle(ctx context.Context, in CallbackInput) (CallbackOutput, error) {
	_ = ctx
	return CallbackOutput{}, nil
}
