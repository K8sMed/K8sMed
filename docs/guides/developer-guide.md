# K8sMed Developer Guide

This guide is intended for developers who want to contribute to K8sMed. It provides detailed information about the project structure, development workflow, and technical architecture.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Architecture Overview](#architecture-overview)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style and Conventions](#code-style-and-conventions)
- [Creating New Analyzers](#creating-new-analyzers)
- [AI Provider Integration](#ai-provider-integration)
- [Debugging Tips](#debugging-tips)
- [Release Process](#release-process)

## Development Environment Setup

### Prerequisites

- Go 1.21+
- Kubernetes cluster (minikube, kind, or a remote cluster)
- kubectl configured to access your cluster
- Git

### Setting Up Your Development Environment

1. **Fork and Clone the Repository**

```bash
# Fork on GitHub first, then clone your fork
git clone https://github.com/YOUR-USERNAME/k8smed.git
cd k8smed

# Add the upstream repository
git remote add upstream https://github.com/k8smed/k8smed.git
```

2. **Install Dependencies**

```bash
# Fetch Go dependencies
go mod download

# Install development tools
make tools
```

3. **Build the Project**

```bash
# Build binary
make build

# Install the kubectl plugin locally for testing
make install
```

4. **Verify Installation**

```bash
kubectl k8smed version
```

## Project Structure

K8sMed is organized into the following directories:

```
.
├── api/              # API definitions and Custom Resource Definitions
├── cmd/              # Command-line interfaces
│   └── kubectl-k8smed/ # The kubectl plugin entry point
├── deploy/           # Deployment manifests and charts
├── docs/             # Documentation
├── examples/         # Example configurations and use cases
├── hack/             # Scripts for development, CI/CD, etc.
├── internal/         # Internal packages not meant for external use
├── pkg/              # Public packages that can be imported by other projects
│   ├── ai/           # AI provider integrations
│   ├── analyzer/     # Resource analyzers
│   ├── collector/    # Kubernetes resource collectors
│   └── config/       # Configuration management
└── tests/            # Integration and end-to-end tests
```

## Architecture Overview

K8sMed follows a modular architecture with these key components:

### Resource Collector

The collector component is responsible for gathering information about Kubernetes resources. Its main interfaces are defined in `pkg/collector/collector.go`.

Key concepts:
- `ResourceCollector`: Interface for collecting resources from a Kubernetes cluster
- `ResourceData`: Struct containing resource information, status, manifest, events, and logs
- `ResourceInfo`: Basic metadata about a resource (kind, name, namespace)

### Resource Analyzers

Analyzers examine collected resources for potential issues. The analyzer framework is in `pkg/analyzer/analyzer.go`.

Key concepts:
- `Analyzer`: Interface that all analyzers must implement
- `AnalysisContext`: Shared context for analyzers with resources and results
- `AnalysisDetail`: Analysis results with problem details and remediation steps

### AI Interface

The AI component connects to different LLM providers to process analysis results. Its interfaces are in `pkg/ai/llm/client.go`.

Key concepts:
- `Client`: Interface for communicating with LLM providers
- `Request`: Struct representing a prompt to the LLM
- `Response`: Struct containing the LLM's response

### CLI Interface

The kubectl plugin interface is defined in `cmd/kubectl-k8smed/main.go`.

## Development Workflow

1. **Create a Branch**

```bash
git checkout -b feature/your-feature-name
```

2. **Make Changes**

Implement your feature or bug fix.

3. **Write Tests**

Add tests for your changes in the appropriate package's `_test.go` files.

4. **Run Tests Locally**

```bash
make test
```

5. **Lint Your Code**

```bash
make lint
```

6. **Commit Your Changes**

```bash
git add .
git commit -m "Add feature: your feature description"
```

7. **Push to Your Fork**

```bash
git push origin feature/your-feature-name
```

8. **Create a Pull Request**

Open a pull request against the main repository's `main` branch.

## Testing

K8sMed uses Go's standard testing package for unit tests and integration tests.

### Running Tests

```bash
# Run all tests
make test

# Run specific tests
go test ./pkg/analyzer/...

# Run tests with verbose output
go test -v ./...

# Run integration tests
make integration-test
```

### Writing Tests

- Unit tests should be in the same package as the code they're testing, with a `_test.go` suffix
- Integration tests should be in the `tests/` directory
- Use table-driven tests where appropriate
- Mock external dependencies (like Kubernetes API) in unit tests

## Code Style and Conventions

K8sMed follows standard Go coding conventions:

- Use `gofmt` to format your code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Document all exported types, functions, and methods
- Keep functions small and focused on a single responsibility
- Use meaningful variable and function names

## Creating New Analyzers

Analyzers are the core of K8sMed's functionality. To create a new analyzer:

1. Create a new file in `pkg/analyzer/` for your resource type
2. Implement the `Analyzer` interface:

```go
type YourResourceAnalyzer struct{}

func (a *YourResourceAnalyzer) Analyze(ctx context.Context, analysisCtx *AnalysisContext) error {
    // Your analyzer implementation
    return nil
}

func (a *YourResourceAnalyzer) Name() string {
    return "YourResourceAnalyzer"
}

func (a *YourResourceAnalyzer) SupportedKinds() []string {
    return []string{"YourResource"}
}
```

3. Register your analyzer in `pkg/analyzer/registry.go`

### Analyzer Implementation Tips

- Focus on one issue type per analyzer function
- Provide clear descriptions of problems
- Include actionable remediation steps
- Generate specific kubectl commands where possible
- Consider the user's experience and level of expertise

## AI Provider Integration

To add a new AI provider:

1. Create a new file in `pkg/ai/llm/` for your provider
2. Implement the `Client` interface
3. Add the provider to the provider factory in `pkg/ai/llm/factory.go`

## Debugging Tips

### Troubleshooting the kubectl Plugin

```bash
# Run with verbose logging
kubectl k8smed analyze pod my-pod --debug

# Inspect collected resources
kubectl k8smed collect pod my-pod --output json > resources.json

# Test AI prompt without analyzing
kubectl k8smed prompt "Your prompt here" --dry-run
```

### Common Issues

- **Permission errors**: Ensure your kubectl context has sufficient permissions
- **API errors**: Check your AI provider credentials and endpoint configuration
- **Resource not found**: Verify the resource exists and the namespace is correct

## Release Process

1. **Version Bump**

Update the version in `pkg/version/version.go`.

2. **Update CHANGELOG.md**

Document all changes since the last release.

3. **Create a Tag**

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

4. **GitHub Release**

Create a new release on GitHub with the changelog content.

5. **Build Release Assets**

```bash
make release
```

---

Thank you for contributing to K8sMed! If you have questions or need help, please open an issue on GitHub or contact the maintainers directly. 