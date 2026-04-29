## collectionx

`collectionx` provides strongly typed collection data structures for Go, including concurrent variants and non-standard structures such as `MultiMap`, `Table`, `Trie`, and interval types.

## Current capabilities

- **Generics-first** API with explicit method names and predictable semantics.
- **Optional concurrent** variants (`ConcurrentSet`, `ConcurrentMap`, …) when data is shared across goroutines.
- **Non-standard but practical** structures inspired by mature ecosystems (ordered maps, multi-maps, 2D `Table`, prefix `Trie`, interval maps, parent/child `Tree`).

## Package layout

- `github.com/arcgolabs/collectionx/set` — `Set`, `ConcurrentSet`, `MultiSet`, `OrderedSet`
- `github.com/arcgolabs/collectionx/mapping` — `Map`, `ConcurrentMap`, `BiMap`, `OrderedMap`, `MultiMap`, `Table`
- `github.com/arcgolabs/collectionx/list` — `List`, `ConcurrentList`, `Deque`, `RingBuffer`, `PriorityQueue`
- `github.com/arcgolabs/collectionx/interval` — `Range`, `RangeSet`, `RangeMap`
- `github.com/arcgolabs/collectionx/prefix` — `Trie` / `PrefixMap`
- `github.com/arcgolabs/collectionx/tree` — `Tree`, `ConcurrentTree`

## Documentation map

- First program (`Set` + `OrderedMap`): [Getting Started](./getting-started)
- Sets, ordered types, `MultiMap`, `Table`, JSON helpers: [Maps, sets, and tables](./mapping-recipes)
- Lists, intervals, trie, tree: [Lists and structured data](./structured-data)
- Release notes: [collectionx v0.1.3](./release-v0.1.3)
- Release notes: [collectionx v0.1.2](./release-v0.1.2)

## Install / Import

```bash
go get github.com/arcgolabs/collectionx/set@latest
go get github.com/arcgolabs/collectionx/mapping@latest
```

Import the **subpackage** you need (for example `collectionx/set`, `collectionx/mapping`, `collectionx/list`).

## Why use collectionx

Go’s standard library containers are intentionally minimal. `collectionx` focuses on generic, strongly typed APIs, explicit ordering guarantees where they matter, and shared engineering conventions across structures.

## Concurrency-safe types

Use **concurrent** variants only when the same instance is accessed from multiple goroutines:

- `ConcurrentSet`, `ConcurrentMap`, `ConcurrentMultiMap`, `ConcurrentTable`, `ConcurrentList`, `ConcurrentTree`

For single-goroutine use or external locking, prefer the non-concurrent types for lower overhead.

## API style notes

- Many `Values` / `All` / `Row` / `Column` style methods return **snapshots** to avoid accidental mutation leakage.
- `GetOption` helpers use `mo.Option` for nullable-style reads where applicable.
- Prefer constructors even when zero values work, for clarity.

## Serialization

Most structures can be passed directly to standard library serializers:

- `json.Marshal` / `json.Unmarshal`
- `gob.NewEncoder(...).Encode(...)` / `Decode(...)`
- `MarshalBinary` / `UnmarshalBinary`

You do not need to call a separate snapshot helper first. `ToJSON()` is still available when you explicitly want raw JSON bytes for logs, debugging, or custom transport code.

Typical usage:

```go
payload, err := json.Marshal(myCollection)
err = json.Unmarshal(payload, &myCollection)
```

Notes:

- `PriorityQueue` does not support automatic restore because its comparator is runtime configuration, not serialized data.
- `ShardedConcurrentMap` can be serialized directly, but to restore it you should initialize the receiver with `NewShardedConcurrentMap(...)` first so the hash function is available.

## Benchmarks

```bash
go test ./set -run ^$ -bench . -benchmem
go test ./mapping -run ^$ -bench . -benchmem
go test ./list -run ^$ -bench . -benchmem
```

Target one package:

```bash
go test ./mapping -run ^$ -bench . -benchmem
go test ./prefix -run ^$ -bench Trie -benchmem
```

## Practical tips

- Prefer `Table` over nested maps when keys are naturally two-dimensional.
- Prefer `OrderedMap` / `OrderedSet` when stable iteration order matters (tests, APIs, serialization).
- Prefer `Trie` for large prefix searches instead of repeated linear scans over string keys.
- Prefer `MultiSet` when frequency counts are the primary operation.
- Prefer `Tree` for parent/child models (org charts, categories, menus).

## FAQ

**Should I always use concurrent variants?**  
No. Use them only when multiple goroutines share the same instance without external synchronization.

**Are returned slices safe to mutate?**  
Snapshot-style APIs return copies; mutating them does not change internal state.

**Why does `OrderedMap` keep insertion order on value update?**  
By design: updates change values, not key order (similar to insertion-ordered maps elsewhere).

**How does `RangeSet` merge ranges?**  
Half-open ranges `[start, end)` are normalized; adjacent ranges merge (for example `[1,5)` + `[5,8)`).

## Troubleshooting

- **Non-deterministic order** — `Map` / `Set` hash iteration is unordered; use `OrderedMap` / `OrderedSet` when you need stable order.
- **`Trie.KeysWithPrefix` allocations** — returns new slices; narrow the prefix, use `RangePrefix` when available, or avoid building huge snapshots on hot paths.
- **Unbounded `MultiMap` / `Table` growth** — use `Delete`, `DeleteRow`, `DeleteColumn`, `DeleteValueIf`, or lifecycle-driven cleanup.

## Anti-patterns

- Defaulting to `Concurrent*` everywhere.
- Relying on hash-map iteration order in tests or business logic.
- Treating snapshot APIs as live views.
- Using `RangeMap` when a plain `map` and point lookups are enough.

## Integration guide

- **configx** — normalize loaded config into typed maps/lists before binding services.
- **clientx** / **kvx** — shape caches and indexes without one-off container code.
- **dix** — provide collection instances from module providers instead of hidden globals.

## Production notes

- Prefer the smallest structure that matches your invariants.
- Document ordering guarantees at API boundaries (`OrderedMap` vs hash-backed maps).
- For concurrent types, define ownership and lifecycle even when internals are lock-safe.
