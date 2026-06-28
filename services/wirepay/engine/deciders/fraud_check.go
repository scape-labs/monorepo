package deciders

// FraudCheckInput is the input to FraudCheck.
//
// TODO: fill out with velocity, geo, device-fingerprint, etc.
type FraudCheckInput struct {
	AccountID string
	Amount    int64
}

// FraudCheckOutput is the decision.
//
// TODO: replace bool with a typed enum (Allow / Review / Block).
type FraudCheckOutput struct {
	Allow bool
}

// FraudCheck runs the fraud rules against the request.
//
// TODO: implement — call into the fraud/rules package, aggregate verdicts.
func FraudCheck(in FraudCheckInput) FraudCheckOutput {
	_ = in
	return FraudCheckOutput{}
}
