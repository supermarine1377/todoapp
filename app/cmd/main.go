package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/supermarine1377/todoapp/app"
)

func main() {
	// gs はGraceful shotdown をするか示す
	var gs bool
	flag.BoolVar(
		&gs,
		"gracefully-shutdown",
		false,
		"trueであればGraceful shutdownする、falseであればGraceful shutdownしない",
	)
	flag.Parse()

	ctx := context.Background()

	fmt.Println(gs)

	if gs {
		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer stop()
		run(ctx)
	} else {
		run(ctx)
	}
}

func run(ctx context.Context) {
	if err := app.Run(ctx); err != nil {
		os.Exit(1)
	}
}
