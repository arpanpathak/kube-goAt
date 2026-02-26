<p align="center">
  <img src="https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo.svg" alt="kube-goAT Logo" width="120" />
</p>

# 🐐 kube-goAT

> **The Greatest of All Time, RISC-Inspired Infrastructure as Code for Kubernetes.** \
> *Eliminate YAML nesting hell. Embrace fast, type-safe, compiled Go infrastructure.*

[![Go Report Card](https://goreportcard.com/badge/github.com/arpanpathak/kube-goAt)](https://goreportcard.com/report/github.com/arpanpathak/kube-goAt)
[![Go Reference](https://pkg.go.dev/badge/github.com/arpanpathak/kube-goAt.svg)](https://pkg.go.dev/github.com/arpanpathak/kube-goAt)
[![Test Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)]()
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

`kube-goAT` is a radically simplified, AST-first Infrastructure as Code (IaC) library for Kubernetes built entirely in idiomatic Go. Designed as a "Secure by Default" infrastructure compiler, it replaces bloated text manifests with type-safe Go structs and ultra-fast binary reconciliation. 

---

## 📖 Table of Contents
- [The Problem (Why kube-goAT?)](#-the-problem-why-kube-goat)
- [How it Works](#-how-it-works)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Architecture & State Management](#-architecture--state-management)
- [Security by Default](#-security-by-default)
- [Contributing](#-contributing)
- [Code of Conduct](#-code-of-conduct)
- [License](#-license)

---

## 🤔 The Problem (Why kube-goAT?)

Kubernetes is incredibly powerful, but infrastructure code is routinely poorly written.
* **YAML Nesting Hell:** YAML was never meant to be a Turing-complete language. Constantly indenting `spec: template: spec: containers:` is error-prone.
* **Lack of Type Safety:** You don't know your YAML has a typo until you apply it to the cluster and it fails validation.
* **State Management Nightmares:** External IaC tools handle state but are slow, require external databases (S3/DynamoDB), and lack tight integration with the target cluster's native systems.

## ⚡ How it Works

`kube-goAT` is designed from the ground up for **maximum productivity, maintainability, and execution speed**.

1. **Instruction Set vs. Templates:** Treat infrastructure as a set of compiled instructions to be executed by an engine, rather than a loose text document.
2. **Fast Native Binary Serialization:** Eliminate the "Text -> JSON -> Go Struct" parsing tax. `kube-goAT` compiles definitions straight into a microscopic native binary format (`encoding/gob`).
3. **Orthogonal Composition:** Use functional, short, chained methods (`AttachedTo()`, `Port()`) to banish the "Indentation of Doom."
4. **Compile-time Safety:** Required fields (like resource names or ports) are enforced at compile time. 

---

## 📦 Installation

To use `kube-goAT` in your Go project:

```bash
go get github.com/arpanpathak/kube-goAt/...
```

---

## 🚀 Quick Start

Notice how clean and terse the grammar is compared to raw `client-go` or standard YAML.

```go
package main

import (
    "context"
    "log"
    "github.com/arpanpathak/kube-goAT/pkg/dsl"
    "github.com/arpanpathak/kube-goAT/pkg/compiler"
    "github.com/arpanpathak/kube-goAT/pkg/engine"
    "github.com/arpanpathak/kube-goAT/pkg/state"
)

func main() {
    // 1. Definition Phase (Builder SDK)
    app := dsl.NewService("api-gateway", 80, 8080).Label("env", "prod")
    web := dsl.NewDeployment("web-server", "golang:1.24-alpine").
        AttachedTo(app).
        Replicas(3)

    graph := dsl.NewGraph().Add(app).Add(web)

    // 2. Compile Phase (Native Binary AST)
    payload, err := compiler.Compile(graph)
    if err != nil {
        log.Fatalf("Compilation failed: %v", err)
    }
    
    // 3. Execution & State Syncing
    eng, err := engine.NewEngine("/path/to/.kube/config", nil)
    
    // Auto-stores infrastructure memory transparently inside the cluster!
    k8sStore := state.NewKubernetesStore(eng.GetClient(), "default")
    eng.SetStore(k8sStore)

    // Diffs against state, Upserts changes, Deletes removed resources!
    err = eng.Apply(context.Background(), payload, "production-infra-state")
    if err != nil {
        log.Fatalf("Deployment Failed: %v", err)
    }
    log.Println("Successfully rolled out infrastructure.")
}
```

Detailed run-throughs can be found in the [Examples Directory](examples).

---

## 🏛️ Architecture & State Management

`kube-goAT` includes out-of-the-box state storage abstractions to ensure your execution engine operates idempotently. We provide two default integrations:

* **`state.LocalStore`**: For fast local debugging, saves binaries straight to disk.
* **`state.KubernetesStore`**: *The recommended production approach.* Eliminates the need for S3 buckets or DynamoDB tables for state management (unlike Terraform). It safely injects your encoded 500-byte infrastructure state directly into a Kubernetes `Secret` right alongside your resources, ensuring High Availability.

If a resource is removed from your codebase, the Execution Engine detects it missing from the binary payload and forcefully deletes it from the Kubernetes API. Field drift (manual hacking of replicas) triggers automatic Upsert overwrites.

---

## 🛡️ Security by Default
The Kubernetes defaults are famously insecure (root privileges, lack of resource limits). `kube-goAT` acts as a **Compilation Target**, meaning we will be aggressively injecting pre-baked `securityContexts` natively into the execution loops so teams get compliance for free without needing to write tedious security templates.

---

## 🤝 Contributing

We welcome contributions of all shapes and sizes! `kube-goAT` is a growing project and actively looking for help expanding the SDK to cover resource types like `StatefulSets`, `Ingress`, or `CronJobs`.

1. Read our [Contributing Guide](CONTRIBUTING.md) to get started. 
2. Fork the repository, create a branch, and open a PR!
3. If you find a bug or have a syntax suggestion, please open an Issue.

## 💬 Code of Conduct

We are committed to fostering a welcoming, respectful, and harassment-free community. Please review our [Code of Conduct](CODE_OF_CONDUCT.md) before participating.

---

## 📜 License

This project is licensed under the **Apache 2.0** License. See the [LICENSE](LICENSE) file for more information.

*Built with ❤️ for true Computer Science engineering.*
