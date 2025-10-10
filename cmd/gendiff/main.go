package main

import (
	cli_app "code/internal/drivers/cli-app"
	"code/internal/drivers/diff-checker"
	"context"
	"os"
)

func main() {
	cliApp := cli_app.NewCliApp()

	cliApp.AddAction(diff_checker.Handler)

	if err := cliApp.Init(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	args := os.Args

	if err := cliApp.Run(ctx, args); err != nil {
		panic(err)
	}
}
