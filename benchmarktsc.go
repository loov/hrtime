package hrtime

import "time"

type BenchmarkTSC struct {
	Step   int
	Counts []Count
	Start  Count
	Stop   Count
}

func NewBenchmarkTSC(count int) *BenchmarkTSC {
	if count <= 0 {
		panic("must have count at least 0")
	}

	return &BenchmarkTSC{
		Step:   0,
		Counts: make([]Count, count),
		Start:  0,
		Stop:   0,
	}
}

func (bench *BenchmarkTSC) finalize(last Count) {
	if bench.Stop != 0 {
		return
	}

	bench.Start = bench.Counts[0]
	bench.Stop = last
	for i := range bench.Counts[:len(bench.Counts)-1] {
		bench.Counts[i] = bench.Counts[i+1] - bench.Counts[i]
	}
	bench.Counts[len(bench.Counts)-1] = bench.Stop - bench.Counts[len(bench.Counts)-1]
}

func (bench *BenchmarkTSC) Next() bool {
	now := TSC()
	if bench.Step >= len(bench.Counts) {
		bench.finalize(now)
		return false
	}
	bench.Counts[bench.Step] = TSC()
	bench.Step++
	return true
}

func (bench *BenchmarkTSC) Laps() []time.Duration {
	laps := make([]time.Duration, len(bench.Counts))
	for i, v := range bench.Counts {
		laps[i] = v.ApproxDuration()
	}
	return laps
}

func (bench *BenchmarkTSC) Histogram(binCount int) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}

	return NewHistogram(bench.Laps(), binCount)
}

func (bench *BenchmarkTSC) HistogramClamp(binCount int, min, max time.Duration) *Histogram {
	if bench.Stop == 0 {
		panic("benchmarking incomplete")
	}

	laps := make([]time.Duration, 0, len(bench.Counts))
	for _, count := range bench.Counts {
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
