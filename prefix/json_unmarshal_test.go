package prefix

import (
	"encoding/json"
	"testing"
)

func TestTrieJSONRoundTrip(t *testing.T) {
	source := NewTrie[int]()
	source.Put("ab", 1)
	source.Put("abc", 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal trie: %v", err)
	}

	var target Trie[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal trie: %v", err)
	}

	if target.Len() != 2 {
		t.Fatalf("unexpected trie size: %d", target.Len())
	}
	if value, ok := target.Get("ab"); !ok || value != 1 {
		t.Fatalf("unexpected trie value for ab: %v %v", value, ok)
	}
	if value, ok := target.Get("abc"); !ok || value != 2 {
		t.Fatalf("unexpected trie value for abc: %v %v", value, ok)
	}
}
