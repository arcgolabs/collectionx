package bitset

import (
	"fmt"
)

func (b *BitSet) marshalJSONBytes() ([]byte, error) {
	data, err := marshalJSONValue(b.Values())
	if err != nil {
		return nil, fmt.Errorf("marshal bitset JSON: %w", err)
	}
	return data, nil
}

// MarshalJSON implements json.Marshaler.
func (b *BitSet) MarshalJSON() ([]byte, error) {
	data, err := b.marshalJSONBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal bitset: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (b *BitSet) String() string {
	data, err := b.marshalJSONBytes()
	return jsonResultString(data, err, "[]")
}
