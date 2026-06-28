package service_test

import (
	"testing"

	"github.com/scape-labs/monorepo/services/service-dura/internal/service"
)

// TODO: replace with real tests once Service is implemented.
func TestNew_returnsNonNil(t *testing.T) {
	if got := service.New(nil); got == nil {
		t.Fatal("New(nil) returned nil")
	}
}
