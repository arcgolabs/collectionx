package mapping_test

import (
	"encoding/json"
	"testing"

	mapping "github.com/arcgolabs/collectionx/mapping"
	"github.com/stretchr/testify/require"
)

func TestMultiMap_BasicOps(t *testing.T) {
	t.Parallel()

	var m mapping.MultiMap[string, int]

	m.Put("a", 1)
	m.PutAll("a", 2, 3)
	m.Put("b", 10)

	require.Equal(t, 2, m.Len())
	require.Equal(t, 4, m.ValueCount())
	require.Equal(t, []int{1, 2, 3}, m.Get("a"))
	key, values, ok := m.GetFirst()
	require.True(t, ok)
	require.Contains(t, []string{"a", "b"}, key)
	if key == "a" {
		require.Equal(t, []int{1, 2, 3}, values)
	} else {
		require.Equal(t, []int{10}, values)
	}

	m.Set("a", 9)
	require.Equal(t, []int{9}, m.Get("a"))

	removed := m.DeleteValueIf("a", func(value int) bool { return value == 9 })
	require.Equal(t, 1, removed)
	require.False(t, m.ContainsKey("a"))
}

func TestMultiMap_ViewStaysStableAfterWriteAndCopy(t *testing.T) {
	t.Parallel()

	m := mapping.NewMultiMap[string, int]()
	m.PutAll("k", 1, 2)

	view := m.Get("k")
	m.Put("k", 3)
	require.Equal(t, []int{1, 2}, view)
	require.Equal(t, []int{1, 2, 3}, m.Get("k"))

	copyValues := m.GetCopy("k")
	copyValues[0] = 99
	require.Equal(t, []int{1, 2, 3}, m.Get("k"))

	opt := m.GetOption("k")
	require.True(t, opt.IsPresent())
	got, ok := opt.Get()
	require.True(t, ok)
	require.Equal(t, []int{1, 2, 3}, got)

	require.True(t, m.GetOption("missing").IsAbsent())
}

func TestNewMultiMapWithCapacity(t *testing.T) {
	t.Parallel()

	m := mapping.NewMultiMapWithCapacity[string, int](8)
	m.PutAll("k", 1, 2)

	require.Equal(t, 1, m.Len())
	require.Equal(t, 2, m.ValueCount())
	require.Equal(t, []int{1, 2}, m.Get("k"))
}

func TestMultiMap_CloneAndFromAll(t *testing.T) {
	t.Parallel()

	source := mapping.NewMultiMapFromAll(map[string][]int{
		"a": {1, 2},
		"b": {3},
	})
	cloned := source.Clone()

	cloned.Put("a", 4)
	cloned.Delete("b")

	require.Equal(t, []int{1, 2}, source.Get("a"))
	require.Equal(t, []int{1, 2, 4}, cloned.Get("a"))
	require.True(t, source.ContainsKey("b"))
	require.False(t, cloned.ContainsKey("b"))
	require.Equal(t, 2, source.Len())
	require.Equal(t, 3, source.ValueCount())
	require.Equal(t, 3, cloned.ValueCount())
}

func TestMultiMap_FluentOps(t *testing.T) {
	t.Parallel()

	values := mapping.NewMultiMap[string, int]()
	values.PutAll("a", 1, 2)
	values.PutAll("b", 3, 4)
	values.Put("c", 5)

	filtered := values.
		WhereKeys(func(key string, _ []int) bool { return key != "c" }).
		RejectValues(func(_ string, value int) bool { return value%2 == 0 })

	require.Equal(t, []int{1}, filtered.Get("a"))
	require.Equal(t, []int{3}, filtered.Get("b"))
	require.False(t, filtered.ContainsKey("c"))

	flattened := filtered.FlattenValues()
	require.ElementsMatch(t, []int{1, 3}, flattened.Values())

	visited := mapping.NewMultiMap[string, int]()
	foundKey, foundValue, ok := values.
		EachKey(func(key string, entries []int) { visited.Set(key, entries...) }).
		EachValue(func(key string, value int) { visited.Put(key, value*10) }).
		FirstValueWhere(func(_ string, value int) bool { return value > 3 })

	require.True(t, ok)
	require.True(t, foundKey == "b" || foundKey == "c")
	require.True(t, foundValue > 3)
	require.Equal(t, []int{1, 2, 10, 20}, visited.Get("a"))
	require.True(t, values.AnyValueMatch(func(_ string, value int) bool { return value == 5 }))
	require.True(t, values.AllValuesMatch(func(_ string, value int) bool { return value > 0 }))
	require.False(t, values.AllValuesMatch(func(_ string, value int) bool { return value%2 == 0 }))
}

func TestMultiMap_JSONCacheReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	m := mapping.NewMultiMap[string, int]()
	m.Put("a", 1)

	data, err := json.Marshal(m)
	require.NoError(t, err)
	require.Equal(t, `{"a":[1]}`, string(data))
	require.Equal(t, `{"a":[1]}`, m.String())

	data[0] = '['
	fresh, err := json.Marshal(m)
	require.NoError(t, err)
	require.Equal(t, `{"a":[1]}`, string(fresh))

	m.Put("a", 2)
	require.Equal(t, `{"a":[1,2]}`, m.String())
}

func TestMultiMap_ViewAllAndRangeView(t *testing.T) {
	t.Parallel()

	m := mapping.NewMultiMap[string, int]()
	m.PutAll("a", 1, 2)

	seen := false
	m.ViewAll(func(items map[string][]int) {
		require.Equal(t, []int{1, 2}, items["a"])
		seen = true
	})
	require.True(t, seen)

	visited := 0
	m.RangeView(func(key string, values []int) bool {
		require.Equal(t, "a", key)
		require.Equal(t, []int{1, 2}, values)
		visited++
		return true
	})
	require.Equal(t, 1, visited)
}
