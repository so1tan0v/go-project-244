package diff

import (
	cli_app "code/internal/drivers/cli-app"
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initCliApp() *cli_app.CliApp {
	cliApp := cli_app.NewCliApp()

	cliApp.AddAction(Handler)

	if err := cliApp.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return cliApp
}

func getExamplePath() string {
	return path.Join("..", "..", "..", "examples")
}

func TestHandler_Success(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "file1.json"), path.Join(examplePathDir, "file2.json")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandler_MissingArgs(t *testing.T) {
	cliApp := initCliApp()

	args := []string{""}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, "you should pass file paths", err.Error())

		return
	}

	t.Fatal("should be error")
}

func TestHandler_UnsupportedTypes(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "some1.some1"), path.Join(examplePathDir, "some2.some2")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, "files must have the same supported extension: got .some1 and .some2", err.Error())

		return
	}

	t.Fatal("should be error")
}

func TestHandler_NotFound(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "some3.some3"), path.Join(examplePathDir, "some2.some2")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("open %s: no such file or directory", path.Join(examplePathDir, "some3.some3")), err.Error())

		return
	}

	t.Fatal("should be error")
}
