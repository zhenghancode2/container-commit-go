package container

import (
	"context"
	"fmt"

	"container-commit-go/pkg/container/runtime"
	"container-commit-go/pkg/imageutil"
	"container-commit-go/pkg/logger"
)

// CommitOptions 提交容器为镜像的选项
type CommitOptions struct {
	RuntimeClient     runtime.RuntimeClient
	ContainerIDorName string
	ImageRef          string
	PushRepo          string
	PushUser          string
	PushPassword      string
	Insecure          bool
	Concurrency       int
}

// validate 验证提交选项
// TODO: 这里可以提交更多校验项，或者用更加优雅的方式实现
func (opts *CommitOptions) validate() error {
	if opts.RuntimeClient == nil {
		// 默认使用 docker runtime
		dockerCli, err := runtime.NewDockerRuntimeClient("", "")
		if err != nil {
			return err
		}
		opts.RuntimeClient = dockerCli
	}
	if opts.ContainerIDorName == "" {
		return fmt.Errorf("container ID or name is required")
	}
	if opts.ImageRef == "" {
		return fmt.Errorf("image reference is required")
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 1
	}
	return nil
}

// CommitContainer 提交容器为镜像并可选推送到仓库
func CommitContainer(ctx context.Context, opts *CommitOptions) error {
	// 0. 校验参数
	if opts == nil {
		return fmt.Errorf("commit options are required")
	}
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate commit options: %w", err)
	}
	// 1. 获取容器挂载目录
	mergedDir, err := opts.RuntimeClient.GetMergedDir(ctx, opts.ContainerIDorName)
	if err != nil {
		return fmt.Errorf("get merged dir: %w", err)
	}
	logger.Info("Got merged dir", logger.WithString("dir", mergedDir))

	// 2. 构建镜像
	img, err := imageutil.SaveImage(mergedDir)
	if err != nil {
		return fmt.Errorf("build image: %w", err)
	}
	logger.Info("Image built", logger.WithString("image", opts.ImageRef))

	// 4. 可选推送到 Harbor
	if opts.PushRepo != "" {
		pushOpts := &imageutil.PushOptions{
			Username:    opts.PushUser,
			Password:    opts.PushPassword,
			Insecure:    opts.Insecure,
			Concurrency: opts.Concurrency,
		}
		destRef := fmt.Sprintf("%s/%s", opts.PushRepo, opts.ImageRef)
		if err := imageutil.PushImage(ctx, img, destRef, pushOpts); err != nil {
			return fmt.Errorf("push image: %w", err)
		}
		logger.Info("Image pushed", logger.WithString("dest", destRef))
	}

	return nil
}
