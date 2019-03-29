// +build !nohrtime
// +build go1.13

package hrtesting

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

func (bench *BenchmarkTSC) Report() {
	hist := bench.hr.Histogram(1)
	if bench.b.N >= 3 {
		bench.b.ReportMetric(hist.P50, "c/p50")
	}
	if bench.b.N >= 10 {
		bench.b.ReportMetric(hist.P90, "c/p90")
	}
	if bench.b.N >= 100 {
		bench.b.ReportMetric(hist.P99, "c/p99")
	}
}
