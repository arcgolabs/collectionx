package tree

import (
	"fmt"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *Tree[K, V]) MarshalBinary() ([]byte, error) {
	data, err := marshalBinaryValue(t.Nodes())
	if err != nil {
		return nil, fmt.Errorf("marshal tree binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (t *Tree[K, V]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *Tree[K, V]) UnmarshalBinary(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal tree binary: nil receiver")
	}
	var roots []NodeSnapshot[K, V]
	if err := unmarshalBinaryValue(data, &roots); err != nil {
		return fmt.Errorf("unmarshal tree binary: %w", err)
	}
	next := NewTree[K, V]()
	for _, root := range roots {
		if err := appendJSONNode(next, root, nil); err != nil {
			return fmt.Errorf("unmarshal tree binary: %w", err)
		}
	}
	*t = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (t *Tree[K, V]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *ConcurrentTree[K, V]) MarshalBinary() ([]byte, error) {
	return t.Snapshot().MarshalBinary()
}

// GobEncode implements gob.GobEncoder.
func (t *ConcurrentTree[K, V]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *ConcurrentTree[K, V]) UnmarshalBinary(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal concurrent tree binary: nil receiver")
	}
	var snapshot Tree[K, V]
	if err := snapshot.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("unmarshal concurrent tree binary: %w", err)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tree = snapshot.Clone()
	return nil
}

// GobDecode implements gob.GobDecoder.
func (t *ConcurrentTree[K, V]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}
