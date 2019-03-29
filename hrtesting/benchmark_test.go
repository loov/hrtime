package hrtesting_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/loov/hrtime/hrtesting"
)

func BenchmarkReport(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()

	for bench.Next() {
		r := fmt.Sprintf("hello, world %d", 123)
		runtime.KeepAlive(r)
	}
}

func BenchmarkTSCReport(b *testing.B) {
	bench := hrtesting.NewBenchmarkTSC(b)
	defer bench.Report()

	for bench.Next() {
		r := fmt.Sprintf("hello, world %d", 123)
		runtime.KeepAlive(r)
	}
}
