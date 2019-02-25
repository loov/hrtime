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
	lap   int32
	done  int32
	spans []Span
	wait  sync.Mutex
}

// NewStopwatch creates a new concurrent benchmark using Now
func NewStopwatch(count int) *Stopwatch {
	if count <= 0 {
		panic("must have count at least 0")
	}

	stop := &Stopwatch{
		lap:   0,
		spans: make([]Span, count),
	}
	stop.wait.Lock()
	return stop
}

// mustBeCompleted checks whether measurement has been completed.
func (bench *Stopwatch) mustBeCompleted() {
	if int(atomic.LoadInt32(&bench.done)) > len(bench.spans) {
		panic("benchmarking incomplete")
	}
}

// Start starts measuring a new lap.
// It returns the lap number to pass in for Stop.
// It will return false, when all measurements have been made.
func (bench *Stopwatch) Start() (int32, bool) {
	lap := atomic.AddInt32(&bench.lap, 1) - 1
	if int(lap) > len(bench.spans) {
		return -1, false
	}
	bench.spans[lap].Start = Now()
	return lap, true
}

// Stop stops measuring the specified lap.
func (bench *Stopwatch) Stop(lap int32) {
	bench.spans[lap].Finish = Now()

	done := atomic.AddInt32(&bench.done, 1)
	if int(done) == len(bench.spans) {
		bench.finalize()
	} else if int(done) > len(bench.spans) {
		panic("stop called too many times")
	}
}

// finalize finalizes the stopwatch
func (bench *Stopwatch) finalize() {
	bench.wait.Unlock()
}

// Wait waits for all measurements to be completed.
func (bench *Stopwatch) Wait() {
	bench.wait.Lock()
	// intentionally empty
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
