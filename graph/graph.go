package graph

import (
	"errors"
	"slices"

	collectionlist "github.com/arcgolabs/collectionx/list"
	collectionmapping "github.com/arcgolabs/collectionx/mapping"
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
	neighbors collectionmapping.Map[K, struct{}]
	order     collectionlist.List[K]
	visitMark uint64
}

// Graph stores nodes with adjacency relationships.
type Graph[K comparable, V any] struct {
	directed  bool
	nodes     collectionmapping.OrderedMap[K, *graphNode[K, V]]
	edgeCount int

	visitEpoch      uint64
	indegreeScratch collectionmapping.Map[K, int]
	queueScratch    []K
	stackScratch    []K
	orderScratch    []K
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
	node, ok := g.nodes.Get(id)
	if ok {
		node.value = value
		return false
	}

	g.nodes.Set(id, &graphNode[K, V]{value: value})
	return true
}

// GetNode returns node value by id.
func (g *Graph[K, V]) GetNode(id K) (V, bool) {
	var zero V
	if g == nil {
		return zero, false
	}
	node, ok := g.nodes.Get(id)
	if !ok {
		return zero, false
	}
	return node.value, true
}

// SetNodeValue updates node value by id.
func (g *Graph[K, V]) SetNodeValue(id K, value V) bool {
	if g == nil {
		return false
	}
	node, ok := g.nodes.Get(id)
	if !ok {
		return false
	}
	node.value = value
	return true
}

// HasNode reports whether id exists.
func (g *Graph[K, V]) HasNode(id K) bool {
	_, ok := g.nodes.Get(id)
	return ok
}

// DeleteNode removes one node and all incident edges.
func (g *Graph[K, V]) DeleteNode(id K) bool {
	if g == nil {
		return false
	}
	node, ok := g.nodes.Get(id)
	if !ok {
		return false
	}

	if g.directed {
		g.edgeCount -= node.order.Len()
		for i := 0; i < g.nodes.Len(); i++ {
			otherID, otherNode, ok := g.nodes.At(i)
			if !ok || otherID == id {
				continue
			}
			if otherNode.deleteNeighbor(id) {
				g.edgeCount--
			}
		}
	} else {
		for i := 0; i < node.order.Len(); i++ {
			neighborID, ok := node.order.Get(i)
			if !ok {
				continue
			}
			if neighborID != id {
				if otherNode, ok := g.nodes.Get(neighborID); ok {
					otherNode.deleteNeighbor(id)
				}
			}
			g.edgeCount--
		}
	}

	_ = g.nodes.Delete(id)
	return true
}

// AddEdge inserts one edge.
func (g *Graph[K, V]) AddEdge(from, to K) error {
	if g == nil {
		return ErrNodeNotFound
	}
	fromNode, ok := g.nodes.Get(from)
	if !ok {
		return ErrNodeNotFound
	}
	toNode, ok := g.nodes.Get(to)
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
	if g == nil {
		return false
	}
	node, ok := g.nodes.Get(from)
	if !ok {
		return false
	}
	_, ok = node.neighbors.Get(to)
	return ok
}

// DeleteEdge removes one edge and reports whether it existed.
func (g *Graph[K, V]) DeleteEdge(from, to K) bool {
	if g == nil {
		return false
	}
	fromNode, ok := g.nodes.Get(from)
	if !ok || !fromNode.deleteNeighbor(to) {
		return false
	}
	if !g.directed && from != to {
		if toNode, ok := g.nodes.Get(to); ok {
			toNode.deleteNeighbor(from)
		}
	}
	g.edgeCount--
	return true
}

// Neighbors returns adjacent node ids in insertion order.
func (g *Graph[K, V]) Neighbors(id K) []K {
	if g == nil {
		return nil
	}
	node, ok := g.nodes.Get(id)
	if !ok {
		return nil
	}
	values := node.order.Values()
	if len(values) == 0 {
		return nil
	}
	return values
}

// NodeIDs returns node ids in insertion order.
func (g *Graph[K, V]) NodeIDs() []K {
	if g == nil {
		return nil
	}
	keys := g.nodes.Keys()
	if len(keys) == 0 {
		return nil
	}
	return keys
}

