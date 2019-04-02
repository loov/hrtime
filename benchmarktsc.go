package hrtime

import (
	"time"
)

// BenchmarkTSC helps benchmarking using CPU counters.
type BenchmarkTSC struct {
	step   int
	counts []Count
	start  Count
	stop   Count
}

// NewBenchmarkTSC creates a new benchmark using CPU counters.
// Count defines the number of samples to measure.
func NewBenchmarkTSC(count int) *BenchmarkTSC {
	if count <= 0 {
		panic("must have count at least 1")
	}

	return &BenchmarkTSC{
		step:   0,
		counts: make([]Count, count),
		start:  0,
		stop:   0,
	}
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *BenchmarkTSC) mustBeCompleted() {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}
}

// finalize calculates diffs for each lap.
func (bench *BenchmarkTSC) finalize(last Count) {
	if bench.stop != 0 {
		return
	}

	bench.start = bench.counts[0]
	bench.stop = last
	for i := range bench.counts[:len(bench.counts)-1] {
		bench.counts[i] = bench.counts[i+1] - bench.counts[i]
	}
	bench.counts[len(bench.counts)-1] = bench.stop - bench.counts[len(bench.counts)-1]
}

// Next starts measuring the next lap.
// It will return false, when all measurements have been made.
func (bench *BenchmarkTSC) Next() bool {
	now := TSC()
	if bench.step >= len(bench.counts) {
		bench.finalize(now)
		return false
	}
	bench.counts[bench.step] = TSC()
	bench.step++
	return true
}

// Counts returns counts for each lap.
func (bench *BenchmarkTSC) Counts() []Count {
	bench.mustBeCompleted()

	return append(bench.counts[:0:0], bench.counts...)
}

// Laps returns timing for each lap using the approximate conversion of Count.
func (bench *BenchmarkTSC) Laps() []time.Duration {
	bench.mustBeCompleted()

	laps := make([]time.Duration, len(bench.counts))
	for i, v := range bench.counts {
		laps[i] = v.ApproxDuration()
	}
	return laps
}

// Name returns name of the benchmark.
func (bench *BenchmarkTSC) Name() string { return "" }

// Unit returns units it measures.
func (bench *BenchmarkTSC) Unit() string { return "tsc" }

// Float64s returns all measurements as float64s
func (bench *BenchmarkTSC) Float64s() []float64 {
	measurements := make([]float64, len(bench.counts))
	for i := range measurements {
		measurements[i] = float64(bench.counts[i])
	}
	return measurements
}

// Histogram creates an histogram of all the laps.
//
// It creates binCount bins to distribute the data and uses the
// 99.9 percentile as the last bucket range. However, for a nicer output
// it might choose a larger value.
func (bench *BenchmarkTSC) Histogram(binCount int) *Histogram {
	bench.mustBeCompleted()

	opts := defaultOptions
	opts.BinCount = binCount

	return NewDurationHistogram(bench.Laps(), &opts)
}

// HistogramClamp creates an historgram of all the laps clamping minimum and maximum time.
//
// It creates binCount bins to distribute the data and uses the
// maximum as the last bucket.
func (bench *BenchmarkTSC) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	bench.mustBeCompleted()

	laps := make([]time.Duration, 0, len(bench.counts))
	for _, count := range bench.counts {
		lap := count.ApproxDuration()
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
