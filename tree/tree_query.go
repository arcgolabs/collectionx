package tree

import collectionlist "github.com/arcgolabs/collectionx/list"

// Get returns node by id.
func (t *Tree[K, V]) Get(id K) (*Node[K, V], bool) {
	if t == nil || t.nodes == nil {
		return nil, false
	}
	return t.nodes.Get(id)
}

// Has reports whether id exists.
func (t *Tree[K, V]) Has(id K) bool {
	_, ok := t.Get(id)
	return ok
}

// Parent returns parent node by child id.
func (t *Tree[K, V]) Parent(id K) (*Node[K, V], bool) {
	node, ok := t.Get(id)
	if !ok || node.parent == nil {
		return nil, false
	}
	return node.parent, true
}

// Children returns children snapshot by node id.
func (t *Tree[K, V]) Children(id K) []*Node[K, V] {
	node, ok := t.Get(id)
	if !ok {
		return nil
	}
	return node.Children()
}

// Roots returns root nodes snapshot.
func (t *Tree[K, V]) Roots() []*Node[K, V] {
	if t == nil || t.roots == nil {
		return nil
	}
	return t.roots.Values()
}

// Ancestors returns parent chain from direct parent to top root.
func (t *Tree[K, V]) Ancestors(id K) []*Node[K, V] {
	node, ok := t.Get(id)
	if !ok {
		return nil
	}

	ancestors := make([]*Node[K, V], 0)
	for current := node.parent; current != nil; current = current.parent {
		ancestors = append(ancestors, current)
	}
	if len(ancestors) == 0 {
		return nil
	}
	return ancestors
}

// Depth returns the number of edges from one root to the node.
func (t *Tree[K, V]) Depth(id K) (int, bool) {
	node, ok := t.Get(id)
	if !ok {
		return 0, false
	}

	depth := 0
	for current := node.parent; current != nil; current = current.parent {
		depth++
	}
	return depth, true
}

// Siblings returns sibling nodes snapshot excluding the node itself.
func (t *Tree[K, V]) Siblings(id K) []*Node[K, V] {
	node, ok := t.Get(id)
	if !ok {
		return nil
	}

	if node.parent == nil {
		if t.roots == nil {
			return nil
		}

		rootCount := t.roots.Len()
		if rootCount <= 1 {
			return nil
		}

		siblings := make([]*Node[K, V], 0, rootCount-1)
		for i := range rootCount {
			candidate, _ := t.roots.Get(i)
			if candidate != node {
				siblings = append(siblings, candidate)
			}
		}
		if len(siblings) == 0 {
			return nil
		}
		return siblings
	}

	childCount := node.parent.children.Len()
	if childCount <= 1 {
		return nil
	}

	siblings := make([]*Node[K, V], 0, childCount-1)
	for i := range childCount {
		candidate, _ := node.parent.children.Get(i)
		if candidate != node {
			siblings = append(siblings, candidate)
		}
	}
	if len(siblings) == 0 {
		return nil
	}
	return siblings
}

// Descendants returns all descendants in DFS pre-order.
func (t *Tree[K, V]) Descendants(id K) []*Node[K, V] {
	node, ok := t.Get(id)
	if !ok {
		return nil
	}

	if node.children.Len() == 0 {
		return nil
	}

	capacity := t.Len()
	if capacity <= 0 {
		capacity = node.children.Len()
	}
	descendants := make([]*Node[K, V], 0, capacity)
	stack := appendChildrenReverse(make([]*Node[K, V], 0, capacity), node)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		descendants = append(descendants, current)

		stack = appendChildrenReverse(stack, current)
	}

	return descendants
}

// Leaves returns all leaf nodes in DFS pre-order.
func (t *Tree[K, V]) Leaves() []*Node[K, V] {
	if t == nil {
		return nil
	}
	return collectLeavesFromList(t.roots)
}

