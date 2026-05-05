package disjointset

// DisjointSet stores partitioned items with union-find operations.
type DisjointSet[T comparable] struct {
	parent   map[T]T
	rank     map[T]int
	size     map[T]int
	setCount int
}

// New creates an empty disjoint set.
func New[T comparable]() *DisjointSet[T] {
	return &DisjointSet[T]{}
}

// Add inserts one or more singleton items.
func (d *DisjointSet[T]) Add(items ...T) {
	if d == nil || len(items) == 0 {
		return
	}
	d.ensureInit()
	for _, item := range items {
		if _, exists := d.parent[item]; exists {
			continue
		}
		d.parent[item] = item
		d.rank[item] = 0
		d.size[item] = 1
		d.setCount++
	}
}

// Has reports whether item exists.
func (d *DisjointSet[T]) Has(item T) bool {
	if d == nil || d.parent == nil {
		return false
	}
	_, ok := d.parent[item]
	return ok
}

// Find returns the representative item for the set containing item.
func (d *DisjointSet[T]) Find(item T) (T, bool) {
	var zero T
	if d == nil || d.parent == nil {
		return zero, false
	}
	root, ok := d.parent[item]
	if !ok {
		return zero, false
	}

	for {
		parent := d.parent[root]
		if parent == root {
			break
		}
		root = parent
	}

	for item != root {
		parent := d.parent[item]
		d.parent[item] = root
		item = parent
	}
	return root, true
}

// Union merges the sets containing left and right.
// Missing items are created as singleton sets.
// It returns true when two different sets were merged.
func (d *DisjointSet[T]) Union(left, right T) bool {
	if d == nil {
		return false
	}
	d.Add(left, right)

	leftRoot, _ := d.Find(left)
	rightRoot, _ := d.Find(right)
	if leftRoot == rightRoot {
		return false
	}

	leftRank := d.rank[leftRoot]
	rightRank := d.rank[rightRoot]
	if leftRank < rightRank {
		leftRoot, rightRoot = rightRoot, leftRoot
		leftRank, rightRank = rightRank, leftRank
	}

	d.parent[rightRoot] = leftRoot
	d.size[leftRoot] += d.size[rightRoot]
	delete(d.size, rightRoot)
	if leftRank == rightRank {
		d.rank[leftRoot]++
	}
	d.setCount--
	return true
}

// Connected reports whether left and right are in the same set.
func (d *DisjointSet[T]) Connected(left, right T) bool {
	leftRoot, ok := d.Find(left)
	if !ok {
		return false
	}
	rightRoot, ok := d.Find(right)
	if !ok {
		return false
	}
	return leftRoot == rightRoot
}

// SizeOf returns the number of items in the set containing item.
func (d *DisjointSet[T]) SizeOf(item T) int {
	root, ok := d.Find(item)
	if !ok {
		return 0
	}
	return d.size[root]
}

// Len returns the total item count.
func (d *DisjointSet[T]) Len() int {
	if d == nil {
		return 0
	}
	return len(d.parent)
}

// SetCount returns the number of disjoint sets.
func (d *DisjointSet[T]) SetCount() int {
	if d == nil {
		return 0
	}
	return d.setCount
}

// IsEmpty reports whether there are no items.
func (d *DisjointSet[T]) IsEmpty() bool {
	return d.Len() == 0
}

// Clear removes all items.
func (d *DisjointSet[T]) Clear() {
	if d == nil {
		return
	}
	d.parent = nil
	d.rank = nil
	d.size = nil
	d.setCount = 0
}

// Groups returns all current groups keyed by representative item.
func (d *DisjointSet[T]) Groups() map[T][]T {
	if d == nil || len(d.parent) == 0 {
		return map[T][]T{}
	}

	groups := make(map[T][]T, d.setCount)
	for item := range d.parent {
		root, _ := d.Find(item)
		members := groups[root]
		if members == nil {
			members = make([]T, 0, d.size[root])
		}
		groups[root] = append(members, item)
	}
	return groups
}

// MembersOf returns all members in the set containing item.
func (d *DisjointSet[T]) MembersOf(item T) []T {
	root, ok := d.Find(item)
	if !ok {
		return nil
	}
	members := make([]T, 0, d.size[root])
	for current := range d.parent {
		currentRoot, _ := d.Find(current)
		if currentRoot == root {
			members = append(members, current)
		}
	}
	return members
}

// RangeGroups iterates current groups until fn returns false.
func (d *DisjointSet[T]) RangeGroups(fn func(root T, members []T) bool) {
	if d == nil || len(d.parent) == 0 || fn == nil {
		return
	}
	for root, members := range d.Groups() {
		if !fn(root, members) {
			return
		}
	}
}

func (d *DisjointSet[T]) ensureInit() {
	if d.parent == nil {
		d.parent = make(map[T]T)
	}
	if d.rank == nil {
		d.rank = make(map[T]int)
	}
	if d.size == nil {
		d.size = make(map[T]int)
	}
}
