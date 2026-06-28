// Package main is the entrypoint for service-flow.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting service.flow")
	// TODO: kit.New("service.flow").Register(...).Run()
	_ = server.Component{}
}
