package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"container-commit-go/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "container-commit",
		Short: "一个将运行中的容器提交为镜像",
		Long: `container-commit 是一个用于将运行中的容器提交为镜像的命令行工具。
该工具的主要功能包括：
- 提交运行中的容器为镜像
- 支持将镜像推送到指定的镜像仓库`,
	}
)

func SetContext(ctx context.Context) {
	rootCmd.SetContext(ctx)
}

// SetConfig 在根命令上下文中设置配置
func SetConfig(cfg *config.Config) {
	rootCmd.SetContext(WithConfig(rootCmd.Context(), cfg))
}

func Execute() error {
	return rootCmd.Execute()
}
