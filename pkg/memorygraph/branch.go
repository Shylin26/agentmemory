package memorygraph

import "fmt"

type Branch struct {
	Name string
	Head string
}

func (s *Store) CreateBranch(name string, headHash string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.branches[name]; exists {
		return fmt.Errorf("branch already exists: %s", name)
	}

	if _, exists := s.commits[headHash]; !exists {
		return fmt.Errorf("cannot create branch pointing to unknown commit: %s", headHash)
	}

	s.branches[name] = &Branch{
		Name: name,
		Head: headHash,
	}
	return nil
}

func (s *Store) UpdateBranch(name string, newHeadHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.branches[name]; !exists {
		return fmt.Errorf("branch does not exists:%s", name)
	}
	if _, exists := s.commits[newHeadHash]; !exists {
		return fmt.Errorf("cannot update branch to unknown commit :%s", newHeadHash)

	}
	s.branches[name].Head = newHeadHash
	return nil
}
func (s *Store) GetBranch(name string) (*Branch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	branch, ok := s.branches[name]
	if !ok {
		return nil, fmt.Errorf("branch not found :%s", name)
	}
	return branch, nil
}
func (s *Store) ListBranches() []*Branch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	branches := make([]*Branch, 0, len(s.branches))
	for _, b := range s.branches {
		branches = append(branches, b)
	}

	return branches

}
