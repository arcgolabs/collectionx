package bytex_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/arcgolabs/collectionx/bytex"
)

func TestRingBufferOverwriteAndViews(t *testing.T) {
	t.Parallel()

	r := bytex.NewRingBuffer(3)
	if evicted, ok := r.Push('a'); ok || evicted != 0 {
		t.Fatalf("unexpected first push eviction: %q %v", evicted, ok)
	}
	r.Push('b')
	r.Push('c')

	evicted, ok := r.Push('d')
	if !ok || evicted != 'a' {
		t.Fatalf("unexpected eviction: %q %v", evicted, ok)
	}
	if got := r.String(); got != "bcd" {
		t.Fatalf("unexpected ring values: %q", got)
	}
	first, ok := r.GetFirst()
	if !ok || first != 'b' {
		t.Fatalf("unexpected first byte: %q %v", first, ok)
	}
	last, ok := r.GetLast()
	if !ok || last != 'd' {
		t.Fatalf("unexpected last byte: %q %v", last, ok)
	}

	ok = r.ViewSegments(func(first, second []byte) {
		if string(first) != "bc" || string(second) != "d" {
			t.Fatalf("unexpected segments: %q %q", string(first), string(second))
		}
	})
	if !ok {
		t.Fatal("expected ViewSegments to succeed")
	}

	dst := make([]byte, 2)
	if copied := r.CopyTo(dst); copied != 2 || string(dst) != "bc" {
		t.Fatalf("unexpected CopyTo result: copied=%d dst=%q", copied, string(dst))
	}
}

func TestRingBufferWriteReadAndDiscard(t *testing.T) {
	t.Parallel()

	r := bytex.NewRingBuffer(4)
	if n, err := r.WriteString("abcdef"); err != nil || n != 6 {
		t.Fatalf("WriteString() = %d, %v", n, err)
	}
	if got := r.String(); got != "cdef" {
		t.Fatalf("unexpected ring after WriteString: %q", got)
	}
	if !r.IsFull() || r.Capacity() != 4 || r.Len() != 4 {
		t.Fatalf("unexpected capacity state: cap=%d len=%d full=%v", r.Capacity(), r.Len(), r.IsFull())
	}

	read := make([]byte, 2)
	n, err := r.Read(read)
	if err != nil || n != 2 || string(read) != "cd" {
		t.Fatalf("Read() = %d, %v, %q", n, err, string(read))
	}
	if got := r.String(); got != "ef" {
		t.Fatalf("unexpected ring after Read: %q", got)
	}

	if _, err := r.Write([]byte("ghij")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	popped, ok := r.PopN(3)
	if !ok || string(popped) != "ghi" {
		t.Fatalf("unexpected PopN result: %q %v", string(popped), ok)
	}
	if got := r.String(); got != "j" {
		t.Fatalf("unexpected ring after PopN: %q", got)
	}
	if discarded := r.Discard(10); discarded != 1 {
		t.Fatalf("unexpected discarded count: %d", discarded)
	}
	if !r.IsEmpty() {
		t.Fatalf("expected empty ring, got %q", r.String())
	}
	if n, err := r.Read(read); n != 0 || err != io.EOF {
		t.Fatalf("expected EOF read, got %d %v", n, err)
	}

	if _, ok := r.Pop(); ok {
		t.Fatal("expected Pop on empty buffer to fail")
	}
	if _, ok := r.PopN(1); ok {
		t.Fatal("expected PopN out of range to fail")
	}
}

func TestRingBufferIOAndSerialization(t *testing.T) {
	t.Parallel()

	r := bytex.NewRingBuffer(8)
	read, err := r.ReadFrom(bytes.NewBufferString("abcdef"))
	if err != nil || read != 6 {
		t.Fatalf("ReadFrom() = %d, %v", read, err)
	}

	var out bytes.Buffer
	written, err := r.WriteTo(&out)
	if err != nil || written != 6 || out.String() != "abcdef" {
		t.Fatalf("WriteTo() = %d, %v, %q", written, err, out.String())
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(data) != `"YWJjZGVm"` {
		t.Fatalf("unexpected JSON payload: %s", data)
	}
	var jsonRestored bytex.RingBuffer
	if err := json.Unmarshal(data, &jsonRestored); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if jsonRestored.String() != "abcdef" || jsonRestored.Capacity() != 6 {
		t.Fatalf("unexpected JSON restored ring: cap=%d value=%q", jsonRestored.Capacity(), jsonRestored.String())
	}

	raw, err := r.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}
	var binaryRestored bytex.RingBuffer
	if err := binaryRestored.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if binaryRestored.String() != "abcdef" || binaryRestored.Capacity() != 8 {
		t.Fatalf("unexpected binary restored ring: cap=%d value=%q", binaryRestored.Capacity(), binaryRestored.String())
	}

	if _, err := (*bytex.RingBuffer)(nil).Write([]byte("x")); err == nil {
		t.Fatal("expected nil receiver write error")
	}
	if _, err := r.ReadFrom(nil); err == nil {
		t.Fatal("expected nil reader error")
	}
	if _, err := r.WriteTo(nil); err == nil {
		t.Fatal("expected nil writer error")
	}
	if err := binaryRestored.UnmarshalBinary([]byte{1}); err == nil {
		t.Fatal("expected invalid binary length error")
	}
}

func TestRingBufferZeroCapacity(t *testing.T) {
	t.Parallel()

	r := bytex.NewRingBuffer(0)
	if n, err := r.WriteString("abc"); err != nil || n != 3 {
		t.Fatalf("WriteString() = %d, %v", n, err)
	}
	if !r.IsEmpty() || r.String() != "" {
		t.Fatalf("zero capacity buffer should remain empty, got %q", r.String())
	}
}
