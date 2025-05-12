package imageutil

import (
	"container-commit-go/pkg/logger"
	"container-commit-go/pkg/runtime"
	"context"
	"encoding/base64"
	"fmt"
	"io"
)

// PushOptions contains options for pushing an image
type PushOptions struct {
	RuntimeClient runtime.RuntimeClient
	ImageRef      string
	ImageID       string
	Username      string
	Password      string
}

func (opts *PushOptions) Validate() error {
	if opts.RuntimeClient == nil {
		// Default to Docker runtime client
		dockerCli, err := runtime.NewDockerRuntimeClient("", "")
		if err != nil {
			return err
		}
		opts.RuntimeClient = dockerCli
	}
	if opts.ImageRef == "" {
		return fmt.Errorf("image reference is required")
	}
	if opts.ImageID == "" {
		return fmt.Errorf("image ID is required")
	}
	if opts.Username == "" && opts.Password != "" {
		return fmt.Errorf("username is required if password is provided")
	}
	return nil
}

func (opts *PushOptions) generateRuntimePushOptions() *runtime.PushOptions {
	var auth string
	if opts.Username != "" && opts.Password != "" {
		authConfig := fmt.Sprintf(`%s:%s`, opts.Username, opts.Password)
		auth = base64.StdEncoding.EncodeToString([]byte(authConfig))
	}
	return &runtime.PushOptions{
		ImageRef:     opts.ImageRef,
		RegistryAuth: auth,
	}
}

// PushImage pushes an image to a specified repository
func PushImage(ctx context.Context, opts *PushOptions) error {
	if opts == nil {
		return fmt.Errorf("push options are required")
	}
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validating push options: %w", err)
	}
	// 1. 获取镜像大小
	size, err := opts.RuntimeClient.GetImageSize(ctx, opts.ImageID)
	if err != nil {
		return fmt.Errorf("getting image size: %w", err)
	}
	logger.Info("Pushing image", logger.WithString("image", opts.ImageRef))
	rc, err := opts.RuntimeClient.PushImage(ctx, opts.generateRuntimePushOptions())
	if err != nil {
		return fmt.Errorf("pushing image: %w", err)
	}
	defer rc.Close()
	// 2. 读取推送结果
	_, err = io.Copy(io.Discard, rc)
	if err != nil {
		return fmt.Errorf("reading push result: %w", err)
	}
	logger.Info("Image pushed", logger.WithString("image", opts.ImageRef), logger.WithInt("size", int(size)))
	return nil
}
