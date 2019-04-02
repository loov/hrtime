// +build !nohrtime

package hrtesting

import (
	"math"
	"testing"

	"github.com/loov/hrtime"
)

// Benchmark wraps *testing.B to measure more details using hrtime.Benchmark
type Benchmark struct {
	hr hrtime.Benchmark
	b  *testing.B
}

// NewBenchmark creates a new *hrtime.Benchmark using *testing.B parameters.
func NewBenchmark(b *testing.B) *Benchmark {
	bench := &Benchmark{
		hr: *hrtime.NewBenchmark(b.N),
		b:  b,
	}
	bench.b.StopTimer()
	return bench
}

// Next starts measuring the next lap.
// It will return false, when all measurements have been made.
func (bench *Benchmark) Next() bool {
	bench.b.StartTimer()
	result := bench.hr.Next()
	if !result {
		bench.b.StopTimer()
	}
	return result
}

// Name returns benchmark name.
func (bench *Benchmark) Name() string { return bench.b.Name() }

// Unit returns units it measures.
func (bench *Benchmark) Unit() string { return bench.hr.Unit() }

// Float64s returns all measurements.
func (bench *Benchmark) Float64s() []float64 { return bench.hr.Float64s() }

// BenchmarkTSC wraps *testing.B to measure more details using hrtime.BenchmarkTSC
type BenchmarkTSC struct {
	hr hrtime.BenchmarkTSC
	b  *testing.B
}

// NewBenchmarkTSC creates a new *hrtime.BenchmarkTSC using *testing.B parameters.
func NewBenchmarkTSC(b *testing.B) *BenchmarkTSC {
	bench := &BenchmarkTSC{
		hr: *hrtime.NewBenchmarkTSC(b.N),
		b:  b,
	}
	bench.b.StopTimer()
	return bench
}

// Name returns benchmark name.
func (bench *BenchmarkTSC) Name() string { return bench.b.Name() }

// Next starts measuring the next lap.
// It will return false, when all measurements have been made.
func (bench *BenchmarkTSC) Next() bool {
	bench.b.StartTimer()
	result := bench.hr.Next()
	if !result {
		bench.b.StopTimer()
	}
	return result
}

// Unit returns units it measures.
func (bench *BenchmarkTSC) Unit() string { return bench.hr.Unit() }

// Float64s returns all measurements as float64s
func (bench *BenchmarkTSC) Float64s() []float64 { return bench.hr.Float64s() }

func truncate(v float64, digits int) float64 {
	if digits == 0 {
		return 0
	}

	scale := math.Pow(10, math.Floor(math.Log10(v))+1-float64(digits))
	return scale * math.Trunc(v/scale)
}
