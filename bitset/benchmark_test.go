package bitset_test

import (
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
