package prefix

import (
	"fmt"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *Trie[V]) MarshalBinary() ([]byte, error) {
	data, err := marshalBinaryValue(t.All())
	if err != nil {
		return nil, fmt.Errorf("marshal trie binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (t *Trie[V]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *Trie[V]) UnmarshalBinary(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal trie binary: nil receiver")
	}
	var items map[string]V
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal trie binary: %w", err)
	}
	next := NewTrie[V]()
	for key, value := range items {
		next.Put(key, value)
	}
	*t = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (t *Trie[V]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}
