package memorygraph

import (
	"testing"
)

func TestNewCommit_ProducesValidHash(t *testing.T) {
	commit, err := NewCommit(nil, "agent-1", []byte("hello world"))
	if err != nil {
		t.Fatalf("NewCommit returned an error:%v", err)
	}
	if commit.Hash == "" {
		t.Error("expected a non-empty hash,got empty string")

	}
	if commit.Author != "agent-1" {
		t.Errorf("expected author 'agent-1',got %q", commit.Author)
	}
}
