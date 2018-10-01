package hrtime

func rdtscpAsm() uint64
func rdtscAsm() uint64
func cpuidAsm(op1, op2 uint32) (eax, ebx, ecx, edx uint32)

func initCPU() {
	cpuid = cpuidAsm
}

// RDTSCP returns Read Time-Stamp Counter value using RDTSCP asm instruction.
func RDTSCP() uint64 { return rdtscpAsm() }

// RDTSC returns Read Time-Stamp Counter value using RDTSC asm instruction.
func RDTSC() uint64 { return rdtscAsm() }
