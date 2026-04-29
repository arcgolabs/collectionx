package graph_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/graph"
	"github.com/stretchr/testify/require"
)

func TestGraphBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	source := buildSerializationGraph()
	data, err := source.MarshalBinary()
	require.NoError(t, err)

	var target graph.Graph[int, string]
	require.NoError(t, target.UnmarshalBinary(data))
	assertGraphShape(t, &target)

	data, err = source.GobEncode()
	require.NoError(t, err)

	var gobTarget graph.Graph[int, string]
	require.NoError(t, gobTarget.GobDecode(data))
	assertGraphShape(t, &gobTarget)
}
