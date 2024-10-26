package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/supermarine1377/todoapp/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := app.Run(ctx); err != nil {
		err = fmt.Errorf("failed to run server: %w", err)
		panic(err)
	}
}
