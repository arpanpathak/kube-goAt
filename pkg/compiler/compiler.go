package compiler

import (
	"github.com/arpanpathak/kube-goAT/pkg/dsl"
)

// Compile turns a GraphBuilder into a serialized binary payload.
// This decouples the DSL formulation from the final gob encoding if needed.
func Compile(g *dsl.GraphBuilder) ([]byte, error) {
	dag := g.Build()
	return dag.Serialize()
}
