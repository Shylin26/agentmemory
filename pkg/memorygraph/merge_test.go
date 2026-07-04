package memorygraph

import (
	"testing"
	"time"
)

func TestMerge_Conflicted(t *testing.T) {
	store := NewStore()
	root, err := NewCommit(nil, "system", []byte("root"))
	if err != nil {
		t.Fatalf("failed to create root commit :%v", err)
	}
	if err := store.Put(root); err != nil {
		t.Fatalf("failed to store root commit :%v", err)
	}
	if err := store.CreateBranch("branch-a", root.Hash); err != nil {
		t.Fatalf("failed to create branch-a:%v", err)
	}
	if err := store.CreateBranch("branch-b", root.Hash); err != nil {
		t.Fatalf("failed to create branch-b: %v", err)
	}
	commitA, err := NewCommit(
		[]string{root.Hash},
		"system",
		[]byte("price is rising"),
	)
	if err != nil {
		t.Fatalf("failed to create CommitA:%v", err)
	}
	if err := store.Put(commitA); err != nil {
		t.Fatalf("failed to store commitA:%v", err)
	}
	if err := store.UpdateBranch("branch-a", commitA.Hash); err != nil {
		t.Fatalf("failed to update branch-a:%v", err)
	}
	time.Sleep(time.Millisecond)
	commitB, err := NewCommit(
		[]string{root.Hash},
		"system",
		[]byte("price is falling"),
	)
	if err != nil {
		t.Fatalf("failed to create commitB: %v", err)
	}

	if err := store.Put(commitB); err != nil {
		t.Fatalf("failed to store commitB: %v", err)
	}

	if err := store.UpdateBranch("branch-b", commitB.Hash); err != nil {
		t.Fatalf("failed to update branch-b: %v", err)
	}
	result, err := store.Merge("branch-a", "branch-b", "merge-bot")
	if err != nil {
		t.Fatalf("merge failed :%v", err)
	}
	if !result.Conflicted {
		t.Errorf("expected winner payload %q, got %q", "price is falling", result.Commit.Payload)
	}
	if string(result.Discarded) != "price is rising" {
		t.Errorf("expected discarded payload %q, got %q", "price is rising", result.Discarded)
	}

	if len(result.Commit.ParentHash) != 2 {
		t.Fatalf("expected 2 parents, got %d", len(result.Commit.ParentHash))
	}

}
func TestMerge_NoConflict(t *testing.T) {
	store := NewStore()

	root, err := NewCommit(nil, "system", []byte("root"))
	if err != nil {
		t.Fatalf("failed to create root commit: %v", err)
	}

	if err := store.Put(root); err != nil {
		t.Fatalf("failed to store root commit: %v", err)
	}

	if err := store.CreateBranch("branch-a", root.Hash); err != nil {
		t.Fatalf("failed to create branch-a: %v", err)
	}

	if err := store.CreateBranch("branch-b", root.Hash); err != nil {
		t.Fatalf("failed to create branch-b: %v", err)
	}

	commitA, err := NewCommit(
		[]string{root.Hash},
		"system",
		[]byte("stable price"),
	)
	if err != nil {
		t.Fatalf("failed to create commitA: %v", err)
	}

	if err := store.Put(commitA); err != nil {
		t.Fatalf("failed to store commitA: %v", err)
	}

	if err := store.UpdateBranch("branch-a", commitA.Hash); err != nil {
		t.Fatalf("failed to update branch-a: %v", err)
	}

	time.Sleep(time.Millisecond)

	commitB, err := NewCommit(
		[]string{root.Hash},
		"system",
		[]byte("stable price"),
	)
	if err != nil {
		t.Fatalf("failed to create commitB: %v", err)
	}

	if err := store.Put(commitB); err != nil {
		t.Fatalf("failed to store commitB: %v", err)
	}

	if err := store.UpdateBranch("branch-b", commitB.Hash); err != nil {
		t.Fatalf("failed to update branch-b: %v", err)
	}

	result, err := store.Merge("branch-a", "branch-b", "merge-bot")
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	if result.Conflicted {
		t.Fatal("expected merge without conflict")
	}

	if len(result.Discarded) != 0 {
		t.Fatalf("expected no discarded payload, got %q", result.Discarded)
	}

	if string(result.Commit.Payload) != "stable price" {
		t.Fatalf("expected payload %q, got %q", "stable price", result.Commit.Payload)
	}

	if len(result.Commit.ParentHash) != 2 {
		t.Fatalf("expected 2 parents, got %d", len(result.Commit.ParentHash))
	}
}
