package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/api/reqcontext"
)

// RequestContextLogger returns a Gin middleware that:
//  1. Generates a unique request ID.
//  2. Builds a handler.RequestContext and stores it in the request context
//     so that downstream (framework-agnostic) handlers can retrieve it.
//  3. Logs the completed request with latency and status.
func RequestContextLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		header := ""
		for k, h := range c.Request.Header {
			for _, v := range h {
				header += k + "=" + v + ";"
			}
		}

		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		requestID, err := uuid.NewV4()
		if err != nil {
			logger.Error("cannot generate a random UUIDv4", slog.Any("error", err))
			requestID = uuid.Nil
		}

		// Build the framework-agnostic RequestContext.
		reqCtx := &reqcontext.RequestContext{
			Context:  c,
			ReqUUID:  requestID,
			RemoteIP: ip,
			Logger:   logger.With(slog.String("request-id", requestID.String())),
		}

		// Store it so that handler.GetRequestContext(r.Context()) works.
		c.Request = c.Request.WithContext(reqcontext.WithRequestContext(c.Request.Context(), reqCtx))
		c.Next()

		latency := time.Since(start)
		logger.Info("request completed",
			slog.String("request-id", requestID.String()),
			slog.Int("status", c.Writer.Status()),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("header", header),
			slog.String("query", query),
			slog.String("ip", ip),
			slog.Duration("latency", latency),
			slog.String("user_agent", userAgent),
		)
	}
}
