package hrtime

import (
	"testing"
)

func BenchmarkRDTSCP(b *testing.B) {
	if !TSCSupported() {
		b.Skip("Cycle counting not supported")
	}

	for i := 0; i < b.N; i++ {
		rdtscp()
	}
}

func BenchmarkRDTSC(b *testing.B) {
	if !TSCSupported() {
		b.Skip("Cycle counting not supported")
	}
	for i := 0; i < b.N; i++ {
		rdtsc()
	}
}
