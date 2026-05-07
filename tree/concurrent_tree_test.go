package tree_test

import (
	"sync"
	"testing"

	tree "github.com/arcgolabs/collectionx/tree"
	"github.com/stretchr/testify/require"
)

func TestConcurrentTree_BasicOperations(t *testing.T) {
	tr := tree.NewConcurrentTree[int, string]()

	require.NoError(t, tr.AddRoot(1, "root"))
	require.NoError(t, tr.AddRoot(9, "root-b"))
	require.NoError(t, tr.AddChild(1, 2, "child-a"))
	require.NoError(t, tr.AddChild(2, 3, "child-b"))
	require.NoError(t, tr.AddChild(1, 4, "child-c"))

	require.True(t, tr.Has(3))
	require.Equal(t, 5, tr.Len())
	require.True(t, tr.SetValue(3, "child-b-updated"))

	node, ok := tr.Get(3)
	require.True(t, ok)
	require.Equal(t, "child-b-updated", node.Value())

	parent, ok := tr.Parent(3)
	require.True(t, ok)
	require.Equal(t, 2, parent.ID())

	depth, ok := tr.Depth(3)
	require.True(t, ok)
	require.Equal(t, 2, depth)

	require.Equal(t, []int{2, 3, 4}, nodeIDs(tr.Descendants(1)))
	require.Equal(t, []int{4}, nodeIDs(tr.Siblings(2)))
	require.Equal(t, []int{9}, nodeIDs(tr.Siblings(1)))
	require.Equal(t, []int{3, 4, 9}, nodeIDs(tr.Leaves()))
}

func TestConcurrentTree_SnapshotIsolation(t *testing.T) {
	tr := tree.NewConcurrentTree[int, string]()
	require.NoError(t, tr.AddRoot(1, "root"))
	require.NoError(t, tr.AddChild(1, 2, "child-a"))

	snapshot := tr.Snapshot()
	require.Equal(t, 2, snapshot.Len())

	require.NoError(t, tr.AddChild(1, 3, "child-b"))
	require.Equal(t, 2, snapshot.Len())
	require.False(t, snapshot.Has(3))

	require.True(t, snapshot.SetValue(2, "snapshot-only"))
	node, _ := tr.Get(2)
	require.Equal(t, "child-a", node.Value())
}

func TestConcurrentTree_SnapshotAPIs(t *testing.T) {
	tr := tree.NewConcurrentTree[int, string]()
	require.NoError(t, tr.AddRoot(1, "root"))
	require.NoError(t, tr.AddChild(1, 2, "child"))

	require.Equal(t, []tree.NodeSnapshot[int, string]{
		{
			ID:    1,
			Value: "root",
			Children: []tree.NodeSnapshot[int, string]{
				{ID: 2, Value: "child"},
			},
		},
	}, tr.Nodes())
	require.Equal(t, []tree.Entry[int, string]{
		tree.RootEntry(1, "root"),
		tree.ChildEntry(2, 1, "child"),
	}, tr.Entries())
}

func TestConcurrentTree_ParallelAddChildren(t *testing.T) {
	tr := tree.NewConcurrentTree[int, int]()
	require.NoError(t, tr.AddRoot(0, 0))

	const workers = 12
	const each = 120

	var wg sync.WaitGroup
	wg.Add(workers)
	var firstErr error
	var firstErrOnce sync.Once
	for w := range workers {
		go func() {
			defer wg.Done()
			base := w * each
			for i := range each {
				id := base + i + 1
				if err := tr.AddChild(0, id, id); err != nil {
					firstErrOnce.Do(func() {
						firstErr = err
					})
				}
			}
		}()
	}
	wg.Wait()

	require.NoError(t, firstErr)
	require.Equal(t, 1+workers*each, tr.Len())
	require.True(t, tr.Has(workers*each))
	require.Equal(t, workers*each, len(tr.Descendants(0)))
}

func TestBuildConcurrent(t *testing.T) {
	entries := []tree.Entry[int, string]{
		tree.RootEntry(1, "root"),
		tree.ChildEntry(2, 1, "a"),
		tree.ChildEntry(3, 2, "b"),
	}

	tr, err := tree.BuildConcurrent(entries)
	require.NoError(t, err)
	require.Equal(t, 3, tr.Len())
	require.Equal(t, []int{1}, nodeIDs(tr.Roots()))
}

func TestBuildConcurrent_WithCycle(t *testing.T) {
	entries := []tree.Entry[int, string]{
		tree.ChildEntry(1, 2, "a"),
		tree.ChildEntry(2, 1, "b"),
	}

	_, err := tree.BuildConcurrent(entries)
	require.ErrorIs(t, err, tree.ErrCycleDetected)
}

func TestConcurrentTree_QueryMissingAndEmptyCases(t *testing.T) {
	tr := tree.NewConcurrentTree[int, string]()

	depth, ok := tr.Depth(42)
	require.False(t, ok)
	require.Zero(t, depth)
	require.Nil(t, tr.Siblings(42))
	require.Nil(t, tr.Leaves())

	require.NoError(t, tr.AddRoot(1, "root"))

	depth, ok = tr.Depth(1)
	require.True(t, ok)
	require.Zero(t, depth)
	require.Nil(t, tr.Siblings(1))
	require.Equal(t, []int{1}, nodeIDs(tr.Leaves()))
}

func TestConcurrentTree_RangeBFS(t *testing.T) {
	tr := tree.NewConcurrentTree[int, string]()

	require.NoError(t, tr.AddRoot(1, "root-a"))
	require.NoError(t, tr.AddRoot(2, "root-b"))
	require.NoError(t, tr.AddChild(1, 3, "a-1"))
	require.NoError(t, tr.AddChild(1, 4, "a-2"))
	require.NoError(t, tr.AddChild(2, 5, "b-1"))
	require.NoError(t, tr.AddChild(3, 6, "a-1-1"))

	var visited []int
	tr.RangeBFS(func(node *tree.Node[int, string]) bool {
		visited = append(visited, node.ID())
		return true
	})
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, visited)

	visited = visited[:0]
	tr.RangeBFS(func(node *tree.Node[int, string]) bool {
		visited = append(visited, node.ID())
		return node.ID() != 4
	})
	require.Equal(t, []int{1, 2, 3, 4}, visited)
}
