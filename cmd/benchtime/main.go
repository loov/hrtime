package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"runtime"
	"runtime/debug"
	"sort"
	"time"
	_ "unsafe"

	"github.com/loov/hrtime"
	"github.com/loov/plot"
)

var (
	samples = flag.Int("samples", 1e6, "measurements per line")
	warmup  = flag.Int("warmup", 256, "warmup count")
	mintime = flag.Float64("min", 0, "minimum ns time to consider")
	maxtime = flag.Float64("max", 100, "maximum ns time to consider")

	svg    = flag.String("svg", "results.svg", "")
	kernel = flag.Float64("kernel", 1, "kernel width")
	width  = flag.Float64("width", 1500, "svg width")
	height = flag.Float64("height", 150, "svg single-plot height")
)

func init() {
	runtime.LockOSThread()
}

//go:linkname nanotime runtime.nanotime
func nanotime() int64

func main() {
	flag.Parse()

	N := *samples

	times := make([]time.Time, N)
	timenanos := make([]int64, N)
	deltas := make([]time.Duration, N)
	nanotimes := make([]int64, N)
	qpcs := make([]int64, N)
	rdtsc := make([]uint64, N)
	rdtscp := make([]uint64, N)

	var beforRDTSC, afterRDTSC, beforRDTSCP, afterRDTSCP time.Time

	fmt.Println("benchmarking")
	debug.SetGCPercent(-1)
	{
		fmt.Println("benchmarking time.Now")
		runtime.GC()
		for i := range times {
			times[i] = time.Now()
		}

		fmt.Println("benchmarking time.UnixNano")
		runtime.GC()
		for i := range timenanos {
			timenanos[i] = time.Now().UnixNano()
		}

		fmt.Println("benchmarking time.Since")
		runtime.GC()
		start := time.Now()
		for i := range deltas {
			deltas[i] = time.Since(start)
		}

		fmt.Println("benchmarking nanotime")
		runtime.GC()
		for i := range nanotimes {
			nanotimes[i] = nanotime()
		}

		if runtime.GOOS == "windows" {
			fmt.Println("benchmarking QPC")
			runtime.GC()
			for i := range qpcs {
				qpcs[i] = QPC()
			}
		}

		if hrtime.TSCSupported() {
			fmt.Println("benchmarking RDTSC")
			runtime.GC()

			beforRDTSC = time.Now()
			for i := range rdtsc {
				rdtsc[i] = hrtime.RDTSC()
			}
			afterRDTSC = time.Now()

			fmt.Println("benchmarking RDTSCP")
			runtime.GC()
			beforRDTSCP = time.Now()
			for i := range rdtscp {
				rdtscp[i] = hrtime.RDTSC()
			}
			afterRDTSCP = time.Now()

			runtime.GC()
		}
	}
	debug.SetGCPercent(100)

	rdtscCalibration := float64(afterRDTSC.Sub(beforRDTSC).Nanoseconds()) / float64(rdtsc[N-1]-rdtsc[0])
	rdtscpCalibration := float64(afterRDTSCP.Sub(beforRDTSCP).Nanoseconds()) / float64(rdtscp[N-1]-rdtscp[0])

	offset := *warmup
	ns_times := make([]float64, N-offset)
	ns_timenanos := make([]float64, N-offset)
	ns_deltas := make([]float64, N-offset)
	ns_nanotimes := make([]float64, N-offset)
	ns_qpcs := make([]float64, N-offset)
	ns_rdtsc := make([]float64, N-offset)
	ns_rdtscp := make([]float64, N-offset)

	qpcmul := 1e9 / float64(QPCFrequency())
	for i := range ns_times {
		ns_times[i] = float64(times[offset+i].Sub(times[offset+i-1]).Nanoseconds())
		ns_timenanos[i] = float64(timenanos[offset+i] - timenanos[offset+i-1])
		ns_deltas[i] = float64((deltas[offset+i] - deltas[offset+i-1]).Nanoseconds())
		ns_nanotimes[i] = float64(nanotimes[offset+i] - nanotimes[offset+i-1])
		ns_qpcs[i] = float64(qpcs[offset+i]-qpcs[offset+i-1]) * qpcmul
		ns_rdtsc[i] = float64(rdtsc[offset+i]-rdtsc[offset+i-1]) * rdtscCalibration
		ns_rdtscp[i] = float64(rdtsc[offset+i]-rdtsc[offset+i-1]) * rdtscpCalibration
	}

	var timings = []*Timing{}
	timings = append(timings,
		&Timing{Name: "time.Now", Measured: ns_times},
		&Timing{Name: "time.UnixNano", Measured: ns_timenanos},
		&Timing{Name: "time.Since", Measured: ns_deltas},
		&Timing{Name: "nanotime", Measured: ns_nanotimes},
	)

	if runtime.GOOS == "windows" {
		timings = append(timings,
			&Timing{Name: "QPC", Measured: ns_qpcs},
		)
	}

	if hrtime.TSCSupported() {
		timings = append(timings,
			&Timing{Name: "RDTSC", Measured: ns_rdtsc},
			&Timing{Name: "RDTSCP", Measured: ns_rdtscp},
		)
	}

	p := plot.New()

	p.X.Min = *mintime
	p.X.Max = *maxtime
	p.X.MajorTicks = 10
	p.X.MinorTicks = 10

	stack := plot.NewVStack()
	stack.Margin = plot.R(5, 5, 5, 5)
	p.Add(stack)

	for _, timing := range timings {
		timing.Prepare(*mintime, *maxtime)
		density := plot.NewDensity("ns", timing.Sanitized)
		density.Kernel = *kernel
		density.Class = timing.Name
		density.Stroke = color.NRGBA{0, 0, 0, 255}
		density.Fill = color.NRGBA{0, 0, 0, 50}

		flex := plot.NewHFlex()

		flex.Add(130, plot.NewTextbox(
			timing.Name,
			fmt.Sprintf("Measured = %v", len(timing.Measured)),
			fmt.Sprintf("Zeros = %v", timing.Zero),
			fmt.Sprintf("Underlimit = %v", timing.Underlimit),
			fmt.Sprintf("Overlimit = %v", timing.Overlimit),
			fmt.Sprintf("99.9 = %v", int(timing.P999)),
			fmt.Sprintf("99.99 = %v", int(timing.P9999)),
			fmt.Sprintf("Max = %v", int(timing.Max)),
		))
		flex.AddGroup(0,
			plot.NewGrid(),
			density,
			plot.NewTickLabels(),
		)

		stack.Add(flex)
	}

	svgcanvas := plot.NewSVG(*width, *height*float64(len(timings)))
	p.Draw(svgcanvas)

	err := ioutil.WriteFile(*svg, svgcanvas.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
}

type Timing struct {
	Name      string
	Measured  []float64
	Sanitized []float64

	Zero       int
	Underlimit int
	Overlimit  int

	P999  float64
	P9999 float64
	Max   float64
}

func (t *Timing) Prepare(min, max float64) {
	t.Sanitized = append(t.Measured[:0:0], t.Measured...)
	sort.Float64s(t.Sanitized)

	t.P999 = t.Sanitized[(len(t.Sanitized)-1)*999/1000]
	t.P9999 = t.Sanitized[(len(t.Sanitized)-1)*9999/10000]
	t.Max = t.Sanitized[len(t.Sanitized)-1]

	for i, v := range t.Sanitized {
		if v <= 0 {
			t.Zero++
		} else if v <= min {
			t.Sanitized[i] = min
			t.Underlimit++
		} else {
			break
		}
	}

	tail := len(t.Sanitized) - 1
	for ; tail >= 0; tail-- {
		if t.Sanitized[tail] < max {
			break
		}
		t.Sanitized[tail] = max
		t.Overlimit++
	}

	t.Sanitized = t.Sanitized[t.Zero:]
}
