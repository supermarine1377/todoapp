package main

import (
	"context"
	"fmt"

	"github.com/supermarine1377/todoapp/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		err = fmt.Errorf("failed to run server: %w", err)
		panic(err)
	}
}
