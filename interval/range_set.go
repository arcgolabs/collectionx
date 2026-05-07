package interval

import (
	"cmp"
	"slices"
	"sort"
)

// RangeSet is a normalized set of half-open ranges [start, end).
// Internal ranges are kept sorted and non-overlapping.
type RangeSet[T cmp.Ordered] struct {
	ranges []Range[T]

	rangesCache []Range[T]
	rangesDirty bool
	jsonCache   []byte
	stringCache string
	jsonDirty   bool
}

// NewRangeSet creates an empty range set.
func NewRangeSet[T cmp.Ordered]() *RangeSet[T] {
	return &RangeSet[T]{
		ranges: make([]Range[T], 0),
	}
}

// Add inserts one range and merges overlaps/adjacent ranges.
func (s *RangeSet[T]) Add(start, end T) bool {
	return s.AddRange(Range[T]{Start: start, End: end})
}

// AddRange inserts one range and merges overlaps/adjacent ranges.
func (s *RangeSet[T]) AddRange(r Range[T]) bool {
	if s == nil || !r.IsValid() {
		return false
	}
	s.ranges = addRangeToSet(s.ranges, r)
	s.invalidateDerivedCaches()
	return true
}

// Remove removes interval part from the set.
func (s *RangeSet[T]) Remove(start, end T) bool {
	if s == nil || len(s.ranges) == 0 {
		return false
	}
	cut := Range[T]{Start: start, End: end}
	if !cut.IsValid() {
		return false
	}
	next, changed := removeRangeFromSet(s.ranges, cut)
	if changed {
		s.ranges = next
		s.invalidateDerivedCaches()
	}
	return changed
}

// Contains reports whether value is in any range.
func (s *RangeSet[T]) Contains(value T) bool {
	_, ok := s.Containing(value)
	return ok
}

// Containing returns the stored range containing value.
func (s *RangeSet[T]) Containing(value T) (Range[T], bool) {
	var zero Range[T]
	if s == nil || len(s.ranges) == 0 {
		return zero, false
	}
	index := sort.Search(len(s.ranges), func(i int) bool {
		return s.ranges[i].End > value
	})
	if index < len(s.ranges) && s.ranges[index].Contains(value) {
		return s.ranges[index], true
	}
	return zero, false
}

// Overlaps reports whether input range overlaps any stored range.
func (s *RangeSet[T]) Overlaps(start, end T) bool {
	return len(s.Overlapping(start, end)) > 0
}

// Overlapping returns stored ranges that overlap the input range.
func (s *RangeSet[T]) Overlapping(start, end T) []Range[T] {
	if s == nil || len(s.ranges) == 0 {
		return nil
	}
	input := Range[T]{Start: start, End: end}
	if !input.IsValid() {
		return nil
	}
	index := sort.Search(len(s.ranges), func(i int) bool {
		return s.ranges[i].End > input.Start
	})
	if index == len(s.ranges) {
		return nil
	}

	overlaps := make([]Range[T], 0, 4)
	for ; index < len(s.ranges); index++ {
		current := s.ranges[index]
		if current.Start >= input.End {
			break
		}
		overlaps = append(overlaps, current)
	}
	if len(overlaps) == 0 {
		return nil
	}
	return overlaps
}

// Bounds returns the smallest range covering all stored ranges.
func (s *RangeSet[T]) Bounds() (Range[T], bool) {
	var zero Range[T]
	if s == nil || len(s.ranges) == 0 {
		return zero, false
	}
	return Range[T]{
		Start: s.ranges[0].Start,
		End:   s.ranges[len(s.ranges)-1].End,
	}, true
}

// Ranges returns copied normalized ranges.
func (s *RangeSet[T]) Ranges() []Range[T] {
	if s == nil || len(s.ranges) == 0 {
		return nil
	}
	if !s.rangesDirty && len(s.rangesCache) > 0 {
		return slices.Clone(s.rangesCache)
	}
	s.rangesCache = slices.Clone(s.ranges)
	s.rangesDirty = false
	return slices.Clone(s.rangesCache)
}

