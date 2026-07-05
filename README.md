# AgentMemory

**AgentMemory** is a Git-inspired memory store for AI agents — commit,
branch, and merge evolving agent knowledge with full version history,
instead of overwriting it.

## Why

Most multi-agent systems manage shared state with a single mutable store:
every write overwrites whatever was there before. That means there's no way
to inspect how a shared belief evolved, branch a hypothesis without
corrupting the shared truth, or recover cleanly when two agents write
contradictory conclusions at the same time.

AgentMemory treats agent memory the way Git treats source code: every write
is an immutable, content-addressed commit; branches are named, movable
pointers over commit history; and merges combine divergent branches with an
explicit, auditable conflict-resolution strategy — never a silent
overwrite.

## How it works

- **Commit** — an immutable record: a SHA-256 hash of its content, one or
  more parent hashes, an author, a timestamp, and a payload.
- **Store** — concurrency-safe storage for commits and branches, built on
  Go's `sync.RWMutex` for safe concurrent reads and writes.
- **Branch** — a named pointer to the current "tip" commit, movable via
  `UpdateBranch`.
- **Merge** — combines two branches into a new commit with both heads as
  parents. Conflicts are resolved with last-write-wins in v1 (see
  [docs/0001](docs/0001-deferred-semantic-merge.md) for why, and what's
  next).

## Performance

Benchmarked on an Apple M4:

| Operation | Latency |
|---|---|
| `NewCommit` | ~1,039 ns/op |
| `Merge` (2 branches) | ~1,923 ns/op |
| `Merge`, store with 1,000 unrelated branches | ~2,041 ns/op |

Merge latency stays effectively flat as unrelated store size grows 100x —
consistent with Go's hash-map-backed storage rather than any accidental
linear scan. See `pkg/memorygraph/*_bench_test.go` for the full benchmark
suite.

## Usage

```go
store := memorygraph.NewStore()

root, _ := memorygraph.NewCommit(nil, "agent-1", []byte("initial belief"))
store.Put(root)
store.CreateBranch("main", root.Hash)

next, _ := memorygraph.NewCommit([]string{root.Hash}, "agent-1", []byte("updated belief"))
store.Put(next)
store.UpdateBranch("main", next.Hash)

branch, _ := store.GetBranch("main")
fmt.Println(branch.Head) // points at `next`
```

## Status

Phase 1 (this repo): commit/branch/merge core, concurrency-safe, benchmarked.
Next: a cache-aware inference scheduler that uses branch topology to route
and batch concurrent agent calls — coming in a future repo.

## Testing

```bash
go test ./...          # run all tests
go test -race ./...    # run with the race detector
go test -bench=. -run=^$ ./...  # run benchmarks
```