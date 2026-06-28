// Package providers wraps the third-party payment rails (EcoCash, M-Pesa,
// VisaNet, …) behind a uniform Provider interface so the effectors don't
// have to special-case them.
package providers

import "context"

// EcoCash is the EcoCash rail adapter.
//
// TODO: implement.
type EcoCash struct {
	BaseURL string
	APIKey  string
}

// Charge submits a charge request to EcoCash.
//
// TODO: implement.
func (p *EcoCash) Charge(ctx context.Context, amount int64, phone string) error {
	_ = ctx
	_ = amount
	_ = phone
	// TODO: implement.
	return nil
}
