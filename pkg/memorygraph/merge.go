package memorygraph

import (
	"bytes"
	"fmt"
)

type MergeResult struct {
	Commit     *Commit
	Conflicted bool
	Discarded  []byte
}

func (s *Store) Merge(branchA string, branchB string, author string) (*MergeResult, error) {
	bA, err := s.GetBranch(branchA)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch %s:%w", branchA, err)
	}
	bB, err := s.GetBranch(branchB)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch %s: %w", branchB, err)
	}
	commitA, err := s.Get(bA.Head)
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit for %s: %w", branchA, err)
	}
	commitB, err := s.Get(bB.Head)
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit for %s: %w", branchB, err)
	}
	var winner, loser *Commit
	if commitA.Timestamp.After(commitB.Timestamp) {
		winner = commitA
		loser = commitB
	} else {
		winner = commitB
		loser = commitA
	}
	conflicted := !bytes.Equal(commitA.Payload, commitB.Payload)
	mergeCommit, err := NewCommit(
		[]string{commitA.Hash, commitB.Hash},
		author,
		winner.Payload,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create merge commit: %w", err)
	}
	if err := s.Put(mergeCommit); err != nil {
		return nil, fmt.Errorf("failed to store merge commit: %w", err)
	}
	result := &MergeResult{
		Commit:     mergeCommit,
		Conflicted: conflicted,
	}
	if conflicted {
		result.Discarded = loser.Payload

	}
	return result, nil

}
