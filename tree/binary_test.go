package tree

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestTreeBinaryRoundTrip(t *testing.T) {
	source := NewTree[int, string]()
	if err := source.AddRoot(1, "root"); err != nil {
		t.Fatalf("add root: %v", err)
	}
	if err := source.AddChild(1, 2, "child"); err != nil {
		t.Fatalf("add child: %v", err)
	}

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target Tree[int, string]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	if node, ok := target.Get(2); !ok || node.Value() != "child" {
		t.Fatalf("unexpected node")
	}
}

func TestConcurrentTreeGobRoundTrip(t *testing.T) {
	source := NewConcurrentTree[int, string]()
	if err := source.AddRoot(1, "root"); err != nil {
		t.Fatalf("add root: %v", err)
	}

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target ConcurrentTree[int, string]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	if node, ok := target.Get(1); !ok || node.Value() != "root" {
		t.Fatalf("unexpected node")
	}
}
