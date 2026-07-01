package app

import (
	"net/http"

	"github.com/scape-labs/platform/kit/server"
)

// Routes groups every HTTP route emicro exposes.
//
// TODO: list every route — public + authenticated + admin — and the handler
// that serves it.
type Routes struct {
	// TODO: handler dependencies.
}

// Register installs the routes onto the kit router.
//
// TODO: implement.
func (r *Routes) Register(router server.Router) error {
	_ = router
	// TODO: register routes.
	return nil
}

// Healthz is the liveness probe.
//
// TODO: replace with a real check (DB ping, AMQP ping, etc.).
func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Readyz is the readiness probe.
//
// TODO: replace with a real check.
func Readyz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}
