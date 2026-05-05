package bitset_test

import (
	"encoding/json"
	"testing"

	"github.com/arcgolabs/collectionx/bitset"
)

const benchBitSetBits = 1 << 12

func buildBenchBitSet() *bitset.BitSet {
	b := &bitset.BitSet{}
	for i := 0; i < benchBitSetBits; i += 2 {
		b.Set(i)
	}
	return b
}

func BenchmarkBitSetSetContains(b *testing.B) {
	s := buildBenchBitSet()
	mask := benchBitSetBits - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		bit := i & mask
		s.Set(bit)
		_ = s.Contains(bit)
	}
}

func BenchmarkBitSetRange(b *testing.B) {
	s := buildBenchBitSet()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		s.Range(func(_ int) bool { return true })
	}
}

func BenchmarkBitSetUnion(b *testing.B) {
	left := buildBenchBitSet()
	right := &bitset.BitSet{}
	for i := 1; i < benchBitSetBits; i += 2 {
		right.Set(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		out := left.Union(right)
		if out.Len() != benchBitSetBits {
			b.Fatalf("unexpected union length: %d", out.Len())
		}
	}
}

func BenchmarkBitSetSymmetricDifference(b *testing.B) {
	left := buildBenchBitSet()
	right := &bitset.BitSet{}
	for i := 1; i < benchBitSetBits; i += 2 {
		right.Set(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		out := left.SymmetricDifference(right)
		if out.Len() != benchBitSetBits {
			b.Fatalf("unexpected symmetric difference length: %d", out.Len())
		}
	}
}

func BenchmarkBitSetIntersects(b *testing.B) {
	left := buildBenchBitSet()
	right := &bitset.BitSet{}
	right.AddRange(benchBitSetBits-64, benchBitSetBits)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if !left.Intersects(right) {
			b.Fatal("Intersects() returned false")
		}
	}
}

func BenchmarkBitSetIsSubsetOf(b *testing.B) {
	left := buildBenchBitSet()
	right := &bitset.BitSet{}
	right.AddRange(0, benchBitSetBits)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if !left.IsSubsetOf(right) {
			b.Fatal("IsSubsetOf() returned false")
		}
	}
}

func BenchmarkBitSetAddRange(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := &bitset.BitSet{}
		s.AddRange(128, benchBitSetBits)
		if s.Len() != benchBitSetBits-128 {
			b.Fatalf("unexpected length: %d", s.Len())
		}
	}
}

func BenchmarkBitSetRemoveRange(b *testing.B) {
	s := &bitset.BitSet{}
	s.AddRange(0, benchBitSetBits)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		removed := s.RemoveRange(128, benchBitSetBits-128)
		if removed != benchBitSetBits-256 {
			b.Fatalf("unexpected removed count: %d", removed)
		}
		s.AddRange(128, benchBitSetBits-128)
	}
}

func BenchmarkBitSetNextSet(b *testing.B) {
	s := buildBenchBitSet()
	span := benchBitSetBits - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		if _, ok := s.NextSet(i % span); !ok {
			b.Fatal("NextSet returned false")
		}
	}
}

func BenchmarkBitSetMarshalJSON(b *testing.B) {
	s := buildBenchBitSet()

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

func BenchmarkBitSetMarshalBinary(b *testing.B) {
	s := buildBenchBitSet()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		data, err := s.MarshalBinary()
		if err != nil {
			b.Fatalf("MarshalBinary() error = %v", err)
		}
		if len(data) == 0 {
			b.Fatal("MarshalBinary() returned empty data")
		}
	}
}
