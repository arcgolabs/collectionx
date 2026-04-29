package disjointset

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (d *DisjointSet[T]) UnmarshalJSON(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal disjoint set JSON: nil receiver")
	}

	var groups [][]T
	if err := json.Unmarshal(data, &groups); err != nil {
		return fmt.Errorf("unmarshal disjoint set JSON: %w", err)
	}

	*d = *New[T]()
	for _, members := range groups {
		d.restoreGroup(members)
	}
	return nil
}
