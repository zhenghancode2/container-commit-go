# container-commit-go

一个命令行工具，用于将运行中的容器导出为镜像，并可推送到 Harbor 等镜像仓库，基于 go-containerregistry 实现。

## 环境要求

- Go 1.24+
- Docker 正常运行
- 可访问 Harbor 或兼容的镜像仓库

## 安装

```bash
git clone https://github.com/zhenghancode2/container-commit-go.git
cd container-commit-go
go build -o bin/container-commit main.go
```

## 使用方法

```bash
# 提交运行中的容器为镜像
container-commit commit <container-id> <image-ref>

# 提交并推送到远程仓库
container-commit commit <container-id> <image-ref> -r <harbor_repo> -u <username> -p <password>

# 允许不安全的推送
container-commit commit <container-id> <image-ref> -r <harbor_repo> -u <username> -p <password> --insecure
```

### 参数说明

- `<container-id>`: 容器 ID 或名称（必填）
- `<image-ref>`: 新镜像名（如 myrepo/myimage:tag）（必填）
- `-r, --repo`: Harbor 仓库地址（可选，指定则推送）
- `-u, --user`: 仓库用户名（可选）
- `-p, --password`: 仓库密码（可选）
- `--insecure`: 允许不安全的仓库连接（可选，默认 false）
