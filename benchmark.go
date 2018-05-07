package hrtime

import (
	"time"
)

type Benchmark struct {
	Step      int
	Durations []time.Duration
	Start     time.Duration
	Stop      time.Duration
}

func NewBenchmark(count int) *Benchmark {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &Benchmark{
		Step:      0,
		Durations: make([]time.Duration, count),
		Start:     0,
		Stop:      0,
	}
}

func (bench *Benchmark) finalize(last time.Duration) {
	if bench.Stop != 0 {
		return
	}

	bench.Start = bench.Durations[0]
	bench.Stop = last
	for i := range bench.Durations[:len(bench.Durations)-1] {
		bench.Durations[i] = bench.Durations[i+1] - bench.Durations[i]
	}
	bench.Durations[len(bench.Durations)-1] = bench.Stop - bench.Durations[len(bench.Durations)-1]
}

func (bench *Benchmark) Next() bool {
	now := Now()
	if bench.Step >= len(bench.Durations) {
		bench.finalize(now)
		return false
	}
	bench.Durations[bench.Step] = Now()
	bench.Step++
	return true
}

func (bench *Benchmark) Laps() []time.Duration {
	return append(bench.Durations[:0:0], bench.Durations...)
}

func (bench *Benchmark) Histogram(binCount int) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}
	return NewHistogram(bench.Durations, binCount)
}

func (bench *Benchmark) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}
	laps := make([]time.Duration, 0, len(bench.Durations))
	for _, lap := range bench.Durations {
		if lap < min {
			laps = append(laps, min)
		} else if lap > max {
			laps = append(laps, max)
		} else {
			laps = append(laps, lap)
		}
	}
	return NewHistogram(laps, binCount)
}
