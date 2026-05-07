package tree

import (
	"fmt"
)

func (t *Tree[K, V]) marshalJSONBytes() ([]byte, error) {
	return marshalTreeJSON("tree", t.Nodes())
}

// MarshalJSON implements json.Marshaler.
func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
	data, err := t.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal tree JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (t *Tree[K, V]) String() string {
	data, err := t.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func (t *ConcurrentTree[K, V]) marshalJSONBytes() ([]byte, error) {
	return t.Snapshot().marshalJSONBytes()
}

// MarshalJSON implements json.Marshaler.
func (t *ConcurrentTree[K, V]) MarshalJSON() ([]byte, error) {
	data, err := t.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal concurrent tree JSON: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (t *ConcurrentTree[K, V]) String() string {
	data, err := t.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}

func marshalTreeJSON[T any](kind string, value T) ([]byte, error) {
	data, err := marshalJSONValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s JSON: %w", kind, err)
	}

	return data, nil
}
