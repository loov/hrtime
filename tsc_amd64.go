package hrtime

func rdtscpAsm() uint64
func cpuidAsm(op1, op2 uint32) (eax, ebx, ecx, edx uint32)

func initCPU() {
	rdtscp = func() uint64 { return uint64(Now()) }
	cpuid = cpuidAsm
}
