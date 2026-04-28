package list

type iterable[T any] interface {
	Len() int
	Range(fn func(index int, item T) bool)
}

// MapList transforms list items into a new List.
func MapList[T any, R any](items iterable[T], mapper func(index int, item T) R) *List[R] {
	if items == nil || mapper == nil || items.Len() == 0 {
		return NewList[R]()
	}
	mapped := NewListWithCapacity[R](items.Len())
	items.Range(func(index int, item T) bool {
		mapped.Add(mapper(index, item))
		return true
	})
	return mapped
}

// FilterList keeps items that match predicate and returns them as a new List.
func FilterList[T any](items iterable[T], predicate func(index int, item T) bool) *List[T] {
	if items == nil || predicate == nil || items.Len() == 0 {
		return NewList[T]()
	}
	filtered := NewListWithCapacity[T](items.Len())
	items.Range(func(index int, item T) bool {
		if predicate(index, item) {
			filtered.Add(item)
		}
		return true
	})
	return filtered
}

// RejectList drops items that match predicate and returns the rest as a new List.
func RejectList[T any](items iterable[T], predicate func(index int, item T) bool) *List[T] {
	if items == nil || predicate == nil || items.Len() == 0 {
		return NewList[T]()
	}
	rejected := NewListWithCapacity[T](items.Len())
	items.Range(func(index int, item T) bool {
		if !predicate(index, item) {
			rejected.Add(item)
		}
		return true
	})
	return rejected
}

// FilterMapList transforms matching items into a new List.
func FilterMapList[T any, R any](items iterable[T], mapper func(index int, item T) (R, bool)) *List[R] {
	if items == nil || mapper == nil || items.Len() == 0 {
		return NewList[R]()
	}
	mapped := NewListWithCapacity[R](items.Len())
	items.Range(func(index int, item T) bool {
		value, ok := mapper(index, item)
		if ok {
			mapped.Add(value)
		}
		return true
	})
	return mapped
}

// FlatMapList expands each item into zero or more output items.
func FlatMapList[T any, R any](items iterable[T], mapper func(index int, item T) []R) *List[R] {
	if items == nil || mapper == nil || items.Len() == 0 {
		return NewList[R]()
	}
	mapped := NewListWithCapacity[R](items.Len())
	items.Range(func(index int, item T) bool {
		mapped.Add(mapper(index, item)...)
		return true
	})
	return mapped
}

// FindList returns the first item matching predicate.
func FindList[T any](items iterable[T], predicate func(index int, item T) bool) (T, bool) {
	var zero T
	if items == nil || predicate == nil || items.Len() == 0 {
		return zero, false
	}
	found := zero
	ok := false
	items.Range(func(index int, item T) bool {
		if !predicate(index, item) {
			return true
		}
		found = item
		ok = true
		return false
	})
	return found, ok
}

// ReduceList folds items into a single accumulator.
func ReduceList[T any, R any](items iterable[T], initial R, reducer func(acc R, index int, item T) R) R {
	if items == nil || reducer == nil || items.Len() == 0 {
		return initial
	}
	acc := initial
	items.Range(func(index int, item T) bool {
		acc = reducer(acc, index, item)
		return true
	})
	return acc
}

// ReduceErrList folds items into a single accumulator and stops on the first error.
func ReduceErrList[T any, R any](items iterable[T], initial R, reducer func(acc R, index int, item T) (R, error)) (R, error) {
	if items == nil || reducer == nil || items.Len() == 0 {
		return initial, nil
	}
	acc := initial
	var resultErr error
	items.Range(func(index int, item T) bool {
		next, err := reducer(acc, index, item)
		if err != nil {
			resultErr = err
			return false
		}
		acc = next
		return true
	})
	return acc, resultErr
}
