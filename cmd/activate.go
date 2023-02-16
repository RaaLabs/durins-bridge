package cmd

import (
	"fmt"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"raalabs.tech/durins-bridge/docker"
	"raalabs.tech/durins-bridge/proxy"
)

var activate = &cobra.Command{
	Use:   "activate",
	Short: "Accepts sockets created by systemd",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, err := cmd.Flags().GetString("docker-socket")
		if err != nil {
			return err
		}

		slog.Info("Activating brige", "to", host)

		listeners, err := activation.Listeners()
		if err != nil {
			return err
		}
		if len(listeners) != 1 {
			return fmt.Errorf("expected exactly 1 file descriptor - got %d", len(listeners))
		}

		client, err := docker.CreateDockerClientFromFlag(host, cmd.Context())
		if err != nil {
			return err
		}

		proxy := proxy.Server{
			Listener: listeners[0],
			Client:   client,
		}

		return proxy.Serve(cmd.Context())
	},
}

func init() {
	activate.Flags().StringP("docker-socket", "s", "/var/run/docker.sock", "The path of the original Docker socket")
}
