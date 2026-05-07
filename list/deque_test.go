package list_test

import (
	"testing"

	list "github.com/arcgolabs/collectionx/list"
	"github.com/stretchr/testify/require"
)

func TestDeque_PushPop(t *testing.T) {
	t.Parallel()

	var d list.Deque[int]
	d.PushBack(2, 3)
	d.PushFront(1)
	require.Equal(t, []int{1, 2, 3}, d.Values())
	first, ok := d.GetFirst()
	require.True(t, ok)
	require.Equal(t, 1, first)
	last, ok := d.GetLast()
	require.True(t, ok)
	require.Equal(t, 3, last)
	require.True(t, d.GetFirstOption().IsPresent())
	require.True(t, d.GetLastOption().IsPresent())

	v, ok := d.PopFront()
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = d.PopBack()
	require.True(t, ok)
	require.Equal(t, 3, v)

	require.Equal(t, []int{2}, d.Values())
}

func TestDeque_GrowAndGet(t *testing.T) {
	t.Parallel()

	var d list.Deque[int]
	for i := range 100 {
		d.PushBack(i)
	}
	require.Equal(t, 100, d.Len())
	value, ok := d.Get(99)
	require.True(t, ok)
	require.Equal(t, 99, value)
}
