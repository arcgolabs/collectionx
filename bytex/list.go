package bytex

import (
	"bytes"
	"errors"
	"io"
	"slices"
)

const readFromBufferSize = 32 * 1024

var (
	errNilList   = errors.New("bytex: nil list receiver")
	errNilReader = errors.New("bytex: nil reader")
	errNilWriter = errors.New("bytex: nil writer")
)

// List is a byte-specialized list backed by a slice.
// Zero value is ready to use.
type List struct {
	items []byte
}

// NewList creates a byte list and copies optional items.
func NewList(items ...byte) *List {
	return NewListWithCapacity(len(items), items...)
}

// NewListFromString creates a byte list from a string.
func NewListFromString(value string) *List {
	l := NewListWithCapacity(len(value))
	_, _ = l.WriteString(value)
	return l
}

// NewListWithCapacity creates a byte list with preallocated capacity.
func NewListWithCapacity(capacity int, items ...byte) *List {
	if capacity < len(items) {
		capacity = len(items)
	}
	if capacity <= 0 {
		return &List{}
	}
	l := &List{items: make([]byte, 0, capacity)}
	l.Add(items...)
	return l
}

// WrapList creates a byte list using items as the backing slice without copying.
// The caller must not mutate items concurrently with the list.
func WrapList(items []byte) *List {
	return &List{items: items}
}

// Add appends one or more bytes.
func (l *List) Add(items ...byte) {
	if l == nil || len(items) == 0 {
		return
	}
	l.items = append(l.items, items...)
}

// Write appends p and implements io.Writer.
func (l *List) Write(p []byte) (int, error) {
	if l == nil {
		return 0, errNilList
	}
	l.items = append(l.items, p...)
	return len(p), nil
}

// WriteByte appends one byte.
func (l *List) WriteByte(value byte) error {
	if l == nil {
		return errNilList
	}
	l.items = append(l.items, value)
	return nil
}

// WriteString appends a string without an intermediate []byte allocation.
func (l *List) WriteString(value string) (int, error) {
	if l == nil {
		return 0, errNilList
	}
	l.items = append(l.items, value...)
	return len(value), nil
}

// ReadFrom appends all bytes from r and implements io.ReaderFrom.
func (l *List) ReadFrom(r io.Reader) (int64, error) {
	if l == nil {
		return 0, errNilList
	}
	if r == nil {
		return 0, errNilReader
	}

	var total int64
	buf := make([]byte, readFromBufferSize)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			l.items = append(l.items, buf[:n]...)
			total += int64(n)
		}
		if errors.Is(err, io.EOF) {
			return total, nil
		}
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, io.ErrNoProgress
		}
	}
}

// WriteTo writes all bytes to w and implements io.WriterTo.
func (l *List) WriteTo(w io.Writer) (int64, error) {
	if w == nil {
		return 0, errNilWriter
	}
	if l == nil || len(l.items) == 0 {
		return 0, nil
	}

	n, err := w.Write(l.items)
	if n < len(l.items) && err == nil {
		err = io.ErrShortWrite
	}
	return int64(n), err
}

// Reader returns a read-only reader over the current backing bytes without copying.
// Do not mutate the list concurrently with the returned reader.
func (l *List) Reader() *bytes.Reader {
	return bytes.NewReader(listBytes(l))
}

// Grow reserves capacity for n more bytes.
func (l *List) Grow(n int) {
	if l == nil || n <= 0 {
		return
	}
	l.items = slices.Grow(l.items, n)
}

// Get returns the byte at index.
func (l *List) Get(index int) (byte, bool) {
	if l == nil || index < 0 || index >= len(l.items) {
		return 0, false
	}
	return l.items[index], true
}

// GetFirst returns the first byte.
func (l *List) GetFirst() (byte, bool) {
	return l.Get(0)
}

// GetLast returns the last byte.
func (l *List) GetLast() (byte, bool) {
	if l == nil || len(l.items) == 0 {
		return 0, false
	}
	return l.items[len(l.items)-1], true
}

// Set replaces the byte at index.
func (l *List) Set(index int, value byte) bool {
	if l == nil || index < 0 || index >= len(l.items) {
		return false
	}
	l.items[index] = value
	return true
}

