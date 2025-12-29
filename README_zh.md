# Mihomo in Docker

## 项目介绍

本项目旨在将 Mihomo (Clash 的一个高性能分支) 容器化，通过 Docker 简化其部署流程，并支持自动生成配置文件，便于在 Kubernetes 或普通容器环境中快速运行。

### 目标用户：
- 需要使用代理工具进行流量管理的技术用户
- 希望在容器环境（如 Docker、Kubernetes）中部署 Mihomo 的运维或开发者

### 解决的核心问题：
- 简化 Mihomo 在不同架构下的安装与配置
- 实现配置的自动化生成，减少手动编辑 YAML 文件的复杂性
- 提供可复用、轻量化的镜像用于生产部署

### 系统功能：
- 自动检测主机架构并下载对应版本的 Mihomo 二进制文件
- 使用 Go 编写的配置生成器，根据环境变量生成 Mihomo 配置文件
- 支持通过启动脚本自动运行 Mihomo 服务
- 容器内暴露常用端口（API、Proxy、UI）

## 功能特性

- **架构自适应**: 在构建时自动识别系统架构（amd64/arm64），下载对应的 Mihomo 二进制
- **配置自动生成**: 利用 Go 程序解析环境变量或模板生成 config.yaml
- **智能节点过滤**: 自动根据关键词过滤代理节点，支持为特定地区创建自动故障转移组
- **轻量运行**: 基于 Alpine Linux 构建最终镜像，仅包含必要组件
- **多阶段构建**: 分离构建、下载和运行阶段，减小镜像体积
- 无 root 权限依赖，适合安全强化场景
- 支持 CI/CD 流水线构建和部署
- 可集成至 Kubernetes 部署

## 快速开始

### 前提条件
- Docker 20+

### 构建镜像

```bash
# 克隆项目
git clone https://github.com/flm210/mihomo-in-docker.git
cd mihomo-in-docker

# 构建镜像
docker build -t mihomo-in-docker .
```

### 运行容器

```bash
# 使用示例配置运行容器
docker run -p 9090:9090 -p 10801:10801 \
  -e "SUBSCRIPTION_URL=https://example.com/subscription" \
  -e "FILTER_KEYWORDS=HK,SG" \
  mihomo-in-docker
```

### 环境变量

- `SUBSCRIPTION_URL`: 代理节点的订阅 URL
- `FILTER_KEYWORDS`: 用于过滤节点的关键字（以逗号分隔）。程序会根据这些关键字自动创建故障转移组。
- `OUTPUT_FILE`: 输出文件名（默认: config.yaml）
- `OUTPUT_FORMAT`: 输出格式（yaml 或 json，默认: yaml）

## 配置

配置生成器支持：
- 订阅 URL 处理
- 节点关键字过滤（自动创建故障转移组）
- YAML/JSON 输出格式
- 支持多种代理协议（VMess、VLESS、Trojan、Shadowsocks 等）
- 根据过滤关键字自动生成带自动故障转移功能的代理组

## 高级用法：自动节点切换

- 如果设置 `FILTER_KEYWORDS=HK`，程序将筛选出名称中包含"HK"（香港）的所有节点并创建一个自动故障转移组。
- 如果设置 `FILTER_KEYWORDS=HK,SG,JP`，它将创建一个包含所有匹配这些关键词的节点的单一组（香港、新加坡、日本），并在它们之间自动故障转移。

### Kubernetes 示例：特定地区部署

如果你想在 K8S 上启动一个专门用于 HK（香港）节点的 Pod，只需将过滤器设置为包含机场中香港对应节点名称的关键词：

```yaml
env:
- name: SUBSCRIPTION_URL
  value: "https://your-provider.com/subscription"
- name: FILTER_KEYWORDS
  value: "香港,HK,hk"  # 与你的提供商香港节点匹配的关键词
```

这将自动过滤并创建仅包含香港节点的故障转移组，并根据连接情况在它们之间自动切换。

## 端口

- `9090`: 外部控制器 API
- `10801`: 混合代理端口（HTTP 和 SOCKS）

## Kubernetes 部署

本项目在 [k8s-deployment.yaml](k8s-deployment.yaml) 文件中提供了一个完整的 Kubernetes 部署配置。配置包括：

- 带有安全上下文的 Deployment
- 暴露所有必要端口的服务
- 用于外部访问的 Ingress

### Kubernetes 部署的前提条件

- 一个正在运行的 Kubernetes 集群
- 配置好连接到集群的 kubectl
- 已安装的入口控制器（例如 nginx-ingress）

### 部署到 Kubernetes

1. **在 [k8s-deployment.yaml](k8s-deployment.yaml) 中更新配置**：
   - 将 `https://example.com/subscription` 替换为您的实际订阅 URL
   - 根据需要调整 [FILTER_KEYWORDS]指定基于地区的过滤器（例如 "HK" 表示香港节点）
   - 更新 Ingress 部分中的主机名（例如 `mihomo.example.com`、`proxy.example.com`）

2. **部署资源**：

```bash
kubectl apply -f k8s-deployment.yaml
```

3. **验证部署**：

```bash
kubectl get pods -l app=mihomo
kubectl get svc -l app=mihomo
kubectl get ingress mihomo-ingress
```

### 访问服务

部署后，您可以访问：
- 通过 Ingress 中指定的主机名访问 API 和仪表板（端口 9090）
- 通过代理主机名访问代理服务（端口 10801）

### Kubernetes 中的安全功能

Kubernetes 部署包括安全最佳实践：
- 以非 root 用户身份运行
- 删除所有功能
- 使用只读根文件系统
- 使用安全上下文限制访问

## 贡献

欢迎贡献！请随时提交 Pull Request。