package bytex

import (
	"math/bits"
)

const byteSetWordCount = 4

// Set is a compact set for byte values.
// Zero value is ready to use.
type Set struct {
	words [byteSetWordCount]uint64
	count int
}

// NewSet creates a byte set and fills it with optional values.
func NewSet(values ...byte) *Set {
	s := &Set{}
	s.Add(values...)
	return s
}

// Add inserts one or more byte values.
func (s *Set) Add(values ...byte) {
	if s == nil || len(values) == 0 {
		return
	}
	for _, value := range values {
		wordIndex, mask := byteSetMask(value)
		before := s.words[wordIndex]
		if before&mask == 0 {
			s.words[wordIndex] = before | mask
			s.count++
		}
	}
}

// AddRange inserts all values in [start, end) after clamping to the byte domain.
func (s *Set) AddRange(start, end int) int {
	if s == nil {
		return 0
	}
	start, end, ok := normalizeByteRange(start, end)
	if !ok {
		return 0
	}
	startWord := start / 64
	endWord := (end - 1) / 64
	startOffset := start % 64
	endOffset := (end-1)%64 + 1

	if startWord == endWord {
		return s.addWordMask(startWord, wordMaskBetweenOffsets(startOffset, endOffset))
	}

	added := s.addWordMask(startWord, ^uint64(0)<<startOffset)
	for wordIndex := startWord + 1; wordIndex < endWord; wordIndex++ {
		added += s.addWordMask(wordIndex, ^uint64(0))
	}
	added += s.addWordMask(endWord, wordMaskBetweenOffsets(0, endOffset))
	return added
}

// Set inserts value and reports whether it was newly added.
func (s *Set) Set(value byte) bool {
	if s == nil {
		return false
	}
	wordIndex, mask := byteSetMask(value)
	if s.words[wordIndex]&mask != 0 {
		return false
	}
	s.words[wordIndex] |= mask
	s.count++
	return true
}

// Remove deletes value and reports whether it existed.
func (s *Set) Remove(value byte) bool {
	if s == nil {
		return false
	}
	wordIndex, mask := byteSetMask(value)
	if s.words[wordIndex]&mask == 0 {
		return false
	}
	s.words[wordIndex] &^= mask
	s.count--
	return true
}

// RemoveRange deletes all values in [start, end) after clamping to the byte domain.
func (s *Set) RemoveRange(start, end int) int {
	if s == nil || s.count == 0 {
		return 0
	}
	start, end, ok := normalizeByteRange(start, end)
	if !ok {
		return 0
	}
	startWord := start / 64
	endWord := (end - 1) / 64
	startOffset := start % 64
	endOffset := (end-1)%64 + 1

	if startWord == endWord {
		return s.removeWordMask(startWord, wordMaskBetweenOffsets(startOffset, endOffset))
	}

	removed := s.removeWordMask(startWord, ^uint64(0)<<startOffset)
	for wordIndex := startWord + 1; wordIndex < endWord; wordIndex++ {
		removed += s.removeWordMask(wordIndex, ^uint64(0))
	}
	removed += s.removeWordMask(endWord, wordMaskBetweenOffsets(0, endOffset))
	return removed
}

// Contains reports whether value exists.
func (s *Set) Contains(value byte) bool {
	if s == nil {
		return false
	}
	wordIndex, mask := byteSetMask(value)
	return s.words[wordIndex]&mask != 0
}

// Len returns the number of distinct byte values.
func (s *Set) Len() int {
	if s == nil {
		return 0
	}
	return s.count
}

// IsEmpty reports whether the set has no values.
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

// Clear removes all values.
func (s *Set) Clear() {
	if s == nil {
		return
	}
	s.words = [byteSetWordCount]uint64{}
	s.count = 0
}

// Clone returns a copy of the set.
func (s *Set) Clone() *Set {
	if s == nil {
		return &Set{}
	}
	return &Set{
		words: s.words,
		count: s.count,
	}
}

// Values returns all byte values in ascending order.
func (s *Set) Values() []byte {
	if s == nil || s.count == 0 {
		return nil
	}
	out := make([]byte, 0, s.count)
	s.Range(func(value byte) bool {
		out = append(out, value)
		return true
	})
	return out
}

// Snapshot returns all byte values in ascending order.
func (s *Set) Snapshot() []byte {
	return s.Values()
}

