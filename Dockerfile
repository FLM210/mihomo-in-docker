# 使用官方 Golang 镜像作为构建环境 - 构建配置生成器
FROM golang:1.25.5-alpine AS config-builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o config-generator .

# 使用 alpine 镜像下载 mihomo
FROM alpine:latest AS mihomo-downloader

# 安装必要的包
RUN apk --no-cache add ca-certificates wget jq gzip

# 下载并解压 mihomo
RUN ARCH=$(uname -m) && \
    case $ARCH in \
        x86_64) BINARY="mihomo-linux-amd64" ;; \
        aarch64) BINARY="mihomo-linux-arm64" ;; \
        *) echo "不支持的架构: $ARCH"; exit 1 ;; \
    esac && \
    API_URL="https://api.github.com/repos/metacubex/mihomo/releases/latest" && \
    DOWNLOAD_URL=$(wget -qO- $API_URL | jq -r ".assets[] | select(.name | test(\"$BINARY\")) | select(.name | endswith(\".gz\")) | .browser_download_url") && \
    wget -qO- $DOWNLOAD_URL | gunzip -c > /mihomo && \
    chmod +x /mihomo

# 最终运行阶段 - 从 Alpine 开始，只包含必要的二进制文件
FROM alpine:latest

# 安装 ca-certificates 以支持 HTTPS 请求
RUN apk --no-cache add ca-certificates bash

# 设置工作目录
WORKDIR /app

# 从构建阶段复制配置生成器和 mihomo 二进制文件
COPY --from=config-builder /app/config-generator .
COPY --from=mihomo-downloader /mihomo /usr/local/bin/mihomo

# 创建默认配置目录
RUN mkdir -p /root/.config/mihomo

# 复制启动脚本
COPY start.sh .

# 确保启动脚本有执行权限
RUN chmod +x start.sh

# 暴露端口（mihomo 默认端口）
EXPOSE 9090 10801

# 运行启动脚本
ENTRYPOINT ["./start.sh"]