// Package testpg spins up a per-test Postgres database for wirepay's integration
// tests. Tests should call Spawn(t) in TestMain to get a ready *sql.DB and a
// cleanup hook.
//
// TODO: implement — fork from template DB, run migrations, return *sql.DB.
package testpg

import (
	"context"
	"database/sql"
	"testing"
)

// Spawn returns a *sql.DB connected to a fresh test database and registers a
// cleanup that drops it when t finishes.
//
// TODO: implement.
func Spawn(t *testing.T) *sql.DB {
	t.Helper()
	_ = context.Background()
	// TODO: connect to template, fork, migrate, return.
	return nil
}
