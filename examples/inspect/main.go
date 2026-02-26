package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/arpanpathak/kube-goAT/pkg/ast"
	"github.com/arpanpathak/kube-goAT/pkg/compiler"
	"github.com/arpanpathak/kube-goAT/pkg/dsl"
	"github.com/arpanpathak/kube-goAT/pkg/state"
)

func main() {
	appService := dsl.NewService("api-gateway", 80, 8080).
		Label("env", "prod")

	webDeployment := dsl.NewDeployment("web-server", "golang:1.24-alpine").
		Replicas(3).
		AttachedTo(appService)

	graph := dsl.NewGraph().
		Add(appService).
		Add(webDeployment)

	// Build the graph internally to inspect
	dag := graph.Build()

	// Convert to nice JSON to visualize graphically what memory looks like
	prettyJSON, _ := json.MarshalIndent(dag, "", "  ")
	fmt.Println("🚀 Here is your Compiled Abstract Syntax Tree (DAG) 🚀")
	fmt.Println(string(prettyJSON))

	// Compile to the native GO binary
	binaryPayload, err := compiler.Compile(graph)
	if err != nil {
		log.Fatalf("Compile error: %v", err)
	}
	fmt.Printf("\n📦 Binary Serialized File Size: %d bytes (Extremely Compact!)\n", len(binaryPayload))

	// Manually push into the state-store to populate it for the user
	store := state.NewLocalStore("./state-store")
	err = store.Save(context.Background(), "web-server-infra", binaryPayload)
	if err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}

	fmt.Println("\n✅ Successfully forced the compiled Binary AST into ./state-store/web-server-infra.gob")

	// Validate it's saved by loading it right back from the state store!
	data, _ := store.Load(context.Background(), "web-server-infra")
	restored, _ := ast.Deserialize(data)
	fmt.Printf("🔍 Restored exactly %d nodes directly from the binary state file!\n", len(restored.Nodes))
}
