package engine

import (
	"context"
	"fmt"
	"log"

	"github.com/arpanpathak/kube-goAT/pkg/ast"
	"github.com/arpanpathak/kube-goAT/pkg/state"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Engine struct {
	client kubernetes.Interface
	store  state.Store
}

func NewEngine(kubeconfig string, store state.Store) (*Engine, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Engine{client: clientset, store: store}, nil
}

// GetClient returns the underlying Kubernetes interface.
func (e *Engine) GetClient() kubernetes.Interface {
	return e.client
}

// SetStore dynamically updates the state storage backend.
func (e *Engine) SetStore(store state.Store) {
	e.store = store
}

// Apply takes a binary Gob AST, compares it to the tracked state, and creates/updates K8s resources.
func (e *Engine) Apply(ctx context.Context, payload []byte, stateKey string) error {
	// Deserialization of the "RISC" binary instructions.
	dag, err := ast.Deserialize(payload)
	if err != nil {
		return fmt.Errorf("failed to deserialize AST: %w", err)
	}

	// State Check Guardrails
	existingState, err := e.store.Load(ctx, stateKey)
	if err == nil {
		log.Printf("[Engine] Loaded existing state for %s (%d bytes)", stateKey, len(existingState))
	} else {
		log.Printf("[Engine] No existing state found for %s, creating new.", stateKey)
	}

	// Execution Loop
	for _, node := range dag.Nodes {
		switch node.Kind {
		case "Service":
			if err := e.applyService(ctx, node); err != nil {
				return err
			}
		case "Deployment":
			if err := e.applyDeployment(ctx, node); err != nil {
				return err
			}
		default:
			log.Printf("[WARNING] Unsupported node kind: %s", node.Kind)
		}
	}

	// Finalize State Record
	return e.store.Save(ctx, stateKey, payload)
}

func (e *Engine) applyService(ctx context.Context, node *ast.Node) error {
	var port, targetPort int32
	if p, ok := node.Properties["port"].(int32); ok {
		port = p
	}
	if p, ok := node.Properties["targetPort"].(int32); ok {
		targetPort = p
	}

	labels := make(map[string]string)
	if l, ok := node.Properties["labels"].(map[string]string); ok {
		for k, v := range l {
			labels[k] = v
		}
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: node.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{Port: port, TargetPort: intstr.FromInt32(targetPort)},
			},
		},
	}

	_, err := e.client.CoreV1().Services(node.Namespace).Get(ctx, node.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = e.client.CoreV1().Services(node.Namespace).Create(ctx, svc, metav1.CreateOptions{})
		log.Printf("[Engine] Created Service: %s", node.Name)
		return err
	} else if err != nil {
		return err
	}

	// Upsert update logic omitted for initial concise implementation
	log.Printf("[Engine] Service %s already exists.", node.Name)
	return nil
}

func (e *Engine) applyDeployment(ctx context.Context, node *ast.Node) error {
	image := "nginx:latest"
	if img, ok := node.Properties["image"].(string); ok {
		image = img
	}

	var replicas int32 = 1
	if r, ok := node.Properties["replicas"].(int32); ok {
		replicas = r
	}

	labels := make(map[string]string)
	if l, ok := node.Properties["labels"].(map[string]string); ok {
		for k, v := range l {
			labels[k] = v
		}
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: node.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "app", Image: image},
					},
				},
			},
		},
	}

	_, err := e.client.AppsV1().Deployments(node.Namespace).Get(ctx, node.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = e.client.AppsV1().Deployments(node.Namespace).Create(ctx, dep, metav1.CreateOptions{})
		log.Printf("[Engine] Created Deployment: %s", node.Name)
		return err
	} else if err != nil {
		return err
	}

	log.Printf("[Engine] Deployment %s already exists.", node.Name)
	return nil
}
