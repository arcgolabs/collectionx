---
title: 'collectionx Maps, Sets, and Tables'
linkTitle: 'maps-sets'
description: 'Recipes for Set, Ordered types, MultiMap, Table, and serialization helpers'
weight: 3
---

## Maps, sets, and tables

Patterns for **`collectionx/set`** and **`collectionx/mapping`**: deduplication, stable iteration order, one-to-many keys, 2D tables, and serialization helpers.

Each section is a standalone `package main` you can paste into its own file.

## 1) Deduplicate with `Set`

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/set"
)

func main() {
	s := set.NewSet[string]()
	s.Add("A", "A", "B")
	fmt.Println(s.Len())
	fmt.Println(s.Contains("B"))
}
```

## 2) Insertion order: `OrderedSet` / `OrderedMap`

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/mapping"
	"github.com/arcgolabs/collectionx/set"
)

func main() {
	os := set.NewOrderedSet[int]()
	os.Add(3, 1, 3, 2)
	fmt.Println(os.Values())

	om := mapping.NewOrderedMap[string, int]()
	om.Set("x", 1)
	om.Set("y", 2)
	om.Set("x", 9)
	fmt.Println(om.Keys())
	fmt.Println(om.Values())
}
```

## 3) One-to-many: `MultiMap`

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/mapping"
)

func main() {
	mm := mapping.NewMultiMap[string, int]()
	mm.PutAll("tag", 1, 2, 3)
	fmt.Println(mm.Get("tag"))
	owned := mm.GetCopy("tag")
	fmt.Println("copy len", len(owned))
	fmt.Println(mm.ValueCount())
	removed := mm.DeleteValueIf("tag", func(v int) bool { return v%2 == 0 })
	fmt.Println(removed, mm.Get("tag"))
}
```

## 4) Two-dimensional keys: `Table`

```go
package main

import (
	"fmt"

	"github.com/arcgolabs/collectionx/mapping"
)

func main() {
	t := mapping.NewTable[string, string, int]()
	t.Put("r1", "c1", 10)
	t.Put("r1", "c2", 20)
	t.Put("r2", "c1", 30)

	v, ok := t.Get("r1", "c2")
	fmt.Println(v, ok)
	fmt.Println(t.Row("r1"))
	fmt.Println(t.Column("c1"))
}
```

## 5) Serialize directly

Most structures can be passed directly to `encoding/json`, `encoding/gob`, or binary codecs without calling a separate snapshot helper.

```go
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/arcgolabs/collectionx/set"
)

func main() {
	s := set.NewOrderedSet[string]("a", "b")

	payload, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(payload))

	var restored set.OrderedSet[string]
	if err := json.Unmarshal(payload, &restored); err != nil {
		panic(err)
	}
	fmt.Println(restored.Values())

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s); err != nil {
		panic(err)
	}

	var restoredGob set.OrderedSet[string]
	if err := gob.NewDecoder(&buf).Decode(&restoredGob); err != nil {
		panic(err)
	}
	fmt.Println(restoredGob.Values())
}
```

`ToJSON()` and `String()` are still available when you explicitly want JSON bytes or log-friendly output.

## 6) Caveats

- `PriorityQueue` cannot be automatically restored from serialized data because its comparator is runtime configuration.
- `ShardedConcurrentMap` can be serialized directly, but restore should target an already-initialized receiver created with `NewShardedConcurrentMap(...)`.

## Related

- Minimal first steps: [Getting Started](./getting-started)
- Deque, intervals, trie, tree: [Lists and structured data](./structured-data)