// GetFirst returns the first normalized range by start.
func (s *RangeSet[T]) GetFirst() (Range[T], bool) {
	var zero Range[T]
	if s == nil || len(s.ranges) == 0 {
		return zero, false
	}
	return s.ranges[0], true
}

// GetLast returns the last normalized range by start.
func (s *RangeSet[T]) GetLast() (Range[T], bool) {
	var zero Range[T]
	if s == nil || len(s.ranges) == 0 {
		return zero, false
	}
	return s.ranges[len(s.ranges)-1], true
}

// Len returns number of normalized ranges.
func (s *RangeSet[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.ranges)
}

// IsEmpty reports whether set has no ranges.
func (s *RangeSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Clear removes all ranges.
func (s *RangeSet[T]) Clear() {
	if s == nil {
		return
	}
	s.ranges = nil
	s.rangesCache = nil
	s.rangesDirty = false
	s.jsonCache = nil
	s.stringCache = ""
	s.jsonDirty = false
}

// Range iterates normalized ranges until fn returns false.
func (s *RangeSet[T]) Range(fn func(r Range[T]) bool) {
	if s == nil || fn == nil {
		return
	}
	for _, r := range s.ranges {
		if !fn(r) {
			return
		}
	}
}

func addRangeToSet[T cmp.Ordered](ranges []Range[T], input Range[T]) []Range[T] {
	if len(ranges) == 0 {
		return append(ranges, input)
	}

	first := findFirstRangeEndingAtOrAfter(ranges, input.Start)
	if first == len(ranges) {
		return append(ranges, input)
	}
	if ranges[first].Start > input.End {
		return replaceSliceRange(ranges, first, first, input)
	}

	merged, end := mergeSetRanges(ranges, first, input)
	return replaceSliceRange(ranges, first, end, merged)
}

func findFirstRangeEndingAtOrAfter[T cmp.Ordered](ranges []Range[T], point T) int {
	return sort.Search(len(ranges), func(i int) bool {
		return ranges[i].End >= point
	})
}

func mergeSetRanges[T cmp.Ordered](ranges []Range[T], first int, merged Range[T]) (Range[T], int) {
	end := first
	for ; end < len(ranges); end++ {
		current := ranges[end]
		if current.Start > merged.End {
			break
		}
		if current.Start < merged.Start {
			merged.Start = current.Start
		}
		if current.End > merged.End {
			merged.End = current.End
		}
	}
	return merged, end
}

func removeRangeFromSet[T cmp.Ordered](ranges []Range[T], cut Range[T]) ([]Range[T], bool) {
	oldLen := len(ranges)
	first := findFirstRangeEndingAfter(ranges, cut.Start)
	if first == len(ranges) {
		return ranges, false
	}

	changed := false
	write := first

	for index := first; index < len(ranges); index++ {
		current := ranges[index]
		if current.Start >= cut.End {
			if write != index {
				copy(ranges[write:], ranges[index:])
			}
			return ranges[:write+len(ranges[index:])], changed
		}

		changed = true
		if current.Start < cut.Start {
			ranges[write] = Range[T]{Start: current.Start, End: cut.Start}
			write++
		}
		if cut.End < current.End {
			if write == len(ranges) {
				ranges = append(ranges, Range[T]{})
			}
			ranges[write] = Range[T]{Start: cut.End, End: current.End}
			write++
			tailCount := oldLen - (index + 1)
			if tailCount > 0 {
				copy(ranges[write:], ranges[index+1:oldLen])
				write += tailCount
			}
			return ranges[:write], changed
		}
	}

	return ranges[:write], changed
}

func findFirstRangeEndingAfter[T cmp.Ordered](ranges []Range[T], point T) int {
	return sort.Search(len(ranges), func(i int) bool {
		return ranges[i].End > point
	})
}

func (s *RangeSet[T]) invalidateDerivedCaches() {
	if s == nil {
		return
	}
	s.rangesCache = nil
	s.rangesDirty = true
	s.jsonCache = nil
	s.stringCache = ""
	s.jsonDirty = true
}
