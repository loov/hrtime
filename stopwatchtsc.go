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

// Duration returns the duration of the count span.
func (span *SpanTSC) Count() Count { return span.Finish - span.Start }

// StopwatchTSC allows concurrent benchmarking using Now
type StopwatchTSC struct {
	lap   int32
	done  int32
	spans []SpanTSC
	wait  sync.Mutex
}

// NewStopwatchTSC creates a new concurrent benchmark using Now
func NewStopwatchTSC(count int) *StopwatchTSC {
	if count <= 0 {
		panic("must have count at least 0")
	}

	stop := &StopwatchTSC{
		lap:   0,
		spans: make([]SpanTSC, count),
	}
	stop.wait.Lock()
	return stop
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *StopwatchTSC) mustBeCompleted() {
	if int(atomic.LoadInt32(&bench.done)) > len(bench.spans) {
		panic("benchmarking incomplete")
	}
}

// Start starts measuring a new lap.
// It returns the lap number to pass in for Stop.
// It will return false, when all measurements have been made.
func (bench *StopwatchTSC) Start() (int32, bool) {
	lap := atomic.AddInt32(&bench.lap, 1) - 1
	if int(lap) > len(bench.spans) {
		return -1, false
	}
	bench.spans[lap].Start = TSC()
	return lap, true
}

// Stop stops measuring the specified lap.
func (bench *StopwatchTSC) Stop(lap int32) {
	bench.spans[lap].Finish = TSC()

	done := atomic.AddInt32(&bench.done, 1)
	if int(done) == len(bench.spans) {
		bench.finalize()
	} else if int(done) > len(bench.spans) {
		panic("stop called too many times")
	}
}

// finalize finalizes the stopwatchTSC
func (bench *StopwatchTSC) finalize() {
	bench.wait.Unlock()
}

// Wait waits for all measurements to be completed.
func (bench *StopwatchTSC) Wait() {
	bench.wait.Lock()
	// intentionally empty
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
