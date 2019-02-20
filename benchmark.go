package hrtime

import (
	"time"
)

// Benchmark helps benchmarking using time.
type Benchmark struct {
	step  int
	laps  []time.Duration
	start time.Duration
	stop  time.Duration
}

// NewBenchmark creates a new benchmark using time.
// Count defines the number of samples to measure.
func NewBenchmark(count int) *Benchmark {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &Benchmark{
		step:  0,
		laps:  make([]time.Duration, count),
		start: 0,
		stop:  0,
	}
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *Benchmark) mustBeCompleted() {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}
}

// finalize calculates diffs for each lap.
func (bench *Benchmark) finalize(last time.Duration) {
	if bench.stop != 0 {
		return
	}

	bench.start = bench.laps[0]
	bench.stop = last
	for i := range bench.laps[:len(bench.laps)-1] {
		bench.laps[i] = bench.laps[i+1] - bench.laps[i]
	}
	bench.laps[len(bench.laps)-1] = bench.stop - bench.laps[len(bench.laps)-1]
}

// Next starts measuring the next lap.
// It will return false, when all measurements have been made.
func (bench *Benchmark) Next() bool {
	now := Now()
	if bench.step >= len(bench.laps) {
		bench.finalize(now)
		return false
	}
	bench.laps[bench.step] = Now()
	bench.step++
	return true
}

// Laps returns timing for each lap.
func (bench *Benchmark) Laps() []time.Duration {
	bench.mustBeCompleted()
	return append(bench.laps[:0:0], bench.laps...)
}

// Histogram creates an histogram of all the laps.
func (bench *Benchmark) Histogram(binCount int) *Histogram {
	bench.mustBeCompleted()

	opts := defaultOptions
	opts.BinCount = binCount

	return NewDurationHistogram(bench.laps, &opts)
}

// HistogramClamp creates an historgram of all the laps clamping minimum and maximum time.
func (bench *Benchmark) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	bench.mustBeCompleted()

	laps := make([]time.Duration, 0, len(bench.laps))
	for _, lap := range bench.laps {
		if lap < min {
			laps = append(laps, min)
		} else {
			laps = append(laps, lap)
		}
	}

	opts := defaultOptions
	opts.BinCount = binCount
	opts.ClampMaximum = float64(max.Nanoseconds())
	opts.ClampPercentile = 0

	return NewDurationHistogram(laps, &opts)
}
