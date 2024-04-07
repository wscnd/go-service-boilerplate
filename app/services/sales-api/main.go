package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/wscnd/go-service-boilerplate/foundation/logger"
)

var build_ref = "develop"

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
	// make the service obey the docker runtime requests/limits
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ___________________________________________________________________________
	// APP STARTING
	log.Info(ctx, "starting service", "version", build_ref)

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
