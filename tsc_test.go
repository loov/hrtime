package hrtime_test

import (
	"testing"

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

	t.Logf("Conversion 1000000 count = %v", hrtime.Count(1000000).ApproxNanos())

	const N = 8 << 10

	startnano := hrtime.Now()
	start := hrtime.TSC()
	for i := 0; i < N; i++ {
		empty()
	}
	stop := hrtime.TSC()
	stopnano := hrtime.Now()

	loopTime := stop - start - 2*hrtime.TSCOverhead()
	wallTime := stopnano - startnano - 2*hrtime.Overhead() - 2*hrtime.TSCOverhead().ApproxNanos()

	approxConversionDrift := wallTime - loopTime.ApproxNanos()
	if approxConversionDrift < 0 {
		approxConversionDrift *= -1
	}
	if approxConversionDrift > 2*hrtime.Overhead()+hrtime.Nano(500) {
		t.Errorf("drift: too large %v (loopTime:%v, wallTime:%v)", approxConversionDrift.Duration(), loopTime.ApproxNanos(), wallTime)
	}

	// we expect each call to take at least 2 nanos
	if loopTime.ApproxNanos() < 2*N {
		t.Errorf("slow: loop time took %v", loopTime)
	}
	// we expect no call to take more than 20 nanos
	if loopTime.ApproxNanos() > 20*N {
		t.Errorf("fast: loop time took %v", loopTime)
	}
}
