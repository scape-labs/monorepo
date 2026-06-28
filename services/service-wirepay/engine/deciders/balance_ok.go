// Package deciders holds the per-request deciders that decide whether a given
// effect may run.
//
// TODO: implement each decider as a pure function: (Request, State) -> Decision.
package deciders

// BalanceOKInput is the input to BalanceOK.
//
// TODO: fill out with the real fields once the ledger schema is stable.
type BalanceOKInput struct {
	AccountID string
	Amount    int64 // in minor units
}

// BalanceOKOutput is the decision.
//
// TODO: replace bool with a typed enum (OK / Insufficient / Unknown).
type BalanceOKOutput struct {
	OK bool
}

// BalanceOK checks whether the account has enough balance to cover Amount.
//
// TODO: implement — query the ledger, return OK if balance >= amount.
func BalanceOK(in BalanceOKInput) BalanceOKOutput {
	_ = in
	return BalanceOKOutput{}
}
