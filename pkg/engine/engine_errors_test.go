package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/arpanpathak/kube-goAT/pkg/dsl"
	"github.com/arpanpathak/kube-goAT/pkg/state"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestEngineApply_ServiceAPIError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("get", "services", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated API error")
	})
	store := state.NewLocalStore(t.TempDir())
	eng := &Engine{client: client, store: store}

	svc := dsl.NewService("test-svc", 80, 8080)
	payload, _ := dsl.NewGraph().Add(svc).Build().Serialize()
	err := eng.Apply(context.Background(), payload, "key")
	if err == nil {
		t.Error("Expected simulated API error")
	}
}

func TestEngineApply_ServiceCreateError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "services", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated Create error")
	})
	store := state.NewLocalStore(t.TempDir())
	eng := &Engine{client: client, store: store}

	svc := dsl.NewService("test-svc", 80, 8080)
	payload, _ := dsl.NewGraph().Add(svc).Build().Serialize()
	err := eng.Apply(context.Background(), payload, "key")
	if err == nil {
		t.Error("Expected simulated Create error")
	}
}

func TestEngineApply_DeploymentAPIError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("get", "deployments", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated API error")
	})
	store := state.NewLocalStore(t.TempDir())
	eng := &Engine{client: client, store: store}

	dep := dsl.NewDeployment("test-dep", "nginx")
	payload, _ := dsl.NewGraph().Add(dep).Build().Serialize()
	err := eng.Apply(context.Background(), payload, "key")
	if err == nil {
		t.Error("Expected simulated API error")
	}
}

func TestEngineApply_DeploymentCreateError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "deployments", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated Create error")
	})
	store := state.NewLocalStore(t.TempDir())
	eng := &Engine{client: client, store: store}

	dep := dsl.NewDeployment("test-dep", "nginx")
	payload, _ := dsl.NewGraph().Add(dep).Build().Serialize()
	err := eng.Apply(context.Background(), payload, "key")
	if err == nil {
		t.Error("Expected simulated Create error")
	}
}
