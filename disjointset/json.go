package disjointset

import (
	"fmt"
)

func (d *DisjointSet[T]) marshalJSONBytes() ([]byte, error) {
	data, err := marshalJSONValue(d.groupsSnapshot())
	if err != nil {
		return nil, fmt.Errorf("marshal disjoint set JSON: %w", err)
	}
	return data, nil
}

// MarshalJSON implements json.Marshaler.
func (d *DisjointSet[T]) MarshalJSON() ([]byte, error) {
	data, err := d.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal disjoint set: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (d *DisjointSet[T]) String() string {
	data, err := d.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func (d *DisjointSet[T]) groupsSnapshot() [][]T {
	if d == nil || len(d.parent) == 0 {
		return nil
	}
	groups := d.Groups()
	out := make([][]T, 0, len(groups))
	for _, members := range groups {
		copyMembers := make([]T, len(members))
		copy(copyMembers, members)
		out = append(out, copyMembers)
	}
	return out
}
