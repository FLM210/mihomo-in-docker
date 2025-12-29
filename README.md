# Mihomo in Docker

## Introduction

This project aims to containerize Mihomo (a high-performance fork of Clash) using Docker, simplifying its deployment process and supporting automatic configuration file generation for easy operation in Kubernetes or common container environments.

### Target Users:
- Technical users who need to use proxy tools for traffic management
- Operations or developers who want to deploy Mihomo in container environments (such as Docker, Kubernetes)

### Core Problems Solved:
- Simplifying installation and configuration of Mihomo across different architectures
- Achieving automatic configuration generation to reduce complexity of manually editing YAML files
- Providing reusable, lightweight images for production deployment

### System Features:
- Automatically detects host architecture and downloads the corresponding version of Mihomo binary
- Uses Go-written configuration generator to generate Mihomo configuration files based on environment variables
- Supports automatic operation of Mihomo service through startup scripts
- Exposes common ports (API, Proxy, UI) within containers

## Features

- **Architecture Adaptation**: Automatically identifies system architecture (amd64/arm64) during build and downloads corresponding Mihomo binary
- **Automatic Configuration Generation**: Utilizes Go program to parse environment variables or templates to generate config.yaml
- **Lightweight Operation**: Built on Alpine Linux for final image, containing only necessary components
- **Multi-stage Build**: Separates build, download, and runtime stages to reduce image size
- No root permission dependency, suitable for security-enhanced scenarios
- Supports CI/CD pipeline builds and deployments
- Can be integrated into Kubernetes deployment

## Quick Start

### Prerequisites
- Docker 20+

### Building the Image

```bash
# Clone the project
git clone https://github.com/flm210/mihomo-in-docker.git
cd mihomo-in-docker

# Build the image
docker build -t mihomo-in-docker .
```

### Running the Container

```bash
# Run container with example configuration
docker run -p 9090:9090 -p 10801:10801 \
  -e "SUBSCRIPTION_URL=https://example.com/subscription" \
  -e "FILTER_KEYWORDS=HK,SG" \
  mihomo-in-docker
```

### Environment Variables

- `SUBSCRIPTION_URL`: Subscription URL for proxy nodes
- `FILTER_KEYWORDS`: Keywords to filter proxy nodes (comma separated)
- `OUTPUT_FILE`: Output filename (default: config.yaml)
- `OUTPUT_FORMAT`: Output format (yaml or json, default: yaml)

## Configuration

The configuration generator supports:
- Subscription URL processing
- Keyword filtering for nodes
- YAML/JSON output formats
- Support for various proxy protocols (VMess, VLESS, Trojan, Shadowsocks, etc.)

## Ports

- `9090`: External controller API
- `10801`: Mixed proxy port (HTTP and SOCKS)

## Kubernetes Deployment

This project provides a complete Kubernetes deployment configuration in the [k8s-deployment.yaml](k8s-deployment.yaml) file. The configuration includes:

- A Deployment with security contexts
- A Service exposing all necessary ports
- An Ingress for external access

### Prerequisites for Kubernetes Deployment

- A running Kubernetes cluster
- kubectl configured to connect to your cluster
- An ingress controller installed (e.g., nginx-ingress)

### Deploying to Kubernetes

1. **Update the configuration** in [k8s-deployment.yaml](k8s-deployment.yaml):
   - Replace `https://example.com/subscription` with your actual subscription URL
   - Adjust the `FILTER_KEYWORDS` as needed
   - Update the hostnames in the Ingress section (e.g., `mihomo.example.com`, `proxy.example.com`)

2. **Deploy the resources**:

```bash
kubectl apply -f k8s-deployment.yaml
```

3. **Verify the deployment**:

```bash
kubectl get pods -l app=mihomo
kubectl get svc -l app=mihomo
kubectl get ingress mihomo-ingress
```

### Accessing the Service

After deployment, you can access:
- The API and dashboard via the hostname specified in the Ingress (port 9090)
- The proxy service via the proxy hostname (port 10801)

### Security Features in Kubernetes

The Kubernetes deployment includes security best practices:
- Runs as a non-root user
- Drops all capabilities
- Uses read-only root filesystem
- Restricts access with security contexts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
