package analyzer

import (
	"context"
	"fmt"
	"strings"

	"github.com/k8smed/k8smed/pkg/collector"
)

// AnalysisContext contains all the information needed for analysis
type AnalysisContext struct {
	// Original user query
	Query string

	// Collected resources data
	Resources []collector.ResourceData

	// Analysis details (filled by analyzers)
	Details []AnalysisDetail
}

// AnalysisDetail represents a single finding or observation during analysis
type AnalysisDetail struct {
	// Type of the issue/finding (error, warning, info)
	Type string

	// Short description of the issue
	Title string

	// Detailed explanation
	Description string

	// Resource related to this finding
	Resource collector.ResourceInfo

	// Suggested remediation steps
	Remediation []string

	// Commands that could help fix the issue
	RemediationCommands []string
}

// Analyzer interface defines methods for analyzing Kubernetes resources
type Analyzer interface {
	// Name returns the analyzer name
	Name() string

	// Description returns the analyzer description
	Description() string

	// Analyze performs analysis on the provided context
	Analyze(ctx context.Context, analysisCtx *AnalysisContext) error
}

// PodAnalyzer analyzes pod-related issues
type PodAnalyzer struct{}

// Name implements the Analyzer interface
func (a *PodAnalyzer) Name() string {
	return "PodAnalyzer"
}

// Description implements the Analyzer interface
func (a *PodAnalyzer) Description() string {
	return "Analyzes pod issues like CrashLoopBackOff, ImagePullBackOff, etc."
}

// Analyze implements the Analyzer interface
func (a *PodAnalyzer) Analyze(ctx context.Context, analysisCtx *AnalysisContext) error {
	// Filter for pod resources
	for _, resource := range analysisCtx.Resources {
		if resource.Resource.Kind == "Pod" {
			// Check container states for common issues
			a.checkContainerStates(resource, analysisCtx)

			// Check pod events for issues
			a.checkEvents(resource, analysisCtx)

			// Check pod logs if available
			if len(resource.Logs) > 0 {
				a.checkLogs(resource, analysisCtx)
			}

			// Check pod status
			a.checkPodStatus(resource, analysisCtx)
		}
	}

	return nil
}

// checkContainerStates checks for common container state issues
func (a *PodAnalyzer) checkContainerStates(resource collector.ResourceData, analysisCtx *AnalysisContext) {
	// Loop through all possible container states in the status map
	for key, value := range resource.Status {
		if strings.HasSuffix(key, ".state") && value == "waiting" {
			// Extract container index and name
			parts := strings.Split(key, ".")
			if len(parts) < 2 {
				continue
			}

			containerPrefix := parts[0] + "." + parts[1] + "."
			containerName := resource.Status[containerPrefix+"name"]
			waitReason := resource.Status[containerPrefix+"reason"]
			waitMessage := resource.Status[containerPrefix+"message"]

			// Skip if no name or reason
			if containerName == "" || waitReason == "" {
				continue
			}

			// Handle specific waiting reasons
			switch waitReason {
			case "CrashLoopBackOff":
				detail := AnalysisDetail{
					Type:        "error",
					Title:       "Container in CrashLoopBackOff",
					Description: "Container " + containerName + " is crash looping: " + waitMessage,
					Resource:    resource.Resource,
					Remediation: []string{
						"Check container logs for errors",
						"Verify the container command is correct",
						"Check if container has appropriate resources",
						"Ensure configuration files exist and are correct",
					},
					RemediationCommands: []string{
						"kubectl logs " + resource.Resource.Name + " -c " + containerName + " -n " + resource.Resource.Namespace,
						"kubectl describe pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
					},
				}
				analysisCtx.Details = append(analysisCtx.Details, detail)

			case "ImagePullBackOff", "ErrImagePull":
				detail := AnalysisDetail{
					Type:        "error",
					Title:       "Image pull failure",
					Description: "Container " + containerName + " cannot pull its image: " + waitMessage,
					Resource:    resource.Resource,
					Remediation: []string{
						"Verify the image name and tag are correct",
						"Check if private registry requires authentication",
						"Ensure image pull secrets are configured",
						"Verify network connectivity to the registry",
					},
					RemediationCommands: []string{
						"kubectl describe pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
					},
				}
				analysisCtx.Details = append(analysisCtx.Details, detail)

			case "CreateContainerConfigError":
				detail := AnalysisDetail{
					Type:        "error",
					Title:       "Container configuration error",
					Description: "Container " + containerName + " has configuration errors: " + waitMessage,
					Resource:    resource.Resource,
					Remediation: []string{
						"Check if referenced ConfigMaps exist",
						"Check if referenced Secrets exist",
						"Verify volume mounts are correctly configured",
						"Check container environment variables",
					},
					RemediationCommands: []string{
						"kubectl describe pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
						"kubectl get configmaps -n " + resource.Resource.Namespace,
						"kubectl get secrets -n " + resource.Resource.Namespace,
					},
				}
				analysisCtx.Details = append(analysisCtx.Details, detail)
			}
		}
	}
}

