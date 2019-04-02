// +build !nohrtime
// +build go1.13

package hrtesting

// Report reports the 50-th, 90-th and 99-th percentile as a metric.
func (bench *Benchmark) Report() {
	hist := bench.hr.Histogram(1)
	if bench.b.N >= 3 {
		bench.b.ReportMetric(hist.P50, "ns/p50")
	}
	if bench.b.N >= 10 {
		bench.b.ReportMetric(hist.P90, "ns/p90")
	}
	if bench.b.N >= 100 {
		bench.b.ReportMetric(hist.P99, "ns/p99")
	}
}

// Report reports the 50-th, 90-th and 99-th percentile as a metric.
func (bench *BenchmarkTSC) Report() {
	hist := bench.hr.Histogram(1)
	if bench.b.N >= 3 {
		bench.b.ReportMetric(hist.P50, "tsc/p50")
	}
	if bench.b.N >= 10 {
		bench.b.ReportMetric(hist.P90, "tsc/p90")
	}
	if bench.b.N >= 100 {
		bench.b.ReportMetric(hist.P99, "tsc/p99")
	}
}
