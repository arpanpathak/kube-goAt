# Basic Local Execution Example

This example demonstrates how to compile `kube-goAT` infrastructure instructions and execute them against a live Kubernetes cluster using the `LocalStore` state tracker (which writes the state guardrail binary directly to `./state-store` on your hard drive).

## Prerequisites

You must have a local Kubernetes cluster running and active in your `~/.kube/config`. You can use `minikube` or `kind`.

### Using Minikube (Recommended)

```bash
# Start a local node
minikube start

# Verify connection
kubectl cluster-info
```

### Using Kind

```bash
# Create a kind cluster
kind create cluster

# Verify connection
kubectl cluster-info
```

## Running the Example

Once your cluster is active and serving traffic locally, run the SDK engine:

```bash
go run main.go
```

## Expected Output

The first time you run the script, it will create the Service and Deployment, and store the binary state footprint inside `./state-store`.

```
Compiled AST into 328 bytes
2026/02/26 13:00:00 Applying Infrastructure...
2026/02/26 13:00:00 [Engine] No existing state found for web-server-infra, creating new.
2026/02/26 13:00:00 [Engine] Created Service: api-gateway
2026/02/26 13:00:00 [Engine] Created Deployment: web-server
2026/02/26 13:00:00 Execution completed successfully.
```

If you run the script a **second time**, the Engine will diff the incoming compiled binary AST against the physical `.gob` file sitting in `./state-store/web-server-infra.gob` and realize the infrastructure is already matching!

```
Compiled AST into 328 bytes
2026/02/26 13:01:00 Applying Infrastructure...
2026/02/26 13:01:00 [Engine] Loaded existing state for web-server-infra (328 bytes)
2026/02/26 13:01:00 [Engine] Service api-gateway already exists.
2026/02/26 13:01:00 [Engine] Deployment web-server already exists.
2026/02/26 13:00:00 Execution completed successfully.
```

## Cleaning Up

```bash
kubectl delete deployment web-server
kubectl delete service api-gateway
rm -rf state-store/
```
