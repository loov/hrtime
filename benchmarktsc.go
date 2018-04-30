package hrtime

import "time"

type BenchmarkTSC struct {
	Step  int
	Laps  []Count
	Start Count
	Stop  Count
}

func NewBenchmarkTSC(count int) *BenchmarkTSC {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &BenchmarkTSC{
		Step:  0,
		Laps:  make([]Count, count),
		Start: 0,
		Stop:  0,
	}
}

func (bench *BenchmarkTSC) finalize(last Count) {
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

func (bench *BenchmarkTSC) Next() bool {
	now := TSC()
	if bench.Step >= len(bench.Laps) {
		bench.finalize(now)
		return false
	}
	bench.Laps[bench.Step] = TSC()
	bench.Step++
	return true
}

func (bench *BenchmarkTSC) Histogram(binCount int) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}

	laps := make([]time.Duration, len(bench.Laps))
	for i, v := range bench.Laps {
		laps[i] = v.ApproxDuration()
	}

	return NewHistogram(laps, binCount)
}

func (bench *BenchmarkTSC) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}
	laps := make([]time.Duration, 0, len(bench.Laps))
	for _, count := range bench.Laps {
		lap := count.ApproxDuration()
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