// checkEvents examines pod events for issues
func (a *PodAnalyzer) checkEvents(resource collector.ResourceData, analysisCtx *AnalysisContext) {
	// Skip if no events
	if len(resource.Events) == 0 {
		return
	}

	// Look for specific event patterns
	for _, event := range resource.Events {
		lowerEvent := strings.ToLower(event)

		// Look for OOMKilled issues
		if strings.Contains(lowerEvent, "oomkilled") {
			detail := AnalysisDetail{
				Type:        "error",
				Title:       "Container terminated due to OOMKilled",
				Description: "A container was terminated because it exceeded its memory limits: " + event,
				Resource:    resource.Resource,
				Remediation: []string{
					"Increase memory limits for the container",
					"Optimize the application to use less memory",
					"Check for memory leaks in the application",
				},
				RemediationCommands: []string{
					"kubectl describe pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
					"kubectl get pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace + " -o yaml",
				},
			}
			analysisCtx.Details = append(analysisCtx.Details, detail)
		}

		// Look for eviction issues
		if strings.Contains(lowerEvent, "evict") {
			detail := AnalysisDetail{
				Type:        "warning",
				Title:       "Pod was evicted",
				Description: "The pod was evicted from its node: " + event,
				Resource:    resource.Resource,
				Remediation: []string{
					"Check node resource pressure (CPU, memory, disk)",
					"Use node affinity to schedule on larger nodes",
					"Add resource quotas to prevent resource exhaustion",
				},
				RemediationCommands: []string{
					"kubectl describe node <node-name>",
					"kubectl top nodes",
				},
			}
			analysisCtx.Details = append(analysisCtx.Details, detail)
		}
	}
}

// checkLogs looks for error patterns in logs
func (a *PodAnalyzer) checkLogs(resource collector.ResourceData, analysisCtx *AnalysisContext) {
	// Combine all logs to search for patterns
	combinedLogs := strings.Join(resource.Logs, "\n")
	lowerLogs := strings.ToLower(combinedLogs)

	// Look for common error patterns
	if strings.Contains(lowerLogs, "exception") || strings.Contains(lowerLogs, "error") {
		// Extract a snippet around the error
		errorSnippet := a.extractErrorSnippet(combinedLogs)

		detail := AnalysisDetail{
			Type:        "warning",
			Title:       "Errors detected in logs",
			Description: "The pod logs contain errors or exceptions: " + errorSnippet,
			Resource:    resource.Resource,
			Remediation: []string{
				"Review application logs for detailed error information",
				"Check application configuration",
				"Verify external dependencies are available",
			},
			RemediationCommands: []string{
				"kubectl logs " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
			},
		}
		analysisCtx.Details = append(analysisCtx.Details, detail)
	}

	// Look for connection issues
	if strings.Contains(lowerLogs, "connection refused") ||
		strings.Contains(lowerLogs, "cannot connect") ||
		strings.Contains(lowerLogs, "dial tcp") {
		detail := AnalysisDetail{
			Type:        "warning",
			Title:       "Connection issues detected",
			Description: "The logs show connection problems to other services",
			Resource:    resource.Resource,
			Remediation: []string{
				"Verify the service endpoints are correct",
				"Check network policies allow the connection",
				"Ensure the target service is running",
			},
			RemediationCommands: []string{
				"kubectl get svc -n " + resource.Resource.Namespace,
				"kubectl get endpoints -n " + resource.Resource.Namespace,
				"kubectl get networkpolicies -n " + resource.Resource.Namespace,
			},
		}
		analysisCtx.Details = append(analysisCtx.Details, detail)
	}
}

