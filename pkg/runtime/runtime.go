package runtime

import (
	"context"
	"fmt"
	"io"
)

type RuntimeClient interface {
	// CommitContainer 提交容器为镜像
	CommitContainer(ctx context.Context, opts *CommitOptions) (string, error)
	// PushImage 推送镜像到仓库
	PushImage(ctx context.Context, opts *PushOptions) (io.ReadCloser, error)
	// GetImageSize 获取镜像大小
	GetImageSize(ctx context.Context, imageID string) (int64, error)
}

type CommitOptions struct {
	ContainerIDorName string
	ImageRef          string
	Message           string
	Author            string
}

func (opts *CommitOptions) Validate() error {
	if opts.ContainerIDorName == "" {
		return fmt.Errorf("container ID or name is required")
	}
	if opts.ImageRef == "" {
		return fmt.Errorf("image reference is required")
	}
	if opts.Message == "" {
		opts.Message = defaultCommitMessage
	}
	if opts.Author == "" {
		opts.Author = defaultAuthor
	}
	return nil
}

const (
	defaultCommitMessage = "default commit message: commit from container-commit"
	defaultAuthor        = "unknown by container-commit"
)

type PushOptions struct {
	ImageRef     string
	RegistryAuth string
}

func (opts *PushOptions) Validate() error {
	if opts.ImageRef == "" {
		return fmt.Errorf("image reference is required")
	}
	if opts.RegistryAuth == "" {
		return fmt.Errorf("registry auth is required")
	}
	return nil
}