// RemoveAt removes and returns the byte at index.
func (l *List) RemoveAt(index int) (byte, bool) {
	if l == nil || index < 0 || index >= len(l.items) {
		return 0, false
	}
	value := l.items[index]
	copy(l.items[index:], l.items[index+1:])
	l.items[len(l.items)-1] = 0
	l.items = l.items[:len(l.items)-1]
	return value, true
}

// PopFirst removes and returns the first byte.
func (l *List) PopFirst() (byte, bool) {
	return l.RemoveAt(0)
}

// PopLast removes and returns the last byte.
func (l *List) PopLast() (byte, bool) {
	if l == nil || len(l.items) == 0 {
		return 0, false
	}
	lastIndex := len(l.items) - 1
	value := l.items[lastIndex]
	l.items[lastIndex] = 0
	l.items = l.items[:lastIndex]
	return value, true
}

// RemoveRange removes bytes in [start, end) and returns the removed count.
func (l *List) RemoveRange(start, end int) int {
	if l == nil || start < 0 || end > len(l.items) || start >= end {
		return 0
	}
	removed := end - start
	copy(l.items[start:], l.items[end:])
	clear(l.items[len(l.items)-removed:])
	l.items = l.items[:len(l.items)-removed]
	return removed
}

// DrainPrefix removes and returns a copy of the first n bytes.
func (l *List) DrainPrefix(n int) ([]byte, bool) {
	if l == nil || n < 0 || n > len(l.items) {
		return nil, false
	}
	out := slices.Clone(l.items[:n])
	l.RemoveRange(0, n)
	return out, true
}

// DrainSuffix removes and returns a copy of the last n bytes.
func (l *List) DrainSuffix(n int) ([]byte, bool) {
	if l == nil || n < 0 || n > len(l.items) {
		return nil, false
	}
	start := len(l.items) - n
	out := slices.Clone(l.items[start:])
	clear(l.items[start:])
	l.items = l.items[:start]
	return out, true
}

// Len returns the byte count.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Cap returns the backing slice capacity.
func (l *List) Cap() int {
	if l == nil {
		return 0
	}
	return cap(l.items)
}

// IsEmpty reports whether the list has no bytes.
func (l *List) IsEmpty() bool {
	return l.Len() == 0
}

// Clear removes all bytes and releases the backing slice.
func (l *List) Clear() {
	if l == nil {
		return
	}
	l.items = nil
}

// Reset removes all bytes while retaining capacity.
func (l *List) Reset() {
	if l == nil {
		return
	}
	clear(l.items)
	l.items = l.items[:0]
}

// Values returns a copy of bytes.
func (l *List) Values() []byte {
	return l.Bytes()
}

// Bytes returns a copy of bytes.
func (l *List) Bytes() []byte {
	if l == nil || len(l.items) == 0 {
		return nil
	}
	return slices.Clone(l.items)
}

// Snapshot returns a copy of bytes.
func (l *List) Snapshot() []byte {
	return l.Bytes()
}

// ViewBytes passes the internal backing slice to fn without copying.
// The slice must be treated as read-only and must not be retained.
func (l *List) ViewBytes(fn func(items []byte)) {
	if l == nil || fn == nil {
		return
	}
	fn(l.items)
}

// Slice returns a copy of bytes in [start, end).
func (l *List) Slice(start, end int) (*List, bool) {
	if l == nil || start < 0 || end > len(l.items) || start > end {
		return nil, false
	}
	return NewList(l.items[start:end]...), true
}

// ViewSlice passes the internal subslice [start, end) to fn without copying.
// The slice must be treated as read-only and must not be retained.
func (l *List) ViewSlice(start, end int, fn func(items []byte)) bool {
	if l == nil || fn == nil || start < 0 || end > len(l.items) || start > end {
		return false
	}
	fn(l.items[start:end])
	return true
}

// Range iterates bytes from left to right until fn returns false.
func (l *List) Range(fn func(index int, value byte) bool) {
	if l == nil || fn == nil {
		return
	}
	for index, value := range l.items {
		if !fn(index, value) {
			return
		}
	}
}

// Clone returns a copy of the list.
func (l *List) Clone() *List {
	if l == nil {
		return &List{}
	}
	return NewList(l.items...)
}

// String returns the bytes interpreted as a string.
func (l *List) String() string {
	if l == nil || len(l.items) == 0 {
		return ""
	}
	return string(l.items)
}

// Equal reports whether two lists contain the same bytes.
func (l *List) Equal(other *List) bool {
	return bytes.Equal(listBytes(l), listBytes(other))
}

