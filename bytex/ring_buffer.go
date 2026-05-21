package bytex

import (
	"errors"
	"io"
	"strings"
)

var errNilRingBuffer = errors.New("bytex: nil ring buffer receiver")

// RingBuffer is a fixed-capacity FIFO circular byte buffer.
// When full, writes overwrite the oldest bytes.
type RingBuffer struct {
	buf  []byte
	head int
	size int
	mask int
}

// NewRingBuffer creates a ring buffer with fixed capacity.
// capacity <= 0 creates an empty buffer that accepts writes but stores nothing.
func NewRingBuffer(capacity int) *RingBuffer {
	if capacity < 0 {
		capacity = 0
	}
	return &RingBuffer{
		buf:  make([]byte, capacity),
		mask: ringBufferMask(capacity),
	}
}

// Capacity returns max byte capacity.
func (r *RingBuffer) Capacity() int {
	if r == nil {
		return 0
	}
	return len(r.buf)
}

// Len returns current byte count.
func (r *RingBuffer) Len() int {
	if r == nil {
		return 0
	}
	return r.size
}

// IsEmpty reports whether the buffer has no bytes.
func (r *RingBuffer) IsEmpty() bool {
	return r.Len() == 0
}

// IsFull reports whether the buffer reached capacity.
func (r *RingBuffer) IsFull() bool {
	return r != nil && r.size == len(r.buf) && len(r.buf) > 0
}

// Push writes one byte at the tail.
// If full, the oldest byte is evicted and returned.
func (r *RingBuffer) Push(value byte) (byte, bool) {
	if r == nil || len(r.buf) == 0 {
		return 0, false
	}
	if r.size < len(r.buf) {
		tail := r.wrap(r.head + r.size)
		r.buf[tail] = value
		r.size++
		return 0, false
	}

	evicted := r.buf[r.head]
	r.buf[r.head] = value
	r.head = r.wrap(r.head + 1)
	return evicted, true
}

// WriteByte writes one byte and implements byte-oriented writers.
func (r *RingBuffer) WriteByte(value byte) error {
	if r == nil {
		return errNilRingBuffer
	}
	r.Push(value)
	return nil
}

// Write appends p and implements io.Writer.
// It accepts all bytes and keeps only the last Capacity() bytes.
func (r *RingBuffer) Write(p []byte) (int, error) {
	if r == nil {
		return 0, errNilRingBuffer
	}
	if len(p) == 0 {
		return 0, nil
	}
	capacity := len(r.buf)
	if capacity == 0 {
		return len(p), nil
	}
	if len(p) >= capacity {
		copy(r.buf, p[len(p)-capacity:])
		r.head = 0
		r.size = capacity
		return len(p), nil
	}

	r.dropPrefix(r.size + len(p) - capacity)
	r.writeTailBytes(p)
	return len(p), nil
}

// WriteString appends value without an intermediate []byte allocation.
// It accepts all bytes and keeps only the last Capacity() bytes.
func (r *RingBuffer) WriteString(value string) (int, error) {
	if r == nil {
		return 0, errNilRingBuffer
	}
	if value == "" {
		return 0, nil
	}
	capacity := len(r.buf)
	if capacity == 0 {
		return len(value), nil
	}
	if len(value) >= capacity {
		copy(r.buf, value[len(value)-capacity:])
		r.head = 0
		r.size = capacity
		return len(value), nil
	}

	r.dropPrefix(r.size + len(value) - capacity)
	r.writeTailString(value)
	return len(value), nil
}

