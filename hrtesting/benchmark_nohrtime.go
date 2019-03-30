// +build nohrtime

package hrtesting

import (
	"testing"
)

// Benchmark implements minimal wrapper over *testing.B for disabling hrtesting.
type Benchmark struct {
	b *testing.B
	k int
}

// NewBenchmark creates a hrtime.Benchmark wrapper for *testing.B
func NewBenchmark(b *testing.B) *Benchmark {
	return &Benchmark{b: b, k: 0}
}

// Next starts measuring the next lap.
// It will return false, when all measurements have been made.
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

// Benchmark implements minimal wrapper over *testing.B for disabling hrtesting.
type BenchmarkTSC = Benchmark

// NewBenchmark creates a hrtime.BenchmarkTSC wrapper for *testing.B
func NewBenchmarkTSC(b *testing.B) *Benchmark {
	return NewBenchmark(b)
}
