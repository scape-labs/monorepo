// Package idgen generates 64-bit flake-style ids.
//
// The id is `(timestamp_ms, worker_id, sequence)` packed into 64 bits so it
// sorts chronologically and fits in a Postgres `bigint`.
//
// TODO: implement New, MustNew, Parse, and a worker-id initialiser that reads
// from kit/config.
package idgen

// ID is a 64-bit identifier.
//
// TODO: define as a typed `uint64` to avoid accidental mixing with int64.
type ID uint64

// New returns a fresh ID. Stub returns 0.
//
// TODO: replace with a real generator.
func New() ID { return 0 }

// Parse decodes an ID from a string representation.
//
// TODO: implement.
func Parse(s string) (ID, error) {
	var id ID
	return id, nil
}

// String returns the base-10 string form.
//
// TODO: implement.
func (id ID) String() string { return "0" }
