package list

import (
	"fmt"
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (l *List[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("list", l.Values())
}

// GobEncode implements gob.GobEncoder.
func (l *List[T]) GobEncode() ([]byte, error) {
	return l.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (l *List[T]) UnmarshalBinary(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal list binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal list binary: %w", err)
	}
	*l = *NewListWithCapacity[T](len(items), items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (l *List[T]) GobDecode(data []byte) error {
	return l.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (g *Grid[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("grid", g.Values())
}

// GobEncode implements gob.GobEncoder.
func (g *Grid[T]) GobEncode() ([]byte, error) {
	return g.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (g *Grid[T]) UnmarshalBinary(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal grid binary: nil receiver")
	}
	var rows [][]T
	if err := unmarshalBinaryValue(data, &rows); err != nil {
		return fmt.Errorf("unmarshal grid binary: %w", err)
	}
	*g = *NewGridWithCapacity[T](len(rows), rows...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (g *Grid[T]) GobDecode(data []byte) error {
	return g.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (g *ConcurrentGrid[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("concurrent grid", g.Values())
}

// GobEncode implements gob.GobEncoder.
func (g *ConcurrentGrid[T]) GobEncode() ([]byte, error) {
	return g.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (g *ConcurrentGrid[T]) UnmarshalBinary(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal concurrent grid binary: nil receiver")
	}
	var rows [][]T
	if err := unmarshalBinaryValue(data, &rows); err != nil {
		return fmt.Errorf("unmarshal concurrent grid binary: %w", err)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.core = NewGridWithCapacity[T](len(rows), rows...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (g *ConcurrentGrid[T]) GobDecode(data []byte) error {
	return g.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (l *ConcurrentList[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("concurrent list", l.Values())
}

// GobEncode implements gob.GobEncoder.
func (l *ConcurrentList[T]) GobEncode() ([]byte, error) {
	return l.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (l *ConcurrentList[T]) UnmarshalBinary(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal concurrent list binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent list binary: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.core = NewListWithCapacity[T](len(items), items...)
	l.jsonCache = nil
	l.stringCache = ""
	l.jsonDirty = false
	return nil
}

// GobDecode implements gob.GobDecoder.
func (l *ConcurrentList[T]) GobDecode(data []byte) error {
	return l.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *Deque[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("deque", d.Values())
}

// GobEncode implements gob.GobEncoder.
func (d *Deque[T]) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *Deque[T]) UnmarshalBinary(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal deque binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal deque binary: %w", err)
	}
	*d = *NewDeque(items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (d *Deque[T]) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *ConcurrentDeque[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("concurrent deque", d.Values())
}

// GobEncode implements gob.GobEncoder.
func (d *ConcurrentDeque[T]) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *ConcurrentDeque[T]) UnmarshalBinary(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal concurrent deque binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent deque binary: %w", err)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.deque = NewDeque(items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (d *ConcurrentDeque[T]) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (r *RingBuffer[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("ring buffer", r.Values())
}

// GobEncode implements gob.GobEncoder.
func (r *RingBuffer[T]) GobEncode() ([]byte, error) {
	return r.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (r *RingBuffer[T]) UnmarshalBinary(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal ring buffer binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal ring buffer binary: %w", err)
	}
	next := NewRingBuffer[T](len(items))
	for _, item := range items {
		_ = next.Push(item)
	}
	*r = *next
	return nil
}

// GobDecode implements gob.GobDecoder.
func (r *RingBuffer[T]) GobDecode(data []byte) error {
	return r.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (r *ConcurrentRingBuffer[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("concurrent ring buffer", r.Values())
}

// GobEncode implements gob.GobEncoder.
func (r *ConcurrentRingBuffer[T]) GobEncode() ([]byte, error) {
	return r.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (r *ConcurrentRingBuffer[T]) UnmarshalBinary(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal concurrent ring buffer binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent ring buffer binary: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buffer = NewRingBuffer[T](len(items))
	for _, item := range items {
		_ = r.buffer.Push(item)
	}
	return nil
}

// GobDecode implements gob.GobDecoder.
func (r *ConcurrentRingBuffer[T]) GobDecode(data []byte) error {
	return r.UnmarshalBinary(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (r *RopeList[T]) MarshalBinary() ([]byte, error) {
	return marshalListBinary("rope list", r.Values())
}

// GobEncode implements gob.GobEncoder.
func (r *RopeList[T]) GobEncode() ([]byte, error) {
	return r.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (r *RopeList[T]) UnmarshalBinary(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal rope list binary: nil receiver")
	}
	var items []T
	if err := unmarshalBinaryValue(data, &items); err != nil {
		return fmt.Errorf("unmarshal rope list binary: %w", err)
	}
	*r = *NewRopeList(items...)
	return nil
}

// GobDecode implements gob.GobDecoder.
func (r *RopeList[T]) GobDecode(data []byte) error {
	return r.UnmarshalBinary(data)
}

func marshalListBinary(kind string, value any) ([]byte, error) {
	data, err := marshalBinaryValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s binary: %w", kind, err)
	}
	return data, nil
}
