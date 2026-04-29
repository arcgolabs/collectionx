package tree

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal tree json: nil receiver")
	}

	var roots []jsonNode[K, V]
	if err := json.Unmarshal(data, &roots); err != nil {
		return fmt.Errorf("unmarshal tree json: %w", err)
	}

	next := NewTree[K, V]()
	for _, root := range roots {
		if err := appendJSONNode(next, root, nil); err != nil {
			return fmt.Errorf("unmarshal tree json: %w", err)
		}
	}

	*t = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *ConcurrentTree[K, V]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal concurrent tree json: nil receiver")
	}

	var snapshot Tree[K, V]
	if err := snapshot.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal concurrent tree json: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.tree = snapshot.Clone()
	return nil
}

func appendJSONNode[K comparable, V any](tree *Tree[K, V], node jsonNode[K, V], parent *jsonNode[K, V]) error {
	if parent == nil {
		if err := tree.AddRoot(node.ID, node.Value); err != nil {
			return err
		}
	} else {
		if err := tree.AddChild(parent.ID, node.ID, node.Value); err != nil {
			return err
		}
	}

	for _, child := range node.Children {
		current := node
		if err := appendJSONNode(tree, child, &current); err != nil {
			return err
		}
	}
	return nil
}
