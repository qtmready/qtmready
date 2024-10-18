package ws

import (
	"testing"
)

func BenchmarkIdempotent(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Idempotent()
	}
}

func BenchmarkIdempotentCollision(b *testing.B) {
	b.ResetTimer()

	seen := make(map[string]bool)

	for i := 0; i < b.N; i++ {
		id := Idempotent()
		if seen[id] {
			b.Fatalf("Collision detected: %s", id)
		}

		seen[id] = true
	}
}

func BenchmarkIdempotentParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Idempotent()
		}
	})
}
