// +build !windows

// Fallback to using time.Now
package hrtime

import "time"

var nanoStart = time.Now()

// Now returns current time.Duration with best possible precision
func Now() time.Duration { return time.Since(nanoStart) }

// NanosPrecision returns maximum possible precision for Nanos in nanoseconds
func NanosPrecision() float64 { return float64(time.Nanosecond) }
