package graph

import (
	"fmt"
)

func (g *Graph[K, V]) marshalJSONBytes() ([]byte, error) {
	data, err := marshalJSONValue(g.Snapshot())
	if err != nil {
		return nil, fmt.Errorf("marshal graph JSON: %w", err)
	}
	return data, nil
}

// MarshalJSON implements json.Marshaler.
func (g *Graph[K, V]) MarshalJSON() ([]byte, error) {
	data, err := g.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal graph: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (g *Graph[K, V]) String() string {
	data, err := g.marshalJSONBytes()
	return jsonResultString(data, err, "{}")
}
