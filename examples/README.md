# K8sMed Examples

This directory contains examples and sample files to help you get started with K8sMed.

## Manifest Examples

The `manifests` directory contains example Kubernetes resources that demonstrate common error scenarios that K8sMed can help diagnose.

### Error Pod Examples

The [error-pod.yaml](manifests/error-pod.yaml) file contains three different pods designed to create common error conditions:

1. **ImagePullBackOff Example**
   ```bash
   kubectl apply -f examples/manifests/error-pod.yaml -n default
   kubectl get pods -n default | grep test-error-pod
   kubectl k8smed analyze pod test-error-pod -n default
   ```
   This example creates a pod that references a non-existent container image, resulting in an ImagePullBackOff error.

2. **CrashLoopBackOff Example**
   ```bash
   kubectl apply -f examples/manifests/error-pod.yaml -n default
   kubectl get pods -n default | grep crash-loop-pod
   kubectl k8smed analyze pod crash-loop-pod -n default
   ```
   This example creates a pod that runs a container that crashes after 5 seconds, resulting in a CrashLoopBackOff error.

3. **OOMKilled Example**
   ```bash
   kubectl apply -f examples/manifests/error-pod.yaml -n default
   kubectl get pods -n default | grep oom-kill-pod
   kubectl k8smed analyze pod oom-kill-pod -n default
   ```
   This example creates a pod that attempts to allocate more memory than its limit, resulting in an OOMKilled error.

## Cleaning Up

To remove all the example pods:

```bash
kubectl delete -f examples/manifests/error-pod.yaml -n default
```

## More Examples

For additional examples and use cases, check out the [docs/examples](../docs/examples/) directory, which includes:

- [Basic Usage Guide](../docs/examples/basic-usage.md)
- [Integration with Gemma AI](../docs/examples/gemma-integration.md) 