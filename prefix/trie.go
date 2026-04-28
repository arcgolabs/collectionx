package prefix

import (
	"slices"

	collectionmapping "github.com/arcgolabs/collectionx/mapping"
	"github.com/samber/mo"
)

type trieNode[V any] struct {
	children   collectionmapping.Map[rune, *trieNode[V]]
	childKeys  []rune
	hasValue   bool
	value      V
	valueCount int
}

// Trie is a prefix tree for string keys.
// Zero value is ready to use.
type Trie[V any] struct {
	root *trieNode[V]
	size int
}

type keyValue[V any] struct {
	key   string
	value V
}

// NewTrie creates an empty trie.
func NewTrie[V any]() *Trie[V] {
	return &Trie[V]{}
}

// NewPrefixMap creates an empty prefix map.
// PrefixMap shares the same implementation as Trie.
func NewPrefixMap[V any]() *Trie[V] {
	return NewTrie[V]()
}

// Put stores key -> value.
// Returns true when inserted as a new key, false when updated existing key.
func (t *Trie[V]) Put(key string, value V) bool {
	if t == nil {
		return false
	}
	t.ensureRoot()

	node := t.root
	path := []*trieNode[V]{node}
	for _, ch := range key {
		next, ok := node.children.Get(ch)
		if !ok {
			next = &trieNode[V]{}
			node.children.Set(ch, next)
			node.insertChildKey(ch)
		}
		node = next
		path = append(path, node)
	}

	isNew := !node.hasValue
	node.value = value
	node.hasValue = true
	if isNew {
		for _, current := range path {
			current.valueCount++
		}
		t.size++
	}
	return isNew
}

// Get returns value by exact key.
func (t *Trie[V]) Get(key string) (V, bool) {
	var zero V
	if t == nil || t.root == nil {
		return zero, false
	}
	node, ok := t.findNode(key)
	if !ok || !node.hasValue {
		return zero, false
	}
	return node.value, true
}

// GetOption returns value by exact key as mo.Option.
func (t *Trie[V]) GetOption(key string) mo.Option[V] {
	value, ok := t.Get(key)
	if !ok {
		return mo.None[V]()
	}
	return mo.Some(value)
}

// Has reports whether exact key exists.
func (t *Trie[V]) Has(key string) bool {
	_, ok := t.Get(key)
	return ok
}

// HasPrefix reports whether prefix exists in trie paths.
func (t *Trie[V]) HasPrefix(prefix string) bool {
	if t == nil || t.root == nil {
		return false
	}
	_, ok := t.findNode(prefix)
	return ok
}

// CountPrefix reports how many stored keys start with prefix.
func (t *Trie[V]) CountPrefix(prefix string) int {
	if t == nil || t.root == nil {
		return 0
	}
	node, ok := t.findNode(prefix)
	if !ok {
		return 0
	}
	return node.valueCount
}

// LongestPrefix returns the longest stored key that is a prefix of key.
func (t *Trie[V]) LongestPrefix(key string) (string, V, bool) {
	var zero V
	if t == nil || t.root == nil {
		return "", zero, false
	}

	node := t.root
	lastMatchIndex := -1
	var lastMatchValue V
	matched := false
	if node.hasValue {
		lastMatchValue = node.value
		matched = true
	}
	runes := []rune(key)
	for index, ch := range runes {
		next, ok := node.children.Get(ch)
		if !ok {
			break
		}
		node = next
		if node.hasValue {
			lastMatchIndex = index
			lastMatchValue = node.value
			matched = true
		}
	}
	if !matched {
		return "", zero, false
	}
	if lastMatchIndex < 0 {
		return "", lastMatchValue, true
	}
	return string(runes[:lastMatchIndex+1]), lastMatchValue, true
}

// Delete removes key and returns whether key existed.
func (t *Trie[V]) Delete(key string) bool {
	if t == nil || t.root == nil {
		return false
	}
	runes := []rune(key)
	removed := t.deleteRec(t.root, runes, 0)
	if removed {
		t.size--
	}
	return removed
}

