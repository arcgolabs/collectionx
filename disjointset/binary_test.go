package disjointset_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/disjointset"
	"github.com/stretchr/testify/require"
)

func TestDisjointSetBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	source := buildSerializationDisjointSet()
	data, err := source.MarshalBinary()
	require.NoError(t, err)

	var target disjointset.DisjointSet[int]
	require.NoError(t, target.UnmarshalBinary(data))
	assertDisjointSetShape(t, &target)

	data, err = source.GobEncode()
	require.NoError(t, err)

	var gobTarget disjointset.DisjointSet[int]
	require.NoError(t, gobTarget.GobDecode(data))
	assertDisjointSetShape(t, &gobTarget)
}
