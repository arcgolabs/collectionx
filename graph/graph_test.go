package graph_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/graph"
	"github.com/stretchr/testify/require"
)

func TestDirectedGraph_BasicOpsAndTraversal(t *testing.T) {
	t.Parallel()

	g := graph.NewDirectedGraph[int, string]()
	require.True(t, g.AddNode(1, "a"))
	require.True(t, g.AddNode(2, "b"))
	require.True(t, g.AddNode(3, "c"))
	require.True(t, g.AddNode(4, "d"))
	require.False(t, g.AddNode(1, "aa"))

	value, ok := g.GetNode(1)
	require.True(t, ok)
	require.Equal(t, "aa", value)
	require.Equal(t, []int{1, 2, 3, 4}, g.NodeIDs())

	require.NoError(t, g.AddEdge(1, 2))
	require.NoError(t, g.AddEdge(1, 3))
	require.NoError(t, g.AddEdge(2, 4))
	require.NoError(t, g.AddEdge(3, 4))
	require.Equal(t, 4, g.EdgeCount())
	require.Equal(t, []int{2, 3}, g.Neighbors(1))
	require.True(t, g.HasEdge(1, 2))

	var bfs []int
	g.BFS(1, func(id int, _ string) bool {
		bfs = append(bfs, id)
		return true
	})
	require.Equal(t, []int{1, 2, 3, 4}, bfs)

	var dfs []int
	g.DFS(1, func(id int, _ string) bool {
		dfs = append(dfs, id)
		return true
	})
	require.Equal(t, []int{1, 2, 4, 3}, dfs)

	order, err := g.TopologicalSort()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4}, order)

	require.True(t, g.DeleteEdge(1, 3))
	require.False(t, g.HasEdge(1, 3))
	require.Equal(t, 3, g.EdgeCount())

	require.True(t, g.DeleteNode(2))
	require.False(t, g.HasNode(2))
	require.False(t, g.HasEdge(1, 2))
	require.Equal(t, 1, g.EdgeCount())
}

func TestGraph_TopologicalSortAndErrors(t *testing.T) {
	t.Parallel()

	directed := graph.NewDirectedGraph[int, int]()
	directed.AddNode(1, 1)
	directed.AddNode(2, 2)
	require.NoError(t, directed.AddEdge(1, 2))
	require.NoError(t, directed.AddEdge(2, 1))

	_, err := directed.TopologicalSort()
	require.ErrorIs(t, err, graph.ErrCycleDetected)

	undirected := graph.NewUndirectedGraph[int, int]()
	undirected.AddNode(1, 1)
	undirected.AddNode(2, 2)
	require.NoError(t, undirected.AddEdge(1, 2))
	require.True(t, undirected.HasEdge(2, 1))
	require.Equal(t, 1, undirected.EdgeCount())

	_, err = undirected.TopologicalSort()
	require.ErrorIs(t, err, graph.ErrTopologicalSortRequiresDirected)

	err = directed.AddEdge(1, 9)
	require.ErrorIs(t, err, graph.ErrNodeNotFound)
}

func TestGraph_RangeAndPathExists(t *testing.T) {
	t.Parallel()

	g := graph.NewDirectedGraph[int, string]()
	g.AddNode(1, "a")
	g.AddNode(2, "b")
	g.AddNode(3, "c")
	g.AddNode(4, "d")
	require.NoError(t, g.AddEdge(1, 2))
	require.NoError(t, g.AddEdge(2, 3))
	require.NoError(t, g.AddEdge(1, 4))

	var nodes []int
	g.RangeNodes(func(id int, _ string) bool {
		nodes = append(nodes, id)
		return true
	})
	require.Equal(t, []int{1, 2, 3, 4}, nodes)

	var edges [][2]int
	g.RangeEdges(func(from, to int) bool {
		edges = append(edges, [2]int{from, to})
		return true
	})
	require.Equal(t, [][2]int{{1, 2}, {1, 4}, {2, 3}}, edges)

	require.True(t, g.PathExists(1, 3))
	require.True(t, g.PathExists(1, 4))
	require.False(t, g.PathExists(4, 1))
	require.False(t, g.PathExists(1, 9))
}