// RangeDFS iterates all nodes in DFS pre-order until fn returns false.
func (t *Tree[K, V]) RangeDFS(fn func(node *Node[K, V]) bool) {
	if t == nil || fn == nil {
		return
	}

	if t.roots == nil || t.roots.Len() == 0 {
		return
	}

	for i := range t.roots.Len() {
		root, _ := t.roots.Get(i)
		if !rangeDFSFromRoot(root, fn) {
			return
		}
	}
}

// RangeBFS iterates all nodes in BFS order until fn returns false.
func (t *Tree[K, V]) RangeBFS(fn func(node *Node[K, V]) bool) {
	if t == nil || fn == nil {
		return
	}

	if t.roots == nil || t.roots.Len() == 0 {
		return
	}

	rangeBFSFromRoots(t.roots, fn)
}

// Len returns total node count.
func (t *Tree[K, V]) Len() int {
	if t == nil || t.nodes == nil {
		return 0
	}
	return t.nodes.Len()
}

// IsEmpty reports whether tree has no nodes.
func (t *Tree[K, V]) IsEmpty() bool {
	return t.Len() == 0
}

func rangeDFSFromRoot[K comparable, V any](root *Node[K, V], fn func(node *Node[K, V]) bool) bool {
	if root == nil {
		return true
	}

	stack := []*Node[K, V]{root}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !fn(current) {
			return false
		}

		stack = appendChildrenReverse(stack, current)
	}

	return true
}

func rangeBFSFromRoots[K comparable, V any](roots *collectionlist.List[*Node[K, V]], fn func(node *Node[K, V]) bool) {
	if roots == nil || roots.Len() == 0 {
		return
	}

	queue := make([]*Node[K, V], 0, roots.Len())
	for i := range roots.Len() {
		root, _ := roots.Get(i)
		queue = append(queue, root)
	}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if !fn(current) {
			return
		}

		childCount := current.children.Len()
		for i := range childCount {
			child, _ := current.children.Get(i)
			queue = append(queue, child)
		}
	}
}

func appendChildrenReverse[K comparable, V any](stack []*Node[K, V], node *Node[K, V]) []*Node[K, V] {
	if node == nil {
		return stack
	}

	for i := node.children.Len() - 1; i >= 0; i-- {
		child, _ := node.children.Get(i)
		stack = append(stack, child)
	}

	return stack
}

func collectLeaves[K comparable, V any](roots []*Node[K, V]) []*Node[K, V] {
	if len(roots) == 0 {
		return nil
	}

	leaves := make([]*Node[K, V], 0, len(roots))
	for _, root := range roots {
		_ = rangeDFSFromRoot(root, func(node *Node[K, V]) bool {
			if node.children.Len() == 0 {
				leaves = append(leaves, node)
			}
			return true
		})
	}
	if len(leaves) == 0 {
		return nil
	}
	return leaves
}

func collectLeavesFromList[K comparable, V any](roots *collectionlist.List[*Node[K, V]]) []*Node[K, V] {
	if roots == nil || roots.Len() == 0 {
		return nil
	}

	leaves := make([]*Node[K, V], 0, roots.Len())
	for i := range roots.Len() {
		root, _ := roots.Get(i)
		_ = rangeDFSFromRoot(root, func(node *Node[K, V]) bool {
			if node.children.Len() == 0 {
				leaves = append(leaves, node)
			}
			return true
		})
	}
	if len(leaves) == 0 {
		return nil
	}
	return leaves
}

func rangeDFSRoots[K comparable, V any](roots []*Node[K, V], fn func(node *Node[K, V]) bool) {
	for _, root := range roots {
		if !rangeDFSFromRoot(root, fn) {
			return
		}
	}
}

func rangeBFSRoots[K comparable, V any](roots []*Node[K, V], fn func(node *Node[K, V]) bool) {
	if len(roots) == 0 {
		return
	}

	queue := append(make([]*Node[K, V], 0, len(roots)), roots...)
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if !fn(current) {
			return
		}

		childCount := current.children.Len()
		for i := range childCount {
			child, _ := current.children.Get(i)
			queue = append(queue, child)
		}
	}
}
