package set

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
	collectionmapping "github.com/arcgolabs/collectionx/mapping"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *Set[T]) MarshalBinary() ([]byte, error) {
	return marshalSetBinary("set", s.Values())
}

// GobEncode implements gob.GobEncoder.
func (s *Set[T]) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Set[T]) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal set binary: nil receiver")
	}
	var items []T
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal set binary: %w", err)
	}
	*s = *NewSetWithCapacity[T](len(items), items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (s *Set[T]) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *ConcurrentSet[T]) MarshalBinary() ([]byte, error) {
	return marshalSetBinary("concurrent set", s.Values())
}

// GobEncode implements gob.GobEncoder.
func (s *ConcurrentSet[T]) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ConcurrentSet[T]) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal concurrent set binary: nil receiver")
	}
	var items []T
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent set binary: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.core = NewSetWithCapacity[T](len(items), items...)
	s.jsonCache = nil
	s.stringCache = ""
	s.jsonDirty = false
	return nil
}

// GobDecode implements gob.GobDecoder.
func (s *ConcurrentSet[T]) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *MultiSet[T]) MarshalBinary() ([]byte, error) {
	return marshalSetBinary("multiset", s.AllCounts())
}

// GobEncode implements gob.GobEncoder.
func (s *MultiSet[T]) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MultiSet[T]) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal multiset binary: nil receiver")
	}
	var counts map[T]int
	if err := common.UnmarshalBinaryValue(data, &counts); err != nil {
		return fmt.Errorf("unmarshal multiset binary: %w", err)
	}
	next := &MultiSet[T]{}
	if len(counts) > 0 {
		next.counts = *collectionmapping.NewMapWithCapacity[T, int](len(counts))
	}
	for item, count := range counts {
		if count < 0 {
			return fmt.Errorf("unmarshal multiset binary: negative count for item")
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

// GobDecode implements gob.GobDecoder.
func (s *MultiSet[T]) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *OrderedSet[T]) MarshalBinary() ([]byte, error) {
	return marshalSetBinary("ordered set", s.Values())
}

// GobEncode implements gob.GobEncoder.
func (s *OrderedSet[T]) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OrderedSet[T]) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal ordered set binary: nil receiver")
	}
	var items []T
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal ordered set binary: %w", err)
	}
	*s = *NewOrderedSetWithCapacity[T](len(items), items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (s *OrderedSet[T]) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

func marshalSetBinary(kind string, value any) ([]byte, error) {
	data, err := common.MarshalBinaryValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s binary: %w", kind, err)
	}
	return data, nil
}
