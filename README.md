# K8sMed: AI-Powered Kubernetes First Responder

K8sMed is an open-source, AI-powered troubleshooting assistant designed to act as a first responder for Kubernetes clusters. By continuously monitoring cluster logs, events, and metrics, K8sMed leverages Large Language Models (LLMs) to diagnose issues, provide natural language explanations, and generate actionable remediation commands—all through a simple kubectl plugin and Operator.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Key Features](#key-features)
- [Architecture Overview](#architecture-overview)
- [Differentiators](#differentiators)
- [Roadmap & Milestones](#roadmap--milestones)
- [Installation](#installation)
  - [CLI Plugin Installation](#cli-plugin-installation)
  - [Operator Deployment](#operator-deployment)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Interactive Mode](#interactive-mode)
  - [Anonymization & Data Privacy](#anonymization--data-privacy)
- [Configuration & Customization](#configuration--customization)
- [Contributing](#contributing)
- [Future Documentation](#future-documentation)
- [License](#license)
- [Contact & Support](#contact--support)

---

## Project Overview

Status: Conceptual / Proposal Stage

K8sMed is a proposed open-source project aimed at revolutionizing Kubernetes troubleshooting. The idea is to create a lightweight “first responder” tool—a kubectl plugin (and optionally, a Kubernetes Operator) that leverages AI (using LLMs) to quickly analyze cluster state, diagnose issues, and suggest actionable remediation steps.

This repository contains the initial project details, design proposal, and roadmap. We are in the early stages of development and are inviting contributions, ideas, and feedback from the community to help shape and build K8sMed into a practical tool for DevOps teams and Kubernetes operators.

### Project Vision
**Objective:**  
Develop an AI-powered troubleshooting assistant for Kubernetes that:
- Continuously gathers logs, events, and metrics.
- Uses AI to perform root-cause analysis.
- Provides clear, human-readable explanations and remediation commands.
- Remains privacy-focused with support for anonymization and local AI options.

Why K8sMed?
Kubernetes environments are complex, and troubleshooting issues often require deep expertise. K8sMed is envisioned to reduce mean time to resolution (MTTR) by acting as a “first responder,” enabling both experts and non-experts to quickly understand and fix problems in their clusters.

### Goals

- **Rapid Diagnosis:** Automatically detect anomalies and provide clear, step-by-step remediation instructions.
- **Actionable Insights:** Generate copy-paste-ready `kubectl` commands and YAML patches.
- **Privacy First:** Anonymize sensitive data and support local LLM deployments.
- **Modular Extensibility:** Allow users to add custom analyzers for specific resources and error patterns.
- **Seamless Integration:** Offer both CLI (kubectl plugin) and Operator deployments to fit different operational needs.

---

## Key Features

- **Real-Time Monitoring & Collection:**  
  Fetch logs, events, and metrics using kubectl and integrations with Prometheus, Trivy, etc.

- **AI-Powered Analysis & Explanation:**  
  Leverage cloud-based LLMs (e.g., OpenAI’s GPT-3.5/4) or local solutions (LocalAI/Ollama) to process collected data and generate natural language diagnostics.

- **Actionable Remediation:**  
  Output precise remediation commands or YAML manifest updates to resolve issues immediately.

- **Interactive Troubleshooting Mode:**  
  Engage in a conversational, iterative troubleshooting session via an interactive CLI mode.

- **Privacy & Security:**  
  Built-in anonymization ensures that only essential, non-sensitive data is shared with AI backends. Local deployment of LLMs is supported for air-gapped or high-security environments.

- **Extensible & Modular Design:**  
  Easy-to-extend architecture allows developers to add new analyzers and integrate third-party tools.

---

## Architecture Overview

### Components

1. **Data Collection Layer:**  
   - Utilizes standard kubectl commands and integrations (e.g., Prometheus, Trivy) to gather cluster state, logs, and events.
2. **Preprocessing & Analyzer Module:**  
   - Filters and formats data for AI processing.
   - Contains built-in analyzers for common Kubernetes objects (Pods, Deployments, Nodes, etc.) and supports custom analyzer plugins.
3. **AI Interface Layer:**  
   - Connects to one or more LLM backends.
   - Provides both one-off analysis and an interactive troubleshooting session.
4. **Remediation Module:**  
   - Converts AI output into actionable commands (e.g., `kubectl patch` or YAML updates).
   - Optionally integrates with CI/CD systems for automation.
5. **Operator Mode (Optional):**  
   - Deployed as a Kubernetes Operator, it automates periodic scans and aggregates diagnostic results as custom resources, which can be queried centrally.

---

## Differentiators

- **Focused First Responder:**  
  K8sMed is specifically designed for rapid troubleshooting, not just manifest generation or alerting.
- **Privacy-First Design:**  
  Offers built-in anonymization and supports local LLMs, addressing data security concerns.
- **Lightweight & Developer-Friendly:**  
  Minimal setup with both CLI and Operator options, reducing friction for both experts and non-experts.
- **Actionable Guidance:**  
  Provides detailed, step-by-step remediation commands that are easy to understand and implement.
- **Modular Extensibility:**  
  A pluggable architecture that allows easy addition of new analyzers and integration with external tools.

---

## Roadmap & Milestones

### Phase 1: MVP Development (0–3 months)
- Implement core CLI plugin functionality:
  - Data collection from Kubernetes cluster.
  - Basic set of analyzers (e.g., for Pods, Deployments, Services).
  - Integration with a default LLM (e.g., GPT-3.5-Turbo).
  - Anonymization feature.
- Publish initial version on GitHub as a kubectl plugin.

### Phase 2: Extended Features (3–6 months)
- Develop additional analyzers for complex resources (e.g., StatefulSets, CronJobs, Network Policies).
- Add support for multiple AI backends (local and cloud-based).
- Build interactive troubleshooting mode.
- Begin developing the Operator version for continuous monitoring and centralized reporting.

### Phase 3: Community Integration & Ecosystem Expansion (6–12 months)
- Open up contribution guidelines and foster community involvement.
- Integrate with external tools (e.g., CI/CD, alerting systems).
- Extend documentation and create use-case examples, tutorials, and webinars.
- Enhance remediation module with automated patch application features.

---

## Installation 
#### How it will look like in future
### CLI Plugin Installation

1. **Prerequisites:**
   - Kubernetes cluster with `kubectl` configured.
   - Go 1.18+ installed (for building from source).
   - An LLM API key (if using a cloud-based LLM) or set up a local LLM backend.

2. **Clone and Build:**

   ```bash
   git clone https://github.com/your-org/k8smed.git
   cd k8smed
   go build -o kubectl-k8smed .
   ```

3. **Add to PATH:**
   - Move the binary to a directory in your PATH (e.g., `/usr/local/bin`):

   ```bash
   sudo mv kubectl-k8smed /usr/local/bin/
   ```

4. **Verify Installation:**

   ```bash
   kubectl k8smed version
   ```

### Operator Deployment

1. **Helm Chart:**
   - A Helm chart is provided to deploy K8sMed as an Operator.
   - Configure the chart values (AI backend, scan frequency, etc.) as needed.
   
2. **Deploy the Operator:**

   ```bash
   helm repo add k8smed https://your-org.github.io/k8smed-charts
   helm install k8smed-operator k8smed/k8smed-operator --namespace k8smed-system --create-namespace
   ```

3. **Configure the Custom Resource:**
   - Create a YAML file (e.g., `k8smed-cr.yaml`) that defines your scan configuration:

   ```yaml
   apiVersion: k8smed.io/v1
   kind: K8sMed
   metadata:
     name: default-scan
     namespace: k8smed-system
   spec:
     aiBackend: "openai"
     model: "gpt-3.5-turbo"
     scanInterval: "5m"
     anonymize: true
   ```

4. **Apply the Custom Resource:**

   ```bash
   kubectl apply -f k8smed-cr.yaml
   ```

---

## Usage

### Basic Commands

- **Run a One-Off Analysis:**

  ```bash
  kubectl k8smed analyze "explain why pod myapp-123 is in CrashLoopBackOff"
  ```

- **Generate Remediation Command:**

  ```bash
  kubectl k8smed analyze --explain "increase memory limit for myapp-123"
  ```

### Interactive Mode

- **Start Interactive Troubleshooting:**

  ```bash
  kubectl k8smed interactive
  ```

  In interactive mode, type queries, and the assistant will maintain context for follow-up questions.

### Anonymization & Data Privacy

- Use the `--anonymize` flag to mask sensitive information in queries:

  ```bash
  kubectl k8smed analyze --explain --anonymize "diagnose issues in namespace sensitive-ns"
  ```

- For high-security environments, configure K8sMed to use a local LLM by setting environment variables (e.g., `LOCAL_LLM=true` and `LLM_ENDPOINT=http://localhost:8080`).

---

## Configuration & Customization

- **AI Backend Configuration:**
  - Set your AI backend by exporting the necessary environment variables:

    ```bash
    export OPENAI_API_KEY=your_openai_api_key
    export AI_BACKEND=openai
    export MODEL_NAME=gpt-3.5-turbo
    ```

- **Custom Analyzers:**
  - Developers can add new analyzers by creating Go modules that implement the Analyzer interface. See `docs/developing-analyzers.md` for guidelines.

- **Operator Config Options:**
  - Adjust the scan interval, target namespaces, and filters via the custom resource configuration.

---

## Contributing & Getting Involved

Since K8sMed is at the conceptual stage, we welcome your ideas, feedback, and contributions. Here’s how you can help:

- **Feedback:** Open issues in this repository to discuss ideas, potential challenges, and proposed features.
- **Contributions:** If you’re interested in prototyping or developing parts of the tool (e.g., initial analyzers or AI integration modules), please submit a pull request.
- **Discussion:** Join our community Slack or Discord channel (details below) to discuss the project roadmap and share your thoughts on the design and functionality.

---

## Future Documentation

As we progress, we plan to add:
- Detailed architecture diagrams and design documents.
- User guides and setup instructions for both the CLI plugin and Operator modes.
- Contribution guidelines and a roadmap for future releases.

---

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

---

## Contact & Support

- **Project Lead:** [Md Imran - imranaec@outlook.com] [@narmidm](https://github.com/narmidm)
- **Community Discussions:** Join our Slack or Discord channels for discussions, feedback, and support.
- **Documentation:** For more details, see our [Documentation Portal](https://your-org.github.io/k8smed-docs).

---

## Final Notes

K8sMed aims to become the definitive AI-powered troubleshooting assistant for Kubernetes. With its focused “first responder” design, privacy-first architecture, and actionable remediation guidance, it is poised to transform how DevOps teams manage and troubleshoot Kubernetes clusters. We look forward to your feedback and contributions as we kick off this project!

---
