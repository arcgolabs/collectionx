package list_test

import (
	"sync"
	"testing"

	list "github.com/arcgolabs/collectionx/list"
	"github.com/stretchr/testify/require"
)

func TestConcurrentList_ParallelAdd(t *testing.T) {
	t.Parallel()

	var l list.ConcurrentList[int]

	const workers = 24
	const each = 150

	var wg sync.WaitGroup
	wg.Add(workers)

	for worker := range workers {
		go func() {
			defer wg.Done()
			base := worker * each
			for i := range each {
				l.Add(base + i)
			}
		}()
	}

	wg.Wait()
	require.Equal(t, workers*each, l.Len())
}

func TestConcurrentList_InsertRemoveAndSnapshot(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList(1, 3)
	require.True(t, l.AddAt(1, 2))
	require.Equal(t, []int{1, 2, 3}, l.Values())

	removed, ok := l.RemoveAt(1)
	require.True(t, ok)
	require.Equal(t, 2, removed)

	snapshot := l.Snapshot()
	l.Add(9)
	require.Equal(t, []int{1, 3}, snapshot.Values())
}

func TestConcurrentList_OptionAPIs(t *testing.T) {
	t.Parallel()

	var l list.ConcurrentList[string]
	l.Add("a", "b")

	value, ok := l.GetFirst()
	require.True(t, ok)
	require.Equal(t, "a", value)

	opt := l.GetFirstOption()
	require.True(t, opt.IsPresent())
	value, ok = opt.Get()
	require.True(t, ok)
	require.Equal(t, "a", value)

	last, ok := l.GetLast()
	require.True(t, ok)
	require.Equal(t, "b", last)

	lastOpt := l.GetLastOption()
	require.True(t, lastOpt.IsPresent())
	last, ok = lastOpt.Get()
	require.True(t, ok)
	require.Equal(t, "b", last)

	removed := l.RemoveAtOption(1)
	require.True(t, removed.IsPresent())
	removedValue, ok := removed.Get()
	require.True(t, ok)
	require.Equal(t, "b", removedValue)

	require.True(t, l.GetOption(99).IsAbsent())
	require.True(t, new(list.ConcurrentList[string]).GetFirstOption().IsAbsent())
	require.True(t, new(list.ConcurrentList[string]).GetLastOption().IsAbsent())
}

func TestConcurrentList_Merge(t *testing.T) {
	t.Parallel()

	left := list.NewConcurrentList(1, 2)
	right := list.NewList(3, 4)
	otherConcurrent := list.NewConcurrentList(5, 6)

	left.Merge(right).MergeConcurrent(otherConcurrent).MergeSlice([]int{7, 8})
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, left.Values())
}

func TestNewConcurrentListWithCapacity(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentListWithCapacity[int](8, 1, 2, 3)

	require.Equal(t, []int{1, 2, 3}, l.Values())
	l.Add(4, 5, 6, 7, 8)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, l.Values())
}

func TestConcurrentList_Join(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList("a", "b", "c")
	require.Equal(t, "a,b,c", l.Join(","))
}

func TestConcurrentList_JSONCacheReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList(1, 2, 3)

	data, err := l.ToJSON()
	require.NoError(t, err)
	require.Equal(t, `[1,2,3]`, string(data))
	require.Equal(t, `[1,2,3]`, l.String())

	data[0] = '{'
	fresh, err := l.ToJSON()
	require.NoError(t, err)
	require.Equal(t, `[1,2,3]`, string(fresh))

	require.True(t, l.Set(1, 9))
	require.Equal(t, `[1,9,3]`, l.String())
}

func TestConcurrentList_Reverse(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList(1, 2, 3, 4)
	reversed := l.Reverse()

	require.Same(t, l, reversed)
	require.Equal(t, []int{4, 3, 2, 1}, l.Values())
}

func TestConcurrentList_ChunkAndWindow(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList(1, 2, 3, 4, 5)

	chunks := l.Chunk(2)
	require.Len(t, chunks, 3)
	require.Equal(t, []int{1, 2}, chunks[0].Values())
	require.Equal(t, []int{3, 4}, chunks[1].Values())
	require.Equal(t, []int{5}, chunks[2].Values())

	windows := l.Window(3)
	require.Len(t, windows, 3)
	require.Equal(t, []int{1, 2, 3}, windows[0].Values())
	require.Equal(t, []int{2, 3, 4}, windows[1].Values())
	require.Equal(t, []int{3, 4, 5}, windows[2].Values())

	chunks[0].Set(0, 99)
	windows[0].Set(0, 88)
	require.Equal(t, []int{1, 2, 3, 4, 5}, l.Values())
}

func TestConcurrentList_BinarySearchFunc(t *testing.T) {
	t.Parallel()

	l := list.NewConcurrentList(1, 3, 5, 7, 9)
	compare := func(item, target int) int {
		switch {
		case item < target:
			return -1
		case item > target:
			return 1
		default:
			return 0
		}
	}

	index, ok := l.BinarySearchFunc(7, compare)
	require.True(t, ok)
	require.Equal(t, 3, index)

	index, ok = l.BinarySearchFunc(8, compare)
	require.False(t, ok)
	require.Equal(t, -1, index)
}
