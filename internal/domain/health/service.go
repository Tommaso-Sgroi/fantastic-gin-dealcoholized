package health

import "context"

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
