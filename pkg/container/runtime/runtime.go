package runtime

import "context"

type RuntimeClient interface {
	// GetMergedDir 返回容器的 merged 目录路径
	GetMergedDir(ctx context.Context, containerIDorName string) (string, error)
}
