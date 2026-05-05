package disjointset

import (
	"fmt"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *DisjointSet[T]) MarshalBinary() ([]byte, error) {
	data, err := marshalBinaryValue(d.groupsSnapshot())
	if err != nil {
		return nil, fmt.Errorf("marshal disjoint set binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (d *DisjointSet[T]) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *DisjointSet[T]) UnmarshalBinary(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal disjoint set binary: nil receiver")
	}

	var groups [][]T
	if err := unmarshalBinaryValue(data, &groups); err != nil {
		return fmt.Errorf("unmarshal disjoint set binary: %w", err)
	}

	*d = *New[T]()
	for _, members := range groups {
		d.restoreGroup(members)
	}
	return nil
}

// GobDecode implements gob.GobDecoder.
func (d *DisjointSet[T]) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

func (d *DisjointSet[T]) restoreGroup(members []T) {
	if len(members) == 0 {
		return
	}
	d.Add(members...)
	root := members[0]
	for _, item := range members[1:] {
		d.Union(root, item)
	}
}
