package graph_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/graph"
)

const (
	benchGraphNodes     = 1 << 10
	benchGraphBranching = 4
)

func buildBenchDirectedGraph(tb testing.TB) *graph.Graph[int, int] {
	tb.Helper()
	g := graph.NewDirectedGraph[int, int]()
	for i := range benchGraphNodes {
		g.AddNode(i, i)
	}
	for i := 1; i < benchGraphNodes; i++ {
		parent := (i - 1) / benchGraphBranching
		if err := g.AddEdge(parent, i); err != nil {
			tb.Fatalf("AddEdge(%d, %d) error = %v", parent, i, err)
		}
	}
	return g
}

func BenchmarkGraphBFS(b *testing.B) {
	g := buildBenchDirectedGraph(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		g.BFS(0, func(_ int, _ int) bool { return true })
	}
}

func BenchmarkGraphDFS(b *testing.B) {
	g := buildBenchDirectedGraph(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		g.DFS(0, func(_ int, _ int) bool { return true })
	}
}

func BenchmarkGraphTopologicalSort(b *testing.B) {
	g := buildBenchDirectedGraph(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		order, err := g.TopologicalSort()
		if err != nil {
			b.Fatalf("TopologicalSort() error = %v", err)
		}
		if len(order) != benchGraphNodes {
			b.Fatalf("unexpected order length: %d", len(order))
		}
	}
}

func BenchmarkGraphDeleteNodeReadd(b *testing.B) {
	g := buildBenchDirectedGraph(b)
	nodeSpan := benchGraphNodes - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		nodeID := (i % nodeSpan) + 1
		parent := (nodeID - 1) / benchGraphBranching
		if !g.DeleteNode(nodeID) {
			b.Fatalf("DeleteNode(%d) failed", nodeID)
		}
		if !g.AddNode(nodeID, nodeID) {
			b.Fatalf("AddNode(%d) failed", nodeID)
		}
		if err := g.AddEdge(parent, nodeID); err != nil {
			b.Fatalf("AddEdge(%d, %d) error = %v", parent, nodeID, err)
		}
	}
}
