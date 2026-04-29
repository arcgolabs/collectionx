package bitset_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/bitset"
	"github.com/stretchr/testify/require"
)

func TestBitSetBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	source := bitset.New(1, 3, 64, 128)
	data, err := source.MarshalBinary()
	require.NoError(t, err)

	var target bitset.BitSet
	require.NoError(t, target.UnmarshalBinary(data))
	require.Equal(t, source.Values(), target.Values())

	data, err = source.GobEncode()
	require.NoError(t, err)

	var gobTarget bitset.BitSet
	require.NoError(t, gobTarget.GobDecode(data))
	require.Equal(t, source.Values(), gobTarget.Values())
}
