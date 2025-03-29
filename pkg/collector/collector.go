package collector

import (
	"context"
	"fmt"

	internalcollector "github.com/k8smed/k8smed/internal/collector"
	internalpod "github.com/k8smed/k8smed/internal/collector/pod"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ResourceType represents the type of Kubernetes resource
type ResourceType string

// Define common resource types
const (
	ResourceTypePod         ResourceType = "pod"
	ResourceTypeDeployment  ResourceType = "deployment"
	ResourceTypeService     ResourceType = "service"
	ResourceTypeNode        ResourceType = "node"
	ResourceTypeNamespace   ResourceType = "namespace"
	ResourceTypeConfigMap   ResourceType = "configmap"
	ResourceTypeSecret      ResourceType = "secret"
	ResourceTypeStatefulSet ResourceType = "statefulset"
	ResourceTypeDaemonSet   ResourceType = "daemonset"
	ResourceTypeIngress     ResourceType = "ingress"
	ResourceTypePVC         ResourceType = "persistentvolumeclaim"
	ResourceTypePV          ResourceType = "persistentvolume"
	ResourceTypeEvent       ResourceType = "event"
)

// CollectionOptions provides options for resource collection
type CollectionOptions struct {
	Namespace     string
	ResourceName  string
	LabelSelector string
	IncludeEvents bool
	IncludeLogs   bool
	SinceSeconds  int64
	TailLines     int64
	Limit         int64
}

// ResourceInfo contains basic information about a Kubernetes resource
type ResourceInfo struct {
	Kind      string            `json:"kind"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	// Additional fields can be added as needed
}

// ResourceData represents the collected data for a Kubernetes resource
type ResourceData struct {
	Resource ResourceInfo      `json:"resource"`
	Manifest string            `json:"manifest"` // YAML representation
	Events   []string          `json:"events,omitempty"`
	Logs     []string          `json:"logs,omitempty"`
	Status   map[string]string `json:"status,omitempty"`
	Related  []ResourceInfo    `json:"related,omitempty"`
}

// Collector provides methods to collect information from a Kubernetes cluster
type Collector struct {
	clientset *kubernetes.Clientset
}

// NewCollector creates a new Collector instance
func NewCollector(kubeConfigPath string) (*Collector, error) {
	// Try to build config from the provided kubeconfig path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		// If that fails, try to use in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes client config: %w", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Collector{
		clientset: clientset,
	}, nil
}

// CollectResource collects data for the specified resource
func (c *Collector) CollectResource(ctx context.Context, resourceType ResourceType, options CollectionOptions) (*ResourceData, error) {
	// Convert our options to internal options
	internalOptions := internalcollector.CollectionOptions{
		Namespace:     options.Namespace,
		ResourceName:  options.ResourceName,
		LabelSelector: options.LabelSelector,
		IncludeEvents: options.IncludeEvents,
		IncludeLogs:   options.IncludeLogs,
		SinceSeconds:  options.SinceSeconds,
		TailLines:     options.TailLines,
		Limit:         options.Limit,
	}

	// Use the appropriate collector
	switch resourceType {
	case ResourceTypePod:
		return c.collectPod(ctx, internalOptions)
	case ResourceTypeDeployment:
		return c.collectDeployment(ctx, internalOptions)
	case ResourceTypeService:
		return c.collectService(ctx, internalOptions)
	case ResourceTypeEvent:
		return c.collectEvents(ctx, internalOptions)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// collectPod collects data for the specified pod
func (c *Collector) collectPod(ctx context.Context, options internalcollector.CollectionOptions) (*ResourceData, error) {
	podCollector := internalpod.NewCollector(c.clientset)
	internalData, err := podCollector.Collect(ctx, options)
	if err != nil {
		return nil, err
	}

	// Convert internal data to our format
	return &ResourceData{
		Resource: ResourceInfo{
			Kind:      internalData.Resource.Kind,
			Name:      internalData.Resource.Name,
			Namespace: internalData.Resource.Namespace,
			Labels:    internalData.Resource.Labels,
		},
		Manifest: internalData.Manifest,
		Events:   internalData.Events,
		Logs:     internalData.Logs,
		Status:   internalData.Status,
		Related:  convertRelatedResources(internalData.Related),
	}, nil
}

// convertRelatedResources converts internal resource info to our format
func convertRelatedResources(internalResources []internalcollector.ResourceInfo) []ResourceInfo {
	resources := make([]ResourceInfo, len(internalResources))
	for i, res := range internalResources {
		resources[i] = ResourceInfo{
			Kind:      res.Kind,
			Name:      res.Name,
			Namespace: res.Namespace,
			Labels:    res.Labels,
		}
	}
	return resources
}

func (c *Collector) collectDeployment(_ context.Context, _ internalcollector.CollectionOptions) (*ResourceData, error) {
	// This will be implemented in a dedicated package
	return nil, fmt.Errorf("not implemented")
}

func (c *Collector) collectService(_ context.Context, _ internalcollector.CollectionOptions) (*ResourceData, error) {
	// This will be implemented in a dedicated package
	return nil, fmt.Errorf("not implemented")
}

func (c *Collector) collectEvents(_ context.Context, _ internalcollector.CollectionOptions) (*ResourceData, error) {
	// This will be implemented in a dedicated package
	return nil, fmt.Errorf("not implemented")
}
