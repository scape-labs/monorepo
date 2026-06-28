// Package rules holds the individual fraud rules that make up the reactive
// control network. Each rule is a pure function from (Request, History) -> Verdict.
package rules

// AmountAnomalyRequest is the input to AmountAnomaly.
//
// TODO: fill out with sender history (median, p99, etc.).
type AmountAnomalyRequest struct {
	AccountID string
	Amount    int64
	Median    int64
}

// AmountAnomalyVerdict is the rule's output.
//
// TODO: replace bool with a typed enum (Allow / Review / Block).
type AmountAnomalyVerdict struct {
	Block bool
}

// AmountAnomaly blocks transfers whose amount is > 10x the account's median.
//
// TODO: implement — replace the placeholder constant with a learned threshold.
func AmountAnomaly(req AmountAnomalyRequest) AmountAnomalyVerdict {
	_ = req
	return AmountAnomalyVerdict{}
}
