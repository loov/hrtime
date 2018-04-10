// +build !amd64 gccgo

package tsc

func initCPU() {
	rdtscp = func() uint64 { return 0 }
	cpuid = func(op1, op2 uint32) (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}
}
