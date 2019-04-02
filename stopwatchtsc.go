package hrtime

import (
	"sync"
	"sync/atomic"
	"time"
)

// SpanTSC defines a Count span
type SpanTSC struct {
	Start  Count
	Finish Count
}

// ApproxDuration returns the approximate duration of the span.
func (span *SpanTSC) ApproxDuration() time.Duration { return span.Count().ApproxDuration() }

// Count returns the duration in count of the count span.
func (span *SpanTSC) Count() Count { return span.Finish - span.Start }

// StopwatchTSC allows concurrent benchmarking using TSC
type StopwatchTSC struct {
	nextLap      int32
	lapsMeasured int32
	spans        []SpanTSC
	wait         sync.Mutex
}

// NewStopwatchTSC creates a new concurrent benchmark using TSC
func NewStopwatchTSC(count int) *StopwatchTSC {
	if count <= 0 {
		panic("must have count at least 1")
	}

	bench := &StopwatchTSC{
		nextLap: 0,
		spans:   make([]SpanTSC, count),
	}
	// lock mutex to ensure Wait() blocks until finalize is called
	bench.wait.Lock()
	return bench
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *StopwatchTSC) mustBeCompleted() {
	if int(atomic.LoadInt32(&bench.lapsMeasured)) < len(bench.spans) {
		panic("benchmarking incomplete")
	}
}

// Start starts measuring a new lap.
// It returns the lap number to pass in for Stop.
// It will return -1, when all measurements have been made.
//
// Call to Stop with -1 is ignored.
func (bench *StopwatchTSC) Start() int32 {
	lap := atomic.AddInt32(&bench.nextLap, 1) - 1
	if int(lap) > len(bench.spans) {
		return -1
	}
	bench.spans[lap].Start = TSC()
	return lap
}

// Stop stops measuring the specified lap.
//
// Call to Stop with -1 is ignored.
func (bench *StopwatchTSC) Stop(lap int32) {
	if lap < 0 {
		return
	}
	bench.spans[lap].Finish = TSC()

	lapsMeasured := atomic.AddInt32(&bench.lapsMeasured, 1)
	if int(lapsMeasured) == len(bench.spans) {
		bench.finalize()
	} else if int(lapsMeasured) > len(bench.spans) {
		panic("stop called too many times")
	}
}

// finalize finalizes the stopwatchTSC
func (bench *StopwatchTSC) finalize() {
	// release the initial lock such that Wait can proceed.
	bench.wait.Unlock()
}

// Wait waits for all measurements to be completed.
func (bench *StopwatchTSC) Wait() {
	// lock waits for finalize to be called by the last measurement.
	bench.wait.Lock()
	_ = 1 // intentionally empty block, suppress staticcheck SA2001 warning
	bench.wait.Unlock()
}

// Spans returns measured time-spans.
func (bench *StopwatchTSC) Spans() []SpanTSC {
	bench.mustBeCompleted()
	return append(bench.spans[:0:0], bench.spans...)
}

// ApproxDurations returns measured durations.
func (bench *StopwatchTSC) ApproxDurations() []time.Duration {
	bench.mustBeCompleted()

	durations := make([]time.Duration, len(bench.spans))
	for i, span := range bench.spans {
		durations[i] = span.ApproxDuration()
	}

	return durations
}

// Name returns name of the benchmark.
func (bench *StopwatchTSC) Name() string { return "" }

// Unit returns units it measures.
func (bench *StopwatchTSC) Unit() string { return "tsc" }

// Float64s returns all measurements.
func (bench *StopwatchTSC) Float64s() []float64 {
	measurements := make([]float64, len(bench.spans))
	for i := range measurements {
		measurements[i] = float64(bench.spans[i].Count())
	}
	return measurements
}

// Histogram creates an histogram of all the durations.
//
// It creates binCount bins to distribute the data and uses the
// 99.9 percentile as the last bucket range. However, for a nicer output
// it might choose a larger value.
func (bench *StopwatchTSC) Histogram(binCount int) *Histogram {
	bench.mustBeCompleted()

	opts := defaultOptions
	opts.BinCount = binCount

	return NewDurationHistogram(bench.ApproxDurations(), &opts)
}

// HistogramClamp creates an historgram of all the durations clamping minimum and maximum time.
//
// It creates binCount bins to distribute the data and uses the
// maximum as the last bucket.
func (bench *StopwatchTSC) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	bench.mustBeCompleted()

	durations := make([]time.Duration, 0, len(bench.spans))
	for _, span := range bench.spans {
		duration := span.ApproxDuration()
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
