package state

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesStore implements Store by saving AST state into a Kubernetes Secret.
type KubernetesStore struct {
	client    kubernetes.Interface
	namespace string
}

func NewKubernetesStore(client kubernetes.Interface, ns string) *KubernetesStore {
	return &KubernetesStore{client: client, namespace: ns}
}

// Save writes the binary gob payload to a K8s Secret.
func (k *KubernetesStore) Save(ctx context.Context, key string, data []byte) error {
	secret, err := k.client.CoreV1().Secrets(k.namespace).Get(ctx, key, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: key, Namespace: k.namespace},
			Data:       map[string][]byte{"state.gob": data},
		}
		_, err = k.client.CoreV1().Secrets(k.namespace).Create(ctx, newSecret, metav1.CreateOptions{})
		return err
	} else if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data["state.gob"] = data
	_, err = k.client.CoreV1().Secrets(k.namespace).Update(ctx, secret, metav1.UpdateOptions{})
	return err
}

// Load retrieves the binary gob payload from a K8s Secret.
func (k *KubernetesStore) Load(ctx context.Context, key string) ([]byte, error) {
	secret, err := k.client.CoreV1().Secrets(k.namespace).Get(ctx, key, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret.Data["state.gob"], nil
}
