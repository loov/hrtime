// +build !race

package hrtime_test

import (
	"testing"

	"github.com/loov/hrtime"
)

func TestNowCalibration(t *testing.T) {
	start := hrtime.Now()
	empty()
	stop := hrtime.Now()
	if stop-start < hrtime.Overhead() {
		t.Errorf("measurement: %v %v", stop-start, hrtime.Overhead())
	}
}

func TestNowPrecision(t *testing.T) {
	const N = 8 << 10

	start := hrtime.Now()
	for i := 0; i < N; i++ {
		empty()
	}
	stop := hrtime.Now()
	loopTime := stop - start - 2*hrtime.Overhead()

	// we expect each call to take at least 1 nanosecond
	if loopTime.Nanoseconds() < N {
		t.Errorf("slow: loop time took %d", loopTime)
	}
	// we expect no call to take more than 10 nanoseconds
	if loopTime.Nanoseconds() > 10*N {
		t.Errorf("fast: loop time took %d", loopTime)
	}
}

//go:noinline
func empty() {}
