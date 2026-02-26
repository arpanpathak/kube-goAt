package engine

import (
	"testing"

	"github.com/arpanpathak/kube-goAT/pkg/state"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEngineHelpers(t *testing.T) {
	client := fake.NewSimpleClientset()
	store := state.NewLocalStore(t.TempDir())

	eng := &Engine{client: client}

	if eng.GetClient() != client {
		t.Error("GetClient did not return the expected kubernetes.Interface")
	}

	eng.SetStore(store)
	if eng.store != store {
		t.Error("SetStore did not properly set the store instance")
	}
}
