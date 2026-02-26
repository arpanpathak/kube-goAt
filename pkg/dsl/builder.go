package dsl

import "github.com/arpanpathak/kube-goAT/pkg/ast"

// Builder is the interface implemented by all DSL resource generators.
type Builder interface {
	Build() *ast.Node
	GetName() string
}

// GraphBuilder composes multiple individual resource builders into a complete DAG.
type GraphBuilder struct {
	resources []Builder
}

func NewGraph() *GraphBuilder {
	return &GraphBuilder{}
}

func (g *GraphBuilder) Add(b Builder) *GraphBuilder {
	g.resources = append(g.resources, b)
	return g
}

// Build generates the final acyclic graph representing the infrastructure.
func (g *GraphBuilder) Build() *ast.DAG {
	dag := &ast.DAG{
		Nodes: make(map[string]*ast.Node),
	}
	for _, res := range g.resources {
		node := res.Build()
		dag.Nodes[node.Name] = node
	}
	return dag
}
