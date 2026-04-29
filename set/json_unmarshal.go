package set

import (
	"encoding/json"
	"fmt"

	collectionmapping "github.com/arcgolabs/collectionx/mapping"
)

// UnmarshalJSON implements json.Unmarshaler.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal set json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal set json: %w", err)
	}

	*s = *NewSetWithCapacity[T](len(items), items...)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ConcurrentSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal concurrent set json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent set json: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.core = NewSetWithCapacity[T](len(items), items...)
	s.jsonCache = nil
	s.stringCache = ""
	s.jsonDirty = false
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *MultiSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal multiset json: nil receiver")
	}

	var counts map[T]int
	if err := json.Unmarshal(data, &counts); err != nil {
		return fmt.Errorf("unmarshal multiset json: %w", err)
	}

	next := &MultiSet[T]{}
	if len(counts) > 0 {
		next.counts = *collectionmapping.NewMapWithCapacity[T, int](len(counts))
	}
	for item, count := range counts {
		if count < 0 {
			return fmt.Errorf("unmarshal multiset json: negative count for item")
		}
		if count == 0 {
			continue
		}
		next.counts.Set(item, count)
		next.size += count
	}
	*s = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal ordered set json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal ordered set json: %w", err)
	}

	*s = *NewOrderedSetWithCapacity[T](len(items), items...)
	return nil
}
