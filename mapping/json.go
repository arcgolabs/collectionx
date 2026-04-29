//revive:disable:file-length-limit JSON mapping helpers are kept together to preserve the collection API surface.

package mapping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"

	common "github.com/arcgolabs/collectionx/internal"
)

func (m *Map[K, V]) marshalJSONBytes() ([]byte, error) {
	if m != nil && !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}

	var (
		data []byte
		err  error
	)
	if m == nil || len(m.items) == 0 {
		data, err = marshalMappingJSON(map[K]V{}, "map")
	} else {
		data, err = marshalMappingJSON(m.items, "map")
	}
	if err != nil {
		return nil, err
	}
	if m != nil {
		m.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *Map[K, V]) String() string {
	if m != nil && !m.jsonDirty && m.stringCache != "" {
		return m.stringCache
	}
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *ConcurrentMap[K, V]) marshalJSONBytes() ([]byte, error) {
	if m == nil {
		return marshalMappingJSON(map[K]V{}, "concurrent map")
	}

	m.mu.RLock()
	if !m.jsonDirty && m.jsonCache != nil {
		data := slices.Clone(m.jsonCache)
		m.mu.RUnlock()
		return data, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}

	var (
		data []byte
		err  error
	)
	if m.core == nil || len(m.core.items) == 0 {
		data, err = marshalMappingJSON(map[K]V{}, "concurrent map")
	} else {
		data, err = marshalMappingJSON(m.core.items, "concurrent map")
	}
	if err != nil {
		return nil, err
	}
	m.jsonCache = data
	m.stringCache = string(data)
	m.jsonDirty = false
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal concurrent map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *ConcurrentMap[K, V]) String() string {
	if m == nil {
		return "{}"
	}
	m.mu.RLock()
	if !m.jsonDirty && m.stringCache != "" {
		value := m.stringCache
		m.mu.RUnlock()
		return value
	}
	m.mu.RUnlock()
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *ShardedConcurrentMap[K, V]) marshalJSONBytes() ([]byte, error) {
	return marshalMappingJSON(m.All(), "sharded concurrent map")
}

// MarshalJSON implements json.Marshaler.
func (m *ShardedConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal sharded concurrent map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *ShardedConcurrentMap[K, V]) String() string {
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *BiMap[K, V]) marshalJSONBytes() ([]byte, error) {
	if m == nil || len(m.kv.items) == 0 {
		return marshalMappingJSON(map[K]V{}, "bimap")
	}
	return marshalMappingJSON(m.kv.items, "bimap")
}

