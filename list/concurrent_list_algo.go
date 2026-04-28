package list

// Reverse reverses items in place and returns the receiver for chaining.
func (l *ConcurrentList[T]) Reverse() *ConcurrentList[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.core == nil {
		return l
	}
	l.core.Reverse()
	l.invalidateSerializationCacheLocked()
	return l
}

// Chunk splits a stable snapshot into consecutive fixed-size chunks.
func (l *ConcurrentList[T]) Chunk(size int) []*List[T] {
	return l.Snapshot().Chunk(size)
}

// Window returns all sliding windows from a stable snapshot.
func (l *ConcurrentList[T]) Window(size int) []*List[T] {
	return l.Snapshot().Window(size)
}

// BinarySearchFunc returns the matching index in a sorted list.
func (l *ConcurrentList[T]) BinarySearchFunc(target T, compare func(item, target T) int) (int, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.core == nil {
		return -1, false
	}
	return l.core.BinarySearchFunc(target, compare)
}
