package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

var _ RuntimeClient = &DockerRuntimeClient{}

type DockerRuntimeClient struct {
	sockHost string
	version  string
	cli      *client.Client
}

const (
	defaultSockHost = "unix:///var/run/docker.sock"
	defaultVersion  = "v1.40"
)

func NewDockerRuntimeClient(sockHost, version string) (*DockerRuntimeClient, error) {
	if sockHost == "" {
		sockHost = defaultSockHost
	}
	if version == "" {
		version = defaultVersion
	}
	cli, err := client.NewClientWithOpts(client.WithHost(sockHost), client.WithVersion(version))
	if err != nil {
		return nil, err
	}
	return &DockerRuntimeClient{cli: cli}, nil
}

func (d *DockerRuntimeClient) GetMergedDir(ctx context.Context, containerIDorName string) (string, error) {
	containerJSON, err := d.cli.ContainerInspect(ctx, containerIDorName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	if containerJSON.GraphDriver.Name != "overlay2" {
		return "", fmt.Errorf("unsupported graph driver: %s", containerJSON.GraphDriver.Name)
	}

	mergedDir, ok := containerJSON.GraphDriver.Data["MergedDir"]
	if !ok || mergedDir == "" {
		return "", fmt.Errorf("merged dir not found")
	}

	return mergedDir, nil
}
