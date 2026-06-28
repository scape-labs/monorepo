package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/scape-labs/monorepo/shared/idgen"
)

// PGStore is the Postgres implementation of Store.
//
// TODO: implement every method on Store.
type PGStore struct {
	db *sql.DB
}

// NewPGStore constructs a PGStore.
//
// TODO: validate the *sql.DB is non-nil.
func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{db: db}
}

// Get fetches an entity by id.
//
// TODO: implement.
func (s *PGStore) Get(ctx context.Context, id idgen.ID) (Entity, error) {
	_ = ctx
	return Entity{}, fmt.Errorf("not implemented")
}

// Create persists a new entity.
//
// TODO: implement.
func (s *PGStore) Create(ctx context.Context, e Entity) error {
	_ = ctx
	_ = e
	return fmt.Errorf("not implemented")
}
