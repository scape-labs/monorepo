package providers

import "context"

// VisaNet is the VisaNet rail adapter.
//
// TODO: implement.
type VisaNet struct {
	BaseURL  string
	Merchant string
}

// Charge submits a charge request to VisaNet.
//
// TODO: implement.
func (p *VisaNet) Charge(ctx context.Context, amount int64, card string) error {
	_ = ctx
	_ = amount
	_ = card
	// TODO: implement.
	return nil
}
