# 构建阶段
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# 拷贝 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝项目源代码
COPY . .

# 构建可执行文件
RUN go build -o container-commit main.go

# 运行阶段
FROM debian:stable-slim

WORKDIR /app

# 拷贝可执行文件和配置文件
COPY --from=builder /app/container-commit .
# 如有默认配置文件
COPY config.yaml .

# 设置容器启动命令
CMD ["./container-commit"]