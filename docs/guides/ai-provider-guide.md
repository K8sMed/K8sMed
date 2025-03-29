# AI Provider Guide for K8sMed

This guide provides detailed instructions for configuring K8sMed with various AI providers, including step-by-step examples for using Gemma and other open-source models.

## Table of Contents

- [OpenAI Integration](#openai-integration)
- [LocalAI Integration](#localai-integration)
- [Using Gemma with K8sMed](#using-gemma-with-k8smed)
- [Other Open-Source Models](#other-open-source-models)
- [Kubernetes Integration with Local Models](#kubernetes-integration-with-local-models)
- [Performance Considerations](#performance-considerations)

## OpenAI Integration

OpenAI's GPT models provide high-quality analysis capabilities for Kubernetes troubleshooting.

### Configuration

1. Obtain an API key from [OpenAI Platform](https://platform.openai.com)
2. Configure K8sMed using environment variables:

```bash
export OPENAI_API_KEY="your-api-key-here"
export K8SMED_AI_PROVIDER="openai"
export K8SMED_AI_MODEL="gpt-4"  # or gpt-3.5-turbo
```

### Example Usage

```bash
kubectl-k8smed analyze "why is my pod in CrashLoopBackOff" --explain
```

## LocalAI Integration

LocalAI is an API-compatible alternative to OpenAI that can run open-source models locally.

### Setting Up LocalAI

1. Install and run LocalAI following the instructions at [LocalAI GitHub](https://github.com/go-skynet/LocalAI)
2. Configure K8sMed to use your LocalAI endpoint:

```bash
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_ENDPOINT="http://localhost:8080/v1"  # Adjust port as needed
export K8SMED_AI_MODEL="your-model-name"  # As configured in LocalAI
```

### Example Usage

```bash
kubectl-k8smed analyze pod mypod --explain
```

## Using Gemma with K8sMed

Gemma is Google's lightweight open-source model that can be used with K8sMed through LocalAI or similar providers.

### Setting Up Gemma

1. Install a LocalAI implementation that supports Gemma (e.g., [Ollama](https://ollama.ai/))
2. Pull the Gemma model:

```bash
# Using Ollama
ollama pull gemma:3b
```

3. Configure K8sMed:

```bash
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_ENDPOINT="http://localhost:11434/api/chat"  # Ollama endpoint
export K8SMED_AI_MODEL="gemma:3b"  # Model name as configured
```

### Using Gemma in Kubernetes with ngrok

To use a locally hosted Gemma model with K8sMed deployed in Kubernetes:

1. Run your local Gemma instance (via Ollama, LocalAI, etc.)
2. Expose it using ngrok:

```bash
ngrok http 11434  # Port depends on your implementation
```

3. Create/update ConfigMap in your Kubernetes cluster:

```yaml
apiVersion: v1
data:
  ai_endpoint: "https://your-ngrok-url.ngrok-free.app/v1/chat/completions"  # Note the full path
  ai_model: "gemma-3-4b-it"  # Exact model name is implementation-specific
  ai_provider: "localai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
```

4. Deploy K8sMed with this ConfigMap:

```bash
kubectl apply -f deploy/manifests/
```

5. Test the integration:

```bash
K8SMED_POD=$(kubectl get pods -n k8smed-system -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it -n k8smed-system $K8SMED_POD -- kubectl-k8smed analyze "pod test-pod has ImagePullBackOff"
```

## Other Open-Source Models

K8sMed can work with various open-source models as long as they're exposed through an OpenAI-compatible API.

### Llama 2 / Llama 3

```bash
# Using Ollama with Llama 2 (7B)
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_ENDPOINT="http://localhost:11434/api/chat"
export K8SMED_AI_MODEL="llama2"
```

### Mistral

```bash
# Using Ollama with Mistral
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_ENDPOINT="http://localhost:11434/api/chat"
export K8SMED_AI_MODEL="mistral"
```

## Kubernetes Integration with Local Models

When using K8sMed in Kubernetes with locally hosted models, you have several options:

### Option 1: Use ngrok (As Shown in the Gemma Example)

Pros:
- Quick to set up
- Works with any Kubernetes cluster, including remote ones
- Minimal configuration required

Cons:
- Relies on external service (ngrok)
- URLs change unless using paid ngrok account
- Potential bandwidth limitations

### Option 2: Deploy the Model in Kubernetes

1. Deploy your AI model (e.g., Gemma) within the Kubernetes cluster using a solution like [LocalAI on Kubernetes](https://localai.io/basics/getting_started/kubernetes/)
2. Configure K8sMed to use the in-cluster service:

```yaml
apiVersion: v1
data:
  ai_endpoint: "http://localai-service.ai-namespace.svc.cluster.local/v1/chat/completions"
  ai_model: "gemma-3-4b-it"
  ai_provider: "localai"
kind: ConfigMap
metadata:
  name: k8smed-config
  namespace: k8smed-system
```

Pros:
- More reliable and stable
- No dependency on external services
- Better security (data stays within cluster)

Cons:
- Requires more resources in the cluster
- More complex setup

### Option 3: Use Port Forwarding

For development/testing purposes:

1. Set up port forwarding from your local machine to a pod in the cluster:

```bash
kubectl port-forward -n k8smed-system deployment/k8smed 8080:8080
```

2. Configure your local LocalAI to listen on the forwarded port
3. K8sMed in the cluster can then access the forwarded port

## Performance Considerations

Different models offer different performance characteristics:

| Model | Size | Analysis Quality | Speed | Resource Requirements |
|-------|------|------------------|-------|------------------------|
| GPT-4 | Large | Excellent | Moderate | API-only (cloud) |
| GPT-3.5-Turbo | Medium | Very Good | Fast | API-only (cloud) |
| Gemma-3-4B-IT | Small | Good for common issues | Fast | 4-8GB RAM |
| Llama-2-7B | Medium | Good | Moderate | 8-16GB RAM |
| Mistral-7B | Medium | Good | Moderate | 8-16GB RAM |

### Tips for Better Results with Smaller Models

1. Be more specific in your queries
2. Use the `--explain` flag to get more context for the model
3. For complex issues, consider using the built-in analyzers without LLM integration
4. If available, use models fine-tuned for technical/code understanding (like Gemma-IT)
