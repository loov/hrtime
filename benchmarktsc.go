package hrtime

import (
	"time"
)

type BenchmarkTSC struct {
	step   int
	counts []Count
	start  Count
	stop   Count
}

func NewBenchmarkTSC(count int) *BenchmarkTSC {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &BenchmarkTSC{
		step:   0,
		counts: make([]Count, count),
		start:  0,
		stop:   0,
	}
}

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

func (bench *BenchmarkTSC) Counts() []Count {
	return append(bench.counts[:0:0], bench.counts...)
}

func (bench *BenchmarkTSC) Laps() []time.Duration {
	laps := make([]time.Duration, len(bench.counts))
	for i, v := range bench.counts {
		laps[i] = v.ApproxDuration()
	}
	return laps
}

func (bench *BenchmarkTSC) Histogram(binCount int) *Histogram {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}

	return NewDurationHistogram(bench.Laps(), binCount)
}

func (bench *BenchmarkTSC) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	if bench.stop == 0 {
		panic("benchmarking incomplete")
	}

	laps := make([]time.Duration, 0, len(bench.counts))
	for _, count := range bench.counts {
		lap := count.ApproxDuration()
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
