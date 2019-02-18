package hrtime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/loov/hrtime"
)

func ExampleBenchmark() {
	const numberOfExperiments = 4096
	bench := hrtime.NewBenchmark(numberOfExperiments)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

func ExampleBenchmarkTSC() {
	const numberOfExperiments = 4096
	bench := hrtime.NewBenchmarkTSC(numberOfExperiments)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

func TestBenchmark(t *testing.T) {
	bench := hrtime.NewBenchmark(8)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	t.Log(bench.Histogram(10))
}

func TestBenchmarkTSC(t *testing.T) {
	bench := hrtime.NewBenchmarkTSC(8)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	t.Log(bench.Histogram(10))
}
