package bytex_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/arcgolabs/collectionx/bytex"
)

func TestSetBasicAndSetOps(t *testing.T) {
	t.Parallel()

	s := bytex.NewSet(1, 3, 255, 3)
	if s.Len() != 3 || !s.Contains(255) || s.Contains(2) {
		t.Fatalf("unexpected set state: len=%d values=%v", s.Len(), s.Values())
	}
	first, ok := s.GetFirst()
	if !ok || first != 1 {
		t.Fatalf("unexpected first value: %d %v", first, ok)
	}
	last, ok := s.GetLast()
	if !ok || last != 255 {
		t.Fatalf("unexpected last value: %d %v", last, ok)
	}

	if !s.Remove(3) || s.Remove(3) {
		t.Fatal("unexpected remove result")
	}
	if !slices.Equal(s.Values(), []byte{1, 255}) {
		t.Fatalf("unexpected values: %v", s.Values())
	}

	other := bytex.NewSet(1, 2)
	if !s.Union(other).IsSupersetOf(bytex.NewSet(1, 2, 255)) {
		t.Fatal("unexpected union result")
	}
	if !s.Intersect(other).IsSubsetOf(bytex.NewSet(1)) {
		t.Fatal("unexpected intersect result")
	}
	if !s.Difference(other).IsSubsetOf(bytex.NewSet(255)) {
		t.Fatal("unexpected difference result")
	}
	if !s.SymmetricDifference(other).IsSupersetOf(bytex.NewSet(2, 255)) {
		t.Fatal("unexpected symmetric difference result")
	}
	if !s.Intersects(other) || s.Intersects(bytex.NewSet(7)) {
		t.Fatal("unexpected intersects result")
	}
}

func TestSetRangeAndComplement(t *testing.T) {
	t.Parallel()

	s := bytex.NewSet()
	if added := s.AddRange(62, 67); added != 5 {
		t.Fatalf("unexpected added count: %d", added)
	}
	if !slices.Equal(s.Values(), []byte{62, 63, 64, 65, 66}) {
		t.Fatalf("unexpected range values: %v", s.Values())
	}
	if added := s.AddRange(-2, 2); added != 2 {
		t.Fatalf("unexpected clamped added count: %d", added)
	}
	if removed := s.RemoveRange(1, 65); removed != 4 {
		t.Fatalf("unexpected removed count: %d values=%v", removed, s.Values())
	}
	if !slices.Equal(s.Values(), []byte{0, 65, 66}) {
		t.Fatalf("unexpected values after remove range: %v", s.Values())
	}

	complement := s.Complement()
	if complement.Len() != 253 || complement.Contains(0) || !complement.Contains(1) || complement.Contains(66) {
		t.Fatalf("unexpected complement: len=%d", complement.Len())
	}
	if all := bytex.NewSet().Complement(); all.Len() != 256 {
		t.Fatalf("empty set complement should contain all bytes, got %d", all.Len())
	}
}

func TestSetSerialization(t *testing.T) {
	t.Parallel()

	s := bytex.NewSet(1, 2, 255)
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(data) != `[1,2,255]` {
		t.Fatalf("unexpected JSON payload: %s", data)
	}

	var restored bytex.Set
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if !slices.Equal(restored.Values(), s.Values()) {
		t.Fatalf("unexpected restored values: %v", restored.Values())
	}

	raw, err := s.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}
	var binaryRestored bytex.Set
	if err := binaryRestored.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if !slices.Equal(binaryRestored.Values(), s.Values()) {
		t.Fatal("binary roundtrip mismatch")
	}

	if err := restored.UnmarshalJSON([]byte(`[256]`)); err == nil {
		t.Fatal("expected out-of-range JSON error")
	}
	if err := restored.UnmarshalBinary([]byte{1, 2}); err == nil {
		t.Fatal("expected invalid binary length error")
	}
}
