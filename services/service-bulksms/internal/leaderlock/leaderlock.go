// Package leaderlock gives bulksms a single-leader guard around tickers and
// background jobs. Only one replica at a time should run a given job.
//
// TODO: implement on top of Postgres advisory locks or kit/leaderlock.
package leaderlock

import "context"

// Lock represents an acquired leader lock.
//
// TODO: implement.
type Lock struct {
	key string
}

// Acquire tries to take the lock for `key`.
//
// TODO: implement with Postgres `pg_try_advisory_lock`.
func Acquire(ctx context.Context, key string) (*Lock, error) {
	_ = ctx
	return &Lock{key: key}, nil
}

// Release frees the lock.
//
// TODO: implement.
func (l *Lock) Release(ctx context.Context) error {
	_ = ctx
	return nil
}
