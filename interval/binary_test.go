package interval

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestRangeMapBinaryRoundTrip(t *testing.T) {
	source := NewRangeMap[int, string]()
	source.Put(1, 3, "a")
	source.Put(5, 7, "b")

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target RangeMap[int, string]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	if value, ok := target.Get(6); !ok || value != "b" {
		t.Fatalf("unexpected range map value: %v %v", value, ok)
	}
}

func TestRangeSetGobRoundTrip(t *testing.T) {
	source := NewRangeSet[int]()
	source.Add(1, 3)
	source.Add(5, 7)

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target RangeSet[int]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	if !target.Contains(2) || !target.Contains(6) {
		t.Fatalf("unexpected range set contents: %#v", target.Ranges())
	}
}
