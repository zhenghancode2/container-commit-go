package runtime

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

func (d *DockerRuntimeClient) CommitContainer(ctx context.Context, opts *CommitOptions) (string, error) {
	if opts == nil {
		return "", fmt.Errorf("commit options are required")
	}
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("validating commit options: %w", err)
	}
	commitResp, err := d.cli.ContainerCommit(ctx, opts.ContainerIDorName, container.CommitOptions{
		Reference: opts.ImageRef,
		Comment:   opts.Message,
		Author:    opts.Author,
	})
	if err != nil {
		return "", fmt.Errorf("failed to commit container: %w", err)
	}
	return commitResp.ID, nil
}

func (d *DockerRuntimeClient) PushImage(ctx context.Context, opts *PushOptions) (io.ReadCloser, error) {
	if opts == nil {
		return nil, fmt.Errorf("push options are required")
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating push options: %w", err)
	}
	return d.cli.ImagePush(ctx, opts.ImageRef, image.PushOptions{
		All:          true,
		RegistryAuth: opts.RegistryAuth,
	})
}

func (d *DockerRuntimeClient) GetImageSize(ctx context.Context, imageID string) (int64, error) {
	imageJSON, err := d.cli.ImageInspect(ctx, imageID)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect image: %w", err)
	}
	return imageJSON.Size, nil
}
