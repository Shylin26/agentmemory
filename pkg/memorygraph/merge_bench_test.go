package memorygraph

import (
	"fmt"
	"testing"
)

func BenchmarkMerge(b *testing.B) {
	store := NewStore()

	root, err := NewCommit(nil, "system", []byte("root"))
	if err != nil {
		b.Fatalf("failed to create root commit: %v", err)
	}
	if err := store.Put(root); err != nil {
		b.Fatalf("failed to store root commit: %v", err)
	}

	if err := store.CreateBranch("branch-a", root.Hash); err != nil {
		b.Fatalf("failed to create branch-a: %v", err)
	}
	if err := store.CreateBranch("branch-b", root.Hash); err != nil {
		b.Fatalf("failed to create branch-b: %v", err)
	}

	commitA, err := NewCommit([]string{root.Hash}, "system", []byte("payload-a"))
	if err != nil {
		b.Fatalf("failed to create commitA: %v", err)
	}
	if err := store.Put(commitA); err != nil {
		b.Fatalf("failed to store commitA: %v", err)
	}
	if err := store.UpdateBranch("branch-a", commitA.Hash); err != nil {
		b.Fatalf("failed to update branch-a: %v", err)
	}

	commitB, err := NewCommit([]string{root.Hash}, "system", []byte("payload-b"))
	if err != nil {
		b.Fatalf("failed to create commitB: %v", err)
	}
	if err := store.Put(commitB); err != nil {
		b.Fatalf("failed to store commitB: %v", err)
	}
	if err := store.UpdateBranch("branch-b", commitB.Hash); err != nil {
		b.Fatalf("failed to update branch-b: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := store.Merge("branch-a", "branch-b", "bench-bot")
		if err != nil {
			b.Fatalf("merge failed: %v", err)
		}
	}
}
func BenchmarkMerge_WithExistingBranches(b *testing.B) {
	for _, branchCount := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("branches-%d", branchCount), func(b *testing.B) {
			store := NewStore()

			root, err := NewCommit(nil, "system", []byte("root"))
			if err != nil {
				b.Fatalf("failed to create root commit: %v", err)
			}
			if err := store.Put(root); err != nil {
				b.Fatalf("failed to store root commit: %v", err)
			}

			for i := 0; i < branchCount; i++ {
				name := fmt.Sprintf("noise-branch-%d", i)
				if err := store.CreateBranch(name, root.Hash); err != nil {
					b.Fatalf("failed to create noise branch: %v", err)
				}
			}

			if err := store.CreateBranch("branch-a", root.Hash); err != nil {
				b.Fatalf("failed to create branch-a: %v", err)
			}
			if err := store.CreateBranch("branch-b", root.Hash); err != nil {
				b.Fatalf("failed to create branch-b: %v", err)
			}

			commitA, err := NewCommit([]string{root.Hash}, "system", []byte("payload-a"))
			if err != nil {
				b.Fatalf("failed to create commitA: %v", err)
			}
			if err := store.Put(commitA); err != nil {
				b.Fatalf("failed to store commitA: %v", err)
			}
			if err := store.UpdateBranch("branch-a", commitA.Hash); err != nil {
				b.Fatalf("failed to update branch-a: %v", err)
			}

			commitB, err := NewCommit([]string{root.Hash}, "system", []byte("payload-b"))
			if err != nil {
				b.Fatalf("failed to create commitB: %v", err)
			}
			if err := store.Put(commitB); err != nil {
				b.Fatalf("failed to store commitB: %v", err)
			}
			if err := store.UpdateBranch("branch-b", commitB.Hash); err != nil {
				b.Fatalf("failed to update branch-b: %v", err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := store.Merge("branch-a", "branch-b", "bench-bot")
				if err != nil {
					b.Fatalf("merge failed: %v", err)
				}
			}
		})
	}
}
