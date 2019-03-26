package hrtesting_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/loov/hrtime/hrtesting"
)

func BenchmarkHistogram(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Print()

	for bench.Next() {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))

	}
}
