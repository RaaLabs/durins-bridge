package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var root = &cobra.Command{
	Use:   "durins-bridge",
	Short: "Durins Bridge is the bridge that separates the trusted Docker socket",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		options := slog.HandlerOptions{}

		if level, err := cmd.Flags().GetString("log-level"); err == nil {
			switch level {
			default:
			case "info":
				options.Level = slog.LevelInfo
			case "debug":
				options.Level = slog.LevelDebug
			case "warn":
				options.Level = slog.LevelWarn
			case "error":
				options.Level = slog.LevelError
			}
		}

		handler := options.NewTextHandler(cmd.OutOrStdout())
		slog.SetDefault(slog.New(handler))
	},
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := root.ExecuteContext(ctx); err != nil {
		slog.Error("Failure while running program", err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(create)
	root.PersistentFlags().StringP("log-level", "l", "info", "Sets the minimum log level to print")
	root.SetHelpCommand(&cobra.Command{Hidden: true})
}
