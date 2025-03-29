# Using K8sMed with Gemma: A Step-by-Step Guide

This guide provides detailed instructions for integrating Gemma with K8sMed, based on hands-on experience. It covers both local setup and Kubernetes deployment with ngrok.

## Overview

Gemma is Google's family of lightweight, state-of-the-art open models designed for various tasks, including troubleshooting and technical assistance. This guide demonstrates how to use Gemma with K8sMed for Kubernetes troubleshooting.

## Prerequisites

- K8sMed installed (locally or in Kubernetes)
- Access to a Kubernetes cluster
- A LocalAI/Ollama installation with Gemma model support
- (Optional) ngrok for exposing local models to Kubernetes

## Local Setup

### Step 1: Install and Configure Gemma

1. Set up a LocalAI implementation that supports Gemma. Options include:
   - [Ollama](https://ollama.ai/) (recommended for simplicity)
   - [LocalAI](https://github.com/go-skynet/LocalAI)
   - Custom OpenAI-compatible API server

2. Download and configure the Gemma model:

   ```bash
   # Using Ollama
   ollama pull gemma:3b-instruct
   ```

3. Start your LocalAI server:

   ```bash
   # For Ollama, it typically runs as a service
   # For other implementations, follow their specific instructions
   ```

4. Verify the model is working:

   ```bash
   # With Ollama
   ollama run gemma:3b-instruct "What is Kubernetes?"
   ```

### Step 2: Configure K8sMed for Local Use

Set the required environment variables:

```bash
export K8SMED_AI_PROVIDER="localai"
export K8SMED_AI_MODEL="gemma:3b-instruct"  # Match the model name in your LocalAI server
export K8SMED_AI_ENDPOINT="http://localhost:11434/api/chat"  # Adjust based on your setup
```

### Step 3: Test the Integration

Run a basic analysis:

```bash
kubectl-k8smed analyze "why would a pod have ImagePullBackOff status" --explain
```

## Kubernetes Deployment with ngrok

### Step 1: Expose LocalAI with ngrok

1. Install ngrok if not already installed:

   ```bash
   # macOS with Homebrew
   brew install ngrok
   
   # Linux
   curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | \
     sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null && \
     echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | \
     sudo tee /etc/apt/sources.list.d/ngrok.list && \
     sudo apt update && sudo apt install ngrok
   ```

2. Sign up for an ngrok account at [ngrok.com](https://ngrok.com/) if you don't have one

3. Authenticate ngrok:

   ```bash
   ngrok config add-authtoken YOUR_AUTH_TOKEN
   ```

4. Start ngrok and expose your LocalAI endpoint:

   ```bash
   # For Ollama (default port 11434)
   ngrok http 11434
   ```

5. Note the HTTPS URL provided by ngrok (e.g., `https://abcd-123-456-789-10.ngrok-free.app`)

### Step 2: Configure K8sMed in Kubernetes

1. Create a namespace for K8sMed if it doesn't exist:

   ```bash
   kubectl apply -f deploy/manifests/namespace.yaml
   ```

2. Create/update the ConfigMap:

   ```bash
   cat > deploy/manifests/configmap.yaml << EOF
   apiVersion: v1
   data:
     ai_endpoint: "https://your-ngrok-url.ngrok-free.app/v1/chat/completions"
     ai_model: "gemma-3-4b-it"
     ai_provider: "localai"
   kind: ConfigMap
   metadata:
     name: k8smed-config
     namespace: k8smed-system
   EOF
   ```

   > **Important:** Make sure to include the full path to the chat completions endpoint (`/v1/chat/completions`). The exact model name may vary based on your LocalAI implementation.

3. Create a placeholder Secret (even though LocalAI might not need an API key):

   ```bash
   kubectl create secret generic k8smed-secrets \
     --namespace=k8smed-system \
     --from-literal=openai_api_key=placeholder \
     --dry-run=client -o yaml > deploy/manifests/secret.yaml
   
   kubectl apply -f deploy/manifests/secret.yaml
   ```

4. Apply RBAC configurations:

   ```bash
   kubectl apply -f deploy/manifests/rbac.yaml
   ```

5. Deploy K8sMed:

   ```bash
   kubectl apply -f deploy/manifests/deployment.yaml
   ```

### Step 3: Test the Integration

1. Verify the pod is running:

   ```bash
   kubectl get pods -n k8smed-system
   ```

2. Test K8sMed with Gemma:

   ```bash
   # Get the pod name
   K8SMED_POD=$(kubectl get pods -n k8smed-system -o jsonpath='{.items[0].metadata.name}')
   
   # Run a test query
   kubectl exec -it -n k8smed-system $K8SMED_POD -- kubectl-k8smed analyze "pod test-pod has ImagePullBackOff"
   ```

## Troubleshooting

### Common Issues with Gemma Integration

1. **Error: "connection refused"**
   - Check that your LocalAI server is running
   - Verify the endpoint URL is correct

2. **Error: "no completions returned"**
   - Check that the endpoint path includes `/v1/chat/completions`
   - Verify the model name matches what's available in your LocalAI service

3. **Incomplete or truncated responses**
   - This may be due to terminal limitations; try redirecting output to a file
   - Adjust the max tokens parameter if supported by your LocalAI implementation

4. **ngrok connection issues**
   - Verify the ngrok tunnel is active
   - Check bandwidth limits on free ngrok accounts

## Performance Optimization

Gemma performs best with:

1. **Clear, specific queries** - "Why is my pod in CrashLoopBackOff?" vs. "pod problems"
2. **Using the `--explain` flag** - Provides more context for the model to work with
3. **Gemma-IT variants** - Models fine-tuned for technical understanding perform better

## Example Queries

Here are some effective queries that work well with Gemma and K8sMed:

```bash
# Specific issue analysis
kubectl-k8smed analyze "pod mypod has ImagePullBackOff status" --explain

# General Kubernetes concepts
kubectl-k8smed analyze "explain how Kubernetes networking works with services and pods"

# Troubleshooting workflows
kubectl-k8smed analyze "how to debug a pod that won't start" --explain
```

## Real-World Example

Here's a complete example of deploying and using K8sMed with Gemma, based on our testing:

1. Start local Gemma model with Ollama:
   ```bash
   ollama serve
   ```

2. Expose with ngrok:
   ```bash
   ngrok http 11434
   ```

3. Deploy K8sMed with the ngrok URL:
   ```bash
   # Update ConfigMap with ngrok URL
   kubectl apply -f deploy/manifests/namespace.yaml
   kubectl apply -f deploy/manifests/configmap.yaml
   kubectl apply -f deploy/manifests/secret.yaml
   kubectl apply -f deploy/manifests/rbac.yaml
   kubectl apply -f deploy/manifests/deployment.yaml
   ```

4. Create a test pod with an issue:
   ```bash
   kubectl apply -f - <<EOF
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-error-pod
     namespace: default
   spec:
     containers:
     - name: error-container
       image: nonexistentimage:latest
       imagePullPolicy: IfNotPresent
     restartPolicy: Never
   EOF
   ```

5. Use K8sMed to analyze the issue:
   ```bash
   K8SMED_POD=$(kubectl get pods -n k8smed-system -o jsonpath='{.items[0].metadata.name}')
   kubectl exec -it -n k8smed-system $K8SMED_POD -- kubectl-k8smed analyze "pod test-error-pod has ImagePullBackOff"
   ```

6. Review the analysis and remediation steps provided by Gemma. 