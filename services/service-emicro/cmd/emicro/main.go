// Package main is the entrypoint for service-emicro.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting service.emicro")
	// TODO: kit.New("service.emicro").Register(...).Run()
	_ = server.Component{}
}
