package service

import (
	"context"

	"github.com/scape-labs/monorepo/shared/idgen"
)

// Store is the data-access boundary for bulksms.
//
// TODO: define the full set of methods that the service layer needs.
type Store interface {
	Get(ctx context.Context, id idgen.ID) (Entity, error)
	Create(ctx context.Context, e Entity) error
	// TODO: add Update / Delete / List / Query.
}

// Entity is a placeholder for whatever aggregate root bulksms owns.
//
// TODO: replace with the real type.
type Entity struct {
	ID idgen.ID
	// TODO: fields.
}
