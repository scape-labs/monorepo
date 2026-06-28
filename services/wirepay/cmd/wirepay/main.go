// Package main is the entrypoint for wirepay.
//
// TODO: wire up all components, register routes, run.
package main

import (
	"log/slog"

	"github.com/scape-labs/platform/kit/server"
	// TODO: import other kit packages
)

func main() {
	slog.Info("starting wirepay")
	// TODO: kit.New("wirepay").Register(...).Run()
	_ = server.Component{}
}
