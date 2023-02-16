package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"golang.org/x/exp/slog"
)

func CreateDockerClientFromFlag(path string, ctx context.Context) (*client.Client, error) {
	client, err := client.NewClientWithOpts(
		client.WithHost("unix://"+path),
		client.FromEnv,
	)
	if err != nil {
		return nil, err
	}

	slog.Info("Connecting to Docker on", "socket", client.DaemonHost())

	info, err := client.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get info from Docker host: %w", err)
	}

	slog.Debug("Connected to Docker", "version", info.ServerVersion)

	return client, nil
}
