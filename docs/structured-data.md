---
title: 'collectionx Lists and Structured Data'
linkTitle: 'lists-data'
description: 'Deque, ring buffer, intervals, trie, and tree examples'
weight: 4
---

## Lists and structured data

Examples for **`collectionx/list`**, **`collectionx/bytex`**, **`collectionx/interval`**, **`collectionx/prefix`**, and **`collectionx/tree`**. Each block is a complete `package main`.

## 1) `Deque` and `RingBuffer`

`Push` on a full `RingBuffer` returns `mo.Option[T]` for the evicted element.

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/list"
)

func main() {
	dq := list.NewDeque[int]()
	dq.PushBack(1, 2)
	dq.PushFront(0)
	fmt.Println(dq.Values())

	rb := list.NewRingBuffer[int](2)
	_ = rb.Push(1)
	_ = rb.Push(2)
	ev := rb.Push(3)
	if v, ok := ev.Get(); ok {
		fmt.Println("evicted", v)
	}
}
```

## 2) Byte-specialized data: `bytex.List`, `RingBuffer`, `Set`, and `Counter`

Use `bytex` when the data domain is bytes and you want byte-oriented helpers without giving up direct serializer support.

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/arcgolabs/collectionx/bytex"
)

func main() {
	buf := bytex.NewListFromString("prefix-body")
	drained, _ := buf.DrainPrefix(len("prefix-"))
	fmt.Println(string(drained))
	fmt.Println(buf.String())

	tail := bytex.NewRingBuffer(4)
	tail.WriteString("abcdef")
	fmt.Println(tail.String())

	seen := bytex.NewSet('a', 'b', 'a')
	seen.AddRange('0', '9'+1)
	fmt.Println(seen.Values())

	counts := bytex.NewCounter()
	counts.AddString("banana")
	fmt.Println(counts.MostCommon(2))

	data, _ := json.Marshal(buf)
	fmt.Println(string(data))
}
```

## 3) Intervals: `RangeSet` and `RangeMap`

Half-open ranges `[start, end)` are normalized inside `RangeSet`. `RangeMap.Get` resolves a point to a value.

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/interval"
)

func main() {
	rs := interval.NewRangeSet[int]()
	rs.Add(1, 5)
	rs.Add(5, 8)
	fmt.Println(rs.Ranges())

	rm := interval.NewRangeMap[int, string]()
	rm.Put(0, 10, "A")
	rm.Put(3, 5, "B")
	v, ok := rm.Get(4)
	fmt.Println(v, ok)
}
```

## 4) Prefix map: `Trie`

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/prefix"
)

func main() {
	tr := prefix.NewTrie[int]()
	tr.Put("user:1", 1)
	tr.Put("user:2", 2)
	tr.Put("order:9", 9)

	fmt.Println(tr.KeysWithPrefix("user:"))
}
```

## 5) Hierarchy: `Tree`

```go
package main

import (
	"fmt"
	"log"

	"github.com/arcgolabs/collectionx/tree"
)

func main() {
	org := tree.NewTree[int, string]()
	if err := org.AddRoot(1, "CEO"); err != nil {
		log.Fatal(err)
	}
	if err := org.AddChild(1, 2, "CTO"); err != nil {
		log.Fatal(err)
	}
	if err := org.AddChild(2, 3, "Platform Lead"); err != nil {
		log.Fatal(err)
	}

	parent, ok := org.Parent(3)
	if !ok {
		log.Fatal("parent not found")
	}
	fmt.Println(parent.ID())
	fmt.Println(len(org.Descendants(1)))
}
```

## Related

- [Getting Started](./getting-started)
- [Maps, sets, and tables](./mapping-recipes)

## Serialization notes

The structures on this page can be passed directly to `json.Marshal`, `json.Unmarshal`, `gob` encoders/decoders, and binary codecs without calling a separate snapshot helper.

The main exception is `PriorityQueue`: values can be serialized, but automatic restore is not supported because the comparator is runtime configuration rather than serialized data.
