# NGINX Deployment using Kubernetes Native State Storage

This is an advanced example representing the **production-recommended path**. 

Instead of tracking infrastructure diffs and opcodes using S3 or local files (like the `basic` example), this example uses `state.KubernetesStore` to dynamically inject the extremely compact binary `kube-goAT` state snapshot **directly into the deployed cluster as a Kubernetes `Secret`**.

This guarantees your runtime environment and your infrastructure tracking mechanism stay unified without requiring third-party databases.

## Prerequisites

You must have a local Kubernetes cluster running.

### Using Minikube

```bash
minikube start
```

### Using Kind

```bash
kind create cluster
```

## Running the Example

Execute the builder script:

```bash
go run main.go
```

## Expected Output

On the first run, the Execution Engine will notice no matching `Secret` exists under the name `nginx-infra-state`, assume a clean slate, and issue creation API calls.

```
Compiled AST into 536 bytes
2026/02/25 20:58:11 🚀 Firing Binary OpCodes into Minikube...
2026/02/25 20:58:11 [Engine] No existing state found for nginx-infra-state, creating new.
2026/02/25 20:58:11 [Engine] Created Service: nginx-https-svc
2026/02/25 20:58:11 [Engine] Created Deployment: nginx-server
2026/02/25 20:58:11 ✅ Execution completed successfully! Infrastructure is live in Minikube.
```

## Verify the Black-Box Secret Storage

You can manually inspect the Kubernetes cluster to see what `kube-goAT` secretly slipped in alongside your NGINX resources.

```bash
kubectl get secret nginx-infra-state
```

Terminal Output:
```
NAME                TYPE     DATA   AGE
nginx-infra-state   Opaque   1      2m
```

The underlying data tracked within that Secret is actually the compiled `go.gob` payload of your precise SDK instructions. Running `go run main.go` again will pull this secret memory, diff it inside the Go routine runtime, and smartly execute idempotent `UPDATE` cycles when needed. 

## Cleaning Up

```bash
kubectl delete deployment nginx-server
kubectl delete service nginx-https-svc
kubectl delete secret nginx-infra-state
```
