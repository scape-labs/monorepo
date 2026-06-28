// Package effectors holds the per-request effectors that mutate side-effectful
// state once all deciders allow.
//
// TODO: implement each effector as an idempotent operation; retries must be
// safe.
package effectors

import "context"

// LedgerPost posts a balanced double-entry record to the ledger.
//
// TODO: implement — open a tx, insert two legs, commit.
func LedgerPost(ctx context.Context, debit, credit string, amount int64) error {
	_ = ctx
	_ = debit
	_ = credit
	_ = amount
	// TODO: implement.
	return nil
}
