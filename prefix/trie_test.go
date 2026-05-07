package prefix_test

import (
	"testing"

	prefix "github.com/arcgolabs/collectionx/prefix"
	"github.com/stretchr/testify/require"
)

func TestTrie_BasicOps(t *testing.T) {
	t.Parallel()

	var tr prefix.Trie[int]
	require.True(t, tr.Put("cat", 1))
	require.True(t, tr.Put("car", 2))
	require.False(t, tr.Put("cat", 9))

	value, ok := tr.Get("cat")
	require.True(t, ok)
	require.Equal(t, 9, value)
	require.True(t, tr.HasPrefix("ca"))
	require.False(t, tr.Has("c"))
	require.Equal(t, 2, tr.Len())
	key, firstValue, ok := tr.GetFirst()
	require.True(t, ok)
	require.Equal(t, "car", key)
	require.Equal(t, 2, firstValue)
	key, firstValue, ok = tr.GetFirstWithPrefix("cat")
	require.True(t, ok)
	require.Equal(t, "cat", key)
	require.Equal(t, 9, firstValue)
}

func TestTrie_PrefixAndDelete(t *testing.T) {
	t.Parallel()

	tr := prefix.NewTrie[string]()
	tr.Put("go", "v1")
	tr.Put("gone", "v2")
	tr.Put("good", "v3")
	tr.Put("zoo", "v4")

	require.Equal(t, []string{"go", "gone", "good"}, tr.KeysWithPrefix("go"))
	require.Equal(t, []prefix.Entry[string]{
		{Key: "go", Value: "v1"},
		{Key: "gone", Value: "v2"},
		{Key: "good", Value: "v3"},
	}, tr.EntriesWithPrefix("go"))
	require.Equal(t, 3, tr.CountPrefix("go"))
	matchedKey, matchedValue, ok := tr.LongestPrefix("goodbye")
	require.True(t, ok)
	require.Equal(t, "good", matchedKey)
	require.Equal(t, "v3", matchedValue)
	require.True(t, tr.Delete("gone"))
	require.Equal(t, []string{"go", "good"}, tr.KeysWithPrefix("go"))
	require.Equal(t, 2, tr.CountPrefix("go"))
	require.Equal(t, 2, tr.DeletePrefix("go"))
	require.Equal(t, []string{"zoo"}, tr.KeysWithPrefix("z"))
	require.Equal(t, 0, tr.CountPrefix("go"))
	require.Equal(t, 1, tr.Len())
	require.Equal(t, 0, tr.DeletePrefix("missing"))
}

func TestTrie_EntriesWithPrefix_EmptyOrMissing(t *testing.T) {
	t.Parallel()

	var tr prefix.Trie[int]
	require.Nil(t, tr.EntriesWithPrefix("go"))

	tr.Put("cat", 1)
	require.Nil(t, tr.EntriesWithPrefix("dog"))
}

func TestTrie_LongestPrefix_NotFound(t *testing.T) {
	t.Parallel()

	tr := prefix.NewTrie[int]()
	tr.Put("cat", 1)

	_, _, ok := tr.LongestPrefix("dog")
	require.False(t, ok)
}

func TestTrie_LongestPrefix_EmptyKey(t *testing.T) {
	t.Parallel()

	tr := prefix.NewTrie[int]()
	tr.Put("", 1)
	tr.Put("cat", 2)

	key, value, ok := tr.LongestPrefix("cab")
	require.True(t, ok)
	require.Equal(t, "", key)
	require.Equal(t, 1, value)
}

func TestPrefixMap_AliasConstructor(t *testing.T) {
	t.Parallel()

	pm := prefix.NewPrefixMap[int]()
	pm.Put("ab", 1)
	opt := pm.GetOption("ab")
	require.True(t, opt.IsPresent())
}
