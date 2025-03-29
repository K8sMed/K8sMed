package pod

import (
	"context"
	"fmt"
	"io"

	"github.com/k8smed/k8smed/internal/collector"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Collector implements pod data collection
type Collector struct {
	clientset *kubernetes.Clientset
}

// NewCollector creates a new pod collector
func NewCollector(clientset *kubernetes.Clientset) *Collector {
	return &Collector{
		clientset: clientset,
	}
}

// Collect gathers data about a pod
func (c *Collector) Collect(ctx context.Context, options collector.CollectionOptions) (*collector.ResourceData, error) {
	// Validate options
	if options.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	// Handle resource filtering
	var pod *corev1.Pod
	var err error

	if options.ResourceName != "" {
		// Get single pod by name
		pod, err = c.clientset.CoreV1().Pods(options.Namespace).Get(ctx, options.ResourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get pod %s: %w", options.ResourceName, err)
		}
	} else if options.LabelSelector != "" {
		// Get pods by label selector
		pods, err := c.clientset.CoreV1().Pods(options.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: options.LabelSelector,
			Limit:         options.Limit,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods with selector %s: %w", options.LabelSelector, err)
		}
		if len(pods.Items) == 0 {
			return nil, fmt.Errorf("no pods found with selector %s", options.LabelSelector)
		}
		// Use the first pod for detailed collection
		pod = &pods.Items[0]
	} else {
		return nil, fmt.Errorf("either pod name or label selector is required")
	}

	// Create resource info
	resourceInfo := collector.ResourceInfo{
		Kind:      "Pod",
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Labels:    pod.Labels,
	}

	// Collect pod data
	resourceData := &collector.ResourceData{
		Resource: resourceInfo,
		Status:   extractPodStatus(pod),
		Related:  []collector.ResourceInfo{},
	}

	// Collect events if requested
	if options.IncludeEvents {
		events, err := c.collectEvents(ctx, pod)
		if err != nil {
			// Log the error but continue
			fmt.Printf("Warning: failed to collect events: %v\n", err)
		} else {
			resourceData.Events = events
		}
	}

	// Collect logs if requested
	if options.IncludeLogs {
		logs, err := c.collectLogs(ctx, pod, options)
		if err != nil {
			// Log the error but continue
			fmt.Printf("Warning: failed to collect logs: %v\n", err)
		} else {
			resourceData.Logs = logs
		}
	}

	return resourceData, nil
}

// collectEvents gathers events related to the pod
func (c *Collector) collectEvents(ctx context.Context, pod *corev1.Pod) ([]string, error) {
	// Get events for the pod
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s,involvedObject.kind=Pod",
		pod.Name, pod.Namespace)

	events, err := c.clientset.CoreV1().Events(pod.Namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	// Format events as strings
	eventStrings := make([]string, 0, len(events.Items))
	for _, event := range events.Items {
		eventTime := event.LastTimestamp.Time
		if eventTime.IsZero() {
			eventTime = event.CreationTimestamp.Time
		}

		eventStr := fmt.Sprintf("[%s] %s %s: %s (count: %d)",
			eventTime.Format("2006-01-02 15:04:05"),
			event.Type,
			event.Reason,
			event.Message,
			event.Count,
		)
		eventStrings = append(eventStrings, eventStr)
	}

	return eventStrings, nil
}

// collectLogs gathers logs from the pod's containers
func (c *Collector) collectLogs(ctx context.Context, pod *corev1.Pod, options collector.CollectionOptions) ([]string, error) {
	allLogs := make([]string, 0)

	// Determine which containers to collect logs from
	containers := make([]corev1.Container, 0)
	containers = append(containers, pod.Spec.Containers...)
	containers = append(containers, pod.Spec.InitContainers...)

	// Collect logs from each container
	for _, container := range containers {
		// Set up log options
		logOptions := &corev1.PodLogOptions{
			Container: container.Name,
		}

		if options.TailLines > 0 {
			logOptions.TailLines = &options.TailLines
		}

		if options.SinceSeconds > 0 {
			logOptions.SinceSeconds = &options.SinceSeconds
		}

		// Request logs
		req := c.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, logOptions)
		logsStream, err := req.Stream(ctx)
		if err != nil {
			// Skip this container if logs aren't available
			continue
		}
		defer logsStream.Close()

		// Read logs
		logsBytes, err := io.ReadAll(logsStream)
		if err != nil {
			return nil, fmt.Errorf("failed to read logs: %w", err)
		}

		logs := string(logsBytes)
		if logs != "" {
			// Add container name header to logs
			containerLogs := fmt.Sprintf("=== Logs for container: %s ===\n%s", container.Name, logs)
			allLogs = append(allLogs, containerLogs)
		}
	}

	return allLogs, nil
}

// extractPodStatus extracts important status information from a pod
func extractPodStatus(pod *corev1.Pod) map[string]string {
	status := make(map[string]string)

	// Basic pod information
	status["phase"] = string(pod.Status.Phase)
	status["hostIP"] = pod.Status.HostIP
	status["podIP"] = pod.Status.PodIP
	status["startTime"] = pod.Status.StartTime.String()

	// Add container statuses
	for i, containerStatus := range pod.Status.ContainerStatuses {
		prefix := fmt.Sprintf("container.%d.", i)
		status[prefix+"name"] = containerStatus.Name
		status[prefix+"ready"] = fmt.Sprintf("%v", containerStatus.Ready)
		status[prefix+"restartCount"] = fmt.Sprintf("%d", containerStatus.RestartCount)

		// Check for container state details
		if containerStatus.State.Waiting != nil {
			status[prefix+"state"] = "waiting"
			status[prefix+"reason"] = containerStatus.State.Waiting.Reason
			status[prefix+"message"] = containerStatus.State.Waiting.Message
		} else if containerStatus.State.Running != nil {
			status[prefix+"state"] = "running"
			status[prefix+"startedAt"] = containerStatus.State.Running.StartedAt.String()
		} else if containerStatus.State.Terminated != nil {
			status[prefix+"state"] = "terminated"
			status[prefix+"reason"] = containerStatus.State.Terminated.Reason
			status[prefix+"exitCode"] = fmt.Sprintf("%d", containerStatus.State.Terminated.ExitCode)
			status[prefix+"message"] = containerStatus.State.Terminated.Message
		}
	}

	// Add conditions
	for i, condition := range pod.Status.Conditions {
		prefix := fmt.Sprintf("condition.%d.", i)
		status[prefix+"type"] = string(condition.Type)
		status[prefix+"status"] = string(condition.Status)
		status[prefix+"reason"] = condition.Reason
		status[prefix+"message"] = condition.Message
	}

	return status
}
