package bytex_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/bytex"
)

const benchBytes = 1 << 12

func buildBenchBytes() []byte {
	values := make([]byte, benchBytes)
	for i := range values {
		values[i] = byte(i)
	}
	return values
}

func BenchmarkListWriteString(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		l := bytex.NewListWithCapacity(benchBytes)
		for range 64 {
			if _, err := l.WriteString("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ!!"); err != nil {
				b.Fatalf("WriteString() error = %v", err)
			}
		}
		if l.Len() != benchBytes {
			b.Fatalf("unexpected length: %d", l.Len())
		}
	}
}

func BenchmarkListIndexSequence(b *testing.B) {
	l := bytex.NewList(buildBenchBytes()...)
	needle := []byte{250, 251, 252, 253}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if l.IndexSequence(needle) < 0 {
			b.Fatal("needle not found")
		}
	}
}

func BenchmarkListWriteTo(b *testing.B) {
	l := bytex.NewList(buildBenchBytes()...)
	var out bytes.Buffer

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		out.Reset()
		if _, err := l.WriteTo(&out); err != nil {
			b.Fatalf("WriteTo() error = %v", err)
		}
	}
}

func BenchmarkRingBufferWrite(b *testing.B) {
	values := buildBenchBytes()
	r := bytex.NewRingBuffer(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := r.Write(values); err != nil {
			b.Fatalf("Write() error = %v", err)
		}
	}
}

func BenchmarkRingBufferWriteString(b *testing.B) {
	value := string(buildBenchBytes())
	r := bytex.NewRingBuffer(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := r.WriteString(value); err != nil {
			b.Fatalf("WriteString() error = %v", err)
		}
	}
}

func BenchmarkRingBufferViewSegments(b *testing.B) {
	r := bytex.NewRingBuffer(1024)
	_, _ = r.Write(buildBenchBytes())

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		r.ViewSegments(func(first, second []byte) {
			if len(first)+len(second) != r.Len() {
				b.Fatal("unexpected segment length")
			}
		})
	}
}

func BenchmarkRingBufferString(b *testing.B) {
	r := bytex.NewRingBuffer(1024)
	_, _ = r.Write(buildBenchBytes())

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if len(r.String()) != r.Len() {
			b.Fatal("unexpected string length")
		}
	}
}

func BenchmarkSetAddContains(b *testing.B) {
	values := buildBenchBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		s := bytex.NewSet()
		for _, value := range values {
			s.Set(value)
			_ = s.Contains(value)
		}
		if s.Len() != 256 {
			b.Fatalf("unexpected set length: %d", s.Len())
		}
	}
}

func BenchmarkSetRange(b *testing.B) {
	s := bytex.NewSet(buildBenchBytes()...)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		s.Range(func(_ byte) bool { return true })
	}
}

func BenchmarkSetAddRange(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := bytex.NewSet()
		if added := s.AddRange(0, 256); added != 256 {
			b.Fatalf("unexpected added count: %d", added)
		}
	}
}

func BenchmarkCounterAdd(b *testing.B) {
	values := buildBenchBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		c := bytex.NewCounter()
		c.Add(values...)
		if c.UniqueLen() != 256 {
			b.Fatalf("unexpected unique length: %d", c.UniqueLen())
		}
	}
}

func BenchmarkCounterAddString(b *testing.B) {
	value := string(buildBenchBytes())

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		c := bytex.NewCounter()
		c.AddString(value)
		if c.UniqueLen() != 256 {
			b.Fatalf("unexpected unique length: %d", c.UniqueLen())
		}
	}
}

func BenchmarkCounterMostCommon(b *testing.B) {
	c := bytex.NewCounter(buildBenchBytes()...)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		common := c.MostCommon(8)
		if len(common) != 8 {
			b.Fatalf("unexpected most common length: %d", len(common))
		}
	}
}

func BenchmarkCounterMostCommonAll(b *testing.B) {
	c := bytex.NewCounter(buildBenchBytes()...)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		common := c.MostCommon(256)
		if len(common) != 256 {
			b.Fatalf("unexpected most common length: %d", len(common))
		}
	}
}

func BenchmarkCounterMarshalJSON(b *testing.B) {
	c := bytex.NewCounter(buildBenchBytes()...)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		data, err := json.Marshal(c)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
		if len(data) == 0 {
			b.Fatal("json.Marshal() returned empty data")
		}
	}
}

func BenchmarkSetMarshalJSON(b *testing.B) {
	s := bytex.NewSet(buildBenchBytes()...)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		data, err := json.Marshal(s)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
		if len(data) == 0 {
			b.Fatal("json.Marshal() returned empty data")
		}
	}
}
