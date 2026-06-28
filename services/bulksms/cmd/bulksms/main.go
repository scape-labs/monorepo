// Package main is the entrypoint for bulksms.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting bulksms")
	// TODO: kit.New("bulksms").Register(...).Run()
	_ = server.Component{}
}
