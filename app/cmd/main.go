package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/supermarine1377/todoapp/app"
)

// @title			タスク管理用のAPI
// @version		0.01
// @description	タスク管理用のAPI
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run app: %v\n", err)
		os.Exit(1)
	}
}
