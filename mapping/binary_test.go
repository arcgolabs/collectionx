package mapping

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestOrderedMapBinaryRoundTripPreservesOrder(t *testing.T) {
	source := NewOrderedMap[string, int]()
	source.Set("third", 3)
	source.Set("first", 1)
	source.Set("second", 2)

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target OrderedMap[string, int]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	keys := target.Keys()
	if len(keys) != 3 || keys[0] != "third" || keys[1] != "first" || keys[2] != "second" {
		t.Fatalf("unexpected key order: %#v", keys)
	}
}

func TestShardedConcurrentMapGobRoundTrip(t *testing.T) {
	source := NewShardedConcurrentMap[string, int](8, HashString)
	source.Set("a", 1)
	source.Set("b", 2)

	data, err := source.GobEncode()
	if err != nil {
		t.Fatalf("gob encode bytes: %v", err)
	}

	target := NewShardedConcurrentMap[string, int](8, HashString)
	if err := target.GobDecode(data); err != nil {
		t.Fatalf("gob decode bytes: %v", err)
	}

	if target.Len() != 2 || target.GetOrDefault("a", 0) != 1 || target.GetOrDefault("b", 0) != 2 {
		t.Fatalf("unexpected sharded map contents: %#v", target.All())
	}
}

func TestConcurrentTableGobRoundTrip(t *testing.T) {
	source := NewConcurrentTable[string, string, int]()
	source.Put("r1", "c1", 1)
	source.Put("r2", "c2", 2)

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target ConcurrentTable[string, string, int]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	if !target.Has("r1", "c1") || !target.Has("r2", "c2") || target.Len() != 2 {
		t.Fatalf("unexpected table contents: %#v", target.All())
	}
}
