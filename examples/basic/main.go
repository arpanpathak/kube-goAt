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
	// Very concise, English-like DSL with no verbose nesting.
	// Required fields are validated at compile-time in the constructor.
	appService := dsl.NewService("api-gateway", 80, 8080).
		Label("env", "prod")

	webDeployment := dsl.NewDeployment("web-server", "nginx:latest").
		Replicas(3).
		AttachedTo(appService)

	// Build the graph
	graph := dsl.NewGraph().
		Add(appService).
		Add(webDeployment)

	// 2. Compilation Phase
	binaryPayload, err := compiler.Compile(graph)
	if err != nil {
		log.Fatalf("Compile error: %v", err)
	}
	fmt.Printf("Compiled AST into %d bytes\n", len(binaryPayload))

	// 3. Execution Phase
	// Setup local state directory for demo purposes (can easily swap to state.NewKubernetesStore)
	store := state.NewLocalStore("./state-store")

	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// Create Engine which connects to the local cluster
	eng, err := engine.NewEngine(kubeconfig, store)
	if err != nil {
		log.Printf("Engine init error: %v (is kubernetes running?)", err)
		log.Printf("Successfully compiled AST, exiting before cluster comms...")
		return
	}

	ctx := context.Background()
	log.Println("Applying Infrastructure...")
	if err := eng.Apply(ctx, binaryPayload, "web-server-infra"); err != nil {
		log.Printf("Execution failed: %v", err)
	} else {
		log.Println("Execution completed successfully.")
	}
}
