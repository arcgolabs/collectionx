package tree

// Get returns node by id as a detached node clone with ancestor chain.
func (t *ConcurrentTree[K, V]) Get(id K) (*Node[K, V], bool) {
	if t == nil {
		return nil, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil, false
	}
	node, ok := t.tree.Get(id)
	if !ok {
		return nil, false
	}
	return cloneNodeWithAncestorsShallow(node), true
}

// Has reports whether id exists.
func (t *ConcurrentTree[K, V]) Has(id K) bool {
	if t == nil {
		return false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return false
	}
	return t.tree.Has(id)
}

// Parent returns parent node by child id as a detached node clone with ancestor chain.
func (t *ConcurrentTree[K, V]) Parent(id K) (*Node[K, V], bool) {
	if t == nil {
		return nil, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil, false
	}
	node, ok := t.tree.Get(id)
	if !ok || node.parent == nil {
		return nil, false
	}
	return cloneNodeWithAncestorsShallow(node.parent), true
}

// Children returns children snapshot by node id.
func (t *ConcurrentTree[K, V]) Children(id K) []*Node[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	node, ok := t.tree.Get(id)
	if !ok {
		return nil
	}
	if node.children.Len() == 0 {
		return nil
	}

	children := make([]*Node[K, V], 0, node.children.Len())
	childCount := node.children.Len()
	for index := range childCount {
		child, _ := node.children.Get(index)
		children = append(children, cloneSubtreeDetached(child))
	}
	return children
}

// Roots returns root nodes snapshot.
func (t *ConcurrentTree[K, V]) Roots() []*Node[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	if t.tree.roots == nil || t.tree.roots.Len() == 0 {
		return nil
	}

	roots := make([]*Node[K, V], 0, t.tree.roots.Len())
	rootCount := t.tree.roots.Len()
	for index := range rootCount {
		root, _ := t.tree.roots.Get(index)
		roots = append(roots, cloneSubtreeDetached(root))
	}
	return roots
}

// Ancestors returns parent chain from direct parent to top root.
func (t *ConcurrentTree[K, V]) Ancestors(id K) []*Node[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	node, ok := t.tree.Get(id)
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

	out := make([]*Node[K, V], len(ancestors))
	var parentClone *Node[K, V]
	for i, ancestor := range ancestors {
		currentClone := newNode(ancestor.ID(), ancestor.Value())
		currentClone.parent = parentClone
		if parentClone != nil {
			parentClone.children.Add(currentClone)
		}
		out[i] = currentClone
		parentClone = currentClone
	}
	return out
}

// Depth returns the number of edges from one root to the node.
func (t *ConcurrentTree[K, V]) Depth(id K) (int, bool) {
	if t == nil {
		return 0, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return 0, false
	}
	node, ok := t.tree.Get(id)
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
func (t *ConcurrentTree[K, V]) Siblings(id K) []*Node[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	node, ok := t.tree.Get(id)
	if !ok {
		return nil
	}

	var siblingsSource []*Node[K, V]
	if node.parent == nil {
		if t.tree.roots == nil || t.tree.roots.Len() <= 1 {
			return nil
		}
		siblingsSource = t.tree.roots.Values()
	} else {
		if node.parent.children.Len() <= 1 {
			return nil
		}
		siblingsSource = node.parent.children.Values()
	}

	siblings := make([]*Node[K, V], 0, len(siblingsSource)-1)
	for _, sibling := range siblingsSource {
		if sibling != node {
			siblings = append(siblings, cloneSubtreeDetached(sibling))
		}
	}
	if len(siblings) == 0 {
		return nil
	}
	return siblings
}

// Descendants returns all descendants in DFS pre-order.
func (t *ConcurrentTree[K, V]) Descendants(id K) []*Node[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	if t.tree == nil {
		t.mu.RUnlock()
		return nil
	}
	node, ok := t.tree.Get(id)
	if !ok {
		t.mu.RUnlock()
		return nil
	}
	_, descendants := cloneSubtreeDetachedWithDescendants(node, true)
	t.mu.RUnlock()
	return descendants
}

// Leaves returns all leaf nodes in DFS pre-order.
func (t *ConcurrentTree[K, V]) Leaves() []*Node[K, V] {
	if t == nil {
		return nil
	}

	return collectLeaves(t.snapshotClonedRoots())
}

// RangeDFS iterates all nodes in DFS pre-order until fn returns false.
func (t *ConcurrentTree[K, V]) RangeDFS(fn func(node *Node[K, V]) bool) {
	if t == nil || fn == nil {
		return
	}

	rangeDFSRoots(t.snapshotClonedRoots(), fn)
}

// RangeBFS iterates all nodes in BFS order until fn returns false.
func (t *ConcurrentTree[K, V]) RangeBFS(fn func(node *Node[K, V]) bool) {
	if t == nil || fn == nil {
		return
	}

	rangeBFSRoots(t.snapshotClonedRoots(), fn)
}

// Len returns total node count.
func (t *ConcurrentTree[K, V]) Len() int {
	if t == nil {
		return 0
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return 0
	}
	return t.tree.Len()
}

// IsEmpty reports whether tree has no nodes.
func (t *ConcurrentTree[K, V]) IsEmpty() bool {
	return t.Len() == 0
}

func cloneNodeWithAncestorsShallow[K comparable, V any](node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	targetClone := newNode(node.ID(), node.Value())
	currentClone := targetClone
	for currentParent := node.parent; currentParent != nil; currentParent = currentParent.parent {
		parentClone := newNode(currentParent.ID(), currentParent.Value())
		currentClone.parent = parentClone
		parentClone.children.Add(currentClone)
		currentClone = parentClone
	}
	return targetClone
}

func (t *ConcurrentTree[K, V]) snapshotClonedRoots() []*Node[K, V] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil || t.tree.roots == nil || t.tree.roots.Len() == 0 {
		return nil
	}

	roots := make([]*Node[K, V], 0, t.tree.roots.Len())
	rootCount := t.tree.roots.Len()
	for index := range rootCount {
		root, _ := t.tree.roots.Get(index)
		roots = append(roots, cloneSubtreeDetached(root))
	}
	return roots
}
