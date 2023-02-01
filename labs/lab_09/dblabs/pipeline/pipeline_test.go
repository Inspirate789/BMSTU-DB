package main

import (
	"testing"
)

func BenchmarkRunSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		requests := make([]Request, requestCount)
		for j := 0; j < requestCount; j++ {
			requests[j].ID = j
		}

		p := Pipeline{[]Line{connectDB, insertTerm, selectTerm, deleteTerm, disconnectDB}}
		b.StartTimer()

		_ = p.RunSequential(requests)
	}
}

//func BenchmarkRunSync100(b *testing.B) {
//	benchmarkRunSync(100, b)
//}

func BenchmarkRunParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		requests := make([]Request, requestCount)
		for j := 0; j < requestCount; j++ {
			requests[j].ID = j
		}

		p := Pipeline{[]Line{connectDB, insertTerm, selectTerm, deleteTerm, disconnectDB}}
		b.StartTimer()

		_ = p.RunParallel(requests)
	}
}
