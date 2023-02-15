package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "durins-bridge",
	Short: "Durins Bridge is the bridge that separates the trusted Docker socket",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := root.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	root.AddCommand(create)
	root.SetHelpCommand(&cobra.Command{Hidden: true})
}
