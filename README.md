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
- **Intelligent Node Filtering**: Automatically filters proxy nodes based on keywords, supporting automatic failover groups for specific regions
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
- `FILTER_KEYWORDS`: Keywords to filter proxy nodes (comma separated). The program automatically creates failover groups based on these keywords.
- `OUTPUT_FILE`: Output filename (default: config.yaml)
- `OUTPUT_FORMAT`: Output format (yaml or json, default: yaml)

## Configuration

The configuration generator supports:
- Subscription URL processing
- Keyword filtering for nodes with automatic failover group creation
- YAML/JSON output formats
- Support for various proxy protocols (VMess, VLESS, Trojan, Shadowsocks, etc.)
- Automatic generation of proxy groups with auto-failover capabilities based on filter keywords

## Advanced Usage: Automatic Node Switching

The program automatically creates a proxy group with automatic failover capabilities based on the keywords provided in the `FILTER_KEYWORDS` environment variable. For example:

- If you set `FILTER_KEYWORDS=HK`, the program will filter all nodes containing "HK" (Hong Kong) in their names and create an automatic failover group.
- If you set `FILTER_KEYWORDS=HK,SG,JP`, it will create a single group containing all nodes that match any of these keywords (Hong Kong, Singapore, Japan), with automatic failover between them.

### Kubernetes Example: Region-Specific Deployment

If you want to deploy a pod in K8S specifically for HK (Hong Kong) nodes, simply set the filter to include the corresponding node names for Hong Kong in your provider's subscription:

```yaml
env:
- name: SUBSCRIPTION_URL
  value: "https://your-provider.com/subscription"
- name: FILTER_KEYWORDS
  value: "Hong Kong,HK,hk"  # Keywords that match your provider's Hong Kong nodes
```

This will automatically filter and create failover groups containing only Hong Kong nodes, with automatic switching between them based on connectivity.

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
   - Adjust the `FILTER_KEYWORDS` as needed to specify region-based filters (e.g., "HK" for Hong Kong nodes)
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