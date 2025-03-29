# K8sMed: AI-Powered Kubernetes First Responder

K8sMed is an open-source, AI-powered troubleshooting assistant designed to act as a first responder for Kubernetes clusters. By analyzing cluster logs, events, and metrics, K8sMed leverages Large Language Models (LLMs) to diagnose issues, provide natural language explanations, and generate actionable remediation commands—all through a simple kubectl plugin.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/k8smed/k8smed)](https://goreportcard.com/report/github.com/k8smed/k8smed)

---

## Table of Contents

- [Project Overview](#project-overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Documentation](#documentation)
- [Privacy & Security](#privacy--security)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [License](#license)
- [Contact](#contact)

---

## Project Overview

K8sMed helps Kubernetes administrators and developers troubleshoot issues faster by acting as a "first responder" for cluster problems. The tool analyzes Kubernetes resources, interprets error messages, and generates clear explanations and remediation steps using AI.

### Why K8sMed?

Kubernetes environments are complex, and troubleshooting issues often requires deep expertise. K8sMed reduces mean time to resolution (MTTR) by:

- Providing instant analysis of Kubernetes resources
- Explaining problems in clear, human-readable language
- Generating actionable remediation commands
- Supporting both beginners and experienced Kubernetes users

### Goals

- **Rapid Diagnosis**: Quickly identify issues across different Kubernetes resources
- **Actionable Insights**: Generate precise kubectl commands and YAML patches
- **Privacy First**: Anonymize sensitive data with built-in protection
- **Flexibility**: Support both cloud-based and local AI models
- **Seamless Experience**: Simple kubectl plugin interface

---

## Key Features

- **Comprehensive Analysis**: 
  Analyze pods, deployments, services, and other Kubernetes resources for common issues

- **Multi-Provider AI Support**: 
  Use OpenAI models (GPT-3.5/4) or local alternatives (LocalAI, Ollama) for analysis

- **Anonymization**: 
  Protect sensitive information with built-in data anonymization

- **Actionable Commands**: 
  Receive ready-to-use kubectl commands for quick remediation

- **Context-Aware Analysis**: 
  Intelligent understanding of Kubernetes concepts and relationships between resources

- **Local-First Architecture**: 
  Run entirely in your environment without requiring external services

---

## Architecture

K8sMed follows a modular architecture with these key components:

1. **Resource Collection**: Gathers information about Kubernetes resources
2. **Problem Analysis**: Examines resources for potential issues
3. **AI Processing**: Sends anonymized data to AI for interpretation
4. **Remediation Generation**: Creates actionable commands to fix issues

The tool runs as a kubectl plugin, requiring only kubectl access to your cluster.

---

## Installation

### Prerequisites

- Kubernetes cluster with kubectl access
- Go 1.21+ (for building from source)
- Access to an AI provider (OpenAI account or local AI setup)

### Quick Install

```bash
# Clone the repository
git clone https://github.com/k8smed/k8smed.git
cd k8smed

# Build the binary
make build

# Install the kubectl plugin
make install
```

### Verify Installation

```bash
kubectl k8smed version
```

---

## Usage

### Basic Analysis

```bash
# Analyze a pod with issues
kubectl k8smed analyze pod my-pod-name --namespace default

# Analyze with a specific question
kubectl k8smed query "Why is my pod in CrashLoopBackOff state?"

# Get remediation suggestions
kubectl k8smed suggest pod my-pod-name
```

### Anonymization

```bash
# Enable anonymization to protect sensitive data
kubectl k8smed analyze pod my-pod-name --anonymize
```

### Using Different AI Providers

```bash
# Use OpenAI
export OPENAI_API_KEY=your_api_key
kubectl k8smed analyze pod my-pod-name

# Use LocalAI
export K8SMED_AI_PROVIDER=localai
export K8SMED_AI_ENDPOINT=http://localhost:8080
kubectl k8smed analyze pod my-pod-name
```

---

## Configuration

K8sMed can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `K8SMED_AI_PROVIDER` | AI provider (openai, localai) | openai |
| `K8SMED_AI_MODEL` | Model name to use | gpt-3.5-turbo |
| `K8SMED_AI_ENDPOINT` | API endpoint for LocalAI | - |
| `K8SMED_ANONYMIZE_DEFAULT` | Enable anonymization by default | false |
| `K8SMED_OUTPUT_FORMAT` | Output format (text, json) | text |
| `OPENAI_API_KEY` | OpenAI API key | - |

---

## Examples

### Diagnosing a Pod in CrashLoopBackOff

```bash
kubectl k8smed analyze pod nginx-deployment-665d87f687-abcde
```

Output:
```
📋 K8sMed Analysis:
🔍 Pod nginx-deployment-665d87f687-abcde is in CrashLoopBackOff state

📝 Description:
The container is repeatedly crashing after startup. The exit code 1 suggests
the application is exiting with an error.

✅ Remediation:
1. Check container logs: kubectl logs nginx-deployment-665d87f687-abcde
2. Verify environment variables are set correctly
3. Check if the application can connect to required services
4. Inspect the startup command for errors

💻 Remediation Commands:
kubectl logs nginx-deployment-665d87f687-abcde
kubectl describe pod nginx-deployment-665d87f687-abcde
```

### More Examples

Check out our [examples directory](docs/examples/) for more use cases, including:
- Troubleshooting ImagePullBackOff errors
- Fixing service connectivity issues
- Resolving permission problems
- Debugging deployment rollout issues

---

## Documentation

Detailed documentation is available in the [docs directory](docs/):

- [Deployment Guide](docs/guides/deployment-guide.md)
- [AI Provider Guide](docs/guides/ai-provider-guide.md)
- [Developer Guide](docs/guides/developer-guide.md) - For contributors and developers
- [Gemma Integration](docs/examples/gemma-integration.md)
- [Basic Usage Guide](docs/examples/basic-usage.md)

---

## Privacy & Security

K8sMed takes privacy seriously:

- **Anonymization**: Built-in anonymization replaces sensitive information before sending to AI providers
- **Local AI Support**: Run entirely in your environment with LocalAI or similar tools
- **Minimal Permissions**: Requires only read access to your cluster
- **No Data Storage**: K8sMed doesn't store any cluster information

For sensitive environments, we recommend:
1. Using the `--anonymize` flag
2. Setting up a local AI model
3. Reviewing prompts sent to the AI

---

## Contributing

We welcome contributions to K8sMed! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up your development environment
- Running tests
- Submitting pull requests
- Our code of conduct

For technical details about the codebase, architecture, and development workflows, check out our [Developer Guide](docs/guides/developer-guide.md).

---

## Roadmap

### Current Focus (Q2-Q3 2025)
- Expanding resource analyzers beyond pods
  - Implementing dedicated analyzers for Deployments, Services, and StatefulSets
  - Creating specialized analyzers for Ingress resources and NetworkPolicies
  - Adding support for Custom Resource analysis
- Improving detection accuracy for common Kubernetes issues
  - Building a comprehensive database of error patterns and solutions
  - Enhancing context-awareness for multi-resource related problems
  - Developing specialized analyzers for networking and storage issues
- Enhancing remediation suggestions for complex scenarios
  - Providing tiered remediation options (quick fixes vs. root cause solutions)
  - Supporting YAML patch generation for configuration fixes
  - Adding simulation capabilities to preview remediation effects
- Adding support for more AI providers and models
  - Implementing dedicated connectors for Anthropic Claude and Google Gemini
  - Optimizing prompts for different model capabilities
  - Creating an abstract provider interface for easy extensions

### Next Steps (Q3-Q4 2025)
- **Interactive Mode Development**
  - Building a conversational CLI interface for multi-turn troubleshooting
  - Implementing session context management for follow-up questions
  - Adding support for exploration-based problem solving with AI guidance
- **Plugin Ecosystem**
  - Creating an extension system for community-contributed analyzers
  - Developing a plugin marketplace or registry
  - Publishing a plugin development guide with examples
- **Performance Optimizations**
  - Implementing parallel resource collection and analysis
  - Adding result caching for faster repeat analysis
  - Optimizing token usage for more efficient AI interactions
- **Integration Capabilities**
  - Building connectors for popular monitoring systems (Prometheus, Grafana)
  - Developing webhook support for automated analysis triggering
  - Creating integration points for CI/CD systems

### Future Plans (2024+)
- Operator mode for continuous monitoring
  - Custom resource definitions for scheduled analysis
  - Alert integration for automatic problem detection
  - Historical analysis storage and trending
- AI training on Kubernetes-specific datasets
  - Creating specialized fine-tuned models for Kubernetes troubleshooting
  - Building synthetic problem datasets for improved accuracy
- Advanced visualization capabilities
  - Resource relationship mapping for complex issues
  - Root cause probability visualizations
  - Remediation impact previews

## Getting Involved

We're actively seeking contributors in the following areas:
1. **Analyzer Development**: Help build analyzers for specific Kubernetes resources
2. **AI Integration**: Assist with implementing new AI provider integrations
3. **Documentation**: Improve guides, examples, and tutorials
4. **Testing**: Create test cases and validation frameworks

If you're interested in contributing, check out our [open issues](https://github.com/k8smed/k8smed/issues) labeled with "good first issue" or "help wanted", or reach out through our contact channels.

---

## License

K8sMed is licensed under the [Apache License 2.0](LICENSE).

---

## Contact

- GitHub Issues: [Submit an issue](https://github.com/k8smed/k8smed/issues)
- Project Lead: [Md Imran](https://github.com/narmidm)

---

K8sMed aims to revolutionize Kubernetes troubleshooting with an AI-powered approach that delivers fast, accurate, and actionable insights. We invite you to try it out, provide feedback, and join our community of contributors!
