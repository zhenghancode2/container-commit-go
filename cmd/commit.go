package cmd

import (
	"container-commit-go/config"
	"container-commit-go/pkg/container"
	"container-commit-go/pkg/logger"
	"container-commit-go/pkg/runtime"
	"context"

	"github.com/spf13/cobra"
)

var (
	imageRepo    string
	repoUser     string
	repoPassword string
	insecure     bool
)

var commitCmd = &cobra.Command{
	Use:   "commit [container-id] [image-ref] [flags]",
	Short: "将运行中的容器提交为镜像",
	Long:  `将运行中的容器提交为指定的镜像，若是提供镜像仓库信息并将其推送到指定的镜像仓库，例如 Harbor。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 解析命令行参数
		if len(args) != 2 {
			return cmd.Help()
		}
		containerID := args[0]
		newImageName := args[1]
		// 初始化日志系统
		logger.Init(config.DefaultLogConfig())
		// run container commit
		if err := run(cmd.Context(), containerID, newImageName); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&repoUser, "user", "u", "", "Harbor registry username")
	commitCmd.Flags().StringVarP(&repoPassword, "password", "p", "", "Harbor registry password")
}

func run(ctx context.Context, containerID, newImageName string) error {
	commitOpts := &container.CommitOptions{
		CommitOptions: runtime.CommitOptions{
			ContainerIDorName: containerID,
			ImageRef:          newImageName,
		},
	}
	return container.CommitContainer(ctx, commitOpts)
}
