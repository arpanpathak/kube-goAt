# kube-goAT
An Infrastructure as code to simplify building cloud native distributed system using Kubernates in more pragmatic way for maximum productivity and maintainability. A RISC-inspired, AST-based Infrastructure as Code library for Kubernetes in Go.


**The RISC Architecture for Infrastructure.**

`kube-goAT` is an AST-first Infrastructure as Code (IaC) library for Kubernetes. It replaces "CISC-like" text manifests (YAML/JSON) with a **Reduced Instruction Set** approach, compiling Go definitions into a **binary-serialized Abstract Syntax Tree**.

## Core Philosophy

* **Instruction Set vs. Templates:** Treat infrastructure as a set of opcodes, not a document.
* **Binary Serialization:** Eliminate the "Text -> JSON -> Go Struct" parsing tax.
* **Orthogonal Composition:** Use functional pipelines to avoid the "Indentation of Doom."
* **Compile-time Safety:** Leverage the Go compiler to validate infrastructure structure before deployment.

## Architecture

1.  **DSL Layer:** Idiomatic Go functions to define resource intent.
2.  **Compiler:** Generates a directed acyclic graph (DAG) representing the AST.
3.  **Serialization:** Encodes the AST into a compact binary format (Protobuf).
4.  **Engine:** A reconciliation loop that executes binary instructions directly against the GKE/K8s API.

## Usage

```go
package main

import "[github.com/arpanpathak/kube-goAT/pkg/dsl](https://github.com/username/kube-goAT/pkg/dsl)"

func main() {
    // Composition via reference, not nesting
    app := dsl.NewService("api-gateway").
        WithPort(80, 8080).
        WithLabel("env", "prod")

    dsl.NewDeployment("web-server").
        AttachedTo(app).
        WithImage("golang:1.24-alpine").
        WithReplicas(5).
        Build() // Serializes to Binary AST
}
```

