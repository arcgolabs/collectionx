package tree

// NodeSnapshot is a detached tree node for serialization adapters.
type NodeSnapshot[K comparable, V any] struct {
	ID       K                    `json:"id"`
	Value    V                    `json:"value"`
	Children []NodeSnapshot[K, V] `json:"children,omitempty"`
}

// Nodes returns detached roots with nested children in tree order.
func (t *Tree[K, V]) Nodes() []NodeSnapshot[K, V] {
	if t == nil || t.IsEmpty() {
		return nil
	}

	roots := make([]NodeSnapshot[K, V], 0, t.roots.Len())
	rootCount := t.roots.Len()
	for index := range rootCount {
		root, _ := t.roots.Get(index)
		roots = append(roots, snapshotNode(root))
	}
	return roots
}

// Entries returns detached flat entries in DFS pre-order.
func (t *Tree[K, V]) Entries() []Entry[K, V] {
	if t == nil || t.IsEmpty() {
		return nil
	}

	entries := make([]Entry[K, V], 0, t.Len())
	t.RangeDFS(func(node *Node[K, V]) bool {
		if node.parent == nil {
			entries = append(entries, RootEntry(node.ID(), node.Value()))
			return true
		}
		entries = append(entries, ChildEntry(node.ID(), node.parent.ID(), node.Value()))
		return true
	})
	return entries
}

// Nodes returns detached roots with nested children from a stable snapshot.
func (t *ConcurrentTree[K, V]) Nodes() []NodeSnapshot[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	return t.tree.Nodes()
}

// Entries returns detached flat entries from a stable snapshot.
func (t *ConcurrentTree[K, V]) Entries() []Entry[K, V] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tree == nil {
		return nil
	}
	return t.tree.Entries()
}

func snapshotNode[K comparable, V any](node *Node[K, V]) NodeSnapshot[K, V] {
	if node == nil {
		return NodeSnapshot[K, V]{}
	}

	out := NodeSnapshot[K, V]{
		ID:    node.ID(),
		Value: node.Value(),
	}
	if node.children.Len() == 0 {
		return out
	}

	out.Children = make([]NodeSnapshot[K, V], 0, node.children.Len())
	childCount := node.children.Len()
	for index := range childCount {
		child, _ := node.children.Get(index)
		out.Children = append(out.Children, snapshotNode(child))
	}
	return out
}
