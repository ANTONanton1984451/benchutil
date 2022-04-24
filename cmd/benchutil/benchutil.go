package main

import (
	"benchutil/internal/load"
	"benchutil/internal/meet"
	"benchutil/pkg/cli"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	app, err := cli.NewApp(meet.Command(), load.New())
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	ctx := context.Background()
	if err = app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "run app: %v", err)
	}
}
