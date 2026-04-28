package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/api/handler"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/api/middleware"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/api/reqcontext"
)

// Handler builds the HTTP handler for the API. It wires up middleware,
// creates the business-logic handler, and registers every route.
func (rt *_router) Handler() http.Handler {
	// Create the (framework-agnostic) handler that holds all dependencies.
	h := handler.New(rt.db)

	// Global middleware — runs before every route.
	rt.engine.Use(middleware.RequestContextLogger(rt.baseLogger))

	// --- routes --------------------------------------------------------
	rt.engine.GET("/", rt.wrap(h.GetHelloWorld))
	rt.engine.GET("/liveness", rt.wrap(h.Liveness))

	return rt.engine
}

// wrap adapts a framework-agnostic API handler to Gin.
func (rt *_router) wrap(fn handler.APIHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := reqcontext.GetRequestContext(c)

		if reqCtx.Logger == nil {
			// fallback
			reqCtx.Logger = rt.baseLogger.With(
				slog.String("WARNING", "Request context is missing logger, falling back to default logger"),
				slog.String("request-id", reqCtx.ReqUUID.String()),
			)
		}

		params := make(handler.PathParameters, len(c.Params))
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		fn(c.Writer, c.Request, params, reqCtx)
	}
}
