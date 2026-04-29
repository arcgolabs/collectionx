package graph

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
)

type graphNodeSnapshot[K comparable, V any] struct {
	ID    K `json:"id"`
	Value V `json:"value"`
}

type graphEdgeSnapshot[K comparable] struct {
	From K `json:"from"`
	To   K `json:"to"`
}

type graphSnapshot[K comparable, V any] struct {
	Directed bool                      `json:"directed"`
	Nodes    []graphNodeSnapshot[K, V] `json:"nodes"`
	Edges    []graphEdgeSnapshot[K]    `json:"edges"`
}

func (g *Graph[K, V]) marshalJSONBytes() ([]byte, error) {
	data, err := common.MarshalJSONValue(g.snapshot())
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
	return common.JSONResultString(data, err, "{}")
}

func (g *Graph[K, V]) snapshot() graphSnapshot[K, V] {
	if g == nil || len(g.nodes) == 0 {
		return graphSnapshot[K, V]{Directed: g != nil && g.directed}
	}

	nodes := make([]graphNodeSnapshot[K, V], 0, len(g.order))
	for _, id := range g.order {
		node := g.nodes[id]
		nodes = append(nodes, graphNodeSnapshot[K, V]{ID: id, Value: node.value})
	}

	edges := make([]graphEdgeSnapshot[K], 0, g.edgeCount)
	if g.directed {
		for _, from := range g.order {
			node := g.nodes[from]
			for _, to := range node.order {
				edges = append(edges, graphEdgeSnapshot[K]{From: from, To: to})
			}
		}
		return graphSnapshot[K, V]{Directed: true, Nodes: nodes, Edges: edges}
	}

	emitted := make(map[graphEdgeSnapshot[K]]struct{}, g.edgeCount*2)
	for _, from := range g.order {
		node := g.nodes[from]
		for _, to := range node.order {
			edge := graphEdgeSnapshot[K]{From: from, To: to}
			if _, seen := emitted[edge]; seen {
				continue
			}
			edges = append(edges, edge)
			emitted[edge] = struct{}{}
			emitted[graphEdgeSnapshot[K]{From: to, To: from}] = struct{}{}
		}
	}
	return graphSnapshot[K, V]{Directed: false, Nodes: nodes, Edges: edges}
}
