package tree

import collectionlist "github.com/arcgolabs/collectionx/list"

// Clone returns a deep copy preserving parent-children structure.
func (t *Tree[K, V]) Clone() *Tree[K, V] {
	if t == nil || t.nodes == nil || t.nodes.IsEmpty() {
		return NewTree[K, V]()
	}

	rootCount := 0
	if t.roots != nil {
		rootCount = t.roots.Len()
	}
	cloned := newTreeWithCapacity[K, V](t.nodes.Len(), rootCount)

	type pair struct {
		source *Node[K, V]
		target *Node[K, V]
	}

	stack := make([]pair, 0, rootCount)
	for index := range rootCount {
		root, _ := t.roots.Get(index)
		rootClone := newNode(root.ID(), root.Value())
		cloned.roots.Add(rootClone)
		cloned.nodes.Set(rootClone.ID(), rootClone)
		stack = append(stack, pair{source: root, target: rootClone})
	}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		childCount := current.source.children.Len()
		for index := range childCount {
			sourceChild, _ := current.source.children.Get(index)
			targetChild := newNode(sourceChild.ID(), sourceChild.Value())
			targetChild.parent = current.target
			current.target.children.Add(targetChild)
			cloned.nodes.Set(targetChild.ID(), targetChild)
			stack = append(stack, pair{source: sourceChild, target: targetChild})
		}
	}

	return cloned
}

func cloneSubtreeDetached[K comparable, V any](root *Node[K, V]) *Node[K, V] {
	rootClone, _ := cloneSubtreeDetachedWithDescendants(root, false)
	return rootClone
}

func cloneSubtreeDetachedWithDescendants[K comparable, V any](root *Node[K, V], collectDescendants bool) (*Node[K, V], []*Node[K, V]) {
	if root == nil {
		return nil, nil
	}

	type sourceEntry struct {
		node        *Node[K, V]
		parentIndex int
		childCount  int
	}

	entries := make([]sourceEntry, 0, 16)
	type stackItem struct {
		node        *Node[K, V]
		parentIndex int
	}
	stack := []stackItem{{node: root, parentIndex: -1}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		childCount := current.node.children.Len()
		entryIndex := len(entries)
		entries = append(entries, sourceEntry{
			node:        current.node,
			parentIndex: current.parentIndex,
			childCount:  childCount,
		})

		for index := childCount - 1; index >= 0; index-- {
			child, _ := current.node.children.Get(index)
			stack = append(stack, stackItem{node: child, parentIndex: entryIndex})
		}
	}

	clones := make([]Node[K, V], len(entries))
	for i, entry := range entries {
		clones[i].id = entry.node.ID()
		clones[i].value = entry.node.Value()
		if entry.childCount > 0 {
			clones[i].children = *collectionlist.NewListWithCapacity[*Node[K, V]](entry.childCount)
		}
	}

	var descendants []*Node[K, V]
	if collectDescendants && len(entries) > 1 {
		descendants = make([]*Node[K, V], 0, len(entries)-1)
	}

	for i, entry := range entries {
		currentClone := &clones[i]
		if entry.parentIndex >= 0 {
			parentClone := &clones[entry.parentIndex]
			currentClone.parent = parentClone
			parentClone.children.Add(currentClone)
		}
		if collectDescendants && i > 0 {
			descendants = append(descendants, currentClone)
		}
	}

	return &clones[0], descendants
}
