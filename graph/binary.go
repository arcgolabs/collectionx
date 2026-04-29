package graph

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (g *Graph[K, V]) MarshalBinary() ([]byte, error) {
	data, err := common.MarshalBinaryValue(g.snapshot())
	if err != nil {
		return nil, fmt.Errorf("marshal graph binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (g *Graph[K, V]) GobEncode() ([]byte, error) {
	return g.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (g *Graph[K, V]) UnmarshalBinary(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal graph binary: nil receiver")
	}

	var snap graphSnapshot[K, V]
	if err := common.UnmarshalBinaryValue(data, &snap); err != nil {
		return fmt.Errorf("unmarshal graph binary: %w", err)
	}

	next := &Graph[K, V]{directed: snap.Directed}
	for _, node := range snap.Nodes {
		next.AddNode(node.ID, node.Value)
	}
	for _, edge := range snap.Edges {
		if err := next.AddEdge(edge.From, edge.To); err != nil {
			return fmt.Errorf("unmarshal graph binary: %w", err)
		}
	}
	*g = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (g *Graph[K, V]) GobDecode(data []byte) error {
	return g.UnmarshalBinary(data)
}
