package interval

func replaceSliceRange[T any](items []T, start, end int, replacement ...T) []T {
	removeCount := end - start
	addCount := len(replacement)
	newLen := len(items) - removeCount + addCount

	if newLen > cap(items) {
		next := make([]T, newLen)
		copy(next, items[:start])
		copy(next[start:], replacement)
		copy(next[start+addCount:], items[end:])
		return next
	}

	oldLen := len(items)
	items = items[:newLen]
	tailStart := start + addCount
	tailSrc := end
	if tailStart != tailSrc {
		copy(items[tailStart:], items[tailSrc:oldLen])
	}
	copy(items[start:], replacement)
	return items
}