// GetFirst returns the smallest value.
func (s *Set) GetFirst() (byte, bool) {
	if s == nil || s.count == 0 {
		return 0, false
	}
	for wordIndex, word := range s.words {
		if word != 0 {
			return byte(wordIndex*64 + bits.TrailingZeros64(word)), true
		}
	}
	return 0, false
}

// GetLast returns the largest value.
func (s *Set) GetLast() (byte, bool) {
	if s == nil || s.count == 0 {
		return 0, false
	}
	for i := len(s.words); i > 0; i-- {
		wordIndex := i - 1
		word := s.words[wordIndex]
		if word != 0 {
			return byte(wordIndex*64 + bits.Len64(word) - 1), true
		}
	}
	return 0, false
}

// Range iterates values in ascending order until fn returns false.
func (s *Set) Range(fn func(value byte) bool) {
	if s == nil || fn == nil || s.count == 0 {
		return
	}
	for wordIndex, word := range s.words {
		current := word
		for current != 0 {
			offset := bits.TrailingZeros64(current)
			if !fn(byte(wordIndex*64 + offset)) {
				return
			}
			current &= current - 1
		}
	}
}

// Union returns a new set containing values from both sets.
func (s *Set) Union(other *Set) *Set {
	out := &Set{}
	for i := range out.words {
		word := byteSetWordAt(s, i) | byteSetWordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
	}
	return out
}

// Intersect returns a new set containing shared values.
func (s *Set) Intersect(other *Set) *Set {
	out := &Set{}
	for i := range out.words {
		word := byteSetWordAt(s, i) & byteSetWordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
	}
	return out
}

// Difference returns a new set with values in s but not in other.
func (s *Set) Difference(other *Set) *Set {
	out := &Set{}
	for i := range out.words {
		word := byteSetWordAt(s, i) &^ byteSetWordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
	}
	return out
}

// SymmetricDifference returns values that exist in exactly one set.
func (s *Set) SymmetricDifference(other *Set) *Set {
	out := &Set{}
	for i := range out.words {
		word := byteSetWordAt(s, i) ^ byteSetWordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
	}
	return out
}

// Complement returns a new set containing every byte value not in s.
func (s *Set) Complement() *Set {
	out := &Set{}
	for i := range out.words {
		out.words[i] = ^byteSetWordAt(s, i)
		out.count += bits.OnesCount64(out.words[i])
	}
	return out
}

// Intersects reports whether the two sets share at least one value.
func (s *Set) Intersects(other *Set) bool {
	for i := range byteSetWordCount {
		if byteSetWordAt(s, i)&byteSetWordAt(other, i) != 0 {
			return true
		}
	}
	return false
}

// IsSubsetOf reports whether every value in s exists in other.
func (s *Set) IsSubsetOf(other *Set) bool {
	if s == nil || s.count == 0 {
		return true
	}
	for i, word := range s.words {
		if word&^byteSetWordAt(other, i) != 0 {
			return false
		}
	}
	return true
}

// IsSupersetOf reports whether every value in other exists in s.
func (s *Set) IsSupersetOf(other *Set) bool {
	return other.IsSubsetOf(s)
}

func (s *Set) recount() {
	if s == nil {
		return
	}
	s.count = 0
	for _, word := range s.words {
		s.count += bits.OnesCount64(word)
	}
}

func (s *Set) addWordMask(wordIndex int, mask uint64) int {
	before := s.words[wordIndex]
	after := before | mask
	s.words[wordIndex] = after
	added := bits.OnesCount64(after) - bits.OnesCount64(before)
	s.count += added
	return added
}

func (s *Set) removeWordMask(wordIndex int, mask uint64) int {
	before := s.words[wordIndex]
	after := before &^ mask
	s.words[wordIndex] = after
	removed := bits.OnesCount64(before) - bits.OnesCount64(after)
	s.count -= removed
	return removed
}

func normalizeByteRange(start, end int) (int, int, bool) {
	if start < 0 {
		start = 0
	}
	if end > 256 {
		end = 256
	}
	if start >= end || start >= 256 || end <= 0 {
		return 0, 0, false
	}
	return start, end, true
}

func wordMaskBetweenOffsets(startOffset, endOffset int) uint64 {
	lower := ^uint64(0) << startOffset
	if endOffset == 64 {
		return lower
	}
	upper := uint64(1)<<endOffset - 1
	return lower & upper
}

func byteSetMask(value byte) (int, uint64) {
	return int(value >> 6), uint64(1) << (value & 63)
}

func byteSetWordAt(s *Set, index int) uint64 {
	if s == nil {
		return 0
	}
	return s.words[index]
}
