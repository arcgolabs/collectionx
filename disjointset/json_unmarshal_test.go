package disjointset_test

import (
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/disjointset"
	"github.com/stretchr/testify/require"
)

func TestDisjointSetJSONRoundTrip(t *testing.T) {
	t.Parallel()

	source := buildSerializationDisjointSet()
	data, err := json.Marshal(source)
	require.NoError(t, err)

	var target disjointset.DisjointSet[int]
	require.NoError(t, json.Unmarshal(data, &target))
	assertDisjointSetShape(t, &target)
}

func buildSerializationDisjointSet() *disjointset.DisjointSet[int] {
	ds := disjointset.New[int]()
	ds.Union(1, 2)
	ds.Union(2, 3)
	ds.Union(10, 11)
	ds.Add(99)
	return ds
}

func assertDisjointSetShape(t *testing.T, ds *disjointset.DisjointSet[int]) {
	t.Helper()
	require.Equal(t, 6, ds.Len())
	require.Equal(t, 3, ds.SetCount())
	require.True(t, ds.Connected(1, 3))
	require.True(t, ds.Connected(10, 11))
	require.False(t, ds.Connected(1, 10))
	require.Equal(t, 3, ds.SizeOf(1))
	require.Equal(t, 1, ds.SizeOf(99))
}
