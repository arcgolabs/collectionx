package list

import (
	"encoding/json"
	"testing"
)

func TestListJSONRoundTrip(t *testing.T) {
	source := NewList(1, 2, 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal list: %v", err)
	}

	var target List[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}

	if got := target.Values(); len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("unexpected list values: %#v", got)
	}
}

func TestConcurrentListJSONRoundTrip(t *testing.T) {
	source := NewConcurrentList("a", "b", "c")

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent list: %v", err)
	}

	var target ConcurrentList[string]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent list: %v", err)
	}

	if got := target.Values(); len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected concurrent list values: %#v", got)
	}
}

func TestGridJSONRoundTrip(t *testing.T) {
	source := NewGrid([]int{1, 2}, []int{3})

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal grid: %v", err)
	}

	var target Grid[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal grid: %v", err)
	}

	rows := target.Values()
	if len(rows) != 2 || len(rows[0]) != 2 || rows[0][0] != 1 || rows[0][1] != 2 || len(rows[1]) != 1 || rows[1][0] != 3 {
		t.Fatalf("unexpected grid values: %#v", rows)
	}
}

func TestConcurrentGridJSONRoundTrip(t *testing.T) {
	source := NewConcurrentGrid([]int{1, 2}, []int{3, 4})

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent grid: %v", err)
	}

	var target ConcurrentGrid[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent grid: %v", err)
	}

	rows := target.Values()
	if len(rows) != 2 || len(rows[0]) != 2 || rows[1][1] != 4 {
		t.Fatalf("unexpected concurrent grid values: %#v", rows)
	}
}

func TestDequeJSONRoundTrip(t *testing.T) {
	source := NewDeque(1, 2, 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal deque: %v", err)
	}

	var target Deque[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal deque: %v", err)
	}

	if got := target.Values(); len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("unexpected deque values: %#v", got)
	}
}

func TestConcurrentDequeJSONRoundTrip(t *testing.T) {
	source := NewConcurrentDeque(1, 2, 3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent deque: %v", err)
	}

	var target ConcurrentDeque[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent deque: %v", err)
	}

	if got := target.Values(); len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("unexpected concurrent deque values: %#v", got)
	}
}

func TestRingBufferJSONRoundTrip(t *testing.T) {
	source := NewRingBuffer[int](4)
	_ = source.Push(1)
	_ = source.Push(2)
	_ = source.Push(3)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal ring buffer: %v", err)
	}

	var target RingBuffer[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal ring buffer: %v", err)
	}

	if got := target.Values(); len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("unexpected ring buffer values: %#v", got)
	}
	if target.Capacity() != 3 {
		t.Fatalf("unexpected ring buffer capacity: %d", target.Capacity())
	}
}

func TestConcurrentRingBufferJSONRoundTrip(t *testing.T) {
	source := NewConcurrentRingBuffer[int](4)
	_ = source.Push(1)
	_ = source.Push(2)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal concurrent ring buffer: %v", err)
	}

	var target ConcurrentRingBuffer[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal concurrent ring buffer: %v", err)
	}

	if got := target.Values(); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("unexpected concurrent ring buffer values: %#v", got)
	}
	if target.Capacity() != 2 {
		t.Fatalf("unexpected concurrent ring buffer capacity: %d", target.Capacity())
	}
}

func TestRopeListJSONRoundTrip(t *testing.T) {
	source := NewRopeList(1, 2, 3, 4)

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal rope list: %v", err)
	}

	var target RopeList[int]
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatalf("unmarshal rope list: %v", err)
	}

	if got := target.Values(); len(got) != 4 || got[0] != 1 || got[1] != 2 || got[2] != 3 || got[3] != 4 {
		t.Fatalf("unexpected rope list values: %#v", got)
	}
}
