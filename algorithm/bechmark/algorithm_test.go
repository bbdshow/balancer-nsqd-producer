package bechmark

import (
	"balancer-nsqd-producer/algorithm"
	"testing"
)

// go test -bench=. -benchmem -run=none
func BenchmarkPolling(b *testing.B) {
	polling := algorithm.NewPolling()
	count := 20
	for count > 0 {
		polling.Put(count)
		count--
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		polling.Get()
	}
}

func BenchmarkRandom(b *testing.B) {
	random := algorithm.NewRandom()
	count := 20
	for count > 0 {
		random.Put(count)
		count--
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.Get()
	}
}

func BenchmarkSmoothWeight(b *testing.B) {
	weight := algorithm.NewSmoothWeight()
	count := 20
	for count > 0 {
		weight.Put(count, count+2)
		count--
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		weight.Get()
	}
}
