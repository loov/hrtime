// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/loov/hrtime"
)

func main() {
	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}
