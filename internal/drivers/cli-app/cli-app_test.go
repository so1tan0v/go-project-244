package cli_app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestCliApp_Init(t *testing.T) {
	app := NewCliApp()

	called := false
	app.AddAction(func(ctx context.Context, cmd *cli.Command) error {
		called = true
		return nil
	})

	err := app.Init()
	require.NoError(t, err)

	assert.Equal(t, "gendiff", app.Cli.Name)
	assert.Equal(t, "Compares two configuration files and shows a difference.", app.Cli.Usage)
	assert.Len(t, app.Cli.Flags, 1)

	flag := app.Cli.Flags[0]
	stringFlag, ok := flag.(*cli.StringFlag)
	require.True(t, ok)
	assert.Equal(t, "format", stringFlag.Name)
	assert.Equal(t, []string{"f"}, stringFlag.Aliases)
	assert.Equal(t, "output format", stringFlag.Usage)
	assert.Equal(t, "stylish", stringFlag.Value)

	assert.Equal(t, called, false)
}

func TestCliApp_GenerateFlags(t *testing.T) {
	app := NewCliApp()
	err := app.GenerateFlags()
	require.NoError(t, err)

	assert.Len(t, app.Flags, 1)
	flag, ok := app.Flags[0].(*cli.StringFlag)
	require.True(t, ok)
	assert.Equal(t, "format", flag.Name)
}

func TestCliApp_AddAction(t *testing.T) {
	app := NewCliApp()
	var receivedCtx context.Context
	var receivedCmd *cli.Command

	app.AddAction(func(ctx context.Context, cmd *cli.Command) error {
		receivedCtx = ctx
		receivedCmd = cmd

		return nil
	})

	assert.NotNil(t, app.Action)

	ctx := context.Background()
	dummyCmd := &cli.Command{}
	err := app.Action(ctx, dummyCmd)

	require.NoError(t, err)
	assert.Equal(t, ctx, receivedCtx)
	assert.Equal(t, dummyCmd, receivedCmd)
}

func TestCliApp_Run_Success(t *testing.T) {
	app := NewCliApp()

	formatValue := ""
	called := false

	app.AddAction(func(ctx context.Context, cmd *cli.Command) error {
		called = true
		formatValue = cmd.String("format")
		return nil
	})

	err := app.Init()
	require.NoError(t, err)

	ctx := context.Background()
	args := []string{"gendiff", "--format=plain"}

	err = app.Run(ctx, args)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "plain", formatValue)
}

func TestCliApp_Run_WithAlias(t *testing.T) {
	app := NewCliApp()

	formatValue := ""
	app.AddAction(func(_ context.Context, cmd *cli.Command) error {
		formatValue = cmd.String("format")
		return nil
	})

	err := app.Init()
	require.NoError(t, err)

	err = app.Run(context.Background(), []string{"gendiff", "-f", "json"})
	require.NoError(t, err)
	assert.Equal(t, "json", formatValue)
}
