// +build !amd64 gccgo

package hrtime

func initCPU() {
	cpuid = func(op1, op2 uint32) (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}
}

// RDTSCP returns 0 for unsupported configuration
//
// If a given OS doesn't support the instruction it returns 0.
// Use TSCSupported to check.
func RDTSCP() uint64 { return 0 }

// RDTSC returns 0 for unsupported configuration
//
// If a given OS doesn't support the instruction it returns 0.
// Use TSCSupported to check.
func RDTSC() uint64 { return 0 }
