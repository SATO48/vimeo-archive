package main

import (
	"log/slog"
	"os"

	"github.com/sato48/vimeo-archive/app"
	"github.com/sato48/vimeo-archive/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func main() {
	App, err := app.Boot(
		app.BootstrapFunc(app.S3),
		app.BootstrapFunc(app.Gorm),
		app.BootstrapFunc(app.Vimeo),
		app.BootstrapFunc(cmd.Bootstrap),
	)
	if err != nil {
		slog.Error("unable to register app", "error", err)
		os.Exit(1)
	}

	var cmd *cobra.Command
	if App.Resolve(&cmd) != nil {
		slog.Error("unable to resolve command", "error", err)
		os.Exit(1)
	}

	cmd.Execute()
}
