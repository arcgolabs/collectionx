package list

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestListBinaryRoundTrip(t *testing.T) {
	source := NewList(1, 2, 3)

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target List[int]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	values := target.Values()
	if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestListGobRoundTrip(t *testing.T) {
	source := NewList("a", "b")

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target List[string]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	values := target.Values()
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Fatalf("unexpected values: %#v", values)
	}
}
