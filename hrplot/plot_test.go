package hrplot_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/loov/hrtime/hrplot"
	"github.com/loov/hrtime/hrtesting"
)

func BenchmarkReport(b *testing.B) {

	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()

	defer hrplot.All("ns-all.svg", bench)
	defer hrplot.Density("ns-density.svg", bench)
	defer hrplot.Line("ns-line.svg", bench)
	defer hrplot.Percentiles("ns-percentiles.svg", bench)

	runtime.GC()
	for bench.Next() {
		r := fmt.Sprintf("hello, world %d", 123)
		runtime.KeepAlive(r)
	}
}

func BenchmarkTSCReport(b *testing.B) {
	bench := hrtesting.NewBenchmarkTSC(b)
	defer bench.Report()

	defer hrplot.All("tsc-all.svg", bench)
	defer hrplot.Density("tsc-density.svg", bench)
	defer hrplot.Line("tsc-line.svg", bench)
	defer hrplot.Percentiles("tsc-percentiles.svg", bench)

	runtime.GC()
	for bench.Next() {
		r := fmt.Sprintf("hello, world %d", 123)
		runtime.KeepAlive(r)
	}
}
