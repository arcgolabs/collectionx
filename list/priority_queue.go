package list

import (
	"errors"
)

// ErrNilPriorityQueueComparator reports that NewPriorityQueue received a nil comparator.
var ErrNilPriorityQueueComparator = errors.New("list: priority queue comparator cannot be nil")

// PriorityQueue is a generic heap-based queue.
// less(a, b) should return true when a has higher priority than b.
type PriorityQueue[T any] struct {
	h *priorityQueueHeap[T]
}

// NewPriorityQueue creates a priority queue with comparator and optional items.
func NewPriorityQueue[T any](less func(a, b T) bool, items ...T) (*PriorityQueue[T], error) {
	if less == nil {
		return nil, ErrNilPriorityQueueComparator
	}
	pq := &PriorityQueue[T]{
		h: &priorityQueueHeap[T]{
			items: make([]T, 0, len(items)),
			less:  less,
		},
	}
	pq.h.items = append(pq.h.items, items...)
	pq.h.heapify()
	return pq, nil
}

// Push inserts value.
func (pq *PriorityQueue[T]) Push(value T) {
	if pq == nil || pq.h == nil {
		return
	}
	pq.h.push(value)
}

// Pop removes and returns top-priority value.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	var zero T
	if pq == nil || pq.h == nil || pq.h.Len() == 0 {
		return zero, false
	}
	return pq.h.pop()
}

// Peek returns top-priority value without removing it.
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	var zero T
	if pq == nil || pq.h == nil || pq.h.Len() == 0 {
		return zero, false
	}
	return pq.h.items[0], true
}

// GetFirst returns the current top-priority value without removing it.
func (pq *PriorityQueue[T]) GetFirst() (T, bool) {
	return pq.Peek()
}

// Len returns item count.
func (pq *PriorityQueue[T]) Len() int {
	if pq == nil || pq.h == nil {
		return 0
	}
	return pq.h.Len()
}

// IsEmpty reports whether queue has no items.
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Len() == 0
}

// Clear removes all items.
func (pq *PriorityQueue[T]) Clear() {
	if pq == nil || pq.h == nil {
		return
	}
	pq.h.items = nil
}

// Values returns a snapshot of internal heap order (not sorted order).
func (pq *PriorityQueue[T]) Values() []T {
	if pq == nil || pq.h == nil || pq.h.Len() == 0 {
		return nil
	}
	out := make([]T, len(pq.h.items))
	copy(out, pq.h.items)
	return out
}

// ValuesSorted returns a sorted snapshot by popping from a temporary heap copy.
func (pq *PriorityQueue[T]) ValuesSorted() []T {
	if pq == nil || pq.h == nil || pq.h.Len() == 0 {
		return nil
	}
	temp := &priorityQueueHeap[T]{
		items: make([]T, len(pq.h.items)),
		less:  pq.h.less,
	}
	copy(temp.items, pq.h.items)

	out := make([]T, 0, temp.Len())
	for temp.Len() > 0 {
		value, ok := temp.pop()
		if !ok {
			break
		}
		out = append(out, value)
	}
	return out
}

type priorityQueueHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func (h *priorityQueueHeap[T]) Len() int {
	return len(h.items)
}

func (h *priorityQueueHeap[T]) lessAt(i, j int) bool {
	return h.less(h.items[i], h.items[j])
}

func (h *priorityQueueHeap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *priorityQueueHeap[T]) heapify() {
	for i := h.Len()/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
}

func (h *priorityQueueHeap[T]) push(value T) {
	h.items = append(h.items, value)
	h.siftUp(h.Len() - 1)
}

func (h *priorityQueueHeap[T]) pop() (T, bool) {
	var zero T
	if h.Len() == 0 {
		return zero, false
	}

	last := h.Len() - 1
	h.swap(0, last)
	value := h.items[last]
	h.items[last] = zero
	h.items = h.items[:last]
	if h.Len() > 0 {
		h.siftDown(0)
	}
	return value, true
}

func (h *priorityQueueHeap[T]) siftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if !h.lessAt(index, parent) {
			return
		}
		h.swap(index, parent)
		index = parent
	}
}

func (h *priorityQueueHeap[T]) siftDown(index int) {
	for {
		left := index*2 + 1
		if left >= h.Len() {
			return
		}

		best := left
		right := left + 1
		if right < h.Len() && h.lessAt(right, left) {
			best = right
		}
		if !h.lessAt(best, index) {
			return
		}
		h.swap(index, best)
		index = best
	}
}
