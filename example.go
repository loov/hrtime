// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/loov/hrtime"
)

func main() {
	start := hrtime.Now()
	time.Sleep(1000 * time.Nanosecond)
	fmt.Println(hrtime.Since(start))

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

// Example output:
// 836.599µs
//    710.64µs [    5]
//   854.331µs [  469] ██████████████▌
//   998.022µs [ 1125] ██████████████████████████████████▌
//  1.141713ms [ 1286] ████████████████████████████████████████
//  1.285405ms [  134] ████
//  1.429096ms [  468] ██████████████▌
//  1.572787ms [  473] ██████████████▌
//  1.716479ms [   85] ██▌
//   1.86017ms [   27] ▌
//  2.003861ms [   24] ▌
