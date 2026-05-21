package bytex

import (
	"encoding/binary"
	"fmt"
)

const (
	byteSetBinarySize        = byteSetWordCount * 8
	byteCounterHeaderSize    = 4
	byteCounterEntryBinSize  = 9
	byteRingBufferHeaderSize = 4
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (l *List) MarshalBinary() ([]byte, error) {
	return l.Bytes(), nil
}

// GobEncode implements gob.GobEncoder.
func (l *List) GobEncode() ([]byte, error) {
	return l.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (l *List) UnmarshalBinary(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal byte list binary: nil receiver")
	}
	*l = *NewList(data...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (l *List) GobDecode(data []byte) error {
	return l.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (r *RingBuffer) MarshalBinary() ([]byte, error) {
	values := r.Bytes()
	capacity := r.Capacity()
	data := make([]byte, byteRingBufferHeaderSize+len(values))
	binary.LittleEndian.PutUint32(data, uint32(capacity))
	copy(data[byteRingBufferHeaderSize:], values)
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (r *RingBuffer) GobEncode() ([]byte, error) {
	return r.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (r *RingBuffer) UnmarshalBinary(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal byte ring buffer binary: nil receiver")
	}
	if len(data) == 0 {
		*r = *NewRingBuffer(0)
		return nil
	}
	if len(data) < byteRingBufferHeaderSize {
		return fmt.Errorf("unmarshal byte ring buffer binary: invalid data length %d", len(data))
	}
	capacity := int(binary.LittleEndian.Uint32(data))
	values := data[byteRingBufferHeaderSize:]
	if capacity < len(values) {
		return fmt.Errorf("unmarshal byte ring buffer binary: capacity %d smaller than value length %d", capacity, len(values))
	}
	next := NewRingBuffer(capacity)
	_, _ = next.Write(values)
	*r = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (r *RingBuffer) GobDecode(data []byte) error {
	return r.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *Set) MarshalBinary() ([]byte, error) {
	data := make([]byte, byteSetBinarySize)
	if s == nil {
		return data, nil
	}
	for i, word := range s.words {
		binary.LittleEndian.PutUint64(data[i*8:], word)
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (s *Set) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Set) UnmarshalBinary(data []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal byte set binary: nil receiver")
	}
	if len(data) == 0 {
		s.Clear()
		return nil
	}
	if len(data) != byteSetBinarySize {
		return fmt.Errorf("unmarshal byte set binary: invalid data length %d", len(data))
	}
	for i := range s.words {
		s.words[i] = binary.LittleEndian.Uint64(data[i*8:])
	}
	s.recount()
	return nil
}

// GobDecode implements gob.GobDecoder.
func (s *Set) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (c *Counter) MarshalBinary() ([]byte, error) {
	entries := c.Entries()
	data := make([]byte, byteCounterHeaderSize+len(entries)*byteCounterEntryBinSize)
	binary.LittleEndian.PutUint32(data, uint32(len(entries)))
	offset := byteCounterHeaderSize
	for _, entry := range entries {
		data[offset] = entry.Value
		binary.LittleEndian.PutUint64(data[offset+1:], uint64(entry.Count))
		offset += byteCounterEntryBinSize
	}
	return data, nil
}

// GobEncode implements gob.GobEncoder.
func (c *Counter) GobEncode() ([]byte, error) {
	return c.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (c *Counter) UnmarshalBinary(data []byte) error {
	if c == nil {
		return fmt.Errorf("unmarshal byte counter binary: nil receiver")
	}
	if len(data) == 0 {
		c.Clear()
		return nil
	}
	if len(data) < byteCounterHeaderSize {
		return fmt.Errorf("unmarshal byte counter binary: invalid data length %d", len(data))
	}
	entryCount := binary.LittleEndian.Uint32(data)
	expectedLen := uint64(byteCounterHeaderSize) + uint64(entryCount)*uint64(byteCounterEntryBinSize)
	if uint64(len(data)) != expectedLen {
		return fmt.Errorf("unmarshal byte counter binary: invalid data length %d", len(data))
	}

	c.Clear()
	offset := byteCounterHeaderSize
	for range int(entryCount) {
		value := data[offset]
		count64 := binary.LittleEndian.Uint64(data[offset+1:])
		count := int(count64)
		if uint64(count) != count64 {
			return fmt.Errorf("unmarshal byte counter binary: count overflows int for value %d", value)
		}
		c.AddN(value, count)
		offset += byteCounterEntryBinSize
	}
	return nil
}

// GobDecode implements gob.GobDecoder.
func (c *Counter) GobDecode(data []byte) error {
	return c.UnmarshalBinary(data)
}
