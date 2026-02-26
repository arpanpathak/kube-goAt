package ast

import (
	"reflect"
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	dag := &DAG{
		Nodes: map[string]*Node{
			"test-node": {
				Kind:      "Service",
				Name:      "test-service",
				Namespace: "default",
				Properties: map[string]any{
					"port":   int32(80),
					"labels": map[string]string{"env": "test"},
				},
			},
		},
	}

	payload, err := dag.Serialize()
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}
	if len(payload) == 0 {
		t.Fatalf("Payload is empty")
	}

	decoded, err := Deserialize(payload)
	if err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if decoded.Nodes["test-node"].Name != "test-service" {
		t.Errorf("Expected Name test-service, got %v", decoded.Nodes["test-node"].Name)
	}
	if !reflect.DeepEqual(dag.Nodes["test-node"].Properties, decoded.Nodes["test-node"].Properties) {
		t.Errorf("Properties mismatch after deserialization")
	}
}

func TestDeserialize_Error(t *testing.T) {
	_, err := Deserialize([]byte("invalid garbage data"))
	if err == nil {
		t.Error("Expected error when deserializing garbage data")
	}
}
