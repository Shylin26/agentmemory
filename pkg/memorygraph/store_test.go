package memorygraph

import "testing"

func TestStorePutAndGet(t *testing.T) {

	store := NewStore()

	commit, err := NewCommit(
		nil,
		"Parisha",
		[]byte("Initial commit"),
	)
	if err != nil {
		t.Fatalf("failed to create commit: %v", err)
	}
	err = store.Put(commit)
	if err != nil {
		t.Fatalf("failed to store commit: %v", err)
	}
	retrieved, err := store.Get(commit.Hash)
	if err != nil {
		t.Fatalf("failed to retrieve commit: %v", err)
	}
	if retrieved.Author != commit.Author {
		t.Errorf("expected author %q, got %q", commit.Author, retrieved.Author)
	}
}
