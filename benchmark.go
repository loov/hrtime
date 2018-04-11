package hrtime

import "time"

type Benchmark struct {
	Step  int
	Laps  []time.Duration
	Start time.Duration
	Stop  time.Duration
}

func NewBenchmark(count int) *Benchmark {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &Benchmark{
		Step:  0,
		Laps:  make([]time.Duration, count),
		Start: 0,
		Stop:  0,
	}
}

func (bench *Benchmark) finalize(last time.Duration) {
	if bench.Stop != 0 {
		return
	}

	bench.Start = bench.Laps[0]
	bench.Stop = last
	for i := range bench.Laps[:len(bench.Laps)-1] {
		bench.Laps[i] = bench.Laps[i+1] - bench.Laps[i]
	}
	bench.Laps[len(bench.Laps)-1] = bench.Stop - bench.Laps[len(bench.Laps)-1]
}

func (bench *Benchmark) Next() bool {
	now := Now()
	if bench.Step >= len(bench.Laps) {
		bench.finalize(now)
		return false
	}
	bench.Laps[bench.Step] = Now()
	bench.Step++
	return true
}

func (bench *Benchmark) Histogram(binCount int) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}
	return NewHistogram(bench.Laps, binCount)
}
