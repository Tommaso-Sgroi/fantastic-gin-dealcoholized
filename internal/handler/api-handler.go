package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/reqcontext"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler/route"
)

// Handler registers all route and returns the http.Handler for the API.
// Call this method exactly once after New().
func (rt *_router) Handler() http.Handler {
	h := routes.New(rt.nameService, rt.healthService)

	rt.engine.GET("/name/:id", rt.wrap(h.GetNameByID))
	rt.engine.POST("/name", rt.wrap(h.PostNewName))
	rt.engine.PUT("/name/:id", rt.wrap(h.PutNameByID))
	rt.engine.DELETE("/name/:id", rt.wrap(h.DeleteNameByID))
	rt.engine.POST("/name/:id/replace", rt.wrap(h.PostReplaceNameByID))
	rt.engine.GET("/healthz", rt.wrap(h.Liveness))

	return rt.engine
}

// wrap adapts a framework-agnostic route handler to Gin.
func (rt *_router) wrap(fn routes.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := reqcontext.GetRequestContext(c)
		if reqCtx.Logger == nil {
			reqCtx.Logger = rt.baseLogger.With(
				slog.String("WARNING", "request context is missing logger, falling back to base logger"),
				slog.String("request-id", reqCtx.ReqUUID.String()),
			)
		}
		fn(reqCtx)
	}
}
