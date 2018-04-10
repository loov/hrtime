package hrtime

// Count represents represents Time Stamp Counter value, when available
//
// Count doesn't depend on power throttling, sleeping and similar effects
// making it useful for benchmarking. However it is not reliably convertible to a
// reasonable time-value.
type Count int64

// ApproxNanos returns approximate conversion into Nano-s
func (count Count) ApproxNanos() Nano {
	if ratioCount == 0 {
		calculateTSCConversion()
	}
	return Nano(count) * ratioNano / Nano(ratioCount)
}

// TSC reads the current Time Stamp Counter value
func TSC() Count { return Count(rdtscp()) }

// TSCSince returns count since start
func TSCSince(start Count) Count { return TSC() - start }

// TSCSupported returns whether processor supports giving invariant time stamp counter values
func TSCSupported() bool { return rdtscpInvariant }

// TSCOverhead returns overhead of Count call
func TSCOverhead() Count { return readTSCOverhead }

var (
	rdtscpInvariant = false
	readTSCOverhead Count

	rdtscp func() uint64
	cpuid  func(op1, op2 uint32) (eax, ebx, ecx, edx uint32)

	ratioNano  Nano
	ratioCount Count
)

func calculateTSCOverhead() {
	if !rdtscpInvariant {
		return
	}

	start := TSC()
	for i := 0; i < calibrationCalls; i++ {
		TSC()
	}
	stop := TSC()

	readTSCOverhead = (stop - start) / (calibrationCalls + 1)
}

func calculateTSCConversion() {
	nanostart := Now()
	countstart := TSC()
	for i := 0; i < calibrationCalls; i++ {
		empty()
	}
	nanoend := Now()
	countstop := TSC()

	ratioNano = nanoend - nanostart - Overhead()
	ratioCount = countstop - countstart - TSCOverhead()
}

//go:noinline
func empty() {}
