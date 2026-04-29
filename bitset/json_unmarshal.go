package bitset

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (b *BitSet) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("unmarshal bitset JSON: nil receiver")
	}

	var bits []int
	if err := json.Unmarshal(data, &bits); err != nil {
		return fmt.Errorf("unmarshal bitset JSON: %w", err)
	}

	*b = *New(bits...)
	return nil
}
