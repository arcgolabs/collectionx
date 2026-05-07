package mapping_test

import (
	"encoding/json"
	"testing"

	mapping "github.com/arcgolabs/collectionx/mapping"
	"github.com/stretchr/testify/require"
)

func TestOrderedMap_OrderStable(t *testing.T) {
	t.Parallel()

	var m mapping.OrderedMap[string, int]
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 9) // update should not move
	m.Set("c", 3)

	require.Equal(t, []string{"a", "b", "c"}, m.Keys())
	require.Equal(t, []int{9, 2, 3}, m.Values())
}

func TestOrderedMap_DeleteAndAt(t *testing.T) {
	t.Parallel()

	m := mapping.NewOrderedMap[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")

	require.True(t, m.Delete(2))
	require.Equal(t, []int{1, 3}, m.Keys())

	key, value, ok := m.At(1)
	require.True(t, ok)
	require.Equal(t, 3, key)
	require.Equal(t, "c", value)

	firstKey, firstValue, ok := m.First()
	require.True(t, ok)
	require.Equal(t, 1, firstKey)
	require.Equal(t, "a", firstValue)
	firstKey, firstValue, ok = m.GetFirst()
	require.True(t, ok)
	require.Equal(t, 1, firstKey)
	require.Equal(t, "a", firstValue)

	lastKey, lastValue, ok := m.Last()
	require.True(t, ok)
	require.Equal(t, 3, lastKey)
	require.Equal(t, "c", lastValue)
}

func TestNewOrderedMapWithCapacity(t *testing.T) {
	t.Parallel()

	m := mapping.NewOrderedMapWithCapacity[int, string](8)
	m.Set(1, "a")
	m.Set(2, "b")

	require.Equal(t, []int{1, 2}, m.Keys())
	require.Equal(t, []string{"a", "b"}, m.Values())
}

func TestOrderedMap_ChainMethods(t *testing.T) {
	t.Parallel()

	values := mapping.NewOrderedMap[string, int]()
	values.Set("a", 1)
	values.Set("b", 2)
	values.Set("c", 3)
	values.Set("d", 4)

	filtered := values.
		WhereEntries(func(_ string, value int) bool { return value >= 2 }).
		RejectEntries(func(key string, _ int) bool { return key == "c" }).
		Take(2)
	require.Equal(t, []string{"b", "d"}, filtered.Keys())
	require.Equal(t, []int{2, 4}, filtered.Values())

	dropped := values.Drop(2)
	require.Equal(t, []string{"c", "d"}, dropped.Keys())

	visited := make([]string, 0, 4)
	key, value, ok := values.
		EachEntry(func(key string, value int) { visited = append(visited, key) }).
		FirstEntryWhere(func(_ string, value int) bool { return value > 2 })
	require.True(t, ok)
	require.Equal(t, "c", key)
	require.Equal(t, 3, value)
	require.Equal(t, []string{"a", "b", "c", "d"}, visited)

	require.True(t, values.AllEntryMatch(func(_ string, value int) bool { return value > 0 }))
	require.True(t, values.AnyEntryMatch(func(key string, _ int) bool { return key == "b" }))
}

func TestOrderedMap_ValuesCacheReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	m := mapping.NewOrderedMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	values := m.Values()
	require.Equal(t, []int{1, 2}, values)

	values[0] = 99
	require.Equal(t, []int{1, 2}, m.Values())

	m.Set("a", 3)
	require.Equal(t, []int{3, 2}, m.Values())

	require.True(t, m.Delete("b"))
	require.Equal(t, []int{3}, m.Values())
}

func TestOrderedMap_JSONCacheReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	m := mapping.NewOrderedMap[string, int]()
	m.Set("a", 1)

	data, err := json.Marshal(m)
	require.NoError(t, err)
	require.Equal(t, `{"a":1}`, string(data))
	require.Equal(t, `{"a":1}`, m.String())

	data[0] = '['
	fresh, err := json.Marshal(m)
	require.NoError(t, err)
	require.Equal(t, `{"a":1}`, string(fresh))

	m.Set("b", 2)
	require.Contains(t, m.String(), `"a":1`)
	require.Contains(t, m.String(), `"b":2`)
}

func TestOrderedMap_FirstAndLast_Empty(t *testing.T) {
	t.Parallel()

	var m mapping.OrderedMap[string, int]

	key, value, ok := m.First()
	require.False(t, ok)
	require.Zero(t, key)
	require.Zero(t, value)

	key, value, ok = m.Last()
	require.False(t, ok)
	require.Zero(t, key)
	require.Zero(t, value)
}
