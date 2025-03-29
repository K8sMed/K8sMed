# Contributing to K8sMed

Thank you for considering contributing to K8sMed! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Contribution Workflow](#contribution-workflow)
- [Coding Standards](#coding-standards)
- [Documentation](#documentation)
- [Testing](#testing)
- [Community](#community)

## Code of Conduct

Our project is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [maintainer email].

## Getting Started

### Issues

- Before creating a new issue, please search to see if a similar issue already exists
- Use the provided issue templates when creating a new issue
- Be descriptive and include as much relevant information as possible
- For bugs, include steps to reproduce, expected behavior, and current behavior

### Pull Requests

- Pull requests should reference an existing issue
- Follow the [pull request template](.github/PULL_REQUEST_TEMPLATE.md)
- Keep pull requests focused on a single topic to make review easier

## Development Environment

### Prerequisites

- Go 1.21+
- Docker
- Kubernetes cluster (for testing)
- Access to an AI model (OpenAI account or LocalAI setup)

### Local Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/k8smed.git
   cd k8smed
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/k8smed/k8smed.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Build the project:
   ```bash
   make build
   ```

## Contribution Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bugfix-name
   ```
2. Make your changes
3. Commit your changes:
   ```bash
   git commit -m "Description of changes"
   ```
4. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Create a pull request against the `main` branch of the upstream repository

### Keeping Your Fork Updated

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

## Coding Standards

### Go Code

- Follow standard Go formatting and style conventions
- Use `gofmt` or `goimports` to automatically format your code
- Add comments for public functions, types, and non-obvious code
- Use meaningful variable and function names
- Keep functions small and focused

### Commit Messages

- Use clear, descriptive commit messages
- Start with a short summary line (50 chars or less)
- Optionally follow with a blank line and a more detailed explanation
- Reference issues and pull requests where appropriate

## Documentation

- Update documentation for any new features or changes
- Document public functions, types, and interfaces
- Include examples where appropriate
- Keep the README updated with relevant information
- Document environment variables and configuration options

## Testing

- Add tests for new functionality
- Ensure existing tests pass before submitting a PR
- Run tests locally using `make test`
- Include both unit tests and integration tests where appropriate

### Testing Locally

```bash
# Run all tests
make test

# Run specific tests
go test ./pkg/analyzer/...

# Run tests with verbose output
go test -v ./...
```

## Community

- Join our [Slack/Discord channel](#) for discussions
- Attend community meetings 
- Follow our [Twitter](#) for updates
- Subscribe to our newsletter

---

Thank you for contributing to K8sMed! Your efforts help make Kubernetes troubleshooting easier for everyone. 