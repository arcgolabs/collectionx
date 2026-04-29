package interval

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (s *RangeSet[T]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal range set json: nil receiver")
	}

	var ranges []Range[T]
	if err := json.Unmarshal(data, &ranges); err != nil {
		return fmt.Errorf("unmarshal range set json: %w", err)
	}

	next := NewRangeSet[T]()
	for _, current := range ranges {
		if !current.IsValid() {
			return fmt.Errorf("unmarshal range set json: invalid range")
		}
		next.AddRange(current)
	}
	*s = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *RangeMap[T, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal range map json: nil receiver")
	}

	var entries []RangeEntry[T, V]
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("unmarshal range map json: %w", err)
	}

	next := NewRangeMap[T, V]()
	for _, entry := range entries {
		if !entry.Range.IsValid() {
			return fmt.Errorf("unmarshal range map json: invalid range entry")
		}
		next.Put(entry.Range.Start, entry.Range.End, entry.Value)
	}
	*m = *next
	return nil
}
