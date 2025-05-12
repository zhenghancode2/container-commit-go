package container

import (
	"context"
	"fmt"

	"container-commit-go/pkg/logger"
	"container-commit-go/pkg/runtime"
)

// CommitOptions 提交容器为镜像的选项
type CommitOptions struct {
	RuntimeClient runtime.RuntimeClient
	runtime.CommitOptions
}

// Validate 验证提交选项
// TODO: 这里可以提交更多校验项，或者用更加优雅的方式实现
func (opts *CommitOptions) Validate() error {
	if opts.RuntimeClient == nil {
		// 默认使用 docker runtime
		dockerCli, err := runtime.NewDockerRuntimeClient("", "")
		if err != nil {
			return err
		}
		opts.RuntimeClient = dockerCli
	}
	return opts.CommitOptions.Validate()
}

// CommitContainer 提交容器为镜像并可选推送到仓库
func CommitContainer(ctx context.Context, opts *CommitOptions) error {
	// 0. 校验参数
	if opts == nil {
		return fmt.Errorf("commit options are required")
	}
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate commit options: %w", err)
	}
	// 1. commit 容器为镜像
	logger.Info("Committing container", logger.WithString("container", opts.ContainerIDorName))
	_, err := opts.RuntimeClient.CommitContainer(ctx, &runtime.CommitOptions{
		ContainerIDorName: opts.ContainerIDorName,
		ImageRef:          opts.ImageRef,
	})
	if err != nil {
		return fmt.Errorf("committing container: %w", err)
	}
	logger.Info("Container committed", logger.WithString("image", opts.ImageRef))
	return nil
}
