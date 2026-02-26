package dsl

import "github.com/arpanpathak/kube-goAT/pkg/ast"

type Service struct {
	name       string
	namespace  string
	port       int32
	targetPort int32
	labels     map[string]string
}

// NewService enforces compile-time validation for required fields: name, port, targetPort.
func NewService(name string, port, targetPort int32) *Service {
	return &Service{
		name:       name,
		namespace:  "default",
		port:       port,
		targetPort: targetPort,
		labels:     make(map[string]string),
	}
}

// Label adds a key-value pair. Method is terse to be less bloated.
func (s *Service) Label(key, value string) *Service {
	s.labels[key] = value
	return s
}

func (s *Service) Namespace(ns string) *Service {
	s.namespace = ns
	return s
}

func (s *Service) GetName() string {
	return s.name
}

// Build compiles the declarative builder into a graph Node.
func (s *Service) Build() *ast.Node {
	return &ast.Node{
		Kind:      "Service",
		Name:      s.name,
		Namespace: s.namespace,
		Properties: map[string]any{
			"port":       s.port,
			"targetPort": s.targetPort,
			"labels":     s.labels,
		},
	}
}
