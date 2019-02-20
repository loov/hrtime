// This program demonstrates the basic usage of the package.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/loov/hrtime"
)

func main() {
	start := hrtime.Now()
	time.Sleep(1000 * time.Nanosecond)
	fmt.Println(hrtime.Since(start))

	const numberOfExperiments = 4096
	bench := hrtime.NewBenchmark(numberOfExperiments)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
		if rand.Intn(2000) == 0 {
			time.Sleep(time.Second)
		}
	}
	fmt.Println(bench.Histogram(10))
}

// Example output:
// 22.779µs
//   avg 1.744377ms;  min 827.221µs;  p50 1.986758ms;  max 2.200263ms;
//   p90 2.036338ms;  p99 2.150237ms;  p999 2.182397ms;  p9999 2.200263ms;
//   827.221µs [  89] ██
//         1ms [1064] █████████████████████████
//       1.2ms [  46] █
//       1.4ms [   0]
//       1.6ms [   1]
//       1.8ms [1698] ████████████████████████████████████████
//         2ms [1197] ████████████████████████████
//       2.2ms [   1]
//       2.4ms [   0]
//       2.6ms [   0]
