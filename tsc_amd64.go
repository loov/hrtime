package hrtime

func rdtscpAsm() uint64
func rdtscAsm() uint64
func cpuidAsm(op1, op2 uint32) (eax, ebx, ecx, edx uint32)

func initCPU() {
	cpuid = cpuidAsm
}

func rdtscp() uint64 { return rdtscpAsm() }
func rdtsc() uint64  { return rdtscAsm() }
