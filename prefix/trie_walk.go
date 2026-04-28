package prefix

type trieWalkFrame[V any] struct {
	node       *trieNode[V]
	childIndex int
}

func walkPrefix[V any](node *trieNode[V], prefix string, visit func(path []rune, node *trieNode[V]) bool) bool {
	if node == nil || visit == nil {
		return true
	}

	prefixRunes := []rune(prefix)
	path := make([]rune, 0, len(prefixRunes)+32)
	path = append(path, prefixRunes...)
	stack := make([]trieWalkFrame[V], 1, 32)
	stack[0] = trieWalkFrame[V]{node: node, childIndex: -1}
	for len(stack) > 0 {
		frame := &stack[len(stack)-1]
		if frame.childIndex < 0 {
			if frame.node.hasValue && !visit(path, frame.node) {
				return false
			}
			frame.childIndex = 0
		}

		if frame.childIndex >= len(frame.node.childKeys) {
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				path = path[:len(path)-1]
			}
			continue
		}

		ch := frame.node.childKeys[frame.childIndex]
		frame.childIndex++
		child, _ := frame.node.children.Get(ch)
		path = append(path, ch)
		stack = append(stack, trieWalkFrame[V]{node: child, childIndex: -1})
	}
	return true
}

func walkPrefixValues[V any](node *trieNode[V], visit func(node *trieNode[V]) bool) bool {
	if node == nil || visit == nil {
		return true
	}

	stack := make([]*trieNode[V], 1, 32)
	stack[0] = node
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.hasValue && !visit(current) {
			return false
		}
		for index := len(current.childKeys) - 1; index >= 0; index-- {
			child, _ := current.children.Get(current.childKeys[index])
			stack = append(stack, child)
		}
	}
	return true
}

func (t *Trie[V]) collectPairs(node *trieNode[V], prefix string, out *[]keyValue[V]) {
	_ = walkPrefix(node, prefix, func(path []rune, node *trieNode[V]) bool {
		*out = append(*out, keyValue[V]{
			key:   string(path),
			value: node.value,
		})
		return true
	})
}

func (t *Trie[V]) collectKeys(node *trieNode[V], prefix string, out *[]string) {
	_ = walkPrefix(node, prefix, func(path []rune, _ *trieNode[V]) bool {
		*out = append(*out, string(path))
		return true
	})
}

func (t *Trie[V]) collectValues(node *trieNode[V], prefix string, out *[]V) {
	_ = prefix
	_ = walkPrefixValues(node, func(node *trieNode[V]) bool {
		*out = append(*out, node.value)
		return true
	})
}

func (t *Trie[V]) rangePrefix(node *trieNode[V], prefix string, fn func(key string, value V) bool) bool {
	return walkPrefix(node, prefix, func(path []rune, node *trieNode[V]) bool {
		return fn(string(path), node.value)
	})
}

func (t *Trie[V]) pairsWithPrefix(prefix string) []keyValue[V] {
	if t == nil || t.root == nil {
		return nil
	}
	startNode, ok := t.findNode(prefix)
	if !ok {
		return nil
	}

	out := make([]keyValue[V], 0, startNode.valueCount)
	t.collectPairs(startNode, prefix, &out)
	return out
}
