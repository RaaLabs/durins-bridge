package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

func CreateDockerClientFromFlag(path string, ctx context.Context) (*client.Client, error) {
	client, err := client.NewClientWithOpts(
		client.WithHost("unix://"+path),
		client.FromEnv,
	)
	if err != nil {
		return nil, err
	}

	log.Println("Connecting to Docker on", client.DaemonHost())

	info, err := client.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get info from Docker host: %w", err)
	}

	log.Println("Connected to Docker version", info.ServerVersion)

	return client, nil
}
