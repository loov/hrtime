// +build !race

package hrtime_test

import (
	"testing"
	"time"

	"github.com/loov/hrtime"
)

func TestCountCalibration(t *testing.T) {
	if !hrtime.TSCSupported() {
		t.Skip("Cycle counting not supported")
	}

	start := hrtime.TSC()
	for i := 0; i < 64; i++ {
		empty()
	}
	stop := hrtime.TSC()

	if stop-start < hrtime.TSCOverhead() {
		t.Errorf("overhead is larger than delta: %v-%v=%v overhead:%v", stop, start, stop-start, hrtime.TSCOverhead())
	}
}

func TestCountPrecision(t *testing.T) {
	if !hrtime.TSCSupported() {
		t.Skip("Cycle counting not supported")
	}

	t.Logf("Conversion 1000000 count = %v", hrtime.Count(1000000).ApproxDuration())

	const N = 8 << 10

	startnano := hrtime.Now()
	start := hrtime.TSC()
	for i := 0; i < N; i++ {
		empty()
	}
	stop := hrtime.TSC()
	stopnano := hrtime.Now()

	loopTime := stop - start - 2*hrtime.TSCOverhead()
	wallTime := stopnano - startnano - 2*hrtime.Overhead() - 2*hrtime.TSCOverhead().ApproxDuration()

	approxConversionDrift := wallTime - loopTime.ApproxDuration()
	if approxConversionDrift < 0 {
		approxConversionDrift *= -1
	}
	if approxConversionDrift > 2*hrtime.Overhead()+500*time.Nanosecond {
		t.Logf("drift: too large %v (loopTime:%v, wallTime:%v)", approxConversionDrift, loopTime.ApproxDuration(), wallTime)
	}

	// we expect each call to take at least 1 nanos
	if loopTime.ApproxDuration() < N {
		t.Errorf("slow: loop time took %v", loopTime)
	}
	// we expect no call to take more than 20 nanos
	if loopTime.ApproxDuration() > 20*N {
		t.Errorf("fast: loop time took %v", loopTime)
	}
}

func BenchmarkTSC(b *testing.B) {
	if !hrtime.TSCSupported() {
		b.Skip("Cycle counting not supported")
	}

	for i := 0; i < b.N; i++ {
		hrtime.TSC()
	}
}
