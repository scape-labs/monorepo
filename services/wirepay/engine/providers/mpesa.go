package providers

import "context"

// MPesa is the M-Pesa rail adapter.
//
// TODO: implement.
type MPesa struct {
	BaseURL    string
	ConsumerKey string
}

// Charge submits a charge request to M-Pesa.
//
// TODO: implement.
func (p *MPesa) Charge(ctx context.Context, amount int64, phone string) error {
	_ = ctx
	_ = amount
	_ = phone
	// TODO: implement.
	return nil
}
