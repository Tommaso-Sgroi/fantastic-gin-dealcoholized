package handler

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/database"
)

// PathParameters holds URL path parameters extracted by the router.
// Framework-agnostic: handlers receive these as a map.
type PathParameters map[string]string

// WebAPIHandlerFunc is the signature for all API handlers. It is completely
// agnostic of the underlying HTTP framework (Gin, chi, stdlib, …).
//
// Parameters:
//   - w: http.ResponseWriter for sending responses
//   - r: *http.Request (standard library request)
//   - ps: PathParameters extracted from the URL path
//   - ctx: RequestContext with per-request metadata (UUID, IP, Logger)
//
// Handlers DO NOT depend on HTTP framework specifics. They receive
// application dependencies via the Handler receiver (see Handler struct).
type WebAPIHandlerFunc func(w http.ResponseWriter, r *http.Request, ps PathParameters, ctx *reqcontext.RequestContext)

// WebAPIHandler Handler holds all application-level dependencies that handlers need.
// Each handler method is a receiver on *Handler, so it has access to
// all these dependencies.
//
// Per-request data comes via RequestContext parameter, not through Handler.
type WebAPIHandler struct {
	// DB is the database connection (shared, long-lived).
	DB database.AppDatabase
}

// New creates a Handler with the given application-level dependencies.
func New(db database.AppDatabase) *WebAPIHandler {
	return &WebAPIHandler{
		DB: db,
	}
}
