# Contributing to `kube-goAT`

First off, thank you for considering contributing to `kube-goAT`! It's people like you that make the open-source community such a fantastic place to learn, inspire, and create.

## How Can I Contribute?

### 1. Reporting Bugs
This section guides you through submitting a bug report. Following these guidelines helps maintainers understand your report, reproduce the behavior, and find related reports.
* Use the GitHub Labels (`bug`, `triage`) to highlight the issue.
* Provide a clear and descriptive title.
* Provide a minimal reproducible Go snippet showing the SDK bug or compilation error.

### 2. Suggesting Enhancements
Enhancement suggestions are highly welcome! Whether it's adding a new Kubernetes object type (like Ingress) or tweaking the Execution Engine's diffing algorithm.
* Use the GitHub Labels (`enhancement`).
* Explain why this enhancement would be useful to most users.

### 3. Pull Requests
1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests! The goal of `kube-goAT` is to remain at **100% statement coverage**.
3. Ensure the test suite passes (`go test -v ./...`).
4. Keep the SDK flat, English-like, and avoid nested configuration hell.
5. Format your code with `gofmt`.

## Core Philosophy Reminders for Contributors
* **Compile-time safety first:** Enforce required Kubernetes fields inside the struct constructors.
* **Keep binary serialization small:** Avoid bloating the `ast.Node` properties.
* **Idempotency:** Whenever touching `pkg/engine`, ensure the code can be re-run 50 times in a row without crashing the Kubernetes API or duplicating resources.

We look forward to reviewing your PRs!
