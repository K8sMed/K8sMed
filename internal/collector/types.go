package collector

import (
	"context"
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
	Kind      string
	Name      string
	Namespace string
	Labels    map[string]string
}

// ResourceData represents the collected data for a Kubernetes resource
type ResourceData struct {
	Resource ResourceInfo
	Manifest string
	Events   []string
	Logs     []string
	Status   map[string]string
	Related  []ResourceInfo
}

// ResourceCollector defines the interface for resource collectors
type ResourceCollector interface {
	// Collect gathers data for a specific resource
	Collect(ctx context.Context, options CollectionOptions) (*ResourceData, error)
}
