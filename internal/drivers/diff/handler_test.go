package diff

import (
	cliapp "code/internal/drivers/cli-app"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initCliApp() *cliapp.CliApp {
	cliApp := cliapp.NewCliApp()

	cliApp.AddAction(Handler)

	if err := cliApp.Init(); err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, err)
		if err2 != nil {
			log.Fatal(err2)
		}

		log.Fatal(err)
	}

	return cliApp
}

func getExamplePath() string {
	return path.Join("..", "..", "..", "examples")
}

func TestHandlerSuccess(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "simple", "file1.json"), path.Join(examplePathDir, "simple", "file2.json")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerMissingArgs(t *testing.T) {
	cliApp := initCliApp()

	args := []string{""}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, "you should pass file paths", err.Error())

		return
	}

	t.Fatal("should be error")
}

func TestHandlerUnsupportedTypes(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "some1.some1"), path.Join(examplePathDir, "some2.some2")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, "unsupported extension: .some1", err.Error())

		return
	}

	t.Fatal("should be error")
}

func TestHandlerNotFound(t *testing.T) {
	cliApp := initCliApp()

	examplePathDir := getExamplePath()
	args := []string{"", path.Join(examplePathDir, "some3.some3"), path.Join(examplePathDir, "some2.some2")}
	err := cliApp.Run(context.Background(), args)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("stat %s: no such file or directory", path.Join(examplePathDir, "some3.some3")), err.Error())

		return
	}

	t.Fatal("should be error")
}
