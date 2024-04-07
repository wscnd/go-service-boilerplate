package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/wscnd/go-service-boilerplate/foundation/logger"
)

var build = "develop"

func main() {
	var log *logger.Logger

	log = logger.NewWithEvents(
		os.Stdout,
		logger.LevelInfo,
		"SALES-API",
		func(ctx context.Context) string {
			return "id-123123-123123"
		},
		logger.Events{
			Error: func(ctx context.Context, r logger.Record) {
				log.Info(ctx, "*** alert ***")
			},
		},
	)

	// _______________

	ctx := context.Background()
	log.Info(ctx, fmt.Sprintf("starting service build[%s] CPU[%v]", build, runtime.GOMAXPROCS(0)))

	defer log.Info(ctx, "service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Info(ctx, "stopping service")
}
