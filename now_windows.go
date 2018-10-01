package hrtime

import (
	"syscall"
	"time"
	"unsafe"
)

// precision timing
var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procFreq    = modkernel32.NewProc("QueryPerformanceFrequency")
	procCounter = modkernel32.NewProc("QueryPerformanceCounter")

	qpcFrequency = getFrequency()
)

// getFrequency returns frequency in ticks per second.
func getFrequency() int64 {
	var freq int64
	r1, _, _ := syscall.Syscall(procFreq.Addr(), 1, uintptr(unsafe.Pointer(&freq)), 0, 0)
	if r1 == 0 {
		panic("call failed")
	}
	return freq
}

// Now returns current time.Duration with best possible precision.
func Now() time.Duration {
	var now int64
	syscall.Syscall(procCounter.Addr(), 1, uintptr(unsafe.Pointer(&now)), 0, 0)
	return time.Duration(now) * time.Second / (time.Duration(qpcFrequency) * time.Nanosecond)
}

// NowPrecision returns maximum possible precision for Now in nanoseconds.
func NowPrecision() float64 {
	return float64(time.Second) / (float64(qpcFrequency) * float64(time.Nanosecond))
}