// MarshalJSON implements json.Marshaler.
func (m *BiMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal bimap: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *BiMap[K, V]) String() string {
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *OrderedMap[K, V]) marshalJSONBytes() ([]byte, error) {
	if m != nil && !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}
	data, err := marshalOrderedMapJSON(m)
	if err != nil {
		return nil, err
	}
	if m != nil {
		m.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal ordered map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *OrderedMap[K, V]) String() string {
	if m != nil && !m.jsonDirty && m.stringCache != "" {
		return m.stringCache
	}
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *MultiMap[K, V]) marshalJSONBytes() ([]byte, error) {
	if m != nil && !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}
	data, err := marshalMappingJSON(m.All(), "multimap")
	if err != nil {
		return nil, err
	}
	if m != nil {
		m.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *MultiMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal multimap: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *MultiMap[K, V]) String() string {
	if m != nil && !m.jsonDirty && m.stringCache != "" {
		return m.stringCache
	}
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (m *ConcurrentMultiMap[K, V]) marshalJSONBytes() ([]byte, error) {
	if m == nil {
		return marshalMappingJSON(map[K][]V{}, "concurrent multimap")
	}

	m.mu.RLock()
	if !m.jsonDirty && m.jsonCache != nil {
		data := slices.Clone(m.jsonCache)
		m.mu.RUnlock()
		return data, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}

	var (
		data []byte
		err  error
	)
	if m.core == nil || m.core.items.Len() == 0 {
		data, err = marshalMappingJSON(map[K][]V{}, "concurrent multimap")
	} else {
		data, err = marshalMappingJSON(m.core.All(), "concurrent multimap")
	}
	if err != nil {
		return nil, err
	}
	m.jsonCache = data
	m.stringCache = string(data)
	m.jsonDirty = false
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *ConcurrentMultiMap[K, V]) MarshalJSON() ([]byte, error) {
	data, err := m.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal concurrent multimap: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *ConcurrentMultiMap[K, V]) String() string {
	if m == nil {
		return "{}"
	}
	m.mu.RLock()
	if !m.jsonDirty && m.stringCache != "" {
		value := m.stringCache
		m.mu.RUnlock()
		return value
	}
	m.mu.RUnlock()
	data, err := m.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (t *Table[R, C, V]) marshalJSONBytes() ([]byte, error) {
	if t != nil && !t.jsonDirty && t.jsonCache != nil {
		return slices.Clone(t.jsonCache), nil
	}
	data, err := marshalMappingJSON(t.All(), "table")
	if err != nil {
		return nil, err
	}
	if t != nil {
		t.cacheSerializationData(data)
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (t *Table[R, C, V]) MarshalJSON() ([]byte, error) {
	data, err := t.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal table: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (t *Table[R, C, V]) String() string {
	if t != nil && !t.jsonDirty && t.stringCache != "" {
		return t.stringCache
	}
	data, err := t.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func (t *ConcurrentTable[R, C, V]) marshalJSONBytes() ([]byte, error) {
	if t == nil {
		return marshalMappingJSON(map[R]map[C]V{}, "concurrent table")
	}

	t.mu.RLock()
	if !t.jsonDirty && t.jsonCache != nil {
		data := slices.Clone(t.jsonCache)
		t.mu.RUnlock()
		return data, nil
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.jsonDirty && t.jsonCache != nil {
		return slices.Clone(t.jsonCache), nil
	}

	var (
		data []byte
		err  error
	)
	if t.core == nil || t.core.data.Len() == 0 {
		data, err = marshalMappingJSON(map[R]map[C]V{}, "concurrent table")
	} else {
		data, err = marshalMappingJSON(t.core.All(), "concurrent table")
	}
	if err != nil {
		return nil, err
	}
	t.jsonCache = data
	t.stringCache = string(data)
	t.jsonDirty = false
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (t *ConcurrentTable[R, C, V]) MarshalJSON() ([]byte, error) {
	data, err := t.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal concurrent table: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (t *ConcurrentTable[R, C, V]) String() string {
	if t == nil {
		return "{}"
	}
	t.mu.RLock()
	if !t.jsonDirty && t.stringCache != "" {
		value := t.stringCache
		t.mu.RUnlock()
		return value
	}
	t.mu.RUnlock()
	data, err := t.marshalJSONBytes()
	return common.JSONResultString(data, err, "{}")
}

func marshalMappingJSON(value any, kind string) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s json: %w", kind, err)
	}
	return data, nil
}

func marshalOrderedMapJSON[K comparable, V any](m *OrderedMap[K, V]) ([]byte, error) {
	if m == nil || m.order.Len() == 0 {
		return marshalMappingJSON(map[string]V{}, "ordered map")
	}

	var buffer bytes.Buffer
	buffer.WriteByte('{')

	for index := range m.order.Len() {
		key, value, _ := m.At(index)
		fieldName, err := encodeObjectKey(key)
		if err != nil {
			return nil, fmt.Errorf("marshal ordered map json: %w", err)
		}
		valueData, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshal ordered map json: %w", err)
		}

		if index > 0 {
			buffer.WriteByte(',')
		}
		keyData, err := json.Marshal(fieldName)
		if err != nil {
			return nil, fmt.Errorf("marshal ordered map json: %w", err)
		}
		buffer.Write(keyData)
		buffer.WriteByte(':')
		buffer.Write(valueData)
	}

	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}
