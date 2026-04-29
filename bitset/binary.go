package bitset

import (
	"fmt"

	common "github.com/arcgolabs/collectionx/internal"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (b *BitSet) MarshalBinary() ([]byte, error) {
	data, err := common.MarshalBinaryValue(b.Values())
	if err != nil {
		return nil, fmt.Errorf("marshal bitset binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (b *BitSet) GobEncode() ([]byte, error) {
	return b.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (b *BitSet) UnmarshalBinary(data []byte) error {
	if b == nil {
		return fmt.Errorf("unmarshal bitset binary: nil receiver")
	}

	var bits []int
	if err := common.UnmarshalBinaryValue(data, &bits); err != nil {
		return fmt.Errorf("unmarshal bitset binary: %w", err)
	}

	*b = *New(bits...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (b *BitSet) GobDecode(data []byte) error {
	return b.UnmarshalBinary(data)
}
