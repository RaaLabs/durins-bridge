package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"raalabs.tech/durins-bridge/docker"
	"raalabs.tech/durins-bridge/listening"
	"raalabs.tech/durins-bridge/proxy"
)

var create = &cobra.Command{
	Use:                   "create [flags] socket",
	Short:                 "Creates proxy-sockets",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		host, err := cmd.Flags().GetString("docker-socket")
		if err != nil {
			return err
		}

		slog.Info("Creating brige", "from", path, "to", host)

		listener, err := listening.CreateListenerFromArgument(path)
		if err != nil {
			return err
		}

		client, err := docker.CreateDockerClientFromFlag(host, cmd.Context())
		if err != nil {
			return err
		}

		proxy := proxy.Server{
			Listener: listener,
			Client:   client,
		}

		return proxy.Serve(cmd.Context())
	},
}

func init() {
	create.Flags().StringP("docker-socket", "s", "/var/run/docker.sock", "The path of the original Docker socket")
}
