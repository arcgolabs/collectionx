package graph

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (g *Graph[K, V]) UnmarshalJSON(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal graph JSON: nil receiver")
	}

	var snap Snapshot[K, V]
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("unmarshal graph JSON: %w", err)
	}

	next := &Graph[K, V]{directed: snap.Directed}
	for _, node := range snap.Nodes {
		next.AddNode(node.ID, node.Value)
	}
	for _, edge := range snap.Edges {
		if err := next.AddEdge(edge.From, edge.To); err != nil {
			return fmt.Errorf("unmarshal graph JSON: %w", err)
		}
	}
	*g = *next
	return nil
}
