package memorygraph

import "testing"

func BenchmarkNewCommit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewCommit(nil, "bench-agent", []byte("benchmark payload"))
		if err != nil {
			b.Fatalf("NewCommit failed: %v", err)
		}
	}
}
