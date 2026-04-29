package bitset

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
)

func (b *BitSet) marshalJSONBytes() ([]byte, error) {
	data, err := common.MarshalJSONValue(b.Values())
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
	return common.JSONResultString(data, err, "[]")
}
