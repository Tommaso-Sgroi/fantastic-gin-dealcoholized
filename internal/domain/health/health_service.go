package health

import "context"

// Pinger is the minimal persistence contract needed for health checks.
type Pinger interface {
	Ping() error
}

// HealthService reports the operational status of the application.
type HealthService interface {
	// Liveness returns a non-nil domainerr when the application cannot serve requests
	// (e.g., the database is unreachable).
	Liveness(ctx context.Context) error
}

type healthService struct {
	pinger Pinger
}

// NewHealthService returns a HealthService that pings the provided Pinger.
func NewHealthService(pinger Pinger) HealthService {
	return &healthService{pinger: pinger}
}

func (s *healthService) Liveness(_ context.Context) error {
	return s.pinger.Ping()
}
