<div align="center">
  <h1>🐐 kube-goAT</h1>
  <p><b>The Greatest of All Time, RISC-Inspired Infrastructure as Code for Kubernetes</b></p>
  <p><i>Eliminate YAML nesting hell. Embrace fast, type-safe, compiled Go infrastructure.</i></p>
</div>

---

`kube-goAT` is a radically simplified, AST-first Infrastructure as Code (IaC) library for Kubernetes built entirely in idiomatic Go. 

It replaces "CISC-like" text manifests (YAML/JSON) with a **Reduced Instruction Set** approach, compiling native, type-safe Go structural definitions directly into a **space-optimized binary-serialized Abstract Syntax Tree (`encoding/gob`)**.

## 🤔 The Problem: Why Another IaC?

Kubernetes is incredibly powerful, but infrastructure code is routinely poorly written using YAML or complex text-based DSLs. 
* **YAML Nesting Hell:** YAML was never meant to be a Turing-complete language for describing dynamic orchestration logic. Constantly indenting `spec: template: spec: containers:` is error-prone.
* **Lack of Type Safety:** You often don't know your YAML has a typo until you apply it to the cluster and it fails validation.
* **Bloated Tooling:** Existing Go SDKs (like CDK8s) suffer from massive dependencies and still generate clunky JSON/YAML intermediary files that must be slowly applied by `kubectl`.
* **State Management Nightmares:** Tools like Terraform handle state but are slow, require external databases or paid backends, and lack tight integration with the target cluster's native systems.

## ✨ The Solution: kube-goAT

`kube-goAT` is designed from the ground up for **maximum productivity, maintainability, and execution speed**.

### Core Philosophy

1. **Instruction Set vs. Templates:** Treat infrastructure as a set of compiled instructions to be executed by an engine, rather than a text document to be parsed.
2. **Fast Native Binary Serialization:** Eliminate the "Text -> JSON -> Go Struct" parsing tax. `kube-goAT` compiles definitions straight into a native binary format (`encoding/gob`).
3. **Orthogonal Composition:** Use functional, short, chained methods (like `AttachedTo()`, `Port()`) to banish the "Indentation of Doom."
4. **Compile-time Safety:** Leverage the Go compiler's strictness. Required fields (e.g., resource names and image tags) are enforced in constructor signatures before the code even runs.

## 🏛️ Architecture

`kube-goAT` operates in four distinct phases:

1. **The SDK Layer:** Idiomatic, flat Go functions used to define your resource intent. It emphasizes short builder patterns and explicit graph mappings (`AttachedTo()`).
2. **The Compiler:** Introspects the SDK layout and generates a Directed Acyclic Graph (DAG) representing the AST instructions.
3. **The Serializer:** Encodes the memory DAG into an extremely compact, native Go binary format (`gob`).
4. **The Execution Engine:** A built-in reconciliation loop that directly talks to the Kubernetes API, applying state guardrails using built-in Kubernetes Secrets or local files.

---

## 🚀 Examples

### 1. The Basics: Service & Deployment
Notice how clean and terse the grammar is compared to raw `client-go` or standard YAML.

```go
package main

import (
    "github.com/arpanpathak/kube-goAT/pkg/dsl"
    "github.com/arpanpathak/kube-goAT/pkg/compiler"
)

func main() {
    // 1. Definition Phase 
    // Composition via clean functional references, not deep nesting!
    app := dsl.NewService("api-gateway", 80, 8080).
        Label("env", "prod")

    // The deployment automatically inherits labels and sets up dependencies 
    // by simply referencing the Service object.
    web := dsl.NewDeployment("web-server", "golang:1.24-alpine").
        AttachedTo(app).
        Replicas(5)

    graph := dsl.NewGraph().Add(app).Add(web)

    // 2. Compile to native binary AST
    payload, err := compiler.Compile(graph)
    if err != nil {
        panic(err)
    }
    
    // Payload is now a tiny []byte ready for the Execution Engine!
}
```

### 2. Built-in State Management & Execution
Deploying your infrastructure shouldn't require installing a hulking CLI. You can embed the engine directly in your deployment pipelines.

```go
package main

import (
    "context"
    "log"
    // ... basic example imports ...
    "github.com/arpanpathak/kube-goAT/pkg/engine"
    "github.com/arpanpathak/kube-goAT/pkg/state"
)

func main() {
    // ... Build and Compile step from previous example ...
    payload, _ := compiler.Compile(graph)

    // Setup Kubernetes Execution Engine
    // State is automatically tracked locally (state.NewLocalStore) 
    // OR securely directly inside the cluster (state.NewKubernetesStore)
    store := state.NewLocalStore("./state-store")
    
    // Connects to your standard ~/.kube/config cluster
    eng, err := engine.NewEngine("/path/to/.kube/config", store)
    if err != nil {
        log.Fatalf("Failed to initialize Engine: %v", err)
    }
    
    // Automatically diffs the generated Binary AST against tracked state 
    // and applies ONLY what changed via basic CRUD logic.
    err = eng.Apply(context.Background(), payload, "production-infrastructure")
    if err != nil {
        log.Fatalf("Deployment Failed: %v", err)
    }
    
    log.Println("Successfully rolled out infrastructure.")
}
```

## 📦 State Management Guardrails

`kube-goAT` includes out-of-the-box state storage abstractions to ensure your execution engine operates idempotently without overwriting cluster resources blindly.

* **`state.LocalStore`**: Great for CI environments or local fast iteration, saves binaries straight to disk.
* **`state.KubernetesStore`**: *The recommended production approach.* Eliminates the need for S3 buckets or DynamoDB tables for state management (unlike Terraform). It safely injects your encoded infrastructure state directly into a Kubernetes `Secret` right alongside your resources, ensuring High Availability.

## 🛠 Contributing

If you'd like to help expand the SDK to cover things like `StatefulSets`, `Ingress`, or `CronJobs`, I'd love contributions! `kube-goAT` is meant to grow while remaining staunchly committed to avoiding boilerplate bloat.

## 📜 License
Apache 2.0
