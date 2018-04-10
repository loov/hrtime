package hrtime

import (
	"time"
)

// Nano is a high resolution time value in nanoseconds
//
// using it to measure times longer than 1 hour should be considered unreliable
type Nano int64

// Since returns count since start
func Since(start Nano) Nano { return Now() - start }

// Hours returns time value as a floating point number of hours.
func (nano Nano) Hours() float64 { return float64(nano) / (60 * 60 * 1e9) }

// Minutes returns time value as a floating point number of minutes.
func (nano Nano) Minutes() float64 { return float64(nano) / (60 * 1e9) }

// Seconds returns time value as a floating point number of seconds.
func (nano Nano) Seconds() float64 { return float64(nano) / 1e9 }

// String returns a string representing the duration in the form "72h3m0.5s".
func (nano Nano) String() string { return nano.Duration().String() }

// Nanoseconds converts Nano into an int64
func (nano Nano) Nanoseconds() int64 { return int64(nano) }

// Duration converts Nano value into a time.Duration
func (nano Nano) Duration() time.Duration { return time.Duration(nano) * time.Nanosecond }

var nanoOverhead Nano

// Overhead returns approximate overhead for a call to Now() or Since()
func Overhead() Nano { return nanoOverhead }

func calculateNanosOverhead() {
	start := Now()
	for i := 0; i < calibrationCalls; i++ {
		Now()
	}
	stop := Now()
	nanoOverhead = (stop - start) / (calibrationCalls + 1)
}
