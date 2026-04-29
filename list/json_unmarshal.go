package list

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements json.Unmarshaler.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal list json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal list json: %w", err)
	}

	*l = List[T]{
		items: items,
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (g *Grid[T]) UnmarshalJSON(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal grid json: nil receiver")
	}

	var rows [][]T
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("unmarshal grid json: %w", err)
	}

	*g = *NewGridWithCapacity[T](len(rows), rows...)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (g *ConcurrentGrid[T]) UnmarshalJSON(data []byte) error {
	if g == nil {
		return fmt.Errorf("unmarshal concurrent grid json: nil receiver")
	}

	var rows [][]T
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("unmarshal concurrent grid json: %w", err)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.core = NewGridWithCapacity[T](len(rows), rows...)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *ConcurrentList[T]) UnmarshalJSON(data []byte) error {
	if l == nil {
		return fmt.Errorf("unmarshal concurrent list json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent list json: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.core = NewListWithCapacity[T](len(items), items...)
	l.jsonCache = nil
	l.stringCache = ""
	l.jsonDirty = false
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Deque[T]) UnmarshalJSON(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal deque json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal deque json: %w", err)
	}

	*d = *NewDeque(items...)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *ConcurrentDeque[T]) UnmarshalJSON(data []byte) error {
	if d == nil {
		return fmt.Errorf("unmarshal concurrent deque json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent deque json: %w", err)
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.deque = NewDeque(items...)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RingBuffer[T]) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal ring buffer json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal ring buffer json: %w", err)
	}

	next := NewRingBuffer[T](len(items))
	for _, item := range items {
		_ = next.Push(item)
	}
	*r = *next
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *ConcurrentRingBuffer[T]) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal concurrent ring buffer json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal concurrent ring buffer json: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.buffer = NewRingBuffer[T](len(items))
	for _, item := range items {
		_ = r.buffer.Push(item)
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RopeList[T]) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("unmarshal rope list json: nil receiver")
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal rope list json: %w", err)
	}

	*r = *NewRopeList(items...)
	return nil
}
