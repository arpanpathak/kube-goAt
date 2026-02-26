package state

import (
	"context"
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestKubernetesStore_GetError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated GET error")
	})
	store := NewKubernetesStore(client, "default")
	err := store.Save(context.Background(), "key", []byte("data"))
	if err == nil {
		t.Error("Expected simulated GET error in Save")
	}
}

func TestKubernetesStore_CreateError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated CREATE error")
	})
	store := NewKubernetesStore(client, "default")
	err := store.Save(context.Background(), "key", []byte("data"))
	if err == nil {
		t.Error("Expected simulated CREATE error in Save")
	}
}

func TestKubernetesStore_UpdateError(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.CoreV1().Secrets("default").Create(context.Background(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "key", Namespace: "default"},
	}, metav1.CreateOptions{})

	client.PrependReactor("update", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("simulated UPDATE error")
	})
	store := NewKubernetesStore(client, "default")
	err := store.Save(context.Background(), "key", []byte("data"))
	if err == nil {
		t.Error("Expected simulated UPDATE error in Save")
	}
}
