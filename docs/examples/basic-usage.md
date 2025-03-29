# K8sMed Basic Usage Examples

This document provides some examples of how to use K8sMed to troubleshoot common Kubernetes issues.

## Installation

First, install the K8sMed kubectl plugin:

```bash
# Download the appropriate binary for your platform
# For macOS (Intel):
curl -L https://github.com/k8smed/k8smed/releases/download/v0.1.0-alpha/kubectl-k8smed_darwin_amd64 -o kubectl-k8smed
chmod +x kubectl-k8smed
sudo mv kubectl-k8smed /usr/local/bin/

# For macOS (Apple Silicon):
curl -L https://github.com/k8smed/k8smed/releases/download/v0.1.0-alpha/kubectl-k8smed_darwin_arm64 -o kubectl-k8smed
chmod +x kubectl-k8smed
sudo mv kubectl-k8smed /usr/local/bin/

# For Linux:
curl -L https://github.com/k8smed/k8smed/releases/download/v0.1.0-alpha/kubectl-k8smed_linux_amd64 -o kubectl-k8smed
chmod +x kubectl-k8smed
sudo mv kubectl-k8smed /usr/local/bin/

# Verify installation
kubectl k8smed version
```

## Configuration

K8sMed uses your existing kubectl configuration by default. You can also set the following environment variables:

```bash
# Set your OpenAI API key (required when using OpenAI as the AI provider)
export OPENAI_API_KEY=your-api-key

# Optional configuration overrides
export K8SMED_AI_PROVIDER=openai       # Options: openai, localai
export K8SMED_AI_MODEL=gpt-3.5-turbo   # Default OpenAI model
export K8SMED_AI_ENDPOINT=             # Custom endpoint for LocalAI
export K8SMED_ANONYMIZE_DEFAULT=false  # Whether to anonymize data by default
```

## Basic Examples

### Analyzing a Pod with Issues

```bash
# Check why a pod is having trouble
kubectl k8smed analyze "why is my pod myapp-pod in CrashLoopBackOff"

# Investigate container termination
kubectl k8smed analyze "why is container terminating in pod myapp-pod"

# Analyze with more detail
kubectl k8smed analyze --explain "diagnose why pod myapp-pod won't start"

# Anonymize sensitive information
kubectl k8smed analyze --anonymize "check why pod in namespace customer-data is failing"
```

### Analyzing Multiple Resources

```bash
# Investigate why a deployment isn't scaling correctly
kubectl k8smed analyze "why is my deployment not scaling to 3 replicas"

# Troubleshoot service connectivity
kubectl k8smed analyze "why can't pod myapp-pod connect to my-service"
```

### Using Interactive Mode (Coming Soon)

```bash
# Start an interactive troubleshooting session
kubectl k8smed interactive
```

## Example Output

When you run an analysis query, K8sMed will:

1. Collect relevant data from your Kubernetes cluster
2. Analyze the data using AI to identify issues
3. Provide a detailed explanation and remediation steps

Example output for a pod in CrashLoopBackOff:

```
Analyzing query: why is my pod myapp-pod in CrashLoopBackOff
Explain: false, Anonymize: false

The pod 'myapp-pod' is in CrashLoopBackOff state because the application is exiting with a non-zero status code.
According to the logs, there's a configuration error: the application is trying to connect to a database at
'db-service.default.svc.cluster.local' but it's unable to establish a connection.

Here are some remediation steps:

1. Verify that the database service exists and is running:
   ```
   kubectl get svc db-service -n default
   kubectl describe svc db-service -n default
   ```

2. Check if there are any network policies preventing the connection:
   ```
   kubectl get networkpolicies -n default
   ```

3. Make sure the database credentials in the pod's environment variables or config are correct.

4. If the database is in a different namespace, ensure the full service DNS name is correctly specified.
```

## Command Reference

K8sMed provides the following commands:

- `kubectl k8smed analyze [query]`: Analyze Kubernetes resources and provide insights
  - `--explain` or `-e`: Show detailed explanation including AI model used
  - `--anonymize` or `-a`: Anonymize sensitive information in the query

- `kubectl k8smed interactive`: Start an interactive troubleshooting session (coming soon)

- `kubectl k8smed version`: Display version information

- `kubectl k8smed config`: Display current configuration