// Clone returns a shallow copy preserving node and edge insertion order.
func (g *Graph[K, V]) Clone() *Graph[K, V] {
	if g == nil {
		return NewUndirectedGraph[K, V]()
	}
	out := &Graph[K, V]{directed: g.directed}
	for i := 0; i < g.nodes.Len(); i++ {
		id, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		out.AddNode(id, node.value)
	}
	for i := 0; i < g.nodes.Len(); i++ {
		from, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		for j := 0; j < node.order.Len(); j++ {
			to, ok := node.order.Get(j)
			if !ok {
				continue
			}
			if g.directed || from == to || !out.HasEdge(from, to) {
				_ = out.AddEdge(from, to)
			}
		}
	}
	return out
}

// Degree returns total incident edge count for id.
// For directed graphs this is InDegree(id) + OutDegree(id).
func (g *Graph[K, V]) Degree(id K) int {
	if g == nil {
		return 0
	}
	if !g.directed {
		node, ok := g.nodes.Get(id)
		if !ok {
			return 0
		}
		return node.order.Len()
	}
	return g.InDegree(id) + g.OutDegree(id)
}

// OutDegree returns the number of outgoing edges for id.
func (g *Graph[K, V]) OutDegree(id K) int {
	if g == nil {
		return 0
	}
	node, ok := g.nodes.Get(id)
	if !ok {
		return 0
	}
	return node.order.Len()
}

// InDegree returns the number of incoming edges for id.
func (g *Graph[K, V]) InDegree(id K) int {
	if g == nil {
		return 0
	}
	if _, ok := g.nodes.Get(id); !ok {
		return 0
	}
	inDegree := 0
	for i := 0; i < g.nodes.Len(); i++ {
		_, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		for j := 0; j < node.order.Len(); j++ {
			neighborID, ok := node.order.Get(j)
			if !ok {
				continue
			}
			if neighborID == id {
				inDegree++
			}
		}
	}
	return inDegree
}

// RangeNodes iterates nodes in insertion order until fn returns false.
func (g *Graph[K, V]) RangeNodes(fn func(id K, value V) bool) {
	if g == nil || fn == nil {
		return
	}
	for i := 0; i < g.nodes.Len(); i++ {
		id, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		if !fn(id, node.value) {
			return
		}
	}
}

// RangeEdges iterates edges in insertion order until fn returns false.
func (g *Graph[K, V]) RangeEdges(fn func(from, to K) bool) {
	if g == nil || fn == nil {
		return
	}
	if g.directed {
		for i := 0; i < g.nodes.Len(); i++ {
			from, node, ok := g.nodes.At(i)
			if !ok {
				continue
			}
			for j := 0; j < node.order.Len(); j++ {
				to, ok := node.order.Get(j)
				if !ok {
					continue
				}
				if !fn(from, to) {
					return
				}
			}
		}
		return
	}

	emitted := make(map[EdgeSnapshot[K]]struct{}, g.edgeCount*2)
	for i := 0; i < g.nodes.Len(); i++ {
		from, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		for j := 0; j < node.order.Len(); j++ {
			to, ok := node.order.Get(j)
			if !ok {
				continue
			}
			edge := EdgeSnapshot[K]{From: from, To: to}
			if _, seen := emitted[edge]; seen {
				continue
			}
			if !fn(from, to) {
				return
			}
			emitted[edge] = struct{}{}
			emitted[EdgeSnapshot[K]{From: to, To: from}] = struct{}{}
		}
	}
}

// PathExists reports whether to is reachable from from.
func (g *Graph[K, V]) PathExists(from, to K) bool {
	if g == nil {
		return false
	}
	if from == to {
		return g.HasNode(from)
	}
	found := false
	g.BFS(from, func(id K, _ V) bool {
		if id == to {
			found = true
			return false
		}
		return true
	})
	return found
}

// BFS iterates reachable nodes in breadth-first order until visit returns false.
func (g *Graph[K, V]) BFS(start K, visit func(id K, value V) bool) {
	if g == nil || visit == nil {
		return
	}
	startNode, ok := g.nodes.Get(start)
	if !ok {
		return
	}

	mark := g.nextVisitEpoch()
	startNode.visitMark = mark
	queue := g.queueScratch[:0]
	queue = append(queue, start)
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		node, ok := g.nodes.Get(id)
		if !ok {
			continue
		}
		if !visit(id, node.value) {
			g.queueScratch = queue[:0]
			return
		}
		for i := 0; i < node.order.Len(); i++ {
			neighborID, ok := node.order.Get(i)
			if !ok {
				continue
			}
			neighbor, ok := g.nodes.Get(neighborID)
			if !ok || neighbor.visitMark == mark {
				continue
			}
			neighbor.visitMark = mark
			queue = append(queue, neighborID)
		}
	}
	g.queueScratch = queue[:0]
}

