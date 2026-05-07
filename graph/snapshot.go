package graph

// NodeSnapshot is a detached graph node for serialization adapters.
type NodeSnapshot[K comparable, V any] struct {
	ID    K `json:"id"`
	Value V `json:"value"`
}

// EdgeSnapshot is a detached graph edge for serialization adapters.
type EdgeSnapshot[K comparable] struct {
	From K `json:"from"`
	To   K `json:"to"`
}

// Snapshot is a detached graph representation for serialization adapters.
type Snapshot[K comparable, V any] struct {
	Directed bool                 `json:"directed"`
	Nodes    []NodeSnapshot[K, V] `json:"nodes"`
	Edges    []EdgeSnapshot[K]    `json:"edges"`
}

// Nodes returns detached nodes in insertion order.
func (g *Graph[K, V]) Nodes() []NodeSnapshot[K, V] {
	if g == nil || len(g.nodes) == 0 {
		return nil
	}

	nodes := make([]NodeSnapshot[K, V], 0, len(g.order))
	for _, id := range g.order {
		node := g.nodes[id]
		nodes = append(nodes, NodeSnapshot[K, V]{ID: id, Value: node.value})
	}
	return nodes
}

// Edges returns detached edges in insertion order.
// Undirected graphs include each edge once.
func (g *Graph[K, V]) Edges() []EdgeSnapshot[K] {
	if g == nil || g.edgeCount == 0 {
		return nil
	}

	edges := make([]EdgeSnapshot[K], 0, g.edgeCount)
	g.RangeEdges(func(from, to K) bool {
		edges = append(edges, EdgeSnapshot[K]{From: from, To: to})
		return true
	})
	return edges
}

// Snapshot returns a detached graph representation.
func (g *Graph[K, V]) Snapshot() Snapshot[K, V] {
	return Snapshot[K, V]{
		Directed: g != nil && g.directed,
		Nodes:    g.Nodes(),
		Edges:    g.Edges(),
	}
}
