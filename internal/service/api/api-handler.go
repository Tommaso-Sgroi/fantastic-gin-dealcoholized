package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/handler"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
)

// Handler returns the http.Handler for the API. Routes and middleware are
// registered once in New() — this is a pure getter.
func (rt *_router) Handler() http.Handler {
	h := handler.New(rt.db)

	// --- routes --------------------------------------------------------
	rt.engine.GET("/name/:id", rt.wrap(h.GetNameByID))
	rt.engine.POST("/name", rt.wrap(h.PostNewName))
	rt.engine.PUT("/name/:id", rt.wrap(h.PutNameByID))
	rt.engine.DELETE("/name/:id", rt.wrap(h.DeleteNameByID))

	rt.engine.GET("/healthz", rt.wrap(h.Liveness))

	return rt.engine
}

// wrap adapts a framework-agnostic API handler to Gin.
func (rt *_router) wrap(fn handler.WebAPIHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := reqcontext.GetRequestContext(c)

		if reqCtx.Logger == nil {
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