// EqualBytes reports whether the list contains the same bytes as items.
func (l *List) EqualBytes(items []byte) bool {
	return bytes.Equal(listBytes(l), items)
}

// CompareBytes compares the list with items using bytes.Compare.
func (l *List) CompareBytes(items []byte) int {
	return bytes.Compare(listBytes(l), items)
}

// Contains reports whether value exists.
func (l *List) Contains(value byte) bool {
	return l.Index(value) >= 0
}

// ContainsSequence reports whether sequence exists.
func (l *List) ContainsSequence(sequence []byte) bool {
	return l.IndexSequence(sequence) >= 0
}

// ContainsAny reports whether any UTF-8 code point in chars exists.
func (l *List) ContainsAny(chars string) bool {
	return l.IndexAny(chars) >= 0
}

// Index returns the first index of value, or -1 when absent.
func (l *List) Index(value byte) int {
	if l == nil {
		return -1
	}
	return bytes.IndexByte(l.items, value)
}

// LastIndex returns the last index of value, or -1 when absent.
func (l *List) LastIndex(value byte) int {
	if l == nil {
		return -1
	}
	return bytes.LastIndexByte(l.items, value)
}

// IndexSequence returns the first index of sequence, or -1 when absent.
func (l *List) IndexSequence(sequence []byte) int {
	if l == nil || len(sequence) == 0 {
		return -1
	}
	return bytes.Index(l.items, sequence)
}

// IndexAny returns the first index of any UTF-8 code point in chars, or -1 when absent.
func (l *List) IndexAny(chars string) int {
	if l == nil || chars == "" {
		return -1
	}
	return bytes.IndexAny(l.items, chars)
}

// Count returns the number of occurrences of value.
func (l *List) Count(value byte) int {
	if l == nil || len(l.items) == 0 {
		return 0
	}
	count := 0
	for _, item := range l.items {
		if item == value {
			count++
		}
	}
	return count
}

// CountSequence returns the number of non-overlapping occurrences of sequence.
func (l *List) CountSequence(sequence []byte) int {
	if l == nil || len(sequence) == 0 {
		return 0
	}
	return bytes.Count(l.items, sequence)
}

// HasPrefix reports whether the list starts with prefix.
func (l *List) HasPrefix(prefix []byte) bool {
	return bytes.HasPrefix(listBytes(l), prefix)
}

// HasSuffix reports whether the list ends with suffix.
func (l *List) HasSuffix(suffix []byte) bool {
	return bytes.HasSuffix(listBytes(l), suffix)
}

// CutPrefix returns a copy of the bytes after prefix when prefix is present.
func (l *List) CutPrefix(prefix []byte) (*List, bool) {
	if l == nil || !bytes.HasPrefix(l.items, prefix) {
		return nil, false
	}
	return NewList(l.items[len(prefix):]...), true
}

// CutSuffix returns a copy of the bytes before suffix when suffix is present.
func (l *List) CutSuffix(suffix []byte) (*List, bool) {
	if l == nil || !bytes.HasSuffix(l.items, suffix) {
		return nil, false
	}
	return NewList(l.items[:len(l.items)-len(suffix)]...), true
}

// TrimPrefix removes prefix when present and reports whether it changed the list.
func (l *List) TrimPrefix(prefix []byte) bool {
	if l == nil || len(prefix) == 0 || !bytes.HasPrefix(l.items, prefix) {
		return false
	}
	l.RemoveRange(0, len(prefix))
	return true
}

// TrimSuffix removes suffix when present and reports whether it changed the list.
func (l *List) TrimSuffix(suffix []byte) bool {
	if l == nil || len(suffix) == 0 || !bytes.HasSuffix(l.items, suffix) {
		return false
	}
	l.items = l.items[:len(l.items)-len(suffix)]
	return true
}

// TrimSpace removes leading and trailing Unicode whitespace.
func (l *List) TrimSpace() bool {
	if l == nil || len(l.items) == 0 {
		return false
	}
	trimmed := bytes.TrimSpace(l.items)
	if len(trimmed) == len(l.items) {
		return false
	}
	copy(l.items, trimmed)
	clear(l.items[len(trimmed):])
	l.items = l.items[:len(trimmed)]
	return true
}

func listBytes(l *List) []byte {
	if l == nil {
		return nil
	}
	return l.items
}
