package remediation

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/k8smed/k8smed/pkg/analyzer"
	"github.com/k8smed/k8smed/pkg/collector"
)

// CommandType represents the type of remediation command
type CommandType string

// Define common command types
const (
	CommandTypeKubectl CommandType = "kubectl"
	CommandTypeHelm    CommandType = "helm"
	CommandTypeYAML    CommandType = "yaml"
	CommandTypeBash    CommandType = "bash"
)

// Command represents a command that can be executed to fix an issue
type Command struct {
	Type        CommandType       `json:"type"`
	Description string            `json:"description"`
	Command     string            `json:"command"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// Plan contains a set of commands to fix issues
type Plan struct {
	// Title of the remediation plan
	Title string `json:"title"`

	// Description of what this plan addresses
	Description string `json:"description"`

	// Steps to follow for remediation
	Steps []string `json:"steps"`

	// Commands to execute for remediation
	Commands []Command `json:"commands"`

	// YAML snippets for remediation that require manifest changes
	YAMLSnippets []string `json:"yamlSnippets,omitempty"`

	// Analysis details that led to this remediation
	AnalysisDetails []analyzer.AnalysisDetail `json:"analysisDetails"`
}

// Generator generates remediation plans based on analysis results
type Generator struct {
	// Templates for common remediation commands
	templates map[string]*template.Template
}

// NewGenerator creates a new remediation generator
func NewGenerator() *Generator {
	g := &Generator{
		templates: make(map[string]*template.Template),
	}

	// Register common templates
	g.registerTemplates()

	return g
}

// registerTemplates registers common command templates
func (g *Generator) registerTemplates() {
	// Template for fixing pod resource limits
	g.templates["fix-pod-resources"] = template.Must(template.New("fix-pod-resources").Parse(
		`kubectl patch {{ .ResourceType }} {{ .ResourceName }} -n {{ .Namespace }} --type=json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/{{ .ResourcePath }}", "value": {"{{ .ResourceType }}": "{{ .ResourceValue }}"}}]'`,
	))

	// Template for restarting a deployment
	g.templates["restart-deployment"] = template.Must(template.New("restart-deployment").Parse(
		`kubectl rollout restart deployment {{ .DeploymentName }} -n {{ .Namespace }}`,
	))

	// Template for scaling a deployment
	g.templates["scale-deployment"] = template.Must(template.New("scale-deployment").Parse(
		`kubectl scale deployment {{ .DeploymentName }} -n {{ .Namespace }} --replicas={{ .Replicas }}`,
	))
}

// GeneratePlan generates a remediation plan based on analysis details
func (g *Generator) GeneratePlan(details []analyzer.AnalysisDetail, resources []collector.ResourceData) *Plan {
	if len(details) == 0 {
		return nil
	}

	plan := &Plan{
		Title:           "Remediation Plan",
		Description:     "Steps to address the identified issues",
		Steps:           []string{},
		Commands:        []Command{},
		AnalysisDetails: details,
	}

	// Group details by type for better organization
	errorDetails := filterDetailsByType(details, "error")
	warningDetails := filterDetailsByType(details, "warning")
	infoDetails := filterDetailsByType(details, "info")

	// Add error remediation first (highest priority)
	if len(errorDetails) > 0 {
		plan.Steps = append(plan.Steps, "Fix critical issues:")
		for _, detail := range errorDetails {
			g.addRemediationForDetail(plan, detail, resources)
		}
	}

	// Add warning remediation next
	if len(warningDetails) > 0 {
		plan.Steps = append(plan.Steps, "Address warnings:")
		for _, detail := range warningDetails {
			g.addRemediationForDetail(plan, detail, resources)
		}
	}

	// Add info remediation last
	if len(infoDetails) > 0 {
		plan.Steps = append(plan.Steps, "Consider improvements:")
		for _, detail := range infoDetails {
			g.addRemediationForDetail(plan, detail, resources)
		}
	}

	return plan
}

// addRemediationForDetail adds remediation steps for a specific detail
func (g *Generator) addRemediationForDetail(plan *Plan, detail analyzer.AnalysisDetail, resources []collector.ResourceData) {
	// Add the remediation steps from the detail
	for _, step := range detail.Remediation {
		plan.Steps = append(plan.Steps, fmt.Sprintf("  - %s", step))
	}

	// Add the remediation commands from the detail
	for _, cmdStr := range detail.RemediationCommands {
		// Check if this is a template reference
		if strings.HasPrefix(cmdStr, "template:") {
			templateName := strings.TrimPrefix(cmdStr, "template:")
			template := g.templates[templateName]
			if template != nil {
				// Find the associated resource
				for _, res := range resources {
					if res.Resource.Kind == detail.Resource.Kind &&
						res.Resource.Name == detail.Resource.Name &&
						res.Resource.Namespace == detail.Resource.Namespace {

						// Create a buffer for the rendered template
						var buf bytes.Buffer

						// Execute the template with the resource data
						if err := template.Execute(&buf, res.Resource); err == nil {
							plan.Commands = append(plan.Commands, Command{
								Type:        CommandTypeKubectl,
								Description: fmt.Sprintf("Apply %s template for %s/%s", templateName, detail.Resource.Kind, detail.Resource.Name),
								Command:     buf.String(),
							})
						}
						break
					}
				}
			}
		} else {
			// This is a direct command
			plan.Commands = append(plan.Commands, Command{
				Type:        determineCommandType(cmdStr),
				Description: fmt.Sprintf("Remediation step for %s/%s", detail.Resource.Kind, detail.Resource.Name),
				Command:     cmdStr,
			})
		}
	}
}

// filterDetailsByType filters analysis details by their type
func filterDetailsByType(details []analyzer.AnalysisDetail, detailType string) []analyzer.AnalysisDetail {
	filtered := make([]analyzer.AnalysisDetail, 0)
	for _, detail := range details {
		if detail.Type == detailType {
			filtered = append(filtered, detail)
		}
	}
	return filtered
}

// determineCommandType tries to figure out the type of command
func determineCommandType(cmd string) CommandType {
	cmd = strings.TrimSpace(cmd)

	if strings.HasPrefix(cmd, "kubectl ") {
		return CommandTypeKubectl
	}

	if strings.HasPrefix(cmd, "helm ") {
		return CommandTypeHelm
	}

	if strings.HasPrefix(cmd, "apiVersion:") || strings.HasPrefix(cmd, "kind:") {
		return CommandTypeYAML
	}

	return CommandTypeBash
}
