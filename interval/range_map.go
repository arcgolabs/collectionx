package interval

import (
	"cmp"
	"slices"
	"sort"

	"github.com/samber/mo"
)

// RangeEntry is one range-value pair used by RangeMap.
type RangeEntry[T cmp.Ordered, V any] struct {
	Range Range[T]
	Value V
}

// RangeMap maps half-open ranges [start, end) to values.
// Overlapping Put overrides existing values in the input interval.
// Internal entries are kept sorted and non-overlapping.
type RangeMap[T cmp.Ordered, V any] struct {
	entries []RangeEntry[T, V]

	entriesCache []RangeEntry[T, V]
	entriesDirty bool
	jsonCache    []byte
	stringCache  string
	jsonDirty    bool
}

// NewRangeMap creates an empty range map.
func NewRangeMap[T cmp.Ordered, V any]() *RangeMap[T, V] {
	return &RangeMap[T, V]{
		entries: make([]RangeEntry[T, V], 0),
	}
}

// Put assigns value to [start, end), overriding any overlaps.
func (m *RangeMap[T, V]) Put(start, end T, value V) bool {
	if m == nil {
		return false
	}
	input, ok := newRangeEntry(start, end, value)
	if !ok {
		return false
	}
	m.entries = putRangeEntry(m.entries, input)
	m.invalidateDerivedCaches()
	return true
}

// Get returns value for point query.
func (m *RangeMap[T, V]) Get(point T) (V, bool) {
	var zero V
	entry, ok := m.Containing(point)
	if !ok {
		return zero, false
	}
	return entry.Value, true
}

// Containing returns the stored entry containing point.
func (m *RangeMap[T, V]) Containing(point T) (RangeEntry[T, V], bool) {
	var zero RangeEntry[T, V]
	if m == nil || len(m.entries) == 0 {
		return zero, false
	}
	index := sort.Search(len(m.entries), func(i int) bool {
		return m.entries[i].Range.End > point
	})
	if index < len(m.entries) && m.entries[index].Range.Contains(point) {
		return m.entries[index], true
	}
	return zero, false
}

// GetOption returns value for point query as mo.Option.
func (m *RangeMap[T, V]) GetOption(point T) mo.Option[V] {
	value, ok := m.Get(point)
	if !ok {
		return mo.None[V]()
	}
	return mo.Some(value)
}

// DeleteRange removes mappings in [start, end).
func (m *RangeMap[T, V]) DeleteRange(start, end T) bool {
	if m == nil || len(m.entries) == 0 {
		return false
	}
	input := Range[T]{Start: start, End: end}
	if !input.IsValid() {
		return false
	}
	next, changed := deleteRangeEntries(m.entries, input)
	if changed {
		m.entries = next
		m.invalidateDerivedCaches()
	}
	return changed
}

// Overlapping returns stored entries whose ranges overlap the input range.
func (m *RangeMap[T, V]) Overlapping(start, end T) []RangeEntry[T, V] {
	if m == nil || len(m.entries) == 0 {
		return nil
	}
	input := Range[T]{Start: start, End: end}
	if !input.IsValid() {
		return nil
	}

	index := sort.Search(len(m.entries), func(i int) bool {
		return m.entries[i].Range.End > input.Start
	})
	if index == len(m.entries) {
		return nil
	}

	overlaps := make([]RangeEntry[T, V], 0, 4)
	for ; index < len(m.entries); index++ {
		current := m.entries[index]
		if current.Range.Start >= input.End {
			break
		}
		overlaps = append(overlaps, current)
	}
	if len(overlaps) == 0 {
		return nil
	}
	return overlaps
}

// Bounds returns the smallest range covering all stored entries.
func (m *RangeMap[T, V]) Bounds() (Range[T], bool) {
	var zero Range[T]
	if m == nil || len(m.entries) == 0 {
		return zero, false
	}
	return Range[T]{
		Start: m.entries[0].Range.Start,
		End:   m.entries[len(m.entries)-1].Range.End,
	}, true
}

// Entries returns copied entries sorted by range start.
func (m *RangeMap[T, V]) Entries() []RangeEntry[T, V] {
	if m == nil || len(m.entries) == 0 {
		return nil
	}
	if !m.entriesDirty && len(m.entriesCache) > 0 {
		return slices.Clone(m.entriesCache)
	}
	m.entriesCache = slices.Clone(m.entries)
	m.entriesDirty = false
	return slices.Clone(m.entriesCache)
}

// GetFirst returns the first entry by range start.
func (m *RangeMap[T, V]) GetFirst() (RangeEntry[T, V], bool) {
	var zero RangeEntry[T, V]
	if m == nil || len(m.entries) == 0 {
		return zero, false
	}
	return m.entries[0], true
}

// GetLast returns the last entry by range start.
func (m *RangeMap[T, V]) GetLast() (RangeEntry[T, V], bool) {
	var zero RangeEntry[T, V]
	if m == nil || len(m.entries) == 0 {
		return zero, false
	}
	return m.entries[len(m.entries)-1], true
}

