package memorygraph

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Commit struct {
	Hash       string    `json:"hash"`
	ParentHash []string  `json:"parent_hash"`
	Author     string    `json:"author"`
	Timestamp  time.Time `json:"timestamp"`
	Payload    []byte    `json:"payload"`
}

func computeHash(parentHash []string, author string, timestamp time.Time, payload []byte) (string, error) {
	hashInput := struct {
		ParentHash []string  `json:"parent_hash"`
		Author     string    `json:"author"`
		Timestamp  time.Time `json:"timestamp"`
		Payload    []byte    `json:"payload"`
	}{
		ParentHash: parentHash,
		Author:     author,
		Timestamp:  timestamp,
		Payload:    payload,
	}
	data, err := json.Marshal(hashInput)
	if err != nil {
		return "", fmt.Errorf("failed to marshal commit for hashing:%w", err)
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum), nil
}

func NewCommit(parentHash []string, author string, payload []byte) (*Commit, error) {
	timestamp := time.Now().UTC()
	hash, err := computeHash(parentHash, author, timestamp, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create commit:%w", err)

	}
	return &Commit{
		Hash:       hash,
		ParentHash: parentHash,
		Author:     author,
		Timestamp:  timestamp,
		Payload:    payload,
	}, nil
}
