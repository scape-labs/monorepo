// Package main is the entrypoint for flow-ussd.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting flow-ussd")
	// TODO: kit.New("flow-ussd").Register(...).Run()
	_ = server.Component{}
}
