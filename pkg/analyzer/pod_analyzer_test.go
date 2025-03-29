package analyzer

import (
	"context"
	"testing"

	"github.com/k8smed/k8smed/pkg/collector"
)

func TestPodAnalyzer_CrashLoopBackOff(t *testing.T) {
	// Create a pod analyzer
	analyzer := &PodAnalyzer{}

	// Create a mock resource with a container in CrashLoopBackOff
	resource := collector.ResourceData{
		Resource: collector.ResourceInfo{
			Kind:      "Pod",
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: map[string]string{
			"phase":                    "Running",
			"container.0.name":         "app",
			"container.0.state":        "waiting",
			"container.0.reason":       "CrashLoopBackOff",
			"container.0.message":      "Back-off 5m0s restarting failed container",
			"container.0.ready":        "false",
			"container.0.restartCount": "5",
		},
	}

	// Create an analysis context with the resource
	analysisCtx := &AnalysisContext{
		Query:     "What's wrong with my pod?",
		Resources: []collector.ResourceData{resource},
		Details:   []AnalysisDetail{},
	}

	// Run the analyzer
	err := analyzer.Analyze(context.Background(), analysisCtx)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	// Verify the results
	if len(analysisCtx.Details) == 0 {
		t.Fatal("Expected at least one analysis detail")
	}

	// Find the CrashLoopBackOff detail
	var crashLoopDetail *AnalysisDetail
	for i, detail := range analysisCtx.Details {
		if detail.Title == "Container in CrashLoopBackOff" {
			crashLoopDetail = &analysisCtx.Details[i]
			break
		}
	}

	if crashLoopDetail == nil {
		t.Fatal("Expected to find a CrashLoopBackOff detail")
	}

	// Verify the detail properties
	if crashLoopDetail.Type != "error" {
		t.Errorf("Expected detail type to be 'error', got '%s'", crashLoopDetail.Type)
	}

	if crashLoopDetail.Resource.Name != "test-pod" {
		t.Errorf("Expected resource name to be 'test-pod', got '%s'", crashLoopDetail.Resource.Name)
	}

	// Verify remediation commands exist
	if len(crashLoopDetail.RemediationCommands) == 0 {
		t.Error("Expected remediation commands to be provided")
	}
}

func TestPodAnalyzer_ImagePullBackOff(t *testing.T) {
	// Create a pod analyzer
	analyzer := &PodAnalyzer{}

	// Create a mock resource with a container having image pull issues
	resource := collector.ResourceData{
		Resource: collector.ResourceInfo{
			Kind:      "Pod",
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: map[string]string{
			"phase":               "Pending",
			"container.0.name":    "app",
			"container.0.state":   "waiting",
			"container.0.reason":  "ImagePullBackOff",
			"container.0.message": "Back-off pulling image",
			"container.0.ready":   "false",
		},
	}

	// Create an analysis context with the resource
	analysisCtx := &AnalysisContext{
		Query:     "Why is my pod not starting?",
		Resources: []collector.ResourceData{resource},
		Details:   []AnalysisDetail{},
	}

	// Run the analyzer
	err := analyzer.Analyze(context.Background(), analysisCtx)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	// Verify the results
	if len(analysisCtx.Details) == 0 {
		t.Fatal("Expected at least one analysis detail")
	}

	// Find the ImagePullBackOff detail
	var imagePullDetail *AnalysisDetail
	for i, detail := range analysisCtx.Details {
		if detail.Title == "Image pull failure" {
			imagePullDetail = &analysisCtx.Details[i]
			break
		}
	}

	if imagePullDetail == nil {
		t.Fatal("Expected to find an Image pull failure detail")
	}

	// Verify the detail properties
	if imagePullDetail.Type != "error" {
		t.Errorf("Expected detail type to be 'error', got '%s'", imagePullDetail.Type)
	}

	// Verify remediation steps exist
	if len(imagePullDetail.Remediation) < 2 {
		t.Error("Expected multiple remediation steps to be provided")
	}
}
