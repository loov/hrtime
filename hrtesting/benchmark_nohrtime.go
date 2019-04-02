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

// Name returns benchmark name.
func (bench *Benchmark) Name() string { return bench.b.Name() }

// Unit returns units it measures.
func (bench *Benchmark) Unit() string { return "" }

// Float64s returns all measurements as float64s
func (bench *BenchmarkTSC) Float64s() []float64 { return nil }

// Benchmark implements minimal wrapper over *testing.B for disabling hrtesting.
type BenchmarkTSC = Benchmark

// NewBenchmark creates a hrtime.BenchmarkTSC wrapper for *testing.B
func NewBenchmarkTSC(b *testing.B) *Benchmark {
	return NewBenchmark(b)
}
