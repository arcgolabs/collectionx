package disjointset_test

import (
	"encoding/json"
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

func BenchmarkDisjointSetMembersOf(b *testing.B) {
	ds := buildBenchDisjointSet(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		members := ds.MembersOf(0)
		if len(members) != benchDisjointSetSize {
			b.Fatalf("unexpected member count: %d", len(members))
		}
	}
}

func BenchmarkDisjointSetRangeGroups(b *testing.B) {
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
		groups := 0
		ds.RangeGroups(func(_ int, members []int) bool {
			groups++
			return len(members) > 0
		})
		if groups == 0 {
			b.Fatal("RangeGroups() visited no groups")
		}
	}
}

func BenchmarkDisjointSetMarshalJSON(b *testing.B) {
	ds := buildBenchDisjointSet(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		data, err := json.Marshal(ds)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
		if len(data) == 0 {
			b.Fatal("json.Marshal() returned empty data")
		}
	}
}

func BenchmarkDisjointSetMarshalBinary(b *testing.B) {
	ds := buildBenchDisjointSet(b)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		data, err := ds.MarshalBinary()
		if err != nil {
			b.Fatalf("MarshalBinary() error = %v", err)
		}
		if len(data) == 0 {
			b.Fatal("MarshalBinary() returned empty data")
		}
	}
}
