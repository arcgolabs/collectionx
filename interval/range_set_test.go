package interval_test

import (
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/interval"
	"github.com/stretchr/testify/require"
)

func TestRangeSet_AddMergeAndContains(t *testing.T) {
	t.Parallel()

	s := interval.NewRangeSet[int]()
	require.True(t, s.Add(1, 3))
	require.True(t, s.Add(3, 5)) // adjacent merge
	require.True(t, s.Add(10, 12))
	require.True(t, s.Add(4, 11)) // overlap merge all

	ranges := s.Ranges()
	require.Equal(t, 1, len(ranges))
	require.Equal(t, interval.Range[int]{Start: 1, End: 12}, ranges[0])
	require.True(t, s.Contains(8))
	require.False(t, s.Contains(12))
}

func TestRangeSet_RemoveSplit(t *testing.T) {
	t.Parallel()

	s := interval.NewRangeSet[int]()
	s.Add(0, 10)
	require.True(t, s.Remove(3, 7))
	require.Equal(
		t,
		[]interval.Range[int]{
			{Start: 0, End: 3},
			{Start: 7, End: 10},
		},
		s.Ranges(),
	)
}

func TestRangeSet_BoundariesAndOverlaps(t *testing.T) {
	t.Parallel()

	s := interval.NewRangeSet[int]()
	s.Add(0, 10)
	s.Add(20, 30)

	containing, ok := s.Containing(0)
	require.True(t, ok)
	require.Equal(t, interval.Range[int]{Start: 0, End: 10}, containing)

	require.True(t, s.Contains(0))
	require.True(t, s.Contains(9))
	require.False(t, s.Contains(10))
	require.False(t, s.Contains(19))
	require.True(t, s.Contains(20))

	require.False(t, s.Overlaps(10, 20))
	require.True(t, s.Overlaps(9, 11))
	require.True(t, s.Overlaps(29, 40))
	require.Equal(t, []interval.Range[int]{
		{Start: 0, End: 10},
		{Start: 20, End: 30},
	}, s.Overlapping(5, 25))

	bounds, ok := s.Bounds()
	require.True(t, ok)
	require.Equal(t, interval.Range[int]{Start: 0, End: 30}, bounds)
}

func TestRangeSet_ContainingOverlappingAndBounds_EmptyOrMiss(t *testing.T) {
	t.Parallel()

	s := interval.NewRangeSet[int]()

	containing, ok := s.Containing(5)
	require.False(t, ok)
	require.Equal(t, interval.Range[int]{}, containing)
	require.Nil(t, s.Overlapping(1, 3))

	bounds, ok := s.Bounds()
	require.False(t, ok)
	require.Equal(t, interval.Range[int]{}, bounds)

	require.True(t, s.Add(10, 20))
	containing, ok = s.Containing(20)
	require.False(t, ok)
	require.Equal(t, interval.Range[int]{}, containing)
	require.Nil(t, s.Overlapping(20, 30))
	require.Nil(t, s.Overlapping(30, 20))
}

func TestRangeSet_CachesReturnDefensiveCopies(t *testing.T) {
	t.Parallel()

	s := interval.NewRangeSet[int]()
	require.True(t, s.Add(1, 3))

	ranges := s.Ranges()
	require.Equal(t, []interval.Range[int]{{Start: 1, End: 3}}, ranges)
	ranges[0] = interval.Range[int]{Start: 9, End: 10}
	require.Equal(t, []interval.Range[int]{{Start: 1, End: 3}}, s.Ranges())

	data, err := json.Marshal(s)
	require.NoError(t, err)
	require.Equal(t, `[{"Start":1,"End":3}]`, string(data))
	require.Equal(t, `[{"Start":1,"End":3}]`, s.String())

	data[0] = '{'
	fresh, err := json.Marshal(s)
	require.NoError(t, err)
	require.Equal(t, `[{"Start":1,"End":3}]`, string(fresh))
}
