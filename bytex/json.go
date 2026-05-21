package bytex

import (
	"encoding/json"
	"fmt"
	"math/bits"
	"strconv"
)

var byteDecimalStrings = func() [256]string {
	var values [256]string
	for i := range values {
		values[i] = strconv.Itoa(i)
	}
	return values
}()

// MarshalJSON implements json.Marshaler.
func (l *List) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(listBytes(l))
	if err != nil {
		return nil, fmt.Errorf("marshal byte list JSON: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *List) UnmarshalJSON(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal byte list JSON: nil receiver")
	}
	var items []byte
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal byte list JSON: %w", err)
	}
	*l = *NewList(items...)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (r *RingBuffer) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(r.Bytes())
	if err != nil {
		return nil, fmt.Errorf("marshal byte ring buffer JSON: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RingBuffer) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal byte ring buffer JSON: nil receiver")
	}
	var items []byte
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal byte ring buffer JSON: %w", err)
	}
	next := NewRingBuffer(len(items))
	_, _ = next.Write(items)
	*r = *next
	return nil
}

// MarshalJSON implements json.Marshaler.
func (s *Set) MarshalJSON() ([]byte, error) {
	if s == nil || s.count == 0 {
		return []byte("null"), nil
	}
	data := make([]byte, 0, s.count*4+2)
	data = append(data, '[')
	first := true
	for wordIndex, word := range s.words {
		current := word
		for current != 0 {
			offset := bits.TrailingZeros64(current)
			if first {
				first = false
			} else {
				data = append(data, ',')
			}
			data = append(data, byteDecimalStrings[wordIndex*64+offset]...)
			current &= current - 1
		}
	}
	data = append(data, ']')
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Set) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal byte set JSON: nil receiver")
	}
	var values []int
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("unmarshal byte set JSON: %w", err)
	}
	s.Clear()
	for _, value := range values {
		if value < 0 || value > 255 {
			return fmt.Errorf("unmarshal byte set JSON: value %d out of byte range", value)
		}
		s.Set(byte(value))
	}
	return nil
}

// String implements fmt.Stringer.
func (s *Set) String() string {
	data, err := s.MarshalJSON()
	return jsonResultString(data, err, "[]")
}

// MarshalJSON implements json.Marshaler.
func (c *Counter) MarshalJSON() ([]byte, error) {
	if c == nil || c.unique == 0 {
		return []byte("null"), nil
	}
	data := make([]byte, 0, c.unique*28+2)
	data = append(data, '[')
	first := true
	for value, count := range c.counts {
		if count == 0 {
			continue
		}
		if first {
			first = false
		} else {
			data = append(data, ',')
		}
		data = append(data, `{"value":`...)
		data = append(data, byteDecimalStrings[value]...)
		data = append(data, `,"count":`...)
		data = strconv.AppendInt(data, int64(count), 10)
		data = append(data, '}')
	}
	data = append(data, ']')
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Counter) UnmarshalJSON(data []byte) error {
	if c == nil {
		return fmt.Errorf("unmarshal byte counter JSON: nil receiver")
	}
	var entries []CounterEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("unmarshal byte counter JSON: %w", err)
	}
	c.Clear()
	for _, entry := range entries {
		if entry.Count < 0 {
			return fmt.Errorf("unmarshal byte counter JSON: negative count for value %d", entry.Value)
		}
		c.AddN(entry.Value, entry.Count)
	}
	return nil
}

// String implements fmt.Stringer.
func (c *Counter) String() string {
	data, err := c.MarshalJSON()
	return jsonResultString(data, err, "[]")
}

func jsonResultString(data []byte, err error, fallback string) string {
	if err != nil {
		return fallback
	}
	return string(data)
}
