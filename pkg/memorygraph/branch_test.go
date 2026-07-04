package memorygraph

import "testing"

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
