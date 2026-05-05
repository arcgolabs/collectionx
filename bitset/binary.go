package bitset

import (
	"fmt"
	"slices"
)

type bitSetBinarySnapshot struct {
	Words []uint64
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (b *BitSet) MarshalBinary() ([]byte, error) {
	var snapshot bitSetBinarySnapshot
	if b != nil {
		snapshot.Words = slices.Clone(b.words)
	}
	data, err := marshalBinaryValue(snapshot)
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

	var snapshot bitSetBinarySnapshot
	if err := unmarshalBinaryValue(data, &snapshot); err == nil {
		b.words = slices.Clone(snapshot.Words)
		b.recount()
		return nil
	}

	var bits []int
	if err := unmarshalBinaryValue(data, &bits); err != nil {
		return fmt.Errorf("unmarshal bitset binary: %w", err)
	}

	*b = *New(bits...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (b *BitSet) GobDecode(data []byte) error {
	return b.UnmarshalBinary(data)
}
