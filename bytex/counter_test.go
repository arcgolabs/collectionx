package bytex_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/arcgolabs/collectionx/bytex"
)

func TestCounterBasicOps(t *testing.T) {
	t.Parallel()

	c := bytex.NewCounter('a', 'b', 'a')
	c.AddString("cc")
	c.AddN('c', 3)
	if c.Len() != 8 || c.UniqueLen() != 3 {
		t.Fatalf("unexpected lengths: total=%d unique=%d", c.Len(), c.UniqueLen())
	}
	if c.Count('a') != 2 || c.Count('c') != 5 || c.Contains('z') {
		t.Fatalf("unexpected counts: %v", c.Entries())
	}

	if removed := c.RemoveN('c', 2); removed != 2 || c.Count('c') != 3 {
		t.Fatalf("unexpected remove result: removed=%d count=%d", removed, c.Count('c'))
	}
	if !slices.Equal(c.Distinct(), []byte{'a', 'b', 'c'}) {
		t.Fatalf("unexpected distinct values: %v", c.Distinct())
	}

	first, ok := c.GetFirst()
	if !ok || first != 'a' {
		t.Fatalf("unexpected first value: %q %v", first, ok)
	}
	last, ok := c.GetLast()
	if !ok || last != 'c' {
		t.Fatalf("unexpected last value: %q %v", last, ok)
	}

	common := c.MostCommon(2)
	if len(common) != 2 || common[0].Value != 'c' || common[0].Count != 3 {
		t.Fatalf("unexpected most common entries: %v", common)
	}
	value, count, ok := c.MostCommonValue()
	if !ok || value != 'c' || count != 3 {
		t.Fatalf("unexpected most common value: %q %d %v", value, count, ok)
	}
	least := c.LeastCommon(2)
	if len(least) != 2 || least[0].Value != 'b' || least[0].Count != 1 {
		t.Fatalf("unexpected least common entries: %v", least)
	}
	value, count, ok = c.LeastCommonValue()
	if !ok || value != 'b' || count != 1 {
		t.Fatalf("unexpected least common value: %q %d %v", value, count, ok)
	}
}

func TestCounterMergeSubtractAndSerialization(t *testing.T) {
	t.Parallel()

	c := bytex.NewCounter('a', 'a')
	c.Merge(bytex.NewCounter('b', 'b', 'b'))
	c.Subtract(bytex.NewCounter('a', 'b'))
	if c.Count('a') != 1 || c.Count('b') != 2 {
		t.Fatalf("unexpected counts after merge/subtract: %v", c.Entries())
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	var restored bytex.Counter
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if restored.Count('a') != 1 || restored.Count('b') != 2 {
		t.Fatalf("unexpected JSON restored counts: %v", restored.Entries())
	}

	raw, err := c.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}
	var binaryRestored bytex.Counter
	if err := binaryRestored.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if binaryRestored.Count('a') != 1 || binaryRestored.Count('b') != 2 {
		t.Fatalf("unexpected binary restored counts: %v", binaryRestored.Entries())
	}

	if err := restored.UnmarshalJSON([]byte(`[{"value":1,"count":-1}]`)); err == nil {
		t.Fatal("expected negative count JSON error")
	}
	if err := restored.UnmarshalBinary([]byte{1, 2}); err == nil {
		t.Fatal("expected invalid binary length error")
	}
}
