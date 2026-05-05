package prefix

import (
	"fmt"
)

// All returns all key-value pairs as a copied map.
func (t *Trie[V]) All() map[string]V {
	pairs := t.pairsWithPrefix("")
	if len(pairs) == 0 {
		return map[string]V{}
	}

	out := make(map[string]V, len(pairs))
	for _, item := range pairs {
		out[item.Key] = item.Value
	}
	return out
}

func (t *Trie[V]) marshalJSONBytes() ([]byte, error) {
	data, err := marshalJSONValue(t.All())
	if err != nil {
		return nil, fmt.Errorf("marshal trie json: %w", err)
	}
	return data, nil
}

// MarshalJSON implements json.Marshaler.
func (t *Trie[V]) MarshalJSON() ([]byte, error) {
	data, err := t.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal trie: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (t *Trie[V]) String() string {
	data, err := t.marshalJSONBytes()
	return jsonResultString(data, err, "{}")
}
