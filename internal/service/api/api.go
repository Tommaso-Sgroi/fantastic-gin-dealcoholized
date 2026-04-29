/*
Package api exposes the main API engine. All HTTP APIs are handled here — so-called "business logic" lives in the
handler sub-package, completely agnostic of the HTTP framework (Gin).

To use this package, create a new instance with New() passing a valid Config. The resulting Router will have the
Router.Handler() function that returns an http.Handler usable in a standard http.Server.

Example:

	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	apiserver := http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: router,
	}
	apiserver.ListenAndServe()

See cmd/webapi/main.go for a full usage example.
*/
package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/service/database"
)

// Config provides dependencies and configuration to New.
type Config struct {
	CORS cors.Config

	// Logger where log entries are sent.
	Logger *slog.Logger

	// TrustedProxies is the list of trusted proxy IPs for Gin.
	TrustedProxies []string

	// Database is the application database instance.
	Database database.AppDatabase
}

// Router is the public interface for the API layer.
type Router interface {
	// Handler returns an http.Handler that serves all registered API routes.
	Handler() http.Handler

	// Close terminates any resource opened during the router lifecycle.
	Close() error
}

// New returns a new Router instance.
func New(cfg Config) (Router, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	engine := gin.New()

	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.HandleMethodNotAllowed = true
	if err := engine.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, fmt.Errorf("setting trusted proxies: %w", err)
	}

	engine.Use(cors.New(cfg.CORS))
	engine.Use(gin.Recovery())

	return &router{
		engine:     engine,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

type router struct {
	engine     *gin.Engine
	baseLogger *slog.Logger
	db         database.AppDatabase
}
