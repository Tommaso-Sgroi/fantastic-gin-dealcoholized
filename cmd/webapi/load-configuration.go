package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/conf/v3/yaml"
)

const prefix = "API"

// WebAPIConfiguration describes the web API configuration. This structure is automatically parsed by
// loadConfiguration and values from flags, environment variable or configuration file will be loaded.
type WebAPIConfiguration struct {
	Config struct {
		Path string `conf:"default:/config/config.yml,short:c,help:path to the YAML configuration file"`
	}
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000,short:H,help:host:port the REST API server listens on"`
		DebugHost       string        `conf:"default:0.0.0.0:4000,help:host:port the debug/pprof server listens on"`
		ReadTimeout     time.Duration `conf:"default:5s,help:maximum duration to read an entire HTTP request"`
		WriteTimeout    time.Duration `conf:"default:5s,help:maximum duration to write an HTTP response"`
		ShutdownTimeout time.Duration `conf:"default:5s,help:maximum duration to wait for active connections during graceful shutdown"`
		TrustedProxies  []string      `conf:"help:list of trusted proxy IP addresses or CIDR ranges"`
	}
	CORS struct {
		AllowAllOrigins  bool     `conf:"default:false,help:allow requests from any origin"`
		AllowedOrigins   []string `conf:"default:http://localhost,help:list of allowed origins"`
		AllowedMethods   []string `conf:"default:GET;POST;PUT;PATCH;DELETE;OPTIONS,help:list of allowed HTTP methods"`
		AllowedHeaders   []string `conf:"default:Accept;Content-Type,help:list of allowed HTTP headers"`
		AllowCredentials bool     `conf:"default:false,help:allow cookies and credentials in cross-origin requests"`
		MaxAge           int      `conf:"default:86400,help:max age in seconds for preflight cache"`
	}
	Debug bool `conf:"short:d,help:enable verbose debug logging"`
	DB    struct {
		Filename string `conf:"default:/tmp/decaf.db,help:path to the SQLite database file"`
	}
	Logger struct {
		Level string `conf:"default:info,help:log level (debug, info, warn, error)"`
		File  string `conf:"default:/tmp/dealch.log,help:path to the log file"`
	}
}

// loadConfiguration creates a WebAPIConfiguration starting from flags, environment variables and configuration file.
// It works by loading the configuration file first (if it exists), then applying environment variables, and finally
// command line flags. So, CLI parameters will override the environment, which will override the configuration file.
// Note that the configuration file path can be specified only via CLI or environment variable.
func loadConfiguration() (WebAPIConfiguration, error) {
	var cfg WebAPIConfiguration

	// First pass: parse env/flags to discover the demo file path.
	help, err := conf.ParseWithOptions(prefix, &cfg, conf.WithStrictFlags())
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help) //nolint:forbidigo
			return cfg, conf.ErrHelpWanted
		}
		return cfg, fmt.Errorf("parsing demo: %w", err)
	}

	// If a YAML demo file exists, reparse supplying the YAML data (useful in k8s/compose).
	// YAML is applied before env vars and flags, so env/flags take precedence.
	fp, err := os.Open(cfg.Config.Path)
	if err != nil && !os.IsNotExist(err) {
		return cfg, fmt.Errorf("can't read the demo file, while it exists: %w", err)
	} else if err == nil {
		yamlFile, err := io.ReadAll(fp)
		_ = fp.Close()
		if err != nil {
			return cfg, fmt.Errorf("can't read demo file: %w", err)
		}
		help, err = conf.Parse(prefix, &cfg, yaml.WithData(yamlFile))
		if err != nil {
			if errors.Is(err, conf.ErrHelpWanted) {
				fmt.Println(help) //nolint:forbidigo
				return cfg, conf.ErrHelpWanted
			}
			return cfg, fmt.Errorf("parsing demo with yaml: %w", err)
		}
	}

	return cfg, nil
}
