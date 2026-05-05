package bitset

import (
	"math/bits"
	"slices"
)

// BitSet is a compact set for non-negative integer values.
// Zero value is ready to use.
type BitSet struct {
	words []uint64
	count int
}

// New creates a bitset and fills it with optional bits.
func New(bits ...int) *BitSet {
	b := &BitSet{}
	b.Add(bits...)
	return b
}

// Set enables one bit and reports whether it was newly added.
func (b *BitSet) Set(bit int) bool {
	if b == nil || bit < 0 {
		return false
	}
	wordIndex := bit / 64
	mask := uint64(1) << (bit % 64)
	b.ensureWord(wordIndex)
	if b.words[wordIndex]&mask != 0 {
		return false
	}
	b.words[wordIndex] |= mask
	b.count++
	return true
}

// Add enables one or more bits.
func (b *BitSet) Add(bits ...int) {
	if b == nil || len(bits) == 0 {
		return
	}
	maxWordIndex := -1
	for _, bit := range bits {
		if bit < 0 {
			continue
		}
		wordIndex := bit / 64
		if wordIndex > maxWordIndex {
			maxWordIndex = wordIndex
		}
	}
	if maxWordIndex >= 0 {
		b.ensureWord(maxWordIndex)
	}
	for _, bit := range bits {
		b.Set(bit)
	}
}

// AddRange enables all bits in [start, end).
func (b *BitSet) AddRange(start, end int) {
	if b == nil || start < 0 || start >= end {
		return
	}
	startWord := start / 64
	endWord := (end - 1) / 64
	b.ensureWord(endWord)
	for wordIndex := startWord; wordIndex <= endWord; wordIndex++ {
		mask := rangeWordMask(start, end, wordIndex)
		before := b.words[wordIndex]
		after := before | mask
		b.words[wordIndex] = after
		b.count += bits.OnesCount64(after) - bits.OnesCount64(before)
	}
}

// Remove clears one bit and reports whether it existed.
func (b *BitSet) Remove(bit int) bool {
	if b == nil || bit < 0 {
		return false
	}
	wordIndex := bit / 64
	if wordIndex >= len(b.words) {
		return false
	}
	mask := uint64(1) << (bit % 64)
	if b.words[wordIndex]&mask == 0 {
		return false
	}
	b.words[wordIndex] &^= mask
	b.count--
	b.trimTrailingZeros()
	return true
}

// RemoveRange clears all bits in [start, end).
func (b *BitSet) RemoveRange(start, end int) int {
	if b == nil || start < 0 || start >= end || len(b.words) == 0 {
		return 0
	}
	maxBit := len(b.words) * 64
	if start >= maxBit {
		return 0
	}
	if end > maxBit {
		end = maxBit
	}
	startWord := start / 64
	endWord := (end - 1) / 64
	removed := 0
	for wordIndex := startWord; wordIndex <= endWord; wordIndex++ {
		mask := rangeWordMask(start, end, wordIndex)
		before := b.words[wordIndex]
		after := before &^ mask
		b.words[wordIndex] = after
		removed += bits.OnesCount64(before) - bits.OnesCount64(after)
	}
	b.count -= removed
	b.trimTrailingZeros()
	return removed
}

// Contains reports whether bit exists.
func (b *BitSet) Contains(bit int) bool {
	if b == nil || bit < 0 {
		return false
	}
	wordIndex := bit / 64
	if wordIndex >= len(b.words) {
		return false
	}
	mask := uint64(1) << (bit % 64)
	return b.words[wordIndex]&mask != 0
}

// Len returns total set bit count.
func (b *BitSet) Len() int {
	if b == nil {
		return 0
	}
	return b.count
}

// IsEmpty reports whether there are no set bits.
func (b *BitSet) IsEmpty() bool {
	return b.Len() == 0
}

// Clear removes all bits.
func (b *BitSet) Clear() {
	if b == nil {
		return
	}
	b.words = nil
	b.count = 0
}

// Clone returns a copy of the bitset.
func (b *BitSet) Clone() *BitSet {
	if b == nil || len(b.words) == 0 {
		return &BitSet{}
	}
	return &BitSet{
		words: slices.Clone(b.words),
		count: b.count,
	}
}

// Values returns all set bits in ascending order.
func (b *BitSet) Values() []int {
	if b == nil || b.count == 0 {
		return nil
	}
	out := make([]int, 0, b.count)
	b.Range(func(bit int) bool {
		out = append(out, bit)
		return true
	})
	return out
}

// Range iterates set bits in ascending order until fn returns false.
func (b *BitSet) Range(fn func(bit int) bool) {
	if b == nil || fn == nil || b.count == 0 {
		return
	}
	for wordIndex, word := range b.words {
		current := word
		for current != 0 {
			offset := bits.TrailingZeros64(current)
			if !fn(wordIndex*64 + offset) {
				return
			}
			current &= current - 1
		}
	}
}

// NextSet returns the next set bit at or after bit.
func (b *BitSet) NextSet(bit int) (int, bool) {
	if b == nil || bit < 0 || len(b.words) == 0 {
		return 0, false
	}
	wordIndex := bit / 64
	if wordIndex >= len(b.words) {
		return 0, false
	}
	current := b.words[wordIndex] &^ ((uint64(1) << (bit % 64)) - 1)
	if current != 0 {
		return wordIndex*64 + bits.TrailingZeros64(current), true
	}
	for wordIndex++; wordIndex < len(b.words); wordIndex++ {
		if b.words[wordIndex] != 0 {
			return wordIndex*64 + bits.TrailingZeros64(b.words[wordIndex]), true
		}
	}
	return 0, false
}

// Union returns a new bitset containing bits from both sets.
func (b *BitSet) Union(other *BitSet) *BitSet {
	maxWords := max(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, maxWords)}
	lastNonZero := 0
	for i := range maxWords {
		word := wordAt(b, i) | wordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
		if word != 0 {
			lastNonZero = i + 1
		}
	}
	out.words = out.words[:lastNonZero]
	return out
}

// Intersect returns a new bitset containing shared bits.
func (b *BitSet) Intersect(other *BitSet) *BitSet {
	minWords := min(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, minWords)}
	lastNonZero := 0
	for i := range minWords {
		word := wordAt(b, i) & wordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
		if word != 0 {
			lastNonZero = i + 1
		}
	}
	out.words = out.words[:lastNonZero]
	return out
}

// Difference returns a new bitset with bits in b but not in other.
func (b *BitSet) Difference(other *BitSet) *BitSet {
	if b == nil || len(b.words) == 0 {
		return &BitSet{}
	}
	out := &BitSet{words: make([]uint64, len(b.words))}
	lastNonZero := 0
	for i := range len(b.words) {
		word := b.words[i] &^ wordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
		if word != 0 {
			lastNonZero = i + 1
		}
	}
	out.words = out.words[:lastNonZero]
	return out
}

// SymmetricDifference returns bits that exist in exactly one of the two sets.
func (b *BitSet) SymmetricDifference(other *BitSet) *BitSet {
	maxWords := max(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, maxWords)}
	lastNonZero := 0
	for i := range maxWords {
		word := wordAt(b, i) ^ wordAt(other, i)
		out.words[i] = word
		out.count += bits.OnesCount64(word)
		if word != 0 {
			lastNonZero = i + 1
		}
	}
	out.words = out.words[:lastNonZero]
	return out
}

func (b *BitSet) ensureWord(index int) {
	if index < len(b.words) {
		return
	}
	b.words = append(b.words, make([]uint64, index-len(b.words)+1)...)
}

func (b *BitSet) trimTrailingZeros() {
	for len(b.words) > 0 && b.words[len(b.words)-1] == 0 {
		b.words = b.words[:len(b.words)-1]
	}
}

func (b *BitSet) recount() {
	total := 0
	for _, word := range b.words {
		total += bits.OnesCount64(word)
	}
	b.count = total
	b.trimTrailingZeros()
}

func lenWords(b *BitSet) int {
	if b == nil {
		return 0
	}
	return len(b.words)
}

func wordAt(b *BitSet, index int) uint64 {
	if b == nil || index >= len(b.words) {
		return 0
	}
	return b.words[index]
}

func rangeWordMask(start, end, wordIndex int) uint64 {
	wordStart := wordIndex * 64
	wordEnd := wordStart + 64
	if start > wordStart {
		wordStart = start
	}
	if end < wordEnd {
		wordEnd = end
	}
	startOffset := wordStart - wordIndex*64
	endOffset := wordEnd - wordIndex*64
	lower := ^uint64(0) << startOffset
	if endOffset == 64 {
		return lower
	}
	upper := uint64(1)<<endOffset - 1
	return lower & upper
}
