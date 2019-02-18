package main

import (
	"syscall"
	"unsafe"
)

// precision timing
var (
	modkernel32                   = syscall.NewLazyDLL("kernel32.dll")
	queryPerformanceFrequencyProc = modkernel32.NewProc("QueryPerformanceFrequency")
	queryPerformanceCounterProc   = modkernel32.NewProc("QueryPerformanceCounter")
)

// now returns time.Duration using queryPerformanceCounter
func QPC() int64 {
	var now int64
	syscall.Syscall(queryPerformanceCounterProc.Addr(), 1, uintptr(unsafe.Pointer(&now)), 0, 0)
	return now
}

// QPCFrequency returns frequency in ticks per second
func QPCFrequency() int64 {
	var freq int64
	r1, _, _ := syscall.Syscall(queryPerformanceFrequencyProc.Addr(), 1, uintptr(unsafe.Pointer(&freq)), 0, 0)
	if r1 == 0 {
		panic("call failed")
	}
	return freq
}
