package dsl

import (
	"reflect"
	"testing"
)

func TestServiceDSL(t *testing.T) {
	svc := NewService("my-svc", 80, 8080).
		Label("key", "val").
		Namespace("test-ns")

	if svc.GetName() != "my-svc" {
		t.Errorf("Expected my-svc, got %v", svc.GetName())
	}

	node := svc.Build()
	if node.Kind != "Service" {
		t.Errorf("Expected Kind Service")
	}
	if node.Namespace != "test-ns" {
		t.Errorf("Expected Namespace test-ns")
	}

	props := node.Properties
	if props["port"] != int32(80) {
		t.Errorf("Expected port 80")
	}

	labels, ok := props["labels"].(map[string]string)
	if !ok || labels["key"] != "val" {
		t.Errorf("Labels map missing or incorrect")
	}
}

func TestDeploymentDSL(t *testing.T) {
	svc := NewService("link-svc", 80, 8080).Label("app", "test")
	dep := NewDeployment("my-dep", "nginx:1.0").
		Replicas(2).
		Namespace("prod").
		AttachedTo(svc)

	if dep.GetName() != "my-dep" {
		t.Errorf("Expected my-dep, got %v", dep.GetName())
	}

	node := dep.Build()
	if node.Kind != "Deployment" {
		t.Errorf("Expected Kind Deployment")
	}
	if node.Namespace != "prod" {
		t.Errorf("Expected Namespace prod")
	}

	if !reflect.DeepEqual(node.Dependencies, []string{"link-svc"}) {
		t.Errorf("Expected dependency on link-svc, got %v", node.Dependencies)
	}

	props := node.Properties
	if props["replicas"] != int32(2) {
		t.Errorf("Expected 2 replicas")
	}

	labels, ok := props["labels"].(map[string]string)
	if !ok || labels["app"] != "test" {
		t.Errorf("Labels map missing or attached incorrectly")
	}
}

func TestGraphBuilder(t *testing.T) {
	svc := NewService("svc1", 80, 8080)
	dep := NewDeployment("dep1", "image")

	graph := NewGraph().Add(svc).Add(dep)
	dag := graph.Build()

	if len(dag.Nodes) != 2 {
		t.Errorf("Expected 2 nodes in DAG, got %v", len(dag.Nodes))
	}
	if _, ok := dag.Nodes["svc1"]; !ok {
		t.Errorf("Missing svc1 in DAG")
	}
	if _, ok := dag.Nodes["dep1"]; !ok {
		t.Errorf("Missing dep1 in DAG")
	}
}