// ReadFrom appends all bytes from reader and implements io.ReaderFrom.
func (r *RingBuffer) ReadFrom(reader io.Reader) (int64, error) {
	if r == nil {
		return 0, errNilRingBuffer
	}
	if reader == nil {
		return 0, errNilReader
	}

	var total int64
	buf := make([]byte, readFromBufferSize)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, writeErr := r.Write(buf[:n]); writeErr != nil {
				return total, writeErr
			}
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

// WriteTo writes buffered bytes to writer from oldest to newest.
func (r *RingBuffer) WriteTo(writer io.Writer) (int64, error) {
	if writer == nil {
		return 0, errNilWriter
	}
	if r == nil || r.size == 0 {
		return 0, nil
	}

	first, second := r.segments()
	total, err := writeFull(writer, first)
	if err != nil {
		return total, err
	}
	written, err := writeFull(writer, second)
	total += written
	return total, err
}

// Read copies and removes bytes from the head and implements io.Reader.
func (r *RingBuffer) Read(p []byte) (int, error) {
	if r == nil {
		return 0, errNilRingBuffer
	}
	if len(p) == 0 {
		return 0, nil
	}
	if r.size == 0 {
		return 0, io.EOF
	}
	n := r.CopyTo(p)
	r.Discard(n)
	return n, nil
}

// Pop removes and returns the oldest byte.
func (r *RingBuffer) Pop() (byte, bool) {
	if r == nil || r.size == 0 {
		return 0, false
	}
	value := r.buf[r.head]
	r.Discard(1)
	return value, true
}

// PopN removes and returns a copy of the oldest n bytes.
func (r *RingBuffer) PopN(n int) ([]byte, bool) {
	if r == nil || n < 0 || n > r.size {
		return nil, false
	}
	out := make([]byte, n)
	r.CopyTo(out)
	r.Discard(n)
	return out, true
}

// Discard removes up to n oldest bytes and returns the removed count.
func (r *RingBuffer) Discard(n int) int {
	if r == nil || n <= 0 || r.size == 0 {
		return 0
	}
	if n > r.size {
		n = r.size
	}
	r.clearPrefix(n)
	r.dropPrefix(n)
	return n
}

// Peek returns the oldest byte without removing it.
func (r *RingBuffer) Peek() (byte, bool) {
	if r == nil || r.size == 0 {
		return 0, false
	}
	return r.buf[r.head], true
}

// GetFirst returns the oldest byte without removing it.
func (r *RingBuffer) GetFirst() (byte, bool) {
	return r.Peek()
}

// GetLast returns the newest byte without removing it.
func (r *RingBuffer) GetLast() (byte, bool) {
	if r == nil || r.size == 0 {
		return 0, false
	}
	return r.buf[r.wrap(r.head+r.size-1)], true
}

// Bytes returns a copy of buffered bytes from oldest to newest.
func (r *RingBuffer) Bytes() []byte {
	if r == nil || r.size == 0 {
		return nil
	}
	out := make([]byte, r.size)
	r.CopyTo(out)
	return out
}

// Values returns a copy of buffered bytes from oldest to newest.
func (r *RingBuffer) Values() []byte {
	return r.Bytes()
}

// Snapshot returns a copy of buffered bytes from oldest to newest.
func (r *RingBuffer) Snapshot() []byte {
	return r.Bytes()
}

// CopyTo copies buffered bytes from oldest to newest into dst.
func (r *RingBuffer) CopyTo(dst []byte) int {
	if r == nil || r.size == 0 || len(dst) == 0 {
		return 0
	}
	n := min(len(dst), r.size)
	firstLen := min(n, len(r.buf)-r.head)
	copied := copy(dst, r.buf[r.head:r.head+firstLen])
	if copied < n {
		copied += copy(dst[copied:], r.buf[:n-copied])
	}
	return copied
}

// ViewSegments passes one or two internal segments from oldest to newest without copying.
// The segments must be treated as read-only and must not be retained.
func (r *RingBuffer) ViewSegments(fn func(first, second []byte)) bool {
	if r == nil || fn == nil {
		return false
	}
	if r.size == 0 {
		fn(nil, nil)
		return true
	}
	firstLen := min(r.size, len(r.buf)-r.head)
	first := r.buf[r.head : r.head+firstLen]
	var second []byte
	if remaining := r.size - firstLen; remaining > 0 {
		second = r.buf[:remaining]
	}
	fn(first, second)
	return true
}

// Clone returns a copy with the same capacity and byte order.
func (r *RingBuffer) Clone() *RingBuffer {
	if r == nil {
		return NewRingBuffer(0)
	}
	clone := NewRingBuffer(r.Capacity())
	_, _ = clone.Write(r.Bytes())
	return clone
}

// Clear removes all buffered bytes and retains capacity.
func (r *RingBuffer) Clear() {
	if r == nil {
		return
	}
	clear(r.buf)
	r.head = 0
	r.size = 0
}

// String returns buffered bytes interpreted as a string.
func (r *RingBuffer) String() string {
	if r == nil || r.size == 0 {
		return ""
	}
	first, second := r.segments()
	if len(second) == 0 {
		return string(first)
	}
	var builder strings.Builder
	builder.Grow(r.size)
	builder.Write(first)
	builder.Write(second)
	return builder.String()
}

func (r *RingBuffer) writeTailBytes(p []byte) {
	if len(p) == 0 {
		return
	}
	tail := r.wrap(r.head + r.size)
	firstLen := min(len(p), len(r.buf)-tail)
	copy(r.buf[tail:], p[:firstLen])
	if firstLen < len(p) {
		copy(r.buf, p[firstLen:])
	}
	r.size += len(p)
}

func (r *RingBuffer) writeTailString(value string) {
	if value == "" {
		return
	}
	tail := r.wrap(r.head + r.size)
	firstLen := min(len(value), len(r.buf)-tail)
	copy(r.buf[tail:], value[:firstLen])
	if firstLen < len(value) {
		copy(r.buf, value[firstLen:])
	}
	r.size += len(value)
}

func (r *RingBuffer) clearPrefix(n int) {
	firstLen := min(n, len(r.buf)-r.head)
	clear(r.buf[r.head : r.head+firstLen])
	if firstLen < n {
		clear(r.buf[:n-firstLen])
	}
}

func (r *RingBuffer) segments() ([]byte, []byte) {
	firstLen := min(r.size, len(r.buf)-r.head)
	first := r.buf[r.head : r.head+firstLen]
	var second []byte
	if remaining := r.size - firstLen; remaining > 0 {
		second = r.buf[:remaining]
	}
	return first, second
}

func (r *RingBuffer) dropPrefix(n int) {
	if n <= 0 {
		return
	}
	if n >= r.size {
		r.head = 0
		r.size = 0
		return
	}
	r.head = r.wrap(r.head + n)
	r.size -= n
}

func (r *RingBuffer) wrap(index int) int {
	if r.mask >= 0 {
		return index & r.mask
	}
	return index % len(r.buf)
}

func ringBufferMask(capacity int) int {
	if capacity > 0 && capacity&(capacity-1) == 0 {
		return capacity - 1
	}
	return -1
}

func writeFull(writer io.Writer, data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	n, err := writer.Write(data)
	if n < len(data) && err == nil {
		err = io.ErrShortWrite
	}
	return int64(n), err
}