// DFS iterates reachable nodes in depth-first pre-order until visit returns false.
func (g *Graph[K, V]) DFS(start K, visit func(id K, value V) bool) {
	if g == nil || visit == nil {
		return
	}
	startNode, ok := g.nodes.Get(start)
	if !ok {
		return
	}

	mark := g.nextVisitEpoch()
	startNode.visitMark = mark
	stack := g.stackScratch[:0]
	stack = append(stack, start)
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		node, ok := g.nodes.Get(id)
		if !ok {
			continue
		}
		if !visit(id, node.value) {
			g.stackScratch = stack[:0]
			return
		}
		for index := node.order.Len() - 1; index >= 0; index-- {
			neighborID, ok := node.order.Get(index)
			if !ok {
				continue
			}
			neighbor, ok := g.nodes.Get(neighborID)
			if !ok || neighbor.visitMark == mark {
				continue
			}
			neighbor.visitMark = mark
			stack = append(stack, neighborID)
		}
	}
	g.stackScratch = stack[:0]
}

// TopologicalSort returns node ids in topological order.
func (g *Graph[K, V]) TopologicalSort() ([]K, error) {
	if g == nil || g.nodes.Len() == 0 {
		return nil, nil
	}
	if !g.directed {
		return nil, ErrTopologicalSortRequiresDirected
	}

	if g.indegreeScratch.Len() == 0 {
		g.indegreeScratch = *collectionmapping.NewMapWithCapacity[K, int](g.nodes.Len())
	}

	for i := 0; i < g.nodes.Len(); i++ {
		id, _, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		g.indegreeScratch.Set(id, 0)
	}
	for i := 0; i < g.nodes.Len(); i++ {
		_, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		for j := 0; j < node.order.Len(); j++ {
			neighborID, ok := node.order.Get(j)
			if !ok {
				continue
			}
			count, _ := g.indegreeScratch.Get(neighborID)
			g.indegreeScratch.Set(neighborID, count+1)
		}
	}

	queue := g.queueScratch[:0]
	for i := 0; i < g.nodes.Len(); i++ {
		id, _, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		if count, ok := g.indegreeScratch.Get(id); ok && count == 0 {
			queue = append(queue, id)
		}
	}

	order := g.orderScratch[:0]
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		order = append(order, id)
		node, ok := g.nodes.Get(id)
		if !ok {
			continue
		}
		for j := 0; j < node.order.Len(); j++ {
			neighborID, ok := node.order.Get(j)
			if !ok {
				continue
			}
			count, ok := g.indegreeScratch.Get(neighborID)
			if !ok {
				continue
			}
			next := count - 1
			g.indegreeScratch.Set(neighborID, next)
			if next == 0 {
				queue = append(queue, neighborID)
			}
		}
	}

	if len(order) != g.nodes.Len() {
		g.queueScratch = queue[:0]
		g.orderScratch = order[:0]
		return nil, ErrCycleDetected
	}

	out := slices.Clone(order)
	g.queueScratch = queue[:0]
	g.orderScratch = order[:0]
	return out, nil
}

// Len returns total node count.
func (g *Graph[K, V]) Len() int {
	if g == nil {
		return 0
	}
	return g.nodes.Len()
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
	g.nodes.Clear()
	g.edgeCount = 0
	g.visitEpoch = 0
	g.indegreeScratch.Clear()
	g.queueScratch = nil
	g.stackScratch = nil
	g.orderScratch = nil
}

func (g *Graph[K, V]) nextVisitEpoch() uint64 {
	g.visitEpoch++
	if g.visitEpoch != 0 {
		return g.visitEpoch
	}
	for i := 0; i < g.nodes.Len(); i++ {
		_, node, ok := g.nodes.At(i)
		if !ok {
			continue
		}
		node.visitMark = 0
	}
	g.visitEpoch = 1
	return g.visitEpoch
}

func (n *graphNode[K, V]) addNeighbor(id K) bool {
	if _, exists := n.neighbors.Get(id); exists {
		return false
	}
	n.neighbors.Set(id, struct{}{})
	n.order.Add(id)
	return true
}

func (n *graphNode[K, V]) deleteNeighbor(id K) bool {
	if !n.neighbors.Delete(id) {
		return false
	}
	index := -1
	for i := 0; i < n.order.Len(); i++ {
		item, ok := n.order.Get(i)
		if !ok {
			continue
		}
		if item == id {
			index = i
			break
		}
	}
	if index >= 0 {
		_, _ = n.order.RemoveAt(index)
	}
	return true
}
