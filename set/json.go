package set

import (
	"fmt"
	"slices"
)

func (s *Set[T]) marshalJSONBytes() ([]byte, error) {
	if s != nil && !s.jsonDirty && s.jsonCache != nil {
		return slices.Clone(s.jsonCache), nil
	}
	data, err := marshalSetJSON("set", s.Values())
	if err != nil {
		return nil, err
	}
	if s != nil {
		s.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	data, err := s.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal set JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *Set[T]) String() string {
	if s != nil && !s.jsonDirty && s.stringCache != "" {
		return s.stringCache
	}
	data, err := s.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func (s *ConcurrentSet[T]) marshalJSONBytes() ([]byte, error) {
	if s == nil {
		return marshalSetJSON("concurrent set", []T(nil))
	}

	s.mu.RLock()
	if !s.jsonDirty && s.jsonCache != nil {
		data := slices.Clone(s.jsonCache)
		s.mu.RUnlock()
		return data, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.jsonDirty && s.jsonCache != nil {
		return slices.Clone(s.jsonCache), nil
	}

	var values []T
	if s.core != nil {
		values = s.core.Values()
	}
	data, err := marshalSetJSON("concurrent set", values)
	if err != nil {
		return nil, err
	}
	s.jsonCache = data
	s.stringCache = string(data)
	s.jsonDirty = false
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (s *ConcurrentSet[T]) MarshalJSON() ([]byte, error) {
	data, err := s.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal concurrent set JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *ConcurrentSet[T]) String() string {
	if s == nil {
		return "[]"
	}
	s.mu.RLock()
	if !s.jsonDirty && s.stringCache != "" {
		value := s.stringCache
		s.mu.RUnlock()
		return value
	}
	s.mu.RUnlock()
	data, err := s.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func (s *MultiSet[T]) marshalJSONBytes() ([]byte, error) {
	return marshalSetJSON("multiset", s.AllCounts())
}

// MarshalJSON implements json.Marshaler.
func (s *MultiSet[T]) MarshalJSON() ([]byte, error) {
	data, err := s.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal multiset JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *MultiSet[T]) String() string {
	data, err := s.marshalJSONBytes()
	return jsonResultString(data, err, "{}")
}

func (s *OrderedSet[T]) marshalJSONBytes() ([]byte, error) {
	if s != nil && !s.jsonDirty && s.jsonCache != nil {
		return slices.Clone(s.jsonCache), nil
	}
	data, err := marshalSetJSON("ordered set", s.Values())
	if err != nil {
		return nil, err
	}
	if s != nil {
		s.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (s *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	data, err := s.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal ordered set JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *OrderedSet[T]) String() string {
	if s != nil && !s.jsonDirty && s.stringCache != "" {
		return s.stringCache
	}
	data, err := s.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func marshalSetJSON[T any](kind string, value T) ([]byte, error) {
	data, err := marshalJSONValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s JSON: %w", kind, err)
	}

	return data, nil
}
