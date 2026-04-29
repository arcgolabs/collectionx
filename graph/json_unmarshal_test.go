package graph_test

import (
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/graph"
	"github.com/stretchr/testify/require"
)

func TestGraphCloneAndDegrees(t *testing.T) {
	t.Parallel()

	directed := buildSerializationGraph()
	clone := directed.Clone()
	assertGraphShape(t, clone)
	require.Equal(t, 3, directed.Degree(2))
	require.Equal(t, 1, directed.InDegree(2))
	require.Equal(t, 2, directed.OutDegree(2))
	require.Equal(t, 0, directed.Degree(99))

	undirected := graph.NewUndirectedGraph[int, string]()
	require.True(t, undirected.AddNode(1, "a"))
	require.True(t, undirected.AddNode(2, "b"))
	require.True(t, undirected.AddNode(3, "c"))
	require.NoError(t, undirected.AddEdge(1, 2))
	require.NoError(t, undirected.AddEdge(1, 3))
	require.Equal(t, 2, undirected.Degree(1))
	require.Equal(t, 2, undirected.OutDegree(1))
	require.Equal(t, 2, undirected.InDegree(1))
}

func TestGraphJSONRoundTrip(t *testing.T) {
	t.Parallel()

	source := buildSerializationGraph()
	data, err := json.Marshal(source)
	require.NoError(t, err)

	var target graph.Graph[int, string]
	require.NoError(t, json.Unmarshal(data, &target))
	assertGraphShape(t, &target)
}

func buildSerializationGraph() *graph.Graph[int, string] {
	g := graph.NewDirectedGraph[int, string]()
	g.AddNode(1, "a")
	g.AddNode(2, "b")
	g.AddNode(3, "c")
	g.AddNode(4, "d")
	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)
	_ = g.AddEdge(2, 4)
	_ = g.AddEdge(2, 3)
	return g
}

func assertGraphShape(t *testing.T, g *graph.Graph[int, string]) {
	t.Helper()
	require.True(t, g.IsDirected())
	require.Equal(t, []int{1, 2, 3, 4}, g.NodeIDs())
	require.Equal(t, []int{2, 3}, g.Neighbors(1))
	require.Equal(t, []int{4, 3}, g.Neighbors(2))
	require.Equal(t, 4, g.EdgeCount())
	order, err := g.TopologicalSort()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 4, 3}, order)
}
