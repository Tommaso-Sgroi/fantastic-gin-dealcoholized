package routes

import (
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/fountain"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/health"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// HandlerFunc is the signature for all route handlers.
//
// Handlers receive a RequestContext that wraps the Gin context plus per-request
// metadata (request UUID and request logger).
// Application-level dependencies are provided via the WebAPIHandler receiver.
type HandlerFunc func(*reqcontext.RequestContext)

// WebAPIHandler holds all application-level dependencies shared by route handlers.
// Per-request data arrives via the RequestContext parameter, not through this struct.
type WebAPIHandler struct {
	// Fountains is the domain service for Fountain CRUD operations.
	Fountains fountain.ServiceFountain

	// Health is the domain service for liveness / readiness checks.
	Health health.HealthService
}

// New creates a WebAPIHandler with the given domain services.
func New(fountains fountain.ServiceFountain, health health.HealthService) *WebAPIHandler {
	return &WebAPIHandler{
		Fountains: fountains,
		Health:    health,
	}
}
