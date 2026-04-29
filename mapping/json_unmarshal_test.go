package mapping

import (
	"encoding/json"
	"testing"
)

func TestMapJSONRoundTrip(t *testing.T) {
	source := NewMap[string, int]()
	source.Set("a", 1)
	source.Set("b", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal map: %v", err)
	}

	var target Map[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}

	if target.Len() != 2 || target.GetOrDefault("a", 0) != 1 || target.GetOrDefault("b", 0) != 2 {
		t.Fatalf("unexpected map contents: %#v", target.All())
	}
}

func TestConcurrentMapJSONRoundTrip(t *testing.T) {
	source := NewConcurrentMap[string, int]()
	source.Set("a", 1)
	source.Set("b", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent map: %v", err)
	}

	var target ConcurrentMap[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent map: %v", err)
	}

	if target.Len() != 2 || target.GetOrDefault("a", 0) != 1 || target.GetOrDefault("b", 0) != 2 {
		t.Fatalf("unexpected concurrent map contents: %#v", target.All())
	}
}

func TestShardedConcurrentMapJSONRoundTrip(t *testing.T) {
	source := NewShardedConcurrentMap[string, int](8, HashString)
	source.Set("a", 1)
	source.Set("b", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal sharded map: %v", err)
	}

	target := NewShardedConcurrentMap[string, int](8, HashString)
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal sharded map: %v", err)
	}

	if target.Len() != 2 || target.GetOrDefault("a", 0) != 1 || target.GetOrDefault("b", 0) != 2 {
		t.Fatalf("unexpected sharded map contents: %#v", target.All())
	}
}

func TestShardedConcurrentMapUnmarshalRequiresInitializedReceiver(t *testing.T) {
	var target ShardedConcurrentMap[string, int]
	if err := json.Unmarshal([]byte(`{"a":1}`), &target); err == nil {
		t.Fatal("expected uninitialized sharded map error")
	}
}

func TestBiMapJSONRoundTrip(t *testing.T) {
	source := NewBiMap[string, int]()
	source.Put("a", 1)
	source.Put("b", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal bimap: %v", err)
	}

	var target BiMap[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal bimap: %v", err)
	}

	if target.Len() != 2 {
		t.Fatalf("unexpected bimap length: %d", target.Len())
	}
	if value, ok := target.GetByKey("a"); !ok || value != 1 {
		t.Fatalf("unexpected bimap forward lookup")
	}
	if key, ok := target.GetByValue(2); !ok || key != "b" {
		t.Fatalf("unexpected bimap inverse lookup")
	}
}

func TestOrderedMapJSONRoundTripPreservesOrder(t *testing.T) {
	source := NewOrderedMap[string, int]()
	source.Set("third", 3)
	source.Set("first", 1)
	source.Set("second", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal ordered map: %v", err)
	}

	var target OrderedMap[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal ordered map: %v", err)
	}

	keys := target.Keys()
	if len(keys) != 3 || keys[0] != "third" || keys[1] != "first" || keys[2] != "second" {
		t.Fatalf("unexpected ordered map keys: %#v", keys)
	}
}

func TestMultiMapJSONRoundTrip(t *testing.T) {
	source := NewMultiMap[string, int]()
	source.PutAll("a", 1, 2)
	source.Put("b", 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal multimap: %v", err)
	}

	var target MultiMap[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal multimap: %v", err)
	}

	if target.Len() != 2 || target.ValueCount() != 3 {
		t.Fatalf("unexpected multimap size: keys=%d values=%d", target.Len(), target.ValueCount())
	}
	values := target.GetCopy("a")
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("unexpected multimap values: %#v", values)
	}
}

func TestConcurrentMultiMapJSONRoundTrip(t *testing.T) {
	source := NewConcurrentMultiMap[string, int]()
	source.PutAll("a", 1, 2)
	source.Put("b", 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent multimap: %v", err)
	}

	var target ConcurrentMultiMap[string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent multimap: %v", err)
	}

	if target.Len() != 2 || target.ValueCount() != 3 {
		t.Fatalf("unexpected concurrent multimap size: keys=%d values=%d", target.Len(), target.ValueCount())
	}
	values := target.GetCopy("a")
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("unexpected concurrent multimap values: %#v", values)
	}
}

func TestTableJSONRoundTrip(t *testing.T) {
	source := NewTable[string, string, int]()
	source.Put("r1", "c1", 1)
	source.Put("r1", "c2", 2)
	source.Put("r2", "c1", 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal table: %v", err)
	}

	var target Table[string, string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal table: %v", err)
	}

	if target.Len() != 3 || !target.Has("r1", "c2") || !target.Has("r2", "c1") {
		t.Fatalf("unexpected table contents: %#v", target.All())
	}
}

func TestConcurrentTableJSONRoundTrip(t *testing.T) {
	source := NewConcurrentTable[string, string, int]()
	source.Put("r1", "c1", 1)
	source.Put("r2", "c2", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent table: %v", err)
	}

	var target ConcurrentTable[string, string, int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent table: %v", err)
	}

	if target.Len() != 2 || !target.Has("r1", "c1") || !target.Has("r2", "c2") {
		t.Fatalf("unexpected concurrent table contents: %#v", target.All())
	}
}
