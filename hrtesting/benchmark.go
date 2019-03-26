package hrtesting

import (
	"testing"

	"github.com/loov/hrtime"
)

type B struct {
	hr hrtime.Benchmark
	b  *testing.B
}

func NewBenchmark(b *testing.B) *B {
	bench := &B{
		hr: *hrtime.NewBenchmark(b.N),
		b:  b,
	}
	bench.b.ResetTimer()
	return bench
}

func (bench *B) Next() bool {
	result := bench.hr.Next()
	if !result {
		bench.b.StopTimer()
	}
	return result
}

func (bench *B) Print() {
	bench.b.Log(bench.hr.Histogram(10))
}
