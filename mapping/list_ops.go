package mapping

type iterable[T any] interface {
	Len() int
	Range(fn func(index int, item T) bool)
}

// GroupByList groups items by key and returns them as a new MultiMap.
func GroupByList[T any, K comparable](items iterable[T], keySelector func(index int, item T) K) *MultiMap[K, T] {
	if items == nil || keySelector == nil || items.Len() == 0 {
		return NewMultiMap[K, T]()
	}
	grouped := NewMultiMapWithCapacity[K, T](items.Len())
	items.Range(func(index int, item T) bool {
		grouped.Put(keySelector(index, item), item)
		return true
	})
	return grouped
}

// AssociateList maps each item to a key/value pair and returns them as a new Map.
func AssociateList[T any, K comparable, V any](items iterable[T], mapper func(index int, item T) (K, V)) *Map[K, V] {
	if items == nil || mapper == nil || items.Len() == 0 {
		return NewMap[K, V]()
	}
	associated := NewMapWithCapacity[K, V](items.Len())
	items.Range(func(index int, item T) bool {
		key, value := mapper(index, item)
		associated.Set(key, value)
		return true
	})
	return associated
}
