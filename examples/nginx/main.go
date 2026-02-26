package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arpanpathak/kube-goAT/pkg/compiler"
	"github.com/arpanpathak/kube-goAT/pkg/dsl"
	"github.com/arpanpathak/kube-goAT/pkg/engine"
	"github.com/arpanpathak/kube-goAT/pkg/state"
)

func main() {
	// 1. Definition Phase
	// Define the NGINX Service. We'll expose port 443 (HTTPS) targeting 443
	// For minikube, we'll just use a standard service for now.
	nginxSvc := dsl.NewService("nginx-https-svc", 443, 443).
		Label("app", "nginx-secure")

	// Define the actual NGINX Deployment
	nginxDep := dsl.NewDeployment("nginx-server", "nginx:latest").
		AttachedTo(nginxSvc).
		Replicas(2)

	// Note: In a real-world secure HTTPS NGINX deployment, we would also need
	// a ConfigMap with the nginx.conf and a Secret with the TLS certs mounted
	// into the Pod as volumes. But for demonstrating the SDK's current
	// capabilities (Service + Deployment), we'll deploy the standard image.

	graph := dsl.NewGraph().
		Add(nginxSvc).
		Add(nginxDep)

	// 2. Compilation Phase
	binaryPayload, err := compiler.Compile(graph)
	if err != nil {
		log.Fatalf("Compile error: %v", err)
	}
	fmt.Printf("Compiled AST into %d bytes\n", len(binaryPayload))

	// 3. Execution Phase against Minikube
	// Let's use the Kubernetes Secret state store to demonstrate
	// production-grade state management directly inside Minikube!
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// The engine initializes the Kubernetes client using the local ~/.kube/config (Minikube)
	// We'll pass the client down to the KubernetesStore so the engine can manage its own state!
	eng, err := engine.NewEngine(kubeconfig, nil) // Temporarily pass nil for the store
	if err != nil {
		log.Fatalf("Minikube connection error: %v (is minikube running?)", err)
	}

	// Create the Kubernetes Secret state store using the engine's initialized Kube client
	// We'll store the infrastructure state in the "default" namespace
	k8sStore := state.NewKubernetesStore(eng.GetClient(), "default")
	eng.SetStore(k8sStore)

	ctx := context.Background()
	log.Println("🚀 Firing Binary OpCodes into Minikube...")

	// Apply the infrastructure graph!
	// The state guardrail will use a Kubernetes Secret named "nginx-infra-state"
	if err := eng.Apply(ctx, binaryPayload, "nginx-infra-state"); err != nil {
		log.Fatalf("Execution failed: %v", err)
	} else {
		log.Println("✅ Execution completed successfully! Infrastructure is live in Minikube.")
	}
}
