package main

import (
	cli_app "code/internal/drivers/cli-app"
	handler "code/internal/drivers/gendiff"
	"context"
	"fmt"
	"os"
)

func main() {
	cliApp := cli_app.NewCliApp()

	cliApp.AddAction(handler.Handler)

	if err := cliApp.Init(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	ctx := context.Background()
	args := os.Args

	if err := cliApp.Run(ctx, args); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
