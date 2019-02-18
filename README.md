# hrtime

[![GoDoc](https://godoc.org/github.com/loov/hrtime?status.svg)](http://godoc.org/github.com/loov/hrtime)

**BETA QUALITY**

Package hrtime implements high-resolution timing functions and benchmarking utilities.

For an example measuring `time.Sleep` precision on Mac and Windows.

```
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
```

Output on Mac:

```
12µs
  avg 14.5µs;  min 2µs;  p50 12µs;  max 74µs;
  p90 22µs;  p99 44µs;  p999 69µs;  p9999 74µs;
        2µs [ 229] ██▉
       10µs [3239] ████████████████████████████████████████
       20µs [ 483] ██████
       30µs [  80] █
       40µs [  39] ▌
       50µs [  17] ▏
       60µs [   6] ▏
       70µs [   3] ▏
       80µs [   0]
       90µs [   0]
```

Output on Windows:

```
1.0633ms
  avg 1.48ms;  min 622µs;  p50 1.06ms;  max 2.39ms;
  p90 2.01ms;  p99 2.26ms;  p999 2.36ms;  p9999 2.39ms;
      622µs [  58] █▊
      800µs [ 640] ██████████████████▍
        1ms [1388] ████████████████████████████████████████
      1.2ms [  44] █▎
      1.4ms [  16] ▌
      1.6ms [  62] █▉
      1.8ms [ 969] ████████████████████████████
        2ms [ 803] ███████████████████████▏
      2.2ms [ 116] ███▍
      2.4ms [   0]
```

_The full explanation why it outputs this is out of the scope of this document. However all sleep instructions have a specified granularity and `time.Sleep` actual sleeping time is `requested time ± sleep granularity`. Of course there are other exceptions to that behavior._