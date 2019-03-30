// +build !nohrtime
// +build !go1.13

package hrtesting

import "time"

// Report reports the 50-th, 90-th and 99-th percentile to the log.
func (bench *Benchmark) Report() {
	hist := bench.hr.Histogram(1)
	bench.b.Logf("%6v₅₀ %6v₉₀ %6v₉₉ N=%v",
		time.Duration(truncate(hist.P50, 3)),
		time.Duration(truncate(hist.P90, 3)),
		time.Duration(truncate(hist.P99, 3)),
		bench.b.N,
	)
}

// Report reports the 50-th, 90-th and 99-th percentile to the log.
func (bench *BenchmarkTSC) Report() {
	hist := bench.hr.Histogram(1)
	bench.b.Logf("%6v₅₀ %6v₉₀ %6v₉₉ N=%v",
		truncate(hist.P50, 3),
		truncate(hist.P90, 3),
		truncate(hist.P99, 3),
		bench.b.N,
	)
}
