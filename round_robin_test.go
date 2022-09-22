package zcache

import (
	"runtime"
	"testing"
)

func BenchmarkRoundRobin(b *testing.B) {
	rr := NewRoundRobin[int, int](100000)

	for i := 0; i < b.N; i++ {
		rr.Set(i, i)
	}
}

func BenchmarkBucketRoundRobin(b *testing.B) {

	bucket := NewBucketRoundRobin[int, int](100000, runtime.NumCPU(), func(i int) uint32 { return uint32(i) })

	for i := 0; i < b.N; i++ {
		bucket.Element(i).Set(i, i)
	}

}
