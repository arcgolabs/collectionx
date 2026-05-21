package bytex_test

import (
	"bytes"
	"encoding/json"
	"io"
	"slices"
	"testing"

	"github.com/arcgolabs/collectionx/bytex"
)

func TestListBasicOps(t *testing.T) {
	t.Parallel()

	l := bytex.NewListFromString("prefix-body-suffix")
	if !l.HasPrefix([]byte("prefix-")) || !l.HasSuffix([]byte("-suffix")) {
		t.Fatal("expected prefix and suffix")
	}
	if !l.TrimPrefix([]byte("prefix-")) || !l.TrimSuffix([]byte("-suffix")) {
		t.Fatal("expected trim operations to change the list")
	}
	if got := l.String(); got != "body" {
		t.Fatalf("unexpected list string: %q", got)
	}

	if _, err := l.WriteString("-body"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if got := l.IndexSequence([]byte("dy-bo")); got != 2 {
		t.Fatalf("unexpected sequence index: %d", got)
	}
	if got := l.CountSequence([]byte("body")); got != 2 {
		t.Fatalf("unexpected sequence count: %d", got)
	}
	if got := l.Count('o'); got != 2 {
		t.Fatalf("unexpected byte count: %d", got)
	}

	first, ok := l.GetFirst()
	if !ok || first != 'b' {
		t.Fatalf("unexpected first byte: %q %v", first, ok)
	}
	last, ok := l.GetLast()
	if !ok || last != 'y' {
		t.Fatalf("unexpected last byte: %q %v", last, ok)
	}

	if removed := l.RemoveRange(4, 5); removed != 1 {
		t.Fatalf("unexpected removed count: %d", removed)
	}
	if got := l.String(); got != "bodybody" {
		t.Fatalf("unexpected list after remove: %q", got)
	}
}

func TestListConsumeAndCut(t *testing.T) {
	t.Parallel()

	l := bytex.NewListFromString(" \talpha,beta\n ")
	if !l.TrimSpace() || l.String() != "alpha,beta" {
		t.Fatalf("unexpected TrimSpace result: %q", l.String())
	}
	if got := l.IndexAny(",;"); got != 5 {
		t.Fatalf("unexpected IndexAny result: %d", got)
	}
	if !l.ContainsAny(";a") {
		t.Fatal("expected ContainsAny to match")
	}

	afterPrefix, ok := l.CutPrefix([]byte("alpha,"))
	if !ok || afterPrefix.String() != "beta" {
		t.Fatalf("unexpected CutPrefix result: %q %v", afterPrefix, ok)
	}
	beforeSuffix, ok := l.CutSuffix([]byte(",beta"))
	if !ok || beforeSuffix.String() != "alpha" {
		t.Fatalf("unexpected CutSuffix result: %q %v", beforeSuffix, ok)
	}

	prefix, ok := l.DrainPrefix(6)
	if !ok || string(prefix) != "alpha," || l.String() != "beta" {
		t.Fatalf("unexpected DrainPrefix result: %q %v list=%q", string(prefix), ok, l.String())
	}
	suffix, ok := l.DrainSuffix(2)
	if !ok || string(suffix) != "ta" || l.String() != "be" {
		t.Fatalf("unexpected DrainSuffix result: %q %v list=%q", string(suffix), ok, l.String())
	}

	first, ok := l.PopFirst()
	if !ok || first != 'b' {
		t.Fatalf("unexpected PopFirst result: %q %v", first, ok)
	}
	last, ok := l.PopLast()
	if !ok || last != 'e' || !l.IsEmpty() {
		t.Fatalf("unexpected PopLast result: %q %v list=%q", last, ok, l.String())
	}

	if _, ok := l.DrainPrefix(1); ok {
		t.Fatal("expected DrainPrefix to reject out-of-range length")
	}
	if _, ok := l.PopLast(); ok {
		t.Fatal("expected PopLast on empty list to fail")
	}
}

func TestListViewsAndCopies(t *testing.T) {
	t.Parallel()

	source := []byte{1, 2, 3}
	l := bytex.WrapList(source)
	source[0] = 9
	if got, _ := l.Get(0); got != 9 {
		t.Fatalf("WrapList should use source backing slice, got %d", got)
	}

	values := l.Values()
	values[0] = 1
	if got, _ := l.Get(0); got != 9 {
		t.Fatalf("Values should return a copy, got %d", got)
	}

	var viewed []byte
	l.ViewBytes(func(items []byte) {
		viewed = slices.Clone(items)
	})
	if !bytes.Equal(viewed, []byte{9, 2, 3}) {
		t.Fatalf("unexpected viewed bytes: %v", viewed)
	}

	ok := l.ViewSlice(1, 3, func(items []byte) {
		if !bytes.Equal(items, []byte{2, 3}) {
			t.Fatalf("unexpected view slice: %v", items)
		}
	})
	if !ok {
		t.Fatal("expected ViewSlice to succeed")
	}

	sliced, ok := l.Slice(1, 3)
	if !ok || !sliced.EqualBytes([]byte{2, 3}) {
		t.Fatalf("unexpected copied slice: %v %v", sliced, ok)
	}
}

func TestListIOAndSerialization(t *testing.T) {
	t.Parallel()

	l := bytex.NewList()
	n, err := l.ReadFrom(bytes.NewBufferString("hello"))
	if err != nil || n != 5 {
		t.Fatalf("ReadFrom() = %d, %v", n, err)
	}

	var out bytes.Buffer
	written, err := l.WriteTo(&out)
	if err != nil || written != 5 || out.String() != "hello" {
		t.Fatalf("WriteTo() = %d, %v, %q", written, err, out.String())
	}

	reader := l.Reader()
	peek := make([]byte, 2)
	if n, err := reader.ReadAt(peek, 1); err != nil || n != 2 || string(peek) != "el" {
		t.Fatalf("Reader().ReadAt() = %d, %v, %q", n, err, string(peek))
	}

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(data) != `"aGVsbG8="` {
		t.Fatalf("unexpected JSON payload: %s", data)
	}

	var restored bytex.List
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if !restored.Equal(l) {
		t.Fatalf("unexpected JSON restored bytes: %q", restored.String())
	}

	raw, err := l.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}
	if !bytes.Equal(raw, []byte("hello")) {
		t.Fatalf("unexpected binary payload: %v", raw)
	}
	var binaryRestored bytex.List
	if err := binaryRestored.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if !binaryRestored.Equal(l) {
		t.Fatal("binary roundtrip mismatch")
	}

	if _, err := (*bytex.List)(nil).Write([]byte("x")); err == nil {
		t.Fatal("expected nil receiver write error")
	}
	if _, err := l.ReadFrom(nil); err == nil {
		t.Fatal("expected nil reader error")
	}
	if _, err := l.WriteTo(nil); err == nil {
		t.Fatal("expected nil writer error")
	}
	if _, err := l.WriteTo(shortWriter{}); err != io.ErrShortWrite {
		t.Fatalf("expected short write error, got %v", err)
	}
}

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	return len(p) - 1, nil
}
