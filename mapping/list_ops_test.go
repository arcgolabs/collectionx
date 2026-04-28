package mapping_test

import (
	"fmt"
	"testing"

	list "github.com/arcgolabs/collectionx/list"
	mapping "github.com/arcgolabs/collectionx/mapping"
	"github.com/stretchr/testify/require"
)

func TestListToMappingHelpers(t *testing.T) {
	t.Parallel()

	source := list.NewList(1, 2, 3, 4)

	grouped := mapping.GroupByList(source, func(_ int, item int) string {
		if item%2 == 0 {
			return "even"
		}
		return "odd"
	})
	require.Equal(t, []int{1, 3}, grouped.Get("odd"))
	require.Equal(t, []int{2, 4}, grouped.Get("even"))

	associated := mapping.AssociateList(source, func(index int, item int) (string, string) {
		return fmt.Sprintf("k%d", index), fmt.Sprintf("v%d", item)
	})
	value, ok := associated.Get("k2")
	require.True(t, ok)
	require.Equal(t, "v3", value)
}
