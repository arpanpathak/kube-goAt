# Inspecting the AST Memory Graph

This example bypasses Kubernetes and simply demonstrates what `kube-goAT` is doing under the hood. It compiles an SDK definition of a Service and Deployment, generates the **Directed Acyclic Graph (DAG)** in memory, and serializes it into the native Go binary format (`encoding/gob`).

## What this demonstrates
1. The type-safe, fluent SDK definition structure.
2. What the AST looks like when rendered out to JSON.
3. How incredibly small the native binary format is.
4. How the `state.LocalStore` writes the binary to the local disk.

## Running the Example

No Kubernetes cluster is required to run this inspection script.

```bash
go run main.go
```

## Expected Output

```
🚀 Here is your Compiled Abstract Syntax Tree (DAG) 🚀
{
  "Nodes": {
    "api-gateway": {
      "Kind": "Service",
      "Name": "api-gateway",
      "Namespace": "default",
      "Dependencies": null,
      "Properties": {
        "labels": {
          "env": "prod"
        },
        "port": 80,
        "targetPort": 8080
      }
    },
    "web-server": {
      "Kind": "Deployment",
      "Name": "web-server",
      "Namespace": "default",
      "Dependencies": [
        "api-gateway"
      ],
      "Properties": {
        "image": "golang:1.24-alpine",
        "labels": {
          "env": "prod"
        },
        "replicas": 3
      }
    }
  }
}

📦 Binary Serialized File Size: 508 bytes (Extremely Compact!)

✅ Successfully forced the compiled Binary AST into ./state-store/web-server-infra.gob
🔍 Restored exactly 2 nodes directly from the binary state file!
```
