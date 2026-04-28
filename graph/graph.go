package graph

import (
	"errors"
	"slices"
)

var (
	// ErrNodeNotFound indicates the node does not exist.
	ErrNodeNotFound = errors.New("graph: node not found")
	// ErrCycleDetected indicates the graph contains a cycle.
	ErrCycleDetected = errors.New("graph: cycle detected")
	// ErrTopologicalSortRequiresDirected indicates topological sort needs a directed graph.
	ErrTopologicalSortRequiresDirected = errors.New("graph: topological sort requires directed graph")
)

type graphNode[K comparable, V any] struct {
	value     V
	neighbors map[K]struct{}
	order     []K
}

// Graph stores nodes with adjacency relationships.
type Graph[K comparable, V any] struct {
	directed  bool
	nodes     map[K]*graphNode[K, V]
	order     []K
	edgeCount int
}

// NewDirectedGraph creates an empty directed graph.
func NewDirectedGraph[K comparable, V any]() *Graph[K, V] {
	return &Graph[K, V]{directed: true}
}

// NewUndirectedGraph creates an empty undirected graph.
func NewUndirectedGraph[K comparable, V any]() *Graph[K, V] {
	return &Graph[K, V]{}
}

// IsDirected reports whether the graph is directed.
func (g *Graph[K, V]) IsDirected() bool {
	return g != nil && g.directed
}

// AddNode inserts or updates one node.
// It returns true when inserted as a new node, false when updated existing node.
func (g *Graph[K, V]) AddNode(id K, value V) bool {
	if g == nil {
		return false
	}
	g.ensureInit()
	if node, ok := g.nodes[id]; ok {
		node.value = value
		return false
	}
	g.nodes[id] = &graphNode[K, V]{
		value:     value,
		neighbors: make(map[K]struct{}),
	}
	g.order = append(g.order, id)
	return true
}

// GetNode returns node value by id.
func (g *Graph[K, V]) GetNode(id K) (V, bool) {
	var zero V
	if g == nil || g.nodes == nil {
		return zero, false
	}
	node, ok := g.nodes[id]
	if !ok {
		return zero, false
	}
	return node.value, true
}

// SetNodeValue updates node value by id.
func (g *Graph[K, V]) SetNodeValue(id K, value V) bool {
	if g == nil || g.nodes == nil {
		return false
	}
	node, ok := g.nodes[id]
	if !ok {
		return false
	}
	node.value = value
	return true
}

// HasNode reports whether id exists.
func (g *Graph[K, V]) HasNode(id K) bool {
	_, ok := g.GetNode(id)
	return ok
}

// DeleteNode removes one node and all incident edges.
func (g *Graph[K, V]) DeleteNode(id K) bool {
	if g == nil || g.nodes == nil {
		return false
	}
	node, ok := g.nodes[id]
	if !ok {
		return false
	}

	if g.directed {
		g.edgeCount -= len(node.order)
		for _, otherID := range g.order {
			if otherID == id {
				continue
			}
			if otherNode, ok := g.nodes[otherID]; ok && otherNode.deleteNeighbor(id) {
				g.edgeCount--
			}
		}
	} else {
		for _, neighborID := range node.order {
			if neighborID != id {
				if neighbor, ok := g.nodes[neighborID]; ok {
					neighbor.deleteNeighbor(id)
				}
			}
			g.edgeCount--
		}
	}

	delete(g.nodes, id)
	g.deleteNodeOrder(id)
	return true
}

// AddEdge inserts one edge.
func (g *Graph[K, V]) AddEdge(from, to K) error {
	if g == nil {
		return ErrNodeNotFound
	}
	g.ensureInit()
	fromNode, ok := g.nodes[from]
	if !ok {
		return ErrNodeNotFound
	}
	toNode, ok := g.nodes[to]
	if !ok {
		return ErrNodeNotFound
	}

	if !fromNode.addNeighbor(to) {
		return nil
	}
	if !g.directed && from != to {
		toNode.addNeighbor(from)
	}
	g.edgeCount++
	return nil
}

// HasEdge reports whether edge exists.
func (g *Graph[K, V]) HasEdge(from, to K) bool {
	if g == nil || g.nodes == nil {
		return false
	}
	node, ok := g.nodes[from]
	if !ok {
		return false
	}
	_, ok = node.neighbors[to]
	return ok
}

// DeleteEdge removes one edge and reports whether it existed.
func (g *Graph[K, V]) DeleteEdge(from, to K) bool {
	if g == nil || g.nodes == nil {
		return false
	}
	fromNode, ok := g.nodes[from]
	if !ok || !fromNode.deleteNeighbor(to) {
		return false
	}
	if !g.directed && from != to {
		if toNode, ok := g.nodes[to]; ok {
			toNode.deleteNeighbor(from)
		}
	}
	g.edgeCount--
	return true
}

