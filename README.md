# kube-goAT
An Infrastructure as code to simplify building cloud native distributed system using Kubernates in more pragmatic way for maximum productivity and maintainability. A RISC-inspired, AST-based Infrastructure as Code library for Kubernetes in Go.

**The RISC Architecture for Infrastructure.**

`kube-goAT` is an AST-first Infrastructure as Code (IaC) library for Kubernetes. It replaces "CISC-like" text manifests (YAML/JSON) with a **Reduced Instruction Set** approach, compiling Go definitions into a **binary-serialized Abstract Syntax Tree**.

# But WHY WHY WHY kube-goAT ?
Kubernetes can become really hard, and infrastructure code is poorly written using YAML. YAML was never meant to be a language for describing orchestration logic.

I am creating a project under the Apache 2.0 license called kube-goAT. I'll provide idiomatic, easy-to-use, composable Infrastructure as Code which will avoid deep nesting of Spec: Spec:, enforce rigorous type safety, and make your job easier. I'll provide injectable state storage; the state of infra code can be stored in etcd or cloud storage. This will directly compile code into a highly space-optimized binary serialization, with the potential to reduce deployment time by bypassing slow data representational formats such as intermediary JSON, YAML, or XML. An infra kube execution engine could easily decompose them into simple CRUD API calls and apply a state guardrail by fetching states from a highly available source of truth: low-latency key-value pairs. It’s kind of a RISC Architecture for infrastructure. Instead of a passive pile of unreadable YAML, we convert it into a DAG(Directed Acyclic Graph) and feed into a execution engine. Execution engine will use native routines. 

I've worked at Amazon and used CDK. It's a good philosophy, but cdk8s is unable to solve the spec nesting problem; I call it "YAML nesting hell." Working at other companies, I've seen how clunky other DSLs, such as Terraform, can become—especially with the growing need to span different kinds of complex hardware infrastructures tied to money-making optimization business logic. We need a powerful, Turing-complete, statically typed language to reliably run our application backbone with full confidence. Kubernetes is a great, generalized way to deploy anything in Linux as a pod and easily inject monitoring, cloud provider APIs, load balancing, vertical and horizontal scaling, networking rules, backups, and more.

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

