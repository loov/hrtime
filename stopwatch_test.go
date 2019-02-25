package hrtime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/loov/hrtime"
)

func ExampleStopwatch() {
	const numberOfExperiments = 4096
	bench := hrtime.NewStopwatch(numberOfExperiments)
	for i := 0; i < numberOfExperiments; i++ {
		go func() {
			lap := bench.Start()
			defer bench.Stop(lap)

			time.Sleep(1000 * time.Nanosecond)
		}()
	}
	bench.Wait()
	fmt.Println(bench.Histogram(10))
}

func ExampleStopwatchTSC() {
	const numberOfExperiments = 4096
	bench := hrtime.NewStopwatchTSC(numberOfExperiments)
	for i := 0; i < numberOfExperiments; i++ {
		go func() {
			lap := bench.Start()
			defer bench.Stop(lap)

			time.Sleep(1000 * time.Nanosecond)
		}()
	}
	bench.Wait()
	fmt.Println(bench.Histogram(10))
}

func TestStopwatch(t *testing.T) {
	bench := hrtime.NewStopwatch(8)
	for i := 0; i < 8; i++ {
		go func() {
			lap := bench.Start()
			defer bench.Stop(lap)

			time.Sleep(1000 * time.Nanosecond)
		}()
	}
	bench.Wait()
	t.Log(bench.Histogram(10))
}

func TestStopwatchTSC(t *testing.T) {
	bench := hrtime.NewStopwatchTSC(8)
	for i := 0; i < 8; i++ {
		go func() {
			lap := bench.Start()
			defer bench.Stop(lap)

			time.Sleep(1000 * time.Nanosecond)
		}()
	}
	bench.Wait()
	t.Log(bench.Histogram(10))
}
