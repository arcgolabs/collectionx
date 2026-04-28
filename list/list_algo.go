package list

import (
	"slices"
	"sort"
)

// Reverse reverses items in place and returns the receiver for chaining.
func (l *List[T]) Reverse() *List[T] {
	if l == nil {
		return NewList[T]()
	}
	if len(l.items) < 2 {
		return l
	}
	slices.Reverse(l.items)
	l.invalidateSerializationCache()
	return l
}

// Chunk splits items into consecutive fixed-size chunks.
func (l *List[T]) Chunk(size int) []*List[T] {
	if l == nil || size <= 0 || len(l.items) == 0 {
		return nil
	}

	chunkCount := (len(l.items) + size - 1) / size
	chunks := make([]*List[T], 0, chunkCount)
	for start := 0; start < len(l.items); start += size {
		end := start + size
		if end > len(l.items) {
			end = len(l.items)
		}
		chunks = append(chunks, NewList(l.items[start:end]...))
	}
	return chunks
}

// Window returns all sliding windows of the given size.
func (l *List[T]) Window(size int) []*List[T] {
	if l == nil || size <= 0 || size > len(l.items) {
		return nil
	}

	windowCount := len(l.items) - size + 1
	windows := make([]*List[T], 0, windowCount)
	for start := 0; start+size <= len(l.items); start++ {
		windows = append(windows, NewList(l.items[start:start+size]...))
	}
	return windows
}

// BinarySearchFunc returns the matching index in a sorted list.
func (l *List[T]) BinarySearchFunc(target T, compare func(item, target T) int) (int, bool) {
	if l == nil || compare == nil || len(l.items) == 0 {
		return -1, false
	}

	index := sort.Search(len(l.items), func(i int) bool {
		return compare(l.items[i], target) >= 0
	})
	if index < len(l.items) && compare(l.items[index], target) == 0 {
		return index, true
	}
	return -1, false
}
