package bytex

import "sort"

// CounterEntry is one byte frequency entry.
type CounterEntry struct {
	Value byte `json:"value"`
	Count int  `json:"count"`
}

// Counter counts byte frequencies.
// Zero value is ready to use.
type Counter struct {
	counts [256]int
	total  int
	unique int
}

const counterBulkThreshold = 64
const counterTopSortThreshold = 32

// NewCounter creates a counter and counts optional values.
func NewCounter(values ...byte) *Counter {
	c := &Counter{}
	c.Add(values...)
	return c
}

// Add counts one or more values.
func (c *Counter) Add(values ...byte) {
	if c == nil || len(values) == 0 {
		return
	}
	if len(values) >= counterBulkThreshold {
		c.addBulkBytes(values)
		return
	}
	for _, value := range values {
		if c.counts[value] == 0 {
			c.unique++
		}
		c.counts[value]++
		c.total++
	}
}

// AddString counts all bytes in value without allocating.
func (c *Counter) AddString(value string) {
	if c == nil || value == "" {
		return
	}
	if len(value) >= counterBulkThreshold {
		c.addBulkString(value)
		return
	}
	for i := range len(value) {
		current := value[i]
		if c.counts[current] == 0 {
			c.unique++
		}
		c.counts[current]++
		c.total++
	}
}

// AddN increases value count by n.
func (c *Counter) AddN(value byte, n int) {
	if c == nil || n <= 0 {
		return
	}
	if c.counts[value] == 0 {
		c.unique++
	}
	c.counts[value] += n
	c.total += n
}

// Remove decreases value count by one and reports whether a value was removed.
func (c *Counter) Remove(value byte) bool {
	return c.RemoveN(value, 1) == 1
}

// RemoveN decreases value count by up to n and returns the removed count.
func (c *Counter) RemoveN(value byte, n int) int {
	if c == nil || n <= 0 || c.counts[value] == 0 {
		return 0
	}
	removed := min(n, c.counts[value])
	c.counts[value] -= removed
	c.total -= removed
	if c.counts[value] == 0 {
		c.unique--
	}
	return removed
}

// Count returns the count for value.
func (c *Counter) Count(value byte) int {
	if c == nil {
		return 0
	}
	return c.counts[value]
}

// Contains reports whether value has a non-zero count.
func (c *Counter) Contains(value byte) bool {
	return c.Count(value) > 0
}

// Len returns the total count across all values.
func (c *Counter) Len() int {
	if c == nil {
		return 0
	}
	return c.total
}

// UniqueLen returns the number of distinct values.
func (c *Counter) UniqueLen() int {
	if c == nil {
		return 0
	}
	return c.unique
}

// IsEmpty reports whether there are no counted values.
func (c *Counter) IsEmpty() bool {
	return c.Len() == 0
}

// Clear removes all counts.
func (c *Counter) Clear() {
	if c == nil {
		return
	}
	c.counts = [256]int{}
	c.total = 0
	c.unique = 0
}

// Clone returns a copy of the counter.
func (c *Counter) Clone() *Counter {
	if c == nil {
		return &Counter{}
	}
	return &Counter{
		counts: c.counts,
		total:  c.total,
		unique: c.unique,
	}
}

// Merge adds counts from other and returns the receiver.
func (c *Counter) Merge(other *Counter) *Counter {
	if c == nil {
		return nil
	}
	if other == nil || other.unique == 0 {
		return c
	}
	other.Range(func(value byte, count int) bool {
		c.AddN(value, count)
		return true
	})
	return c
}

// Subtract removes counts from other and returns the receiver.
func (c *Counter) Subtract(other *Counter) *Counter {
	if c == nil {
		return nil
	}
	if other == nil || other.unique == 0 {
		return c
	}
	other.Range(func(value byte, count int) bool {
		c.RemoveN(value, count)
		return true
	})
	return c
}

// Distinct returns values with non-zero counts in ascending order.
func (c *Counter) Distinct() []byte {
	if c == nil || c.unique == 0 {
		return nil
	}
	out := make([]byte, 0, c.unique)
	c.Range(func(value byte, _ int) bool {
		out = append(out, value)
		return true
	})
	return out
}

// Values returns distinct values with non-zero counts in ascending order.
func (c *Counter) Values() []byte {
	return c.Distinct()
}

// Entries returns non-zero count entries in ascending byte order.
func (c *Counter) Entries() []CounterEntry {
	if c == nil || c.unique == 0 {
		return nil
	}
	out := make([]CounterEntry, 0, c.unique)
	c.Range(func(value byte, count int) bool {
		out = append(out, CounterEntry{Value: value, Count: count})
		return true
	})
	return out
}

// Snapshot returns non-zero count entries in ascending byte order.
func (c *Counter) Snapshot() []CounterEntry {
	return c.Entries()
}

// MostCommon returns up to n entries ordered by count descending, then value ascending.
func (c *Counter) MostCommon(n int) []CounterEntry {
	if c == nil || c.unique == 0 || n <= 0 {
		return nil
	}
	return c.topMostCommon(n)
}

