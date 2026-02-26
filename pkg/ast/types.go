package ast

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]any{})
	gob.Register(map[string]string{})
	gob.Register([]any{})
	gob.Register([]string{})
	gob.Register(int32(0))
}

// Node represents a generic Kubernetes resource intent.
type Node struct {
	Kind         string
	Name         string
	Namespace    string
	Dependencies []string
	Properties   map[string]any
}

// DAG represents the complete infrastructure graph.
type DAG struct {
	Nodes map[string]*Node
}

// Serialize converts the DAG to a compact binary format using Gob.
func (d *DAG) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize restores the DAG from gob binary format.
func Deserialize(data []byte) (*DAG, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var d DAG
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}
	return &d, nil
}
