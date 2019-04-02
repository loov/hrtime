package hrtime

import (
	"sync"
	"sync/atomic"
	"time"
)

// Span defines a time.Duration span
type Span struct {
	Start  time.Duration
	Finish time.Duration
}

// Duration returns the duration of the time span.
func (span *Span) Duration() time.Duration {
	return span.Finish - span.Start
}

// Stopwatch allows concurrent benchmarking using Now
type Stopwatch struct {
	nextLap      int32
	lapsMeasured int32
	spans        []Span
	wait         sync.Mutex
}

// NewStopwatch creates a new concurrent benchmark using Now
func NewStopwatch(count int) *Stopwatch {
	if count <= 0 {
		panic("must have count at least 1")
	}

	bench := &Stopwatch{
		nextLap: 0,
		spans:   make([]Span, count),
	}
	// lock mutex to ensure Wait() blocks until finalize is called
	bench.wait.Lock()
	return bench
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *Stopwatch) mustBeCompleted() {
	if int(atomic.LoadInt32(&bench.lapsMeasured)) < len(bench.spans) {
		panic("benchmarking incomplete")
	}
}

// Start starts measuring a new lap.
// It returns the lap number to pass in for Stop.
// It will return -1, when all measurements have been made.
//
// Call to Stop with -1 is ignored.
func (bench *Stopwatch) Start() int32 {
	lap := atomic.AddInt32(&bench.nextLap, 1) - 1
	if int(lap) > len(bench.spans) {
		return -1
	}
	bench.spans[lap].Start = Now()
	return lap
}

// Stop stops measuring the specified lap.
//
// Call to Stop with -1 is ignored.
func (bench *Stopwatch) Stop(lap int32) {
	if lap < 0 {
		return
	}
	bench.spans[lap].Finish = Now()

	lapsMeasured := atomic.AddInt32(&bench.lapsMeasured, 1)
	if int(lapsMeasured) == len(bench.spans) {
		bench.finalize()
	} else if int(lapsMeasured) > len(bench.spans) {
		panic("stop called too many times")
	}
}

// finalize finalizes the stopwatch
func (bench *Stopwatch) finalize() {
	// release the initial lock such that Wait can proceed.
	bench.wait.Unlock()
}

// Wait waits for all measurements to be completed.
func (bench *Stopwatch) Wait() {
	// lock waits for finalize to be called by the last measurement.
	bench.wait.Lock()
	_ = 1 // intentionally empty block, suppress staticcheck SA2001 warning
	bench.wait.Unlock()
}

// Spans returns measured time-spans.
func (bench *Stopwatch) Spans() []Span {
	bench.mustBeCompleted()
	return append(bench.spans[:0:0], bench.spans...)
}

// Durations returns measured durations.
func (bench *Stopwatch) Durations() []time.Duration {
	bench.mustBeCompleted()

	durations := make([]time.Duration, len(bench.spans))
	for i, span := range bench.spans {
		durations[i] = span.Duration()
	}

	return durations
}

// Name returns name of the benchmark.
func (bench *Stopwatch) Name() string { return "" }

// Unit returns units it measures.
func (bench *Stopwatch) Unit() string { return "ns" }

// Float64s returns all measurements.
func (bench *Stopwatch) Float64s() []float64 {
	measurements := make([]float64, len(bench.spans))
	for i := range measurements {
		measurements[i] = float64(bench.spans[i].Duration().Nanoseconds())
	}
	return measurements
}

// Histogram creates an histogram of all the durations.
//
// It creates binCount bins to distribute the data and uses the
// 99.9 percentile as the last bucket range. However, for a nicer output
// it might choose a larger value.
func (bench *Stopwatch) Histogram(binCount int) *Histogram {
	bench.mustBeCompleted()

	opts := defaultOptions
	opts.BinCount = binCount

	return NewDurationHistogram(bench.Durations(), &opts)
}

// HistogramClamp creates an historgram of all the durations clamping minimum and maximum time.
//
// It creates binCount bins to distribute the data and uses the
// maximum as the last bucket.
func (bench *Stopwatch) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	bench.mustBeCompleted()

	durations := make([]time.Duration, 0, len(bench.spans))
	for _, span := range bench.spans {
		duration := span.Duration()
		if duration < min {
			durations = append(durations, min)
		} else {
			durations = append(durations, duration)
		}
	}

	opts := defaultOptions
	opts.BinCount = binCount
	opts.ClampMaximum = float64(max.Nanoseconds())
	opts.ClampPercentile = 0

	return NewDurationHistogram(durations, &opts)
}
