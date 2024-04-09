package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/wscnd/go-service-boilerplate/foundation/logger"

	"github.com/ardanlabs/conf/v3"
)

var build = "develop"

func main() {
	// ___________________________________________________________________________
	// LOGGER INITIALIZATION
	var log *logger.Logger

	log = logger.NewWithEvents(
		os.Stdout,
		logger.LevelInfo,
		"SALES-API",
		func(ctx context.Context) string {
			return "00000000-0000-0000-0000-000000000000"
		},
		logger.Events{
			Error: func(ctx context.Context, r logger.Record) {
				log.Info(ctx, "*** alert ***")
			},
		},
	)

	// ___________________________________________________________________________
	// RUN PROJECT
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// ___________________________________________________________________________
	// GOMAXPROCS
	// make the service obey the k8s runtime requests/limits
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// CONFIGURATION
	// the conf package uses the following ordering: (1) default, (2) environment variables, (3) cmd line flag override
	// https://pkg.go.dev/github.com/ardanlabs/conf/v3

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s,mask"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:4000"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Sales Service Project",
		},
	}

	const prefix = "SALES"
	helpWanted, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(helpWanted)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// ___________________________________________________________________________
	// APP STARTING
	log.Info(ctx, "starting service", "version", cfg.Version.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)


	// ___________________________________________________________________________
	// CLEAN SHUTDOWN
	// make the service obey the container runtime requests/limits
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", shutdown)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", shutdown)
	}

	return nil
}
