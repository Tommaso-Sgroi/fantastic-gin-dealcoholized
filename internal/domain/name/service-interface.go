package name

import "context"

// ServiceName defines the business operations available on Name entities.
// All methods accept a context.Context for cancellation and tracing propagation.
type ServiceName interface {
	// GetByID returns the Name with the given id.
	// Returns ErrNotFound when no such entity exists.
	GetByID(ctx context.Context, id int64) (Name, error)

	// Create validates and persists a new Name, returning the created entity.
	// Returns ErrInvalidInput when the name is empty or exceeds the allowed length.
	Create(ctx context.Context, value string) (Name, error)

	// UpdateByID replaces the value of the Name identified by id.
	// Returns ErrNotFound when no such entity exists.
	// Returns ErrInvalidInput when the new value fails validation.
	UpdateByID(ctx context.Context, id int64, value string) (Name, error)

	// DeleteByID removes the Name identified by id.
	// Returns ErrNotFound when no such entity exists.
	DeleteByID(ctx context.Context, id int64) error
}
