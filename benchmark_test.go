package hrtime_test

import (
	"fmt"
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
