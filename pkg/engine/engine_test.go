package engine

import (
	"context"
	"testing"

	"github.com/arpanpathak/kube-goAT/pkg/ast"
	"github.com/arpanpathak/kube-goAT/pkg/dsl"
	"github.com/arpanpathak/kube-goAT/pkg/state"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEngineApply(t *testing.T) {
	client := fake.NewSimpleClientset()
	store := state.NewLocalStore(t.TempDir())

	// Create engine instance manually to bypass kubeconfig requirement for testing
	eng := &Engine{
		client: client,
		store:  store,
	}

	ctx := context.Background()
	stateKey := "test-engine"

	// 1. Build a mock AST payload
	svc := dsl.NewService("test-svc", 80, 8080).Label("app", "test")
	dep := dsl.NewDeployment("test-dep", "nginx:1.0").Replicas(2).AttachedTo(svc)

	dag := dsl.NewGraph().Add(svc).Add(dep).Build()
	payload, err := dag.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize test DAG: %v", err)
	}

	// 2. Test Apply
	err = eng.Apply(ctx, payload, stateKey)
	if err != nil {
		t.Fatalf("Engine apply failed: %v", err)
	}

	// 3. Verify Kubernetes Actions
	// Service created
	s, err := client.CoreV1().Services("default").Get(ctx, "test-svc", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected service test-svc to be created")
	}
	if s.Labels["app"] != "test" {
		t.Errorf("Service missing label")
	}

	// Deployment created
	d, err := client.AppsV1().Deployments("default").Get(ctx, "test-dep", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected deployment test-dep to be created")
	}
	if *d.Spec.Replicas != 2 {
		t.Errorf("Deployment replicas expected 2, got %d", *d.Spec.Replicas)
	}

	// 4. Test State File Written
	loadedState, err := store.Load(ctx, stateKey)
	if err != nil {
		t.Fatalf("State not saved")
	}
	if string(loadedState) != string(payload) {
		t.Errorf("Saved state mismatch")
	}

	// 5. Test Apply Again (Idempotent update logic branches)
	err = eng.Apply(ctx, payload, stateKey)
	if err != nil {
		t.Fatalf("Engine second apply failed (Idempotency check): %v", err)
	}

	// 6. Test with Deserialization Error
	err = eng.Apply(ctx, []byte("garbage"), stateKey)
	if err == nil {
		t.Errorf("Expected apply to fail on bad payload")
	}
}

func TestEngineApply_UnsupportedKind(t *testing.T) {
	client := fake.NewSimpleClientset()
	store := state.NewLocalStore(t.TempDir())
	eng := &Engine{client: client, store: store}

	dag := &ast.DAG{
		Nodes: map[string]*ast.Node{
			"bad": {Kind: "UnknownKind", Name: "bad", Namespace: "default"},
		},
	}
	payload, _ := dag.Serialize()

	// Should log warning and succeed
	err := eng.Apply(context.Background(), payload, "test-key")
	if err != nil {
		t.Errorf("Expected unsupported kinds to be skipped but apply failed: %v", err)
	}
}
