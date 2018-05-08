package hrtime

import (
	"time"
)

type Benchmark struct {
	step  int
	laps  []time.Duration
	start time.Duration
	stop  time.Duration
}

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

func (bench *Benchmark) Laps() []time.Duration {
	return append(bench.laps[:0:0], bench.laps...)
}

func (bench *Benchmark) Histogram(binCount int) *Histogram {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}
	return NewDurationHistogram(bench.laps, binCount)
}

func (bench *Benchmark) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}
	laps := make([]time.Duration, 0, len(bench.laps))
	for _, lap := range bench.laps {
		if lap < min {
			laps = append(laps, min)
		} else if lap > max {
			laps = append(laps, max)
		} else {
			laps = append(laps, lap)
		}
	}
	return NewDurationHistogram(laps, binCount)
}
