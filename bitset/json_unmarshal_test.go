package bitset_test

import (
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/bitset"
	"github.com/stretchr/testify/require"
)

func TestBitSetJSONRoundTrip(t *testing.T) {
	t.Parallel()

	source := bitset.New(1, 3, 64, 128)
	data, err := json.Marshal(source)
	require.NoError(t, err)

	var target bitset.BitSet
	require.NoError(t, json.Unmarshal(data, &target))
	require.Equal(t, source.Values(), target.Values())
}
