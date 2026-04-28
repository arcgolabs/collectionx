package list_test

import (
	"fmt"
	"testing"

	list "github.com/arcgolabs/collectionx/list"
	"github.com/stretchr/testify/require"
)

func TestListOps(t *testing.T) {
	t.Parallel()

	source := list.NewList(1, 2, 3, 4)

	mapped := list.MapList(source, func(index int, item int) string {
		return fmt.Sprintf("%d:%d", index, item)
	})
	require.Equal(t, []string{"0:1", "1:2", "2:3", "3:4"}, mapped.Values())

	filtered := list.FilterList(source, func(_ int, item int) bool {
		return item%2 == 0
	})
	require.Equal(t, []int{2, 4}, filtered.Values())

	rejected := list.RejectList(source, func(_ int, item int) bool {
		return item%2 == 0
	})
	require.Equal(t, []int{1, 3}, rejected.Values())

	filterMapped := list.FilterMapList(source, func(_ int, item int) (string, bool) {
		if item%2 != 0 {
			return "", false
		}
		return fmt.Sprintf("v%d", item), true
	})
	require.Equal(t, []string{"v2", "v4"}, filterMapped.Values())

	flatMapped := list.FlatMapList(source, func(_ int, item int) []string {
		return []string{fmt.Sprintf("x%d", item), fmt.Sprintf("y%d", item)}
	})
	require.Equal(t, []string{"x1", "y1", "x2", "y2", "x3", "y3", "x4", "y4"}, flatMapped.Values())

	found, ok := list.FindList(source, func(_ int, item int) bool {
		return item > 2
	})
	require.True(t, ok)
	require.Equal(t, 3, found)

	reduced := list.ReduceList(source, 0, func(acc int, _ int, item int) int {
		return acc + item
	})
	require.Equal(t, 10, reduced)

	reducedErr, err := list.ReduceErrList(source, 0, func(acc int, _ int, item int) (int, error) {
		if item == 4 {
			return 0, fmt.Errorf("stop at %d", item)
		}
		return acc + item, nil
	})
	require.Error(t, err)
	require.Equal(t, 6, reducedErr)
	require.ErrorContains(t, err, "stop at 4")
}
