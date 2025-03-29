# Security Policy

## Supported Versions

K8sMed is currently in active development. The following versions are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

We take the security of K8sMed seriously. If you believe you've found a security vulnerability, please follow these steps:

1. **Do not disclose the vulnerability publicly**
2. **Email the core maintainers directly at [k8smed.security@example.com]** with details about the vulnerability
3. Include the following information:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested mitigations (if known)

The maintainers will acknowledge receipt of your report within 48 hours and will provide an estimated timeline for a fix.

## Security Best Practices

When using K8sMed, consider the following security practices:

1. Use the anonymization feature when sharing diagnostics with external AI providers
2. Run K8sMed with minimal required permissions in your Kubernetes cluster
3. For sensitive environments, consider using LocalAI or other self-hosted AI providers
4. Review the AI prompts and responses for any sensitive information before sharing

Thank you for helping keep K8sMed and its users safe!
