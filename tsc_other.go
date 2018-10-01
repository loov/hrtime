// +build !amd64 gccgo

package tsc

func initCPU() {
	cpuid = func(op1, op2 uint32) (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}
}

// RDTSCP returns 0 for unsupported configuration
func RDTSCP() uint64 { return 0 }

// RDTSC returns 0 for unsupported configuration
func RDTSC() uint64 { return 0 }
