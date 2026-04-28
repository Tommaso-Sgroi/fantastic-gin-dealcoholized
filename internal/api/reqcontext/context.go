package reqcontext

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// RequestContext carries per-request metadata that middleware populates and
// handlers consume. This is ONLY for per-request data, NOT for application-level
// dependencies.
//
// Middleware (framework-specific Gin code) creates this and stores it so that
// framework-agnostic handlers can retrieve it without importing Gin.
//
// Application-level dependencies (DB, services, etc.) are provided via the
// Handler struct receiver on handler methods.
type RequestContext struct {
	// Embed the standard context.Context so handlers can use context features
	// like timeouts and cancellation.
	Context *gin.Context

	// ReqUUID is a unique ID for this request (typically a UUIDv4).
	// Use this for tracing and correlation logs.
	ReqUUID uuid.UUID

	// RemoteIP is the client's IP address (may be the actual client or a proxy,
	// depending on trusted proxy configuration).
	RemoteIP string

	// Logger is a per-request logger. Middleware should have already added
	// contextual fields (like request-id). Handlers should use this instead of
	// creating new loggers.
	Logger *slog.Logger
}

// contextKey is an unexported type to avoid collisions in context.Value.
type contextKey string

const reqCtxKey contextKey = "request-ctx"

// WithRequestContext stores a RequestContext in a context.Context.
// Typically called by middleware before passing control to handlers.
func WithRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, reqCtxKey, rc)
}

// GetRequestContext retrieves the RequestContext stored on the underlying
// http.Request context by middleware. If none is found it returns a zero-value
// RequestContext whose embedded Context is the one passed in, so callers always
// get something usable even if middleware didn't run (useful for testing).
func GetRequestContext(ctx *gin.Context) *RequestContext {
	if rc, ok := ctx.Request.Context().Value(reqCtxKey).(*RequestContext); ok {
		return rc
	}
	return &RequestContext{Context: ctx, Logger: slog.Default()}
}
