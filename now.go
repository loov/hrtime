package hrtime

import (
	"time"
)

var nanoOverhead time.Duration

// Overhead returns approximate overhead for a call to Now() or Since()
func Overhead() time.Duration { return nanoOverhead }

// Since returns time.Duration since start
func Since(start time.Duration) time.Duration { return Now() - start }

func calculateNanosOverhead() {
	start := Now()
	for i := 0; i < calibrationCalls; i++ {
		Now()
	}
	stop := Now()
	nanoOverhead = (stop - start) / (calibrationCalls + 1)
}
