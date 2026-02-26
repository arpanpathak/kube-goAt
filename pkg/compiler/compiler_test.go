package compiler

import (
	"testing"

	"github.com/arpanpathak/kube-goAT/pkg/dsl"
)

func TestCompile(t *testing.T) {
	svc := dsl.NewService("comp-svc", 80, 8080)
	dep := dsl.NewDeployment("comp-dep", "nginx").AttachedTo(svc)

	graph := dsl.NewGraph().Add(svc).Add(dep)

	payload, err := Compile(graph)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(payload) == 0 {
		t.Fatal("Compiled payload is empty")
	}
}
