package interval

import (
	"fmt"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *RangeSet[T]) MarshalBinary() ([]byte, error) {
	data, err := marshalBinaryValue(s.Ranges())
	if err != nil {
		return nil, fmt.Errorf("marshal range set binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (s *RangeSet[T]) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RangeSet[T]) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal range set binary: nil receiver")
	}
	var ranges []Range[T]
	if err := unmarshalBinaryValue(data, &ranges); err != nil {
		return fmt.Errorf("unmarshal range set binary: %w", err)
	}
	next := NewRangeSet[T]()
	for _, current := range ranges {
		if !current.IsValid() {
			return fmt.Errorf("unmarshal range set binary: invalid range")
		}
		next.AddRange(current)
	}
	*s = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (s *RangeSet[T]) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *RangeMap[T, V]) MarshalBinary() ([]byte, error) {
	data, err := marshalBinaryValue(m.Entries())
	if err != nil {
		return nil, fmt.Errorf("marshal range map binary: %w", err)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (m *RangeMap[T, V]) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *RangeMap[T, V]) UnmarshalBinary(data []byte) error {
	if m == nil {
		return fmt.Errorf("unmarshal range map binary: nil receiver")
	}
	var entries []RangeEntry[T, V]
	if err := unmarshalBinaryValue(data, &entries); err != nil {
		return fmt.Errorf("unmarshal range map binary: %w", err)
	}
	next := NewRangeMap[T, V]()
	for _, entry := range entries {
		if !entry.Range.IsValid() {
			return fmt.Errorf("unmarshal range map binary: invalid range entry")
		}
		next.Put(entry.Range.Start, entry.Range.End, entry.Value)
	}
	*m = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (m *RangeMap[T, V]) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}
