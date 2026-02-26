package state

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestKubernetesStore(t *testing.T) {
	client := fake.NewSimpleClientset()
	store := NewKubernetesStore(client, "default")
	ctx := context.Background()

	key := "test-k8s-state"
	data := []byte("k8s binary payload")

	// Test Save (Create)
	if err := store.Save(ctx, key, data); err != nil {
		t.Errorf("Save (Create) failed: %v", err)
	}

	// Verify Secret Created
	secret, err := client.CoreV1().Secrets("default").Get(ctx, key, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to retrieve created secret: %v", err)
	}
	if string(secret.Data["state.gob"]) != string(data) {
		t.Errorf("Secret data mismatch: expected %s, got %s", data, secret.Data["state.gob"])
	}

	// Test Save (Update)
	newData := []byte("updated k8s binary payload")
	if err := store.Save(ctx, key, newData); err != nil {
		t.Errorf("Save (Update) failed: %v", err)
	}

	// Test Load
	loaded, err := store.Load(ctx, key)
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if string(loaded) != string(newData) {
		t.Errorf("Expected %s, got %s", newData, loaded)
	}

	// Test Load Non-Existent
	_, err = store.Load(ctx, "does-not-exist")
	if err == nil {
		t.Error("Expected error loading non-existent k8s state")
	}

	// Test Save with existing secret but nil Data map
	nilDataSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "nil-data", Namespace: "default"},
	}
	client.CoreV1().Secrets("default").Create(ctx, nilDataSecret, metav1.CreateOptions{})
	if err := store.Save(ctx, "nil-data", data); err != nil {
		t.Errorf("Save on nil Data secret failed: %v", err)
	}
}
