package fountain

import "context"

// ServiceFountain defines the business operations available on Fountain entities.
// All methods accept a context.Context for cancellation and tracing propagation.
type ServiceFountain interface {
	// List returns all fountains.
	// If lat and lng are non-nil, the result is filtered to fountains within rangeKm
	// kilometres of the given coordinates and sorted by distance.
	// rangeKm is ignored when lat/lng are nil.
	List(ctx context.Context, lat *float64, lng *float64, rangeKm float64) ([]Fountain, error)

	// Create validates and persists a new Fountain, returning the created entity
	// (including the server-assigned ID).
	// Returns ErrInvalidInput when the fountain fails validation.
	Create(ctx context.Context, f Fountain) (Fountain, error)

	// Update replaces all properties of the Fountain identified by id with those
	// provided in f. Returns the updated entity.
	// Returns ErrNotFound when no entity with id exists.
	// Returns ErrInvalidInput when f fails validation.
	Update(ctx context.Context, id uint64, f Fountain) (Fountain, error)

	// Delete removes the Fountain identified by id.
	// Returns ErrNotFound when no entity with id exists.
	Delete(ctx context.Context, id uint64) error
}