// checkPodStatus checks the pod's phase and conditions
func (a *PodAnalyzer) checkPodStatus(resource collector.ResourceData, analysisCtx *AnalysisContext) {
	// Check pod phase
	phase, ok := resource.Status["phase"]
	if !ok {
		return
	}

	switch phase {
	case "Pending":
		// Check for scheduling issues
		for i := 0; ; i++ {
			condTypeKey := fmt.Sprintf("condition.%d.type", i)
			condStatusKey := fmt.Sprintf("condition.%d.status", i)
			condReasonKey := fmt.Sprintf("condition.%d.reason", i)
			condMessageKey := fmt.Sprintf("condition.%d.message", i)

			condType, typeExists := resource.Status[condTypeKey]
			if !typeExists {
				break
			}

			condStatus, statusExists := resource.Status[condStatusKey]
			if !statusExists {
				continue
			}

			if condType == "PodScheduled" && condStatus == "False" {
				reason := resource.Status[condReasonKey]
				message := resource.Status[condMessageKey]

				detail := AnalysisDetail{
					Type:        "error",
					Title:       "Pod scheduling issues",
					Description: "The pod is in a Pending state and cannot be scheduled: " + message,
					Resource:    resource.Resource,
					Remediation: []string{
						"Check cluster resource capacity",
						"Check node taints and affinities",
						"Check resource requests and limits",
					},
				}

				if strings.Contains(reason, "Insufficient") {
					detail.RemediationCommands = []string{
						"kubectl get nodes",
						"kubectl describe nodes <node-name>",
						"kubectl top nodes",
					}
				}

				analysisCtx.Details = append(analysisCtx.Details, detail)
			}
		}

	case "Failed":
		// Pod in Failed state
		detail := AnalysisDetail{
			Type:        "error",
			Title:       "Pod failed",
			Description: "The pod is in a Failed state",
			Resource:    resource.Resource,
			Remediation: []string{
				"Check pod logs for errors",
				"Check pod events for issues",
				"Check if containers are configured correctly",
			},
			RemediationCommands: []string{
				"kubectl logs " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
				"kubectl describe pod " + resource.Resource.Name + " -n " + resource.Resource.Namespace,
			},
		}
		analysisCtx.Details = append(analysisCtx.Details, detail)
	}
}

// extractErrorSnippet extracts a short snippet around the first error or exception in the logs
func (a *PodAnalyzer) extractErrorSnippet(logs string) string {
	// Find the index of "error" or "exception"
	errorIndex := strings.Index(strings.ToLower(logs), "error")
	exceptionIndex := strings.Index(strings.ToLower(logs), "exception")

	var index int
	if errorIndex >= 0 && exceptionIndex >= 0 {
		if errorIndex < exceptionIndex {
			index = errorIndex
		} else {
			index = exceptionIndex
		}
	} else if errorIndex >= 0 {
		index = errorIndex
	} else if exceptionIndex >= 0 {
		index = exceptionIndex
	} else {
		return ""
	}

	// Extract up to 150 characters around the error
	start := index - 50
	if start < 0 {
		start = 0
	}

	end := index + 100
	if end > len(logs) {
		end = len(logs)
	}

	// Get the snippet and truncate at line boundaries
	snippet := logs[start:end]
	lines := strings.Split(snippet, "\n")

	// Take at most 3 lines
	if len(lines) > 3 {
		if len(lines) > 4 {
			lines = append(lines[:2], "...", lines[len(lines)-1])
		} else {
			lines = lines[:3]
		}
	}

	return strings.Join(lines, "\n")
}

// DeploymentAnalyzer analyzes deployment-related issues
type DeploymentAnalyzer struct{}

// Name implements the Analyzer interface
func (a *DeploymentAnalyzer) Name() string {
	return "DeploymentAnalyzer"
}

// Description implements the Analyzer interface
func (a *DeploymentAnalyzer) Description() string {
	return "Analyzes deployment issues like unavailable replicas, rollout failures, etc."
}

// Analyze implements the Analyzer interface
func (a *DeploymentAnalyzer) Analyze(ctx context.Context, analysisCtx *AnalysisContext) error {
	// This is a placeholder implementation

	// For now, just add a dummy detail
	analysisCtx.Details = append(analysisCtx.Details, AnalysisDetail{
		Type:  "info",
		Title: "Deployment analyzer executed",
		Description: "This is a placeholder for the deployment analyzer. In a real implementation, " +
			"it would analyze replica status, rollout status, pod template issues, etc.",
		Resource: collector.ResourceInfo{
			Kind: "Deployment",
		},
		Remediation: []string{
			"Check deployment status using 'kubectl rollout status deployment <name>'",
			"Check deployment events using 'kubectl describe deployment <name>'",
		},
		RemediationCommands: []string{
			"kubectl rollout status deployment ${DEPLOYMENT_NAME}",
			"kubectl describe deployment ${DEPLOYMENT_NAME}",
		},
	})

	return nil
}

// Registry keeps track of available analyzers
type Registry struct {
	analyzers map[string]Analyzer
}

// NewRegistry creates a new registry with the default analyzers
func NewRegistry() *Registry {
	registry := &Registry{
		analyzers: make(map[string]Analyzer),
	}

	// Register default analyzers
	registry.Register(&PodAnalyzer{})
	registry.Register(&DeploymentAnalyzer{})

	return registry
}

// Register adds an analyzer to the registry
func (r *Registry) Register(analyzer Analyzer) {
	r.analyzers[analyzer.Name()] = analyzer
}

// Get returns the analyzer with the given name
func (r *Registry) Get(name string) Analyzer {
	return r.analyzers[name]
}

// GetAll returns all registered analyzers
func (r *Registry) GetAll() []Analyzer {
	analyzers := make([]Analyzer, 0, len(r.analyzers))
	for _, analyzer := range r.analyzers {
		analyzers = append(analyzers, analyzer)
	}
	return analyzers
}
