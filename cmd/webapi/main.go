/*
Webapi is the executable for the main web server.
It builds a web server around APIs from `service/api`.
Webapi connects to external resources needed (database) and starts two web servers: the API web server, and the debug.
Everything is served via the API web server, except debug variables (/debug/vars) and profiler infos (pprof).

Usage:

	webapi [flags]

Flags and configurations are handled automatically by the code in `load-configuration.go`.

Return values (exit codes):

	0
		The program ended successfully (no errors, stopped by signal)

	> 0
		The program ended due to an domainerr

Note that this program will update the schema of the database to the latest version available (embedded in the
executable during the build).
*/
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/gin-contrib/cors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/health"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/domain/name"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/handler"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository"
)

// main is the program entry point. The only purpose of this function is to call run() and set the exit code if there is
// any domainerr
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "domainerr: ", err)
		os.Exit(1)
	}
}

// run executes the program. The body of this function should perform the following steps:
// * reads the configuration
// * creates and configure the logger
// * connects to any external resources (like databases, authenticators, etc.)
// * creates an instance of the service/api package
// * starts the principal web server (using the service/api.Router.Handler() for HTTP handlers)
// * waits for any termination event: SIGTERM signal (UNIX), non-recoverable server domainerr, etc.
// * closes the principal web server
func run() error {
	// Load Configuration and defaults
	cfg, err := loadConfiguration()
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return err
	}

	// Init logging
	var level slog.Level
	switch cfg.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "domainerr":
		level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", cfg.Logger.Level)
	}

	var logWriter io.Writer
	switch cfg.Logger.File {
	case "stderr":
		logWriter = os.Stderr
	case "stdout":
		logWriter = os.Stdout
	default:
		file, err := os.OpenFile(cfg.Logger.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}
		defer func() {
			_ = file.Close()
		}()
		logWriter = file
	}

	logger := slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug, // enable if in debug mode
	}))

	logger.Info("Initializing application...")

	// Start Database
	logger.Info("initializing database support")
	dbconn, err := sql.Open("sqlite3", cfg.DB.Filename)
	if err != nil {
		logger.Error("domainerr opening SQLite DB", slog.Any("err", err))
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Info("database stopping")
		_ = dbconn.Close()
	}()
	db, err := repository.New(dbconn)
	if err != nil {
		logger.Error("domainerr creating AppDatabase", slog.Any("err", err))
		return fmt.Errorf("creating AppDatabase: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("domainerr pinging AppDatabase", slog.Any("err", err))
		return fmt.Errorf("pinging AppDatabase: %w", err)
	}

	// Start (main) API server
	logger.Info("initializing API server")

	// Wire domain services
	nameService := name.NewNameService(db)
	healthService := health.NewHealthService(db)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this domainerr.
	serverErrors := make(chan error, 1)

	// Create the API router
	apirouter, err := handler.New(
		handler.Config{
			Logger: logger,
			CORS: cors.Config{
				AllowAllOrigins:  cfg.CORS.AllowAllOrigins,
				AllowOrigins:     cfg.CORS.AllowedOrigins,
				AllowMethods:     cfg.CORS.AllowedMethods,
				AllowHeaders:     cfg.CORS.AllowedHeaders,
				AllowCredentials: cfg.CORS.AllowCredentials,
				MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
			},
			NameService:    nameService,
			HealthService:  healthService,
			TrustedProxies: cfg.Web.TrustedProxies,
		},
	)
	if err != nil {
		logger.Error("domainerr creating the API server instance", slog.Any("err", err))
		return fmt.Errorf("creating the API server instance: %w", err)
	}
	apihandler := apirouter.Handler()

	apihandler, err = registerWebUI(apihandler)
	if err != nil {
		logger.Error("domainerr registering web UI handler", slog.Any("err", err))
		return fmt.Errorf("registering web UI handler: %w", err)
	}

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           apihandler,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	go func() {
		logger.Info("API listening on", slog.String("addr", apiserver.Addr))
		serverErrors <- apiserver.ListenAndServe()
		logger.Info("stopping API server")
	}()

	// Waiting for shutdown signal or POSIX signals
	select {
	case err := <-serverErrors:
		// Non-recoverable server domainerr
		return fmt.Errorf("server domainerr: %w", err)

	case sig := <-shutdown:
		logger.Info("signal received, start shutdown", slog.String("signal", sig.String()))

		// Asking API server to shut down and load shed.
		err := apirouter.Close()
		if err != nil {
			logger.Warn("graceful shutdown of apirouter domainerr", slog.Any("err", err))
		}

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and load shed.
		err = apiserver.Shutdown(ctx)
		if err != nil {
			logger.Warn("domainerr during graceful shutdown of HTTP server", slog.Any("err", err))
			err = apiserver.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
