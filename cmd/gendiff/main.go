package main

import (
	cli_app "code/internal/drivers/cli-app"
	"context"
	"os"
)

func main() {
	cliApp := cli_app.NewCliApp()

	if err := cliApp.Init(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	args := os.Args

	if err := cliApp.Run(ctx, args); err != nil {
		panic(err)
	}
}
