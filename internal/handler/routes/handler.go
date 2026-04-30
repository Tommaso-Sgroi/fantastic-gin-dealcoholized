package routes

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/health"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/name"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
)

// PathParameters holds URL path parameters extracted by the router.
// Framework-agnostic: handlers receive these as a plain map.
type PathParameters map[string]string

// HandlerFunc is the signature for all route handlers. It is completely
// agnostic of the underlying HTTP framework (Gin, chi, stdlib, …).
//
// Parameters:
//   - w:   http.ResponseWriter for sending the response.
//   - r:   *http.Request — standard library request; use r.Context() for cancellation.
//   - ps:  PathParameters extracted from the URL path.
//   - ctx: RequestContext with per-request metadata (UUID, IP, Logger).
//
// Application-level dependencies are provided via the WebAPIHandler receiver.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, ps PathParameters, ctx *reqcontext.RequestContext)

// WebAPIHandler holds all application-level dependencies shared by route handlers.
// Per-request data arrives via the RequestContext parameter, not through this struct.
type WebAPIHandler struct {
	// Names is the domain service for Name CRUD operations.
	Names name.NameService

	// Health is the domain service for liveness / readiness checks.
	Health health.HealthService
}

// New creates a WebAPIHandler with the given domain services.
func New(names name.NameService, health health.HealthService) *WebAPIHandler {
	return &WebAPIHandler{
		Names:  names,
		Health: health,
	}
}
