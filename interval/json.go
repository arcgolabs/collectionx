package interval

import (
	"fmt"
	"slices"
)

func (s *RangeSet[T]) marshalJSONBytes() ([]byte, error) {
	if s != nil && !s.jsonDirty && s.jsonCache != nil {
		return slices.Clone(s.jsonCache), nil
	}
	data, err := marshalJSONValue(s.Ranges())
	if err != nil {
		return nil, fmt.Errorf("marshal range set json: %w", err)
	}
	if s != nil {
		s.jsonCache = data
		s.stringCache = string(data)
		s.jsonDirty = false
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (s *RangeSet[T]) MarshalJSON() ([]byte, error) {
	data, err := s.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal range set: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *RangeSet[T]) String() string {
	if s != nil && !s.jsonDirty && s.stringCache != "" {
		return s.stringCache
	}
	data, err := s.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func (m *RangeMap[T, V]) marshalJSONBytes() ([]byte, error) {
	if m != nil && !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}
	data, err := marshalJSONValue(m.Entries())
	if err != nil {
		return nil, fmt.Errorf("marshal range map json: %w", err)
	}
	if m != nil {
		m.jsonCache = data
		m.stringCache = string(data)
		m.jsonDirty = false
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *RangeMap[T, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal range map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *RangeMap[T, V]) String() string {
	if m != nil && !m.jsonDirty && m.stringCache != "" {
		return m.stringCache
	}
	data, err := m.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}
