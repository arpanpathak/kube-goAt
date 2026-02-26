package dsl

import "testing"

func TestDeployment_Label(t *testing.T) {
	dep := NewDeployment("test", "img").Label("foo", "bar")
	if dep.labels["foo"] != "bar" {
		t.Error("Label not set correctly")
	}
}
