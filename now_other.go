// +build !windows

package hrtime

import "time"

// Now returns current time.Duration with best possible precision
func Now() time.Duration {
	return time.Duration(time.Now().UnixNano()) * time.Nanosecond
}

// NowPrecision returns maximum possible precision for Nanos in nanoseconds
func NowPrecision() float64 { return float64(1) }
