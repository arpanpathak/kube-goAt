package ast

import "testing"

func TestSerialize_Error(t *testing.T) {
	dag := &DAG{
		Nodes: map[string]*Node{
			"bad": {
				Properties: map[string]any{"chan": make(chan int)},
			},
		},
	}
	_, err := dag.Serialize()
	if err == nil {
		t.Error("Expected error when serializing un-encodable data")
	}
}
