package set

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestOrderedSetBinaryRoundTrip(t *testing.T) {
	source := NewOrderedSet(3, 1, 2)

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target OrderedSet[int]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	values := target.Values()
	if len(values) != 3 || values[0] != 3 || values[1] != 1 || values[2] != 2 {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestSetGobRoundTrip(t *testing.T) {
	source := NewSet("a", "b")

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target Set[string]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	if !target.Contains("a") || !target.Contains("b") || target.Len() != 2 {
		t.Fatalf("unexpected contents: %#v", target.Values())
	}
}