// MostCommonValue returns the most frequent value, breaking ties by smaller value.
func (c *Counter) MostCommonValue() (byte, int, bool) {
	if c == nil || c.unique == 0 {
		return 0, 0, false
	}
	bestValue := byte(0)
	bestCount := 0
	found := false
	for value, count := range c.counts {
		if count > 0 && (!found || count > bestCount) {
			bestValue = byte(value)
			bestCount = count
			found = true
		}
	}
	return bestValue, bestCount, true
}

// LeastCommon returns up to n entries ordered by count ascending, then value ascending.
func (c *Counter) LeastCommon(n int) []CounterEntry {
	if c == nil || c.unique == 0 || n <= 0 {
		return nil
	}
	return c.topLeastCommon(n)
}

// LeastCommonValue returns the least frequent value, breaking ties by smaller value.
func (c *Counter) LeastCommonValue() (byte, int, bool) {
	if c == nil || c.unique == 0 {
		return 0, 0, false
	}
	bestValue := byte(0)
	bestCount := 0
	found := false
	for value, count := range c.counts {
		if count > 0 && (!found || count < bestCount) {
			bestValue = byte(value)
			bestCount = count
			found = true
		}
	}
	return bestValue, bestCount, true
}

// GetFirst returns the smallest value with a non-zero count.
func (c *Counter) GetFirst() (byte, bool) {
	if c == nil || c.unique == 0 {
		return 0, false
	}
	for value, count := range c.counts {
		if count > 0 {
			return byte(value), true
		}
	}
	return 0, false
}

// GetLast returns the largest value with a non-zero count.
func (c *Counter) GetLast() (byte, bool) {
	if c == nil || c.unique == 0 {
		return 0, false
	}
	for value := len(c.counts); value > 0; value-- {
		if c.counts[value-1] > 0 {
			return byte(value - 1), true
		}
	}
	return 0, false
}

// Range iterates non-zero count entries in ascending byte order until fn returns false.
func (c *Counter) Range(fn func(value byte, count int) bool) {
	if c == nil || fn == nil || c.unique == 0 {
		return
	}
	for value, count := range c.counts {
		if count > 0 && !fn(byte(value), count) {
			return
		}
	}
}

func (c *Counter) addBulkBytes(values []byte) {
	var counts [256]int
	for _, value := range values {
		counts[value]++
	}
	c.mergeCounts(&counts, len(values))
}

func (c *Counter) addBulkString(value string) {
	var counts [256]int
	for i := range len(value) {
		counts[value[i]]++
	}
	c.mergeCounts(&counts, len(value))
}

func (c *Counter) mergeCounts(counts *[256]int, total int) {
	for value, count := range counts {
		if count == 0 {
			continue
		}
		if c.counts[value] == 0 {
			c.unique++
		}
		c.counts[value] += count
	}
	c.total += total
}

func (c *Counter) topMostCommon(n int) []CounterEntry {
	limit := min(n, c.unique)
	if limit > counterTopSortThreshold {
		entries := c.Entries()
		sort.Slice(entries, func(i, j int) bool {
			return moreCommonEntry(entries[i], entries[j])
		})
		if limit < len(entries) {
			entries = entries[:limit]
		}
		return entries
	}

	out := make([]CounterEntry, 0, limit)
	for value, count := range c.counts {
		if count <= 0 {
			continue
		}
		candidate := CounterEntry{Value: byte(value), Count: count}
		if len(out) < limit {
			out = append(out, candidate)
			moveMoreCommonTowardFront(out, len(out)-1)
			continue
		}
		if moreCommonEntry(candidate, out[limit-1]) {
			out[limit-1] = candidate
			moveMoreCommonTowardFront(out, limit-1)
		}
	}
	return out
}

func (c *Counter) topLeastCommon(n int) []CounterEntry {
	limit := min(n, c.unique)
	if limit > counterTopSortThreshold {
		entries := c.Entries()
		sort.Slice(entries, func(i, j int) bool {
			return lessCommonEntry(entries[i], entries[j])
		})
		if limit < len(entries) {
			entries = entries[:limit]
		}
		return entries
	}

	out := make([]CounterEntry, 0, limit)
	for value, count := range c.counts {
		if count <= 0 {
			continue
		}
		candidate := CounterEntry{Value: byte(value), Count: count}
		if len(out) < limit {
			out = append(out, candidate)
			moveLessCommonTowardFront(out, len(out)-1)
			continue
		}
		if lessCommonEntry(candidate, out[limit-1]) {
			out[limit-1] = candidate
			moveLessCommonTowardFront(out, limit-1)
		}
	}
	return out
}

func moveMoreCommonTowardFront(entries []CounterEntry, index int) {
	for index > 0 && moreCommonEntry(entries[index], entries[index-1]) {
		entries[index], entries[index-1] = entries[index-1], entries[index]
		index--
	}
}

func moveLessCommonTowardFront(entries []CounterEntry, index int) {
	for index > 0 && lessCommonEntry(entries[index], entries[index-1]) {
		entries[index], entries[index-1] = entries[index-1], entries[index]
		index--
	}
}

func moreCommonEntry(candidate, existing CounterEntry) bool {
	if candidate.Count == existing.Count {
		return candidate.Value < existing.Value
	}
	return candidate.Count > existing.Count
}

func lessCommonEntry(candidate, existing CounterEntry) bool {
	if candidate.Count == existing.Count {
		return candidate.Value < existing.Value
	}
	return candidate.Count < existing.Count
}
