package hrtime

import (
	"syscall"
	"unsafe"
)

// precision timing
var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procFreq    = modkernel32.NewProc("QueryPerformanceFrequency")
	procCounter = modkernel32.NewProc("QueryPerformanceCounter")

	qpcFrequency = getFrequency()
)

// getFrequency returns frequency in ticks per second
func getFrequency() int64 {
	var freq int64
	r1, _, _ := syscall.Syscall(procFreq.Addr(), 1, uintptr(unsafe.Pointer(&freq)), 0, 0)
	if r1 == 0 {
		panic("call failed")
	}
	return freq
}

// Now returns current nanoseconds with best possible precision
func Now() Nano {
	var now int64
	syscall.Syscall(procCounter.Addr(), 1, uintptr(unsafe.Pointer(&now)), 0, 0)
	return Nano(now * 1e9 / qpcFrequency)
}

// NanosPrecision returns maximum possible precision for Now in nanoseconds
func NanosPrecision() float64 { return float64(1e9) / float64(qpcFrequency) }

// NanosFrequency returns counts per second
func NanosFrequency() Nano { return Nano(qpcFrequency) }