// Neighbors returns adjacent node ids in insertion order.
func (g *Graph[K, V]) Neighbors(id K) []K {
	if g == nil || g.nodes == nil {
		return nil
	}
	node, ok := g.nodes[id]
	if !ok || len(node.order) == 0 {
		return nil
	}
	return slices.Clone(node.order)
}

// NodeIDs returns node ids in insertion order.
func (g *Graph[K, V]) NodeIDs() []K {
	if g == nil || len(g.order) == 0 {
		return nil
	}
	return slices.Clone(g.order)
}

// BFS iterates reachable nodes in breadth-first order until visit returns false.
func (g *Graph[K, V]) BFS(start K, visit func(id K, value V) bool) {
	if g == nil || g.nodes == nil || visit == nil {
		return
	}
	if _, ok := g.nodes[start]; !ok {
		return
	}

	visited := make(map[K]struct{}, len(g.nodes))
	queue := []K{start}
	visited[start] = struct{}{}
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		node := g.nodes[id]
		if !visit(id, node.value) {
			return
		}
		for _, neighborID := range node.order {
			if _, seen := visited[neighborID]; seen {
				continue
			}
			visited[neighborID] = struct{}{}
			queue = append(queue, neighborID)
		}
	}
}

// DFS iterates reachable nodes in depth-first pre-order until visit returns false.
func (g *Graph[K, V]) DFS(start K, visit func(id K, value V) bool) {
	if g == nil || g.nodes == nil || visit == nil {
		return
	}
	if _, ok := g.nodes[start]; !ok {
		return
	}

	visited := make(map[K]struct{}, len(g.nodes))
	stack := []K{start}
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, seen := visited[id]; seen {
			continue
		}
		visited[id] = struct{}{}

		node := g.nodes[id]
		if !visit(id, node.value) {
			return
		}
		for i := len(node.order) - 1; i >= 0; i-- {
			neighborID := node.order[i]
			if _, seen := visited[neighborID]; !seen {
				stack = append(stack, neighborID)
			}
		}
	}
}

// TopologicalSort returns node ids in topological order.
func (g *Graph[K, V]) TopologicalSort() ([]K, error) {
	if g == nil || g.nodes == nil || len(g.nodes) == 0 {
		return nil, nil
	}
	if !g.directed {
		return nil, ErrTopologicalSortRequiresDirected
	}

	indegree := make(map[K]int, len(g.nodes))
	for _, id := range g.order {
		indegree[id] = 0
	}
	for _, id := range g.order {
		node := g.nodes[id]
		for _, neighborID := range node.order {
			indegree[neighborID]++
		}
	}

	queue := make([]K, 0, len(g.nodes))
	for _, id := range g.order {
		if indegree[id] == 0 {
			queue = append(queue, id)
		}
	}

	order := make([]K, 0, len(g.nodes))
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		order = append(order, id)
		node := g.nodes[id]
		for _, neighborID := range node.order {
			indegree[neighborID]--
			if indegree[neighborID] == 0 {
				queue = append(queue, neighborID)
			}
		}
	}

	if len(order) != len(g.nodes) {
		return nil, ErrCycleDetected
	}
	return order, nil
}

// Len returns total node count.
func (g *Graph[K, V]) Len() int {
	if g == nil {
		return 0
	}
	return len(g.nodes)
}

// EdgeCount returns total edge count.
func (g *Graph[K, V]) EdgeCount() int {
	if g == nil {
		return 0
	}
	return g.edgeCount
}

// IsEmpty reports whether graph has no nodes.
func (g *Graph[K, V]) IsEmpty() bool {
	return g.Len() == 0
}

// Clear removes all nodes and edges.
func (g *Graph[K, V]) Clear() {
	if g == nil {
		return
	}
	g.nodes = nil
	g.order = nil
	g.edgeCount = 0
}

func (g *Graph[K, V]) ensureInit() {
	if g.nodes == nil {
		g.nodes = make(map[K]*graphNode[K, V])
	}
}

func (g *Graph[K, V]) deleteNodeOrder(id K) {
	for index, current := range g.order {
		if current != id {
			continue
		}
		g.order = append(g.order[:index], g.order[index+1:]...)
		return
	}
}

func (n *graphNode[K, V]) addNeighbor(id K) bool {
	if _, exists := n.neighbors[id]; exists {
		return false
	}
	n.neighbors[id] = struct{}{}
	n.order = append(n.order, id)
	return true
}

func (n *graphNode[K, V]) deleteNeighbor(id K) bool {
	if _, exists := n.neighbors[id]; !exists {
		return false
	}
	delete(n.neighbors, id)
	for index, current := range n.order {
		if current != id {
			continue
		}
		n.order = append(n.order[:index], n.order[index+1:]...)
		break
	}
	return true
}
