package mapping

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal map json: nil receiver")
	}

	var items map[K]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal map json: %w", err)
	}

	*m = *NewMapWithCapacity[K, V](len(items))
	m.SetAll(items)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *ConcurrentMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal concurrent map json: nil receiver")
	}

	var items map[K]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent map json: %w", err)
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

// UnmarshalJSON implements json.Unmarshaler.
// The receiver must already be constructed with a hash function.
func (m *ShardedConcurrentMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal sharded concurrent map json: nil receiver")
	}
	if m.hash == nil || len(m.shards) == 0 {
		return fmt.Errorf("unmarshal sharded concurrent map json: receiver must be initialized with NewShardedConcurrentMap")
	}

	var items map[K]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal sharded concurrent map json: %w", err)
	}

	m.Clear()
	m.SetAll(items)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *BiMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal bimap json: nil receiver")
	}

	var items map[K]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal bimap json: %w", err)
	}

	next := NewBiMap[K, V]()
	seenValues := NewMapWithCapacity[V, struct{}](len(items))
	for key, value := range items {
		if _, exists := seenValues.Get(value); exists {
			return fmt.Errorf("unmarshal bimap json: duplicate value")
		}
		seenValues.Set(value, struct{}{})
		next.Put(key, value)
	}

	*m = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal ordered map json: nil receiver")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("unmarshal ordered map json: %w", err)
	}

	switch value := token.(type) {
	case nil:
		*m = *NewOrderedMap[K, V]()
		return nil
	case json.Delim:
		if value != '{' {
			return fmt.Errorf("unmarshal ordered map json: expected object")
		}
	default:
		return fmt.Errorf("unmarshal ordered map json: expected object")
	}

	next := NewOrderedMap[K, V]()
	for decoder.More() {
		rawKeyToken, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("unmarshal ordered map json: %w", err)
		}

		rawKey, ok := rawKeyToken.(string)
		if !ok {
			return fmt.Errorf("unmarshal ordered map json: expected string key")
		}

		key, err := decodeObjectKey[K](rawKey)
		if err != nil {
			return fmt.Errorf("unmarshal ordered map json: %w", err)
		}

		var value V
		if err := decoder.Decode(&value); err != nil {
			return fmt.Errorf("unmarshal ordered map json: %w", err)
		}
		next.Set(key, value)
	}

	end, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("unmarshal ordered map json: %w", err)
	}
	if delim, ok := end.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("unmarshal ordered map json: expected object end")
	}

	*m = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *MultiMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal multimap json: nil receiver")
	}

	var items map[K][]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal multimap json: %w", err)
	}

	next := NewMultiMapWithCapacity[K, V](len(items))
	for key, values := range items {
		next.Set(key, values...)
	}
	*m = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *ConcurrentMultiMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal concurrent multimap json: nil receiver")
	}

	var items map[K][]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent multimap json: %w", err)
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

// UnmarshalJSON implements json.Unmarshaler.
func (t *Table[R, C, V]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal table json: nil receiver")
	}

	var rows map[R]map[C]V
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("unmarshal table json: %w", err)
	}

	next := NewTable[R, C, V]()
	for rowKey, rowValues := range rows {
		next.SetRow(rowKey, rowValues)
	}
	*t = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *ConcurrentTable[R, C, V]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal concurrent table json: nil receiver")
	}

	var rows map[R]map[C]V
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("unmarshal concurrent table json: %w", err)
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

func decodeObjectKey[K comparable](rawKey string) (K, error) {
	var object map[K]json.RawMessage
	data, err := json.Marshal(map[string]json.RawMessage{
		rawKey: json.RawMessage("null"),
	})
	if err != nil {
		var zero K
		return zero, fmt.Errorf("marshal object key: %w", err)
	}
	if err := json.Unmarshal(data, &object); err != nil {
		var zero K
		return zero, fmt.Errorf("decode object key %q: %w", rawKey, err)
	}
	for key := range object {
		return key, nil
	}
	var zero K
	return zero, fmt.Errorf("decode object key %q: no decoded key", rawKey)
}

func encodeObjectKey[K comparable](key K) (string, error) {
	object := map[K]json.RawMessage{
		key: json.RawMessage("null"),
	}
	data, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("encode object key: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("encode object key: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return "", fmt.Errorf("encode object key: expected object")
	}

	keyToken, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("encode object key: %w", err)
	}
	rawKey, ok := keyToken.(string)
	if !ok {
		return "", fmt.Errorf("encode object key: expected string key")
	}
	return rawKey, nil
}