// DeletePrefix removes all keys that start with prefix and returns removed key count.
func (t *Trie[V]) DeletePrefix(prefix string) int {
	if t == nil || t.root == nil {
		return 0
	}
	if prefix == "" {
		removed := t.size
		t.Clear()
		return removed
	}

	runes := []rune(prefix)
	path := make([]*trieNode[V], 1, len(runes)+1)
	path[0] = t.root
	node := t.root
	for _, ch := range runes {
		next, ok := node.children.Get(ch)
		if !ok {
			return 0
		}
		node = next
		path = append(path, node)
	}

	removed := node.valueCount
	if removed == 0 {
		return 0
	}

	parent := path[len(path)-2]
	lastRune := runes[len(runes)-1]
	parent.children.Delete(lastRune)
	parent.deleteChildKey(lastRune)

	for _, current := range path[:len(path)-1] {
		current.valueCount -= removed
	}
	t.size -= removed
	return removed
}

// Len returns stored key count.
func (t *Trie[V]) Len() int {
	if t == nil {
		return 0
	}
	return t.size
}

// IsEmpty reports whether trie has no keys.
func (t *Trie[V]) IsEmpty() bool {
	return t.Len() == 0
}

// Clear removes all keys.
func (t *Trie[V]) Clear() {
	if t == nil {
		return
	}
	t.root = nil
	t.size = 0
}

// KeysWithPrefix returns all keys that start with prefix.
func (t *Trie[V]) KeysWithPrefix(prefix string) []string {
	if t == nil || t.root == nil {
		return nil
	}

	startNode, ok := t.findNode(prefix)
	if !ok {
		return nil
	}

	out := make([]string, 0, startNode.valueCount)
	t.collectKeys(startNode, prefix, &out)
	return out
}

// ValuesWithPrefix returns all values under prefix.
func (t *Trie[V]) ValuesWithPrefix(prefix string) []V {
	if t == nil || t.root == nil {
		return nil
	}

	startNode, ok := t.findNode(prefix)
	if !ok {
		return nil
	}

	out := make([]V, 0, startNode.valueCount)
	t.collectValues(startNode, prefix, &out)
	return out
}

// RangePrefix iterates keys with prefix in lexicographic key order until fn returns false.
func (t *Trie[V]) RangePrefix(prefix string, fn func(key string, value V) bool) {
	if t == nil || t.root == nil || fn == nil {
		return
	}

	startNode, ok := t.findNode(prefix)
	if !ok {
		return
	}

	t.rangePrefix(startNode, prefix, fn)
}

func (t *Trie[V]) ensureRoot() {
	if t.root == nil {
		t.root = &trieNode[V]{}
	}
}

func (t *Trie[V]) findNode(key string) (*trieNode[V], bool) {
	node := t.root
	for _, ch := range key {
		next, ok := node.children.Get(ch)
		if !ok {
			return nil, false
		}
		node = next
	}
	return node, true
}

func (t *Trie[V]) deleteRec(node *trieNode[V], runes []rune, depth int) bool {
	if node == nil {
		return false
	}
	if depth == len(runes) {
		if !node.hasValue {
			return false
		}
		node.hasValue = false
		node.valueCount--
		var zero V
		node.value = zero
		return true
	}

	ch := runes[depth]
	child, ok := node.children.Get(ch)
	if !ok {
		return false
	}
	removed := t.deleteRec(child, runes, depth+1)
	if !removed {
		return false
	}

	node.valueCount--
	if !child.hasValue && child.children.Len() == 0 {
		node.children.Delete(ch)
		node.deleteChildKey(ch)
	}
	return true
}

func (n *trieNode[V]) insertChildKey(ch rune) {
	index, found := slices.BinarySearch(n.childKeys, ch)
	if found {
		return
	}
	n.childKeys = append(n.childKeys, 0)
	copy(n.childKeys[index+1:], n.childKeys[index:])
	n.childKeys[index] = ch
}

func (n *trieNode[V]) deleteChildKey(ch rune) {
	index, found := slices.BinarySearch(n.childKeys, ch)
	if !found {
		return
	}
	copy(n.childKeys[index:], n.childKeys[index+1:])
	n.childKeys[len(n.childKeys)-1] = 0
	n.childKeys = n.childKeys[:len(n.childKeys)-1]
}
