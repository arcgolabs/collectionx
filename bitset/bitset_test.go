package bitset_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/bitset"
	"github.com/stretchr/testify/require"
)

func TestBitSet_BasicOps(t *testing.T) {
	t.Parallel()

	b := bitset.New(1, 3, 64)
	require.True(t, b.Contains(1))
	require.True(t, b.Contains(64))
	require.False(t, b.Contains(2))
	require.Equal(t, 3, b.Len())
	require.Equal(t, []int{1, 3, 64}, b.Values())

	require.True(t, b.Set(2))
	require.False(t, b.Set(2))
	require.True(t, b.Remove(3))
	require.False(t, b.Remove(3))
	require.False(t, b.Set(-1))
	require.False(t, b.Contains(-1))
	require.Equal(t, []int{1, 2, 64}, b.Values())
}

func TestBitSet_SetOperations(t *testing.T) {
	t.Parallel()

	left := bitset.New(1, 2, 64)
	right := bitset.New(2, 3, 65)

	require.Equal(t, []int{1, 2, 3, 64, 65}, left.Union(right).Values())
	require.Equal(t, []int{2}, left.Intersect(right).Values())
	require.Equal(t, []int{1, 64}, left.Difference(right).Values())
	require.Equal(t, []int{1, 3, 64, 65}, left.SymmetricDifference(right).Values())
}

func TestBitSet_RangeAndClear(t *testing.T) {
	t.Parallel()

	var b bitset.BitSet
	b.Add(-1, 0, 5, 9)

	var visited []int
	b.Range(func(bit int) bool {
		visited = append(visited, bit)
		return bit != 5
	})
	require.Equal(t, []int{0, 5}, visited)

	b.Clear()
	require.True(t, b.IsEmpty())
	require.Nil(t, b.Values())
}
