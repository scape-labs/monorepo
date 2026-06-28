package money_test

import (
	"testing"

	"github.com/scape-labs/monorepo/shared/money"
)

// TODO: replace with real tests once money.Add / money.Sub are implemented.
func TestAmount_zero_value(t *testing.T) {
	var a money.Amount
	if a.Currency != "" {
		t.Fatalf("zero-value Currency should be empty, got %q", a.Currency)
	}
	if a.Minor != 0 {
		t.Fatalf("zero-value Minor should be 0, got %d", a.Minor)
	}
}
