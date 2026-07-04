package memorygraph

import (
	"fmt"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	commits map[string]*Commit
}

func NewStore() *Store {
	return &Store{
		commits: make(map[string]*Commit),
	}
}

func (s *Store) Put(commit *Commit) error {
	if commit == nil {
		return fmt.Errorf("cannot store a nil comment")

	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.commits[commit.Hash] = commit
	return nil
}

func (s *Store) Get(hash string) (*Commit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	commit, ok := s.commits[hash]
	if !ok {
		return nil, fmt.Errorf("commit not found:%s", hash)

	}
	return commit, nil
}
