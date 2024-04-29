package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/wscnd/go-service-boilerplate/apps/tools/admin/commands"
	"github.com/wscnd/go-service-boilerplate/libs/logger"
)

func main() {
	log := logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	if err := run(log); err != nil {
		fmt.Println("msg", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	return commands.GenToken(log)
}
