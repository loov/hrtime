// +build !windows

// Fallback to using time.Now
package hrtime

import "time"

var nanoStart = time.Now()

// Now returns current nanoseconds with best possible precision
func Now() Nano { return Nano(time.Since(nanoStart)) }

// NanosPrecision returns maximum possible precision for Nanos in nanoseconds
func NanosPrecision() float64 { return float64(time.Nanosecond) * 1e9 / float64(time.Second) }

// NanosFrequency returns counts per second
func NanosFrequency() Nano { return Nano(time.Second / time.Nanosecond) }
