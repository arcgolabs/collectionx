package tree

import (
	"encoding/json"
	"testing"
)

func TestTreeJSONRoundTrip(t *testing.T) {
	source := NewTree[int, string]()
	if err := source.AddRoot(1, "root"); err != nil {
		t.Fatalf("add root: %v", err)
	}
	if err := source.AddChild(1, 2, "left"); err != nil {
		t.Fatalf("add child left: %v", err)
	}
	if err := source.AddChild(1, 3, "right"); err != nil {
		t.Fatalf("add child right: %v", err)
	}

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal tree: %v", err)
	}

	var target Tree[int, string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal tree: %v", err)
	}

	if target.Len() != 3 {
		t.Fatalf("unexpected tree size: %d", target.Len())
	}
	if node, ok := target.Get(1); !ok || node.Value() != "root" {
		t.Fatalf("unexpected root node")
	}
	if children := target.Children(1); len(children) != 2 || children[0].ID() != 2 || children[1].ID() != 3 {
		t.Fatalf("unexpected tree children: %#v", children)
	}
}

func TestConcurrentTreeJSONRoundTrip(t *testing.T) {
	source := NewConcurrentTree[int, string]()
	if err := source.AddRoot(1, "root"); err != nil {
		t.Fatalf("add root: %v", err)
	}
	if err := source.AddChild(1, 2, "child"); err != nil {
		t.Fatalf("add child: %v", err)
	}

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent tree: %v", err)
	}

	var target ConcurrentTree[int, string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent tree: %v", err)
	}

	if target.Len() != 2 {
		t.Fatalf("unexpected concurrent tree size: %d", target.Len())
	}
	if node, ok := target.Get(2); !ok || node.Value() != "child" {
		t.Fatalf("unexpected child node")
	}
}

func TestTreeUnmarshalRejectsDuplicateIDs(t *testing.T) {
	var target Tree[int, string]
	err := json.Unmarshal([]byte(`[{"id":1,"value":"root"},{"id":1,"value":"dup"}]`), &target)
	if err == nil {
		t.Fatal("expected duplicate id error")
	}
}
