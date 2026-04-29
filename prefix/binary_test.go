package prefix

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestTrieBinaryRoundTrip(t *testing.T) {
	source := NewTrie[int]()
	source.Put("ab", 1)
	source.Put("abc", 2)

	data, err := source.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary: %v", err)
	}

	var target Trie[int]
	if err := target.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal binary: %v", err)
	}

	if value, ok := target.Get("abc"); !ok || value != 2 {
		t.Fatalf("unexpected trie value: %v %v", value, ok)
	}
}

func TestTrieGobRoundTrip(t *testing.T) {
	source := NewTrie[string]()
	source.Put("k", "v")

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(source); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var target Trie[string]
	if err := gob.NewDecoder(&buffer).Decode(&target); err != nil {
		t.Fatalf("gob decode: %v", err)
	}

	if value, ok := target.Get("k"); !ok || value != "v" {
		t.Fatalf("unexpected trie value: %v %v", value, ok)
	}
}
