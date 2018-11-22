package hrtime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/loov/hrtime"
)

func ExampleBenchmark() {
	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

func ExampleBenchmarkTSC() {
	bench := hrtime.NewBenchmarkTSC(4 << 10)
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
