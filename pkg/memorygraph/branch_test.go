package memorygraph

import (
	"fmt"
	"sync"
	"testing"
)

func TestBranchLifecycle(t *testing.T) {
	store := NewStore()
	rootCommit, err := NewCommit(
		nil,
		"Parisha",
		[]byte("Initial commit"),
	)
	if err != nil {
		t.Fatalf("failed to create commit :%v", err)
	}
	err = store.Put(rootCommit)
	if err != nil {
		t.Fatalf("failed to store root commmit :%v", err)
	}
	err = store.CreateBranch("main", rootCommit.Hash)
	if err != nil {
		t.Fatalf("failed to create a branch :%v", err)
	}
	branch, err := store.GetBranch("main")
	if err != nil {
		t.Fatalf("failed to get branch: %v", err)
	}

	if branch.Head != rootCommit.Hash {
		t.Errorf("expected head %q, got %q", rootCommit.Hash, branch.Head)
	}

	secondCommit, err := NewCommit(
		[]string{rootCommit.Hash},
		"Parisha",
		[]byte("Second commit"),
	)
	if err != nil {
		t.Fatalf("failed to create second commit: %v", err)
	}
	err = store.Put(secondCommit)
	if err != nil {
		t.Fatalf("failed to store second commit :%v", err)
	}
	err = store.UpdateBranch("main", secondCommit.Hash)
	if err != nil {
		t.Fatalf("failed to update branch :%v", err)

	}
	if branch.Head != secondCommit.Hash {
		t.Errorf("expected head %q, got %q", secondCommit.Hash, branch.Head)
	}

}

func TestCreateBranchInvalidCommit(t *testing.T) {
	store := NewStore()
	err := store.CreateBranch("main", "does-not-exists")
	if err == nil {
		t.Fatal("expected an error when creating a branch with a nonexistent commit ")
	}
}

func TestStore_ConcurrentBranchUpdates(t *testing.T) {
	store := NewStore()
	root, err := NewCommit(nil, "system", []byte("root"))
	if err != nil {
		t.Fatalf("failed to create root commit %v", err)
	}
	if err := store.Put(root); err != nil {
		t.Fatalf("failed to store root commit %v", err)
	}
	if err := store.CreateBranch("shared", root.Hash); err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}
	const numGoroutines = 50
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			commit, err := NewCommit([]string{root.Hash}, "system", []byte(fmt.Sprintf("update-%d", i)))
			if err != nil {
				errCh <- fmt.Errorf("goroutine %d: failed to create commit: %w", i, err)
				return
			}
			if err := store.Put(commit); err != nil {
				errCh <- fmt.Errorf("goroutine %d: failed to store commit: %w", i, err)
				return
			}
			if err := store.UpdateBranch("shared", commit.Hash); err != nil {
				errCh <- fmt.Errorf("goroutine %d: failed to update branch: %w", i, err)
				return
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}

	finalBranch, err := store.GetBranch("shared")
	if err != nil {
		t.Fatalf("failed to get final branch state: %v", err)
	}
	if _, err := store.Get(finalBranch.Head); err != nil {
		t.Errorf("final branch head does not point to a valid commit: %v", err)
	}
}
func TestListBranches(t *testing.T) {
	store := NewStore()

	root, err := NewCommit(nil, "system", []byte("root"))
	if err != nil {
		t.Fatalf("failed to create root commit: %v", err)
	}

	if err := store.Put(root); err != nil {
		t.Fatalf("failed to store root commit: %v", err)
	}

	if err := store.CreateBranch("main", root.Hash); err != nil {
		t.Fatalf("failed to create main branch: %v", err)
	}

	if err := store.CreateBranch("dev", root.Hash); err != nil {
		t.Fatalf("failed to create dev branch: %v", err)
	}

	if err := store.CreateBranch("feature", root.Hash); err != nil {
		t.Fatalf("failed to create feature branch: %v", err)
	}

	branches := store.ListBranches()

	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}
}
