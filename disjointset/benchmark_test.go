package disjointset_test

import (
	"testing"

	"github.com/arcgolabs/collectionx/disjointset"
)

const benchDisjointSetSize = 1 << 12

func buildBenchDisjointSet(tb testing.TB) *disjointset.DisjointSet[int] {
	tb.Helper()
	ds := disjointset.New[int]()
	for i := range benchDisjointSetSize {
		ds.Add(i)
	}
	for i := 1; i < benchDisjointSetSize; i++ {
		ds.Union(i-1, i)
	}
	return ds
}

func BenchmarkDisjointSetFind(b *testing.B) {
	ds := buildBenchDisjointSet(b)
	mask := benchDisjointSetSize - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		if _, ok := ds.Find(i & mask); !ok {
			b.Fatalf("Find(%d) failed", i&mask)
		}
	}
}

func BenchmarkDisjointSetConnected(b *testing.B) {
	ds := buildBenchDisjointSet(b)
	mask := (benchDisjointSetSize >> 1) - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		left := i & mask
		right := left + (benchDisjointSetSize >> 1)
		if !ds.Connected(left, right) {
			b.Fatalf("Connected(%d, %d) = false", left, right)
		}
	}
}

func BenchmarkDisjointSetGroups(b *testing.B) {
	ds := disjointset.New[int]()
	for i := range benchDisjointSetSize {
		ds.Add(i)
	}
	for i := 0; i < benchDisjointSetSize; i += 4 {
		ds.Union(i, i+1)
		ds.Union(i+2, i+3)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		groups := ds.Groups()
		if len(groups) == 0 {
			b.Fatal("Groups() returned empty result")
		}
	}
}

func BenchmarkDisjointSetUnionFresh(b *testing.B) {
	const nodesPerRun = 1 << 10

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		ds := disjointset.New[int]()
		for i := range nodesPerRun {
			ds.Add(i)
		}
		for i := 1; i < nodesPerRun; i++ {
			ds.Union(i-1, i)
		}
	}
}
