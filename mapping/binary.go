package mapping

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
)

type orderedMapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *Map[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("map", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *Map[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Map[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal map binary: nil receiver")
	}
	var items map[K]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal map binary: %w", err)
	}
	*m = *NewMapWithCapacity[K, V](len(items))
	m.SetAll(items)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *Map[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *ConcurrentMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("concurrent map", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *ConcurrentMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *ConcurrentMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal concurrent map binary: nil receiver")
	}
	var items map[K]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent map binary: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.core = NewMapWithCapacity[K, V](len(items))
	m.core.SetAll(items)
	m.jsonCache = nil
	m.stringCache = ""
	m.jsonDirty = false
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *ConcurrentMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *ShardedConcurrentMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("sharded concurrent map", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *ShardedConcurrentMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// The receiver must already be constructed with a hash function.
func (m *ShardedConcurrentMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal sharded concurrent map binary: nil receiver")
	}
	if m.hash == nil || len(m.shards) == 0 {
		return fmt.Errorf("unmarshal sharded concurrent map binary: receiver must be initialized with NewShardedConcurrentMap")
	}
	var items map[K]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal sharded concurrent map binary: %w", err)
	}
	m.Clear()
	m.SetAll(items)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *ShardedConcurrentMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *BiMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("bimap", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *BiMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *BiMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal bimap binary: nil receiver")
	}
	var items map[K]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal bimap binary: %w", err)
	}
	next := NewBiMap[K, V]()
	seenValues := NewMapWithCapacity[V, struct{}](len(items))
	for key, value := range items {
		if _, exists := seenValues.Get(value); exists {
			return fmt.Errorf("unmarshal bimap binary: duplicate value")
		}
		seenValues.Set(value, struct{}{})
		next.Put(key, value)
	}
	*m = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *BiMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *OrderedMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("ordered map", m.entriesSnapshot())
}

// GobEncode implements gob.GobEncoder.
func (m *OrderedMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *OrderedMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal ordered map binary: nil receiver")
	}
	var entries []orderedMapEntry[K, V]
	if err := common.UnmarshalBinaryValue(data, &entries); err != nil {
		return fmt.Errorf("unmarshal ordered map binary: %w", err)
	}
	next := NewOrderedMapWithCapacity[K, V](len(entries))
	for _, entry := range entries {
		next.Set(entry.Key, entry.Value)
	}
	*m = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *OrderedMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *MultiMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("multimap", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *MultiMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *MultiMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal multimap binary: nil receiver")
	}
	var items map[K][]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal multimap binary: %w", err)
	}
	next := NewMultiMapWithCapacity[K, V](len(items))
	for key, values := range items {
		next.Set(key, values...)
	}
	*m = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *MultiMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *ConcurrentMultiMap[K, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("concurrent multimap", m.All())
}

// GobEncode implements gob.GobEncoder.
func (m *ConcurrentMultiMap[K, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *ConcurrentMultiMap[K, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal concurrent multimap binary: nil receiver")
	}
	var items map[K][]V
	if err := common.UnmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent multimap binary: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.core = NewMultiMapWithCapacity[K, V](len(items))
	for key, values := range items {
		m.core.Set(key, values...)
	}
	m.jsonCache = nil
	m.stringCache = ""
	m.jsonDirty = false
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *ConcurrentMultiMap[K, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *Table[R, C, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("table", t.All())
}

// GobEncode implements gob.GobEncoder.
func (t *Table[R, C, V]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *Table[R, C, V]) UnmarshalBinary(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal table binary: nil receiver")
	}
	var rows map[R]map[C]V
	if err := common.UnmarshalBinaryValue(data, &rows); err != nil {
		return fmt.Errorf("unmarshal table binary: %w", err)
	}
	next := NewTable[R, C, V]()
	for rowKey, rowValues := range rows {
		next.SetRow(rowKey, rowValues)
	}
	*t = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (t *Table[R, C, V]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *ConcurrentTable[R, C, V]) MarshalBinary() ([]byte, error) {
	return marshalMappingBinary("concurrent table", t.All())
}

// GobEncode implements gob.GobEncoder.
func (t *ConcurrentTable[R, C, V]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *ConcurrentTable[R, C, V]) UnmarshalBinary(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal concurrent table binary: nil receiver")
	}
	var rows map[R]map[C]V
	if err := common.UnmarshalBinaryValue(data, &rows); err != nil {
		return fmt.Errorf("unmarshal concurrent table binary: %w", err)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.core = NewTable[R, C, V]()
	for rowKey, rowValues := range rows {
		t.core.SetRow(rowKey, rowValues)
	}
	t.jsonCache = nil
	t.stringCache = ""
	t.jsonDirty = false
	return nil
}

// GobDecode implements gob.GobDecoder.
func (t *ConcurrentTable[R, C, V]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}

func (m *OrderedMap[K, V]) entriesSnapshot() []orderedMapEntry[K, V] {
	if m == nil || m.order.Len() == 0 {
		return nil
	}
	entries := make([]orderedMapEntry[K, V], 0, m.order.Len())
	m.Range(func(key K, value V) bool {
		entries = append(entries, orderedMapEntry[K, V]{
			Key:   key,
			Value: value,
		})
		return true
	})
	return entries
}

func marshalMappingBinary(kind string, value any) ([]byte, error) {
	data, err := common.MarshalBinaryValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s binary: %w", kind, err)
	}
	return data, nil
}
