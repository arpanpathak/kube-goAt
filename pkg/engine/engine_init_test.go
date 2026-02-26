package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewEngineError(t *testing.T) {
	_, err := NewEngine("invalid-kubeconfig-path", nil)
	if err == nil {
		t.Error("Expected error for invalid kubeconfig")
	}
}

func TestNewEngineSuccess(t *testing.T) {
	kubeconfig := []byte(`
apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: test-token
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	os.WriteFile(path, kubeconfig, 0644)

	_, err := NewEngine(path, nil)
	if err != nil {
		t.Errorf("Expected success for valid kubeconfig, got %v", err)
	}
}

func TestNewEngine_ClientCreateError(t *testing.T) {
	kubeconfig := []byte(`
apiVersion: v1
clusters:
- cluster:
    server: ://invalid-url
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: test-token
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "config_err")
	os.WriteFile(path, kubeconfig, 0644)

	_, err := NewEngine(path, nil)
	if err == nil {
		t.Errorf("Expected error for invalid URL in kubeconfig")
	}
}
