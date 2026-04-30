package health

import "context"

// HealthService reports the operational status of the application.
type HealthService interface {
	// Liveness returns a non-nil error when the application cannot serve requests
	// (e.g., the database is unreachable).
	Liveness(ctx context.Context) error
}

// Pinger is the minimal persistence contract needed for health checks.
type Pinger interface {
	Ping(ctx context.Context) error
}
