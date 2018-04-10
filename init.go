package hrtime

const calibrationCalls = 1 << 10

func init() {
	calculateNanosOverhead()

	initCPU()
	{
		_, _, _, edx := cpuid(0x80000007, 0x0)
		rdtscpInvariant = edx&(1<<8) != 0
	}
	calculateTSCOverhead()
	calculateTSCConversion()
}
