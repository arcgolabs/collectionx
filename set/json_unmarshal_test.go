package set

import (
	"encoding/json"
	"testing"
)

func TestSetJSONRoundTrip(t *testing.T) {
	source := NewSet(1, 2, 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal set: %v", err)
	}

	var target Set[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal set: %v", err)
	}

	if target.Len() != 3 || !target.Contains(1) || !target.Contains(2) || !target.Contains(3) {
		t.Fatalf("unexpected set contents: %#v", target.Values())
	}
}

func TestConcurrentSetJSONRoundTrip(t *testing.T) {
	source := NewConcurrentSet("a", "b", "c")

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent set: %v", err)
	}

	var target ConcurrentSet[string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent set: %v", err)
	}

	if target.Len() != 3 || !target.Contains("a") || !target.Contains("b") || !target.Contains("c") {
		t.Fatalf("unexpected concurrent set contents: %#v", target.Values())
	}
}

func TestOrderedSetJSONRoundTrip(t *testing.T) {
	source := NewOrderedSet(3, 1, 2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal ordered set: %v", err)
	}

	var target OrderedSet[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal ordered set: %v", err)
	}

	values := target.Values()
	if len(values) != 3 || values[0] != 3 || values[1] != 1 || values[2] != 2 {
		t.Fatalf("unexpected ordered set values: %#v", values)
	}
}

func TestMultiSetJSONRoundTrip(t *testing.T) {
	source := NewMultiSet[string]("a", "a", "b")

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal multiset: %v", err)
	}

	var target MultiSet[string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal multiset: %v", err)
	}

	if target.Len() != 3 || target.Count("a") != 2 || target.Count("b") != 1 {
		t.Fatalf("unexpected multiset counts: %#v", target.AllCounts())
	}
}

func TestMultiSetUnmarshalRejectsNegativeCount(t *testing.T) {
	var target MultiSet[string]
	if err := json.Unmarshal([]byte(`{"bad":-1}`), &target); err == nil {
		t.Fatal("expected negative count error")
	}
}
