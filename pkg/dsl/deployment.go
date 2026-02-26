package dsl

import "github.com/arpanpathak/kube-goAT/pkg/ast"

type Deployment struct {
	name      string
	namespace string
	image     string
	replicas  int32
	labels    map[string]string
	dependsOn []string
}

// NewDeployment enforces compile-time validation for required fields: name, image.
func NewDeployment(name, image string) *Deployment {
	return &Deployment{
		name:      name,
		namespace: "default",
		image:     image,
		replicas:  1,
		labels:    make(map[string]string),
	}
}

func (d *Deployment) Replicas(n int32) *Deployment {
	d.replicas = n
	return d
}

func (d *Deployment) Label(key, value string) *Deployment {
	d.labels[key] = value
	return d
}

func (d *Deployment) Namespace(ns string) *Deployment {
	d.namespace = ns
	return d
}

// AttachedTo is an idiomatic way to reference another resource.
// It implicitly links labels and enforces graph dependencies.
func (d *Deployment) AttachedTo(svc *Service) *Deployment {
	for k, v := range svc.labels {
		d.labels[k] = v
	}
	d.dependsOn = append(d.dependsOn, svc.GetName())
	return d
}

func (d *Deployment) GetName() string {
	return d.name
}

// Build compiles the declarative builder into a graph Node.
func (d *Deployment) Build() *ast.Node {
	return &ast.Node{
		Kind:         "Deployment",
		Name:         d.name,
		Namespace:    d.namespace,
		Dependencies: d.dependsOn,
		Properties: map[string]any{
			"image":    d.image,
			"replicas": d.replicas,
			"labels":   d.labels,
		},
	}
}
