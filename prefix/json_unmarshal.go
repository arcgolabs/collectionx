package prefix

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (t *Trie[V]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("unmarshal trie json: nil receiver")
	}

	var items map[string]V
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal trie json: %w", err)
	}

	next := NewTrie[V]()
	for key, value := range items {
		next.Put(key, value)
	}
	*t = *next
	return nil
}
