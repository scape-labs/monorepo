// Package main is the entrypoint for service-dura.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting service.dura")
	// TODO: kit.New("service.dura").Register(...).Run()
	_ = server.Component{}
}