// Len returns non-overlapping entry count.
func (m *RangeMap[T, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.entries)
}

// IsEmpty reports whether map has no entries.
func (m *RangeMap[T, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Clear removes all entries.
func (m *RangeMap[T, V]) Clear() {
	if m == nil {
		return
	}
	m.entries = nil
	m.entriesCache = nil
	m.entriesDirty = false
	m.jsonCache = nil
	m.stringCache = ""
	m.jsonDirty = false
}

// Range iterates entries in start order until fn returns false.
func (m *RangeMap[T, V]) Range(fn func(entry RangeEntry[T, V]) bool) {
	if m == nil || fn == nil {
		return
	}
	for _, entry := range m.entries {
		if !fn(entry) {
			return
		}
	}
}

func newRangeEntry[T cmp.Ordered, V any](start, end T, value V) (RangeEntry[T, V], bool) {
	input := Range[T]{Start: start, End: end}
	if !input.IsValid() {
		return RangeEntry[T, V]{}, false
	}
	return RangeEntry[T, V]{Range: input, Value: value}, true
}

func putRangeEntry[T cmp.Ordered, V any](entries []RangeEntry[T, V], input RangeEntry[T, V]) []RangeEntry[T, V] {
	if len(entries) == 0 {
		return append(entries, input)
	}

	first := findFirstRangeEntryEndingAfter(entries, input.Range.Start)
	if first == len(entries) {
		return append(entries, input)
	}
	if entries[first].Range.Start >= input.Range.End {
		return replaceSliceRange(entries, first, first, input)
	}

	var replacement [3]RangeEntry[T, V]
	replCount := 0
	if left, ok := splitLeftRangeEntry(entries[first], input.Range.Start); ok {
		replacement[replCount] = left
		replCount++
	}
	replacement[replCount] = input
	replCount++

	end, right, ok := splitRightRangeEntry(entries, first, input.Range.End)
	if ok {
		replacement[replCount] = right
		replCount++
	}
	return replaceSliceRange(entries, first, end, replacement[:replCount]...)
}

func findFirstRangeEntryEndingAfter[T cmp.Ordered, V any](entries []RangeEntry[T, V], point T) int {
	return sort.Search(len(entries), func(i int) bool {
		return entries[i].Range.End > point
	})
}

func splitLeftRangeEntry[T cmp.Ordered, V any](entry RangeEntry[T, V], splitAt T) (RangeEntry[T, V], bool) {
	if entry.Range.Start >= splitAt {
		return RangeEntry[T, V]{}, false
	}
	return RangeEntry[T, V]{
		Range: Range[T]{Start: entry.Range.Start, End: splitAt},
		Value: entry.Value,
	}, true
}

func splitRightRangeEntry[T cmp.Ordered, V any](entries []RangeEntry[T, V], first int, end T) (int, RangeEntry[T, V], bool) {
	for index := first; index < len(entries); index++ {
		entry := entries[index]
		if entry.Range.Start >= end {
			return index, RangeEntry[T, V]{}, false
		}
		if entry.Range.End > end {
			return index + 1, RangeEntry[T, V]{
				Range: Range[T]{Start: end, End: entry.Range.End},
				Value: entry.Value,
			}, true
		}
	}
	return len(entries), RangeEntry[T, V]{}, false
}

func deleteRangeEntries[T cmp.Ordered, V any](entries []RangeEntry[T, V], cut Range[T]) ([]RangeEntry[T, V], bool) {
	oldLen := len(entries)
	first := findFirstRangeEntryEndingAfter(entries, cut.Start)
	if first == len(entries) {
		return entries, false
	}

	changed := false
	write := first

	for index := first; index < len(entries); index++ {
		entry := entries[index]
		if entry.Range.Start >= cut.End {
			if write != index {
				copy(entries[write:], entries[index:])
			}
			return entries[:write+len(entries[index:])], changed
		}

		changed = true
		if entry.Range.Start < cut.Start {
			entries[write] = RangeEntry[T, V]{
				Range: Range[T]{Start: entry.Range.Start, End: cut.Start},
				Value: entry.Value,
			}
			write++
		}
		if cut.End < entry.Range.End {
			if write == len(entries) {
				entries = append(entries, RangeEntry[T, V]{})
			}
			entries[write] = RangeEntry[T, V]{
				Range: Range[T]{Start: cut.End, End: entry.Range.End},
				Value: entry.Value,
			}
			write++
			tailCount := oldLen - (index + 1)
			if tailCount > 0 {
				copy(entries[write:], entries[index+1:oldLen])
				write += tailCount
			}
			return entries[:write], changed
		}
	}

	return entries[:write], changed
}

func (m *RangeMap[T, V]) invalidateDerivedCaches() {
	if m == nil {
		return
	}
	m.entriesCache = nil
	m.entriesDirty = true
	m.jsonCache = nil
	m.stringCache = ""
	m.jsonDirty = true
}
