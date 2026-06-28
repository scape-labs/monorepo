// Package logging is the thin wrapper around slog that bulksms uses everywhere.
//
// TODO: configure once at startup (JSON handler, level from env) and expose a
// single Logger for the rest of the service to import.
package logging

import "log/slog"

// L is the service-wide logger.
//
// TODO: initialise in main() and replace this default with the configured one.
var L *slog.Logger = slog.Default()
