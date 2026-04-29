package interval

import (
	"encoding/json"
	"testing"
)

func TestRangeSetJSONRoundTrip(t *testing.T) {
	source := NewRangeSet[int]()
	source.Add(1, 3)
	source.Add(5, 7)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal range set: %v", err)
	}

	var target RangeSet[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal range set: %v", err)
	}

	ranges := target.Ranges()
	if len(ranges) != 2 || ranges[0].Start != 1 || ranges[0].End != 3 || ranges[1].Start != 5 || ranges[1].End != 7 {
		t.Fatalf("unexpected range set values: %#v", ranges)
	}
}

func TestRangeMapJSONRoundTrip(t *testing.T) {
	source := NewRangeMap[int, string]()
	source.Put(1, 3, "a")
	source.Put(5, 7, "b")

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal range map: %v", err)
	}

	var target RangeMap[int, string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal range map: %v", err)
	}

	entries := target.Entries()
	if len(entries) != 2 || entries[0].Range.Start != 1 || entries[0].Range.End != 3 || entries[0].Value != "a" || entries[1].Range.Start != 5 || entries[1].Range.End != 7 || entries[1].Value != "b" {
		t.Fatalf("unexpected range map values: %#v", entries)
	}
}
