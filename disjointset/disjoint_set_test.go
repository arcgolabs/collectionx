package disjointset_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/disjointset"
	"github.com/stretchr/testify/require"
)

func TestDisjointSet_BasicOps(t *testing.T) {
	t.Parallel()

	ds := disjointset.New[int]()
	ds.Add(1, 2, 3)

	require.Equal(t, 3, ds.Len())
	require.Equal(t, 3, ds.SetCount())
	require.True(t, ds.Has(2))
	require.False(t, ds.Connected(1, 2))

	require.True(t, ds.Union(1, 2))
	require.True(t, ds.Connected(1, 2))
	require.Equal(t, 2, ds.SizeOf(1))
	require.Equal(t, 2, ds.SetCount())

	rootA, ok := ds.Find(1)
	require.True(t, ok)
	rootB, ok := ds.Find(2)
	require.True(t, ok)
	require.Equal(t, rootA, rootB)

	require.False(t, ds.Union(1, 2))

	groups := ds.Groups()
	require.Len(t, groups, 2)
}

func TestDisjointSet_UnionAutoAddsAndClear(t *testing.T) {
	t.Parallel()

	var ds disjointset.DisjointSet[string]
	require.True(t, ds.Union("a", "b"))
	require.Equal(t, 2, ds.Len())
	require.Equal(t, 1, ds.SetCount())
	require.True(t, ds.Connected("a", "b"))
	require.False(t, ds.Connected("a", "c"))
	require.Zero(t, ds.SizeOf("missing"))

	ds.Clear()
	require.True(t, ds.IsEmpty())
	require.Equal(t, map[string][]string{}, ds.Groups())
}

func TestDisjointSet_MembersAndRangeGroups(t *testing.T) {
	t.Parallel()

	ds := disjointset.New[int]()
	ds.Union(1, 2)
	ds.Union(2, 3)
	ds.Union(10, 11)

	members := ds.MembersOf(2)
	require.ElementsMatch(t, []int{1, 2, 3}, members)
	require.Nil(t, ds.MembersOf(99))

	groupCount := 0
	totalMembers := 0
	ds.RangeGroups(func(_ int, members []int) bool {
		groupCount++
		totalMembers += len(members)
		return true
	})
	require.Equal(t, ds.SetCount(), groupCount)
	require.Equal(t, ds.Len(), totalMembers)
}
