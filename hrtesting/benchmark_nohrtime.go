// +build nohrtime

package hrtesting

import (
	"testing"
)

// Benchmark implements minimal wrapper over *testing.Benchmark for disabling hrtesting.
type Benchmark struct {
	b *testing.B
	k int
}

// NewBenchmark creates a hrtime.Benchmark wrapper for *testing.B
func NewBenchmark(b *testing.Benchmark) *Benchmark {
	return &Benchmark{b: b, k: 0}
}

// Next is used to loop through all the tests.
func (bench *Benchmark) Next() bool {
	bench.b.StartTimer()
	bench.k++
	next := bench.k <= bench.b.N
	if !next {
		bench.b.StopTimer()
	}
	return next
}

// Report reports the result to the console.
func (bench *Benchmark) Report() {}

// Benchmark implements minimal wrapper over *testing.Benchmark for disabling hrtesting.
type BenchmarkTSC = Benchmark

// NewBenchmark creates a hrtime.BenchmarkTSC wrapper for *testing.B
func NewBenchmarkTSC(b *testing.B) *Benchmark {
	return NewBenchmark(b)
}
