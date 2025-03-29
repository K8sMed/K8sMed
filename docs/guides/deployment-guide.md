# K8sMed Deployment Guide

This guide provides detailed instructions for deploying and configuring K8sMed in both local and Kubernetes environments. It covers different AI provider options, configuration parameters, and troubleshooting tips.

## Table of Contents

- [Overview](#overview)
- [Local Installation](#local-installation)
- [Kubernetes Deployment](#kubernetes-deployment)
- [AI Provider Configuration](#ai-provider-configuration)
  - [OpenAI](#openai)
  - [LocalAI/Ollama](#localaiollama)
  - [Self-hosted Models via ngrok](#self-hosted-models-via-ngrok)
- [Common Usage Patterns](#common-usage-patterns)
- [Troubleshooting](#troubleshooting)
- [FAQs](#faqs)

## Overview

K8sMed is an AI-powered Kubernetes troubleshooting assistant designed to diagnose issues, provide natural language explanations, and generate actionable remediation commands. It can be deployed locally as a CLI tool or within a Kubernetes cluster.

## Local Installation

### Prerequisites

- Go 1.21+ installed
- `kubectl` installed and configured with access to a Kubernetes cluster
- API key for your chosen LLM provider (if using OpenAI or similar)

### Installation from Source

1. Clone the repository:

```bash
git clone https://github.com/k8smed/k8smed.git
cd k8smed
```

2. Build the binary:

```bash
make build
```

This will create a binary at `bin/kubectl-k8smed`.

3. Make the binary available in your PATH:

```bash
# Option 1: Move to a location in your PATH
sudo mv bin/kubectl-k8smed /usr/local/bin/

# Option 2: Add the bin directory to your PATH
export PATH=$PATH:$(pwd)/bin
```

### Configuration

K8sMed can be configured using environment variables:

```bash
# For OpenAI
export OPENAI_API_KEY="your-api-key"
export K8SMED_AI_PROVIDER="openai"
export K8SMED_AI_MODEL="gpt-4"    # Options: gpt-3.5-turbo, gpt-4, etc.

# For LocalAI/Ollama
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_MODEL="llama2"   # Or your model name
export K8SMED_AI_ENDPOINT="http://localhost:11434/v1"  # Your LocalAI/Ollama endpoint
```

### Basic Usage

```bash
# Analyze a specific pod
kubectl-k8smed analyze pod mypod -n mynamespace

# Analyze with detailed explanations
kubectl-k8smed analyze "pod mypod has CrashLoopBackOff" --explain

# Anonymize sensitive information
kubectl-k8smed analyze deployment myapp -n mynamespace --anonymize
```

## Kubernetes Deployment

K8sMed can be deployed within your Kubernetes cluster to provide centralized troubleshooting capabilities.

### Prerequisites

- Kubernetes cluster (1.19+)
- kubectl access with admin privileges
- Docker (for building the image)

### Deployment Steps

1. Build and push the Docker image:

```bash
# Build the image
docker build -t yourusername/k8smed:latest .

# Push to a registry (if deploying to a remote cluster)
docker push yourusername/k8smed:latest
```

2. If using a local cluster like kind or minikube, load the image directly:

```bash
# For kind
kind load docker-image k8smed:latest

# For minikube
minikube image load k8smed:latest
```

3. Configure the deployment resources:

Create the necessary ConfigMap with your AI provider settings:

```bash
cat > deploy/manifests/configmap.yaml << EOF
apiVersion: v1
data:
  ai_endpoint: "https://api.openai.com/v1"  # For OpenAI
  ai_model: "gpt-4"
  ai_provider: "openai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
EOF
```

Create a Secret for your API key:

```bash
kubectl create secret generic k8smed-secrets \
  --namespace=k8smed-system \
  --from-literal=openai_api_key=your-api-key \
  --dry-run=client -o yaml > deploy/manifests/secret.yaml
```

4. Deploy K8sMed to your Kubernetes cluster:

```bash
# Create namespace, RBAC, ConfigMap, Secret and Deployment
kubectl apply -f deploy/manifests/
```

5. Verify the deployment:

```bash
kubectl get pods -n k8smed-system
```

### Using K8sMed in Kubernetes

Once deployed, you can use K8sMed by executing commands in the pod:

```bash
# Get the pod name
K8SMED_POD=$(kubectl get pods -n k8smed-system -o jsonpath='{.items[0].metadata.name}')

# Analyze a pod
kubectl exec -it -n k8smed-system $K8SMED_POD -- kubectl-k8smed analyze pod problematic-pod
```

## AI Provider Configuration

K8sMed supports multiple AI providers, each with their own configuration requirements.

### OpenAI

For OpenAI (GPT-3.5, GPT-4), you need:

1. An API key from OpenAI
2. Configuration in environment variables or ConfigMap:

```yaml
# For local usage (environment variables)
export OPENAI_API_KEY="your-api-key"
export K8SMED_AI_PROVIDER="openai"
export K8SMED_AI_MODEL="gpt-4"  # Options: gpt-3.5-turbo, gpt-4, etc.

# For Kubernetes (ConfigMap)
apiVersion: v1
data:
  ai_endpoint: "https://api.openai.com/v1"
  ai_model: "gpt-4"
  ai_provider: "openai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
```

And a Secret for the API key:

```yaml
apiVersion: v1
data:
  openai_api_key: "base64-encoded-api-key"
kind: Secret
metadata:
  name: k8smed-secrets
  namespace: k8smed-system
```

### LocalAI/Ollama

For LocalAI or Ollama (self-hosted models):

1. Configure with the correct endpoint and model name:

```yaml
# For local usage (environment variables)
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_MODEL="llama2"  # Or your model name
export K8SMED_AI_ENDPOINT="http://localhost:11434/v1"  # Your LocalAI/Ollama endpoint

# For Kubernetes (ConfigMap)
apiVersion: v1
data:
  ai_endpoint: "http://localhost:11434/v1"
  ai_model: "llama2"
  ai_provider: "localai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
```

### Self-hosted Models via ngrok

For Kubernetes access to locally hosted models, use ngrok to expose your LocalAI endpoint:

1. Start your LocalAI/Ollama server locally

2. Expose it with ngrok:

```bash
ngrok http 11434  # Adjust port as needed
```

3. Update the ConfigMap with the ngrok URL:

```yaml
apiVersion: v1
data:
  ai_endpoint: "https://your-ngrok-url.ngrok-free.app/v1/chat/completions"
  ai_model: "your-model-name"
  ai_provider: "localai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
```

> **Note:** For LocalAI with OpenAI-compatible API, you must include the full path to the chat completions endpoint. The exact path may vary based on your LocalAI implementation.

## Common Usage Patterns

### Analyzing Specific Resources

```bash
# Analyze a pod
kubectl-k8smed analyze pod mypod

# Analyze a deployment
kubectl-k8smed analyze deployment mydeployment

# Analyze a service
kubectl-k8smed analyze service myservice
```

### Natural Language Queries

```bash
# Ask about a problem with specific symptoms
kubectl-k8smed analyze "why is my pod in CrashLoopBackOff"

# Get remediation steps for a specific issue
kubectl-k8smed analyze "how to fix ImagePullBackOff in my deployment"
```

### Interactive Mode

```bash
# Start an interactive troubleshooting session
kubectl-k8smed interactive
```

## Troubleshooting

### Common Issues

#### LLM API Connection Issues

If you see errors like `Error getting LLM response: connection refused` or `Error getting LLM response: no completions returned`:

1. Check your API endpoint is correct
2. Verify the API key is valid
3. For LocalAI/Ollama, ensure the server is running
4. For ngrok, check the tunnel is active and the URL is correct

#### Image Pull Issues in Kubernetes

If the K8sMed pod has `ErrImagePull` or `ImagePullBackOff` status:

1. Ensure the image exists and is correctly tagged
2. For local clusters, load the image using `kind load` or `minikube image load`
3. For remote clusters, push the image to a registry accessible to your cluster

#### Permission Issues

If you see `Error: forbidden: User "system:serviceaccount:..." cannot get resource` errors:

1. Verify the RBAC configuration is correctly applied
2. Check that the service account has appropriate permissions
3. For cluster-wide access, the ClusterRole and ClusterRoleBinding must be properly configured

## FAQs

### Q: Can I use K8sMed without an LLM?

A: Yes! While the LLM integration provides enhanced analysis and natural language interaction, K8sMed has built-in analyzers that can identify common issues like CrashLoopBackOff, ImagePullBackOff, and resource constraints without requiring an LLM.

### Q: Which LLM performs best for Kubernetes troubleshooting?

A: GPT-4 and similar high-capability models generally provide the most accurate and detailed analysis. However, smaller models like Gemma-3-4B-IT can still provide useful insights for common Kubernetes issues.

### Q: How can I share K8sMed with my team?

A: The recommended approach is to deploy K8sMed in your Kubernetes cluster with appropriate RBAC permissions. This allows multiple team members to use the tool without individual setup.

### Q: Is my Kubernetes data sent to OpenAI?

A: K8sMed collects resource data from your cluster to send to the configured LLM. If you're using OpenAI, this data will be sent to their API. For sensitive environments, we recommend:

1. Using the `--anonymize` flag to remove sensitive information
2. Self-hosting an LLM using LocalAI/Ollama
3. Using the built-in analyzers without LLM integration

### Q: How to contribute to K8sMed?

A: We welcome contributions! Please check our [CONTRIBUTING.md](../../CONTRIBUTING.md) guide and feel free to open issues or pull requests in the GitHub repository. 