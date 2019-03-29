// Package hrtesting implements wrappers for testing.B that allow to output
// more detailed information in standard go benchmarks.
//
// Since, it is mainly built for convenience it needs to call `b.StartTimer` and other benchmark
// timer functions. As a result, manually calling them can cause unintended measurement results.
//
// Since such benchmarking will have an overhead, it will increase single measurement results.
// To disable that measurement temporarily "-tags nohrtime".
//
// To use this package write your benchmark as:
//
//     func BenchmarkHello(b *testing.B) {
//         bench := hrtesting.NewBenchmark(b)
//         defer bench.Report()
//
//         for bench.Next() {
//             fmt.Sprintf("hello")
//         }
//     }
//
// Only statements in the `for bench.Next() {` loop will be measured.
//
// To use time stamp counters, which are not supported on all platforms:
//
//     func BenchmarkHello(b *testing.B) {
//         bench := hrtesting.NewBenchmarkTSC(b)
//         defer bench.Report()
//
//         for bench.Next() {
//             fmt.Sprintf("hello")
//         }
//     }
//
package hrtesting
