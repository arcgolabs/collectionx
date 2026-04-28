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
	for _, bit := range bits {
		b.Set(bit)
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
			current &^= uint64(1) << offset
		}
	}
}

// Union returns a new bitset containing bits from both sets.
func (b *BitSet) Union(other *BitSet) *BitSet {
	maxWords := max(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, maxWords)}
	for i := range maxWords {
		out.words[i] = wordAt(b, i) | wordAt(other, i)
	}
	out.recount()
	return out
}

// Intersect returns a new bitset containing shared bits.
func (b *BitSet) Intersect(other *BitSet) *BitSet {
	minWords := min(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, minWords)}
	for i := range minWords {
		out.words[i] = wordAt(b, i) & wordAt(other, i)
	}
	out.recount()
	return out
}

// Difference returns a new bitset with bits in b but not in other.
func (b *BitSet) Difference(other *BitSet) *BitSet {
	out := b.Clone()
	if other == nil || len(other.words) == 0 || len(out.words) == 0 {
		return out
	}
	for i := range len(out.words) {
		out.words[i] &^= wordAt(other, i)
	}
	out.recount()
	return out
}

// SymmetricDifference returns bits that exist in exactly one of the two sets.
func (b *BitSet) SymmetricDifference(other *BitSet) *BitSet {
	maxWords := max(lenWords(b), lenWords(other))
	out := &BitSet{words: make([]uint64, maxWords)}
	for i := range maxWords {
		out.words[i] = wordAt(b, i) ^ wordAt(other, i)
	}
	out.recount()
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
