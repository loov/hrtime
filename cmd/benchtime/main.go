package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"time"
	_ "unsafe"

	"github.com/loov/hrtime"
)

var (
	samples = flag.Int("samples", 100e6, "measurements per line")
	warmup  = flag.Int("warmup", 256, "warmup count")
	maxtime = flag.Float64("max", 100, "maximum ns time to consider")

	svg     = flag.String("svg", "results.svg", "")
	density = flag.Int("density", 5000, "points per line for graph")
	kernel  = flag.Float64("kernel", 0.5, "kernel width")
	width   = flag.Float64("width", 1500, "svg width")
	height  = flag.Float64("height", 150, "svg single-plot height")
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

	timings = append(timings,
		&Timing{Name: "RDTSC", Measured: ns_rdtsc},
		&Timing{Name: "RDTSCP", Measured: ns_rdtscp},
	)

	out, err := os.Create(*svg)
	check(err)
	defer out.Close()

	buf := bufio.NewWriter(out)
	defer buf.Flush()

	plot(buf, timings)
}

func plot(w io.Writer, timings []*Timing) {
	fmt.Println("plotting")

	write := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	width, height := *width, *height
	density := *density
	kernel := *kernel

	min := 0.0
	max := *maxtime
	pointstep := (max - min) / float64(density)
	tickstepx := *maxtime / 50
	majtickx := 5

	if kernel < pointstep {
		fmt.Println("kernel to small using:", pointstep)
		kernel = pointstep
	}

	tickstepy := 10.0 / 100
	majticky := 5

	pad := 5.0

	legendwidth := 150.0
	tox := func(v float64) float64 {
		return legendwidth + (v-min)*(width-legendwidth)/(max-min)
	}

	ylog_compress := 50.0
	ylog_mul := 1 / math.Log(ylog_compress+1)
	toy := func(p float64) float64 {
		p = math.Log(p*ylog_compress+1) * ylog_mul
		return height - p*height
	}

	write(`<?xml version="1.0" standalone="no"?>
		<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN" "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
		<svg xmlns="http://www.w3.org/2000/svg"
			width="%.0fpx" height="%.0fpx">`, width+pad*2, height*float64(len(timings))+pad*2)
	defer write(`</svg>`)

	write(`<style>
		/* <![CDATA[ */
		svg { dominant-baseline: hanging; }
		polyline { fill: transparent; }
		.line {
			stroke: #000;
			fill: rgba(0,0,0,0.2);
		}
		text {
			font-family: monospace;
			white-space: pre;
			font-size: 12px;
			text-shadow:
				-1px -1px 0 white,
				 1px -1px 0 white,
				 1px  1px 0 white,
				-1px  1px 0 white;
		}
		.ticklabel {
			font-size: 10px;
			dominant-baseline: text-after-edge;
		}

		.ps { font-size: 9px;}
		/* ]]> */
	  </style>`)

	const nl = "\n"

	y := 0.0
	for _, timing := range timings {
		fmt.Println("plotting ", timing.Name)
		timing.Prepare(max)
		func() {
			write(`<g transform="translate(%.2f,%.2f)">`, pad, pad+y)
			defer write(`</g>` + nl)
			y += height

			write(`<rect x="0" y="0" width="%.2f" height="%.2f" style="fill:#f0f0f0;" />`+nl, width, height)

			var tick int

			tick = 0
			for at := min; at <= max; at += tickstepx {
				write(`<polyline `)
				if tick%majtickx == 0 {
					write(`stroke="#333" stroke-dasharray="1, 1" `)
				} else {
					write(`stroke="#666" stroke-dasharray="2, 5" `)
				}
				tick++

				write(`points="%.2f,%.2f %.2f,%2.f" />`+nl, tox(at), 0.0, tox(at), height)

				write(`<text x="%.0f" y="%.0f" class="ticklabel">%.0f</text>" />`+nl, tox(at), height, at)
			}

			tick = 0
			for p := 0.0; p <= 1; p += tickstepy {
				write(`<polyline `)
				if tick%majticky == 0 {
					write(`stroke="#333" stroke-dasharray="1, 1" `)
				} else {
					write(`stroke="#666" stroke-dasharray="2, 5" `)
				}
				tick++

				write(`points="%.2f,%.2f %.2f,%2.f" />`+nl, 0.0, toy(p), width, toy(p))
			}

			write(`<polyline class="line" points="`)

			write(`%.2f,%.2f `, tox(min), toy(0))
			index := 0
			for at := min; at <= max; at += pointstep {
				total := 0.0
				mul := 1.0 / float64(len(timing.Sanitized))

				low, high := at-kernel, at+kernel
				for ; index <= len(timing.Sanitized); index++ {
					if timing.Sanitized[index] >= low {
						break
					}
				}
				for _, time := range timing.Sanitized[index:] {
					if time > high {
						break
					}
					total += cubicPulse(at, kernel, time) * mul
				}

				write(`%.2f,%.2f `, tox(at), toy(total))
			}
			write(`%.2f,%.2f `, tox(max), toy(0))
			write(`" />` + nl)

			write(`<rect x="10"  y="5" width="%v" height="%v" style="fill:rgba(255,255,255,0.7);" />`+nl, legendwidth-20, height-10)

			write(`<text x="30"  y="15" style="font-weight: bold;">%v</text>`, timing.Name)
			write(`<text x="30"  y="30">measu=%v</text>`, len(timing.Measured))
			write(`<text x="30"  y="45">valid=%v</text>`, len(timing.Sanitized))
			write(`<text x="30"  y="60">zeros=%v</text>`, timing.Zero)
			write(`<text x="30"  y="75">overs=%v</text>`, timing.Over)

			write(`<text class="ps" x="30" y="90">avg  = %8.2f</text>`, timing.Average)
			write(`<text class="ps" x="30" y="100">.5   = %8.2f</text>`, timing.Ps[0])
			write(`<text class="ps" x="30" y="110">.9   = %8.2f</text>`, timing.Ps[1])
			write(`<text class="ps" x="30" y="120">.99  = %8.2f</text>`, timing.Ps[2])
			write(`<text class="ps" x="30" y="130">.999 = %8.2f</text>`, timing.Ps[3])
		}()
	}
}

type Timing struct {
	Name      string
	Measured  []float64
	Sanitized []float64

	Zero int
	Over int

	Average float64
	Ps      []float64
}

func (t *Timing) Prepare(max float64) {
	t.Sanitized = make([]float64, len(t.Measured))

	copy(t.Sanitized, t.Measured)
	sort.Float64s(t.Sanitized)

	for _, v := range t.Sanitized {
		if v <= 0 {
			t.Zero++
		} else {
			break
		}
	}

	avg := 0.0
	frac := 1.0 / float64(len(t.Sanitized[t.Zero:]))
	for _, v := range t.Sanitized[t.Zero:] {
		avg += v * frac
	}
	t.Average = avg
	t.Ps = quant(t.Sanitized, 0.5, 0.9, 0.99, 0.999)

	tail := len(t.Sanitized) - 1
	for ; tail >= 0; tail-- {
		if t.Sanitized[tail] < max {
			break
		}
		t.Sanitized[tail] = max
		t.Over++
	}

	t.Sanitized = t.Sanitized[t.Zero:]
	if len(t.Sanitized) == 0 {
		t.Sanitized = []float64{0}
	}
}

func quant(timings []float64, ps ...float64) []float64 {
	xs := make([]float64, len(ps))
	for i, p := range ps {
		pi := int(p * float64(len(timings)))
		if pi > len(timings) {
			pi = len(timings)
		}
		xs[i] = timings[pi]
	}
	return xs
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func cubicPulse(center, radius, at float64) float64 {
	at = at - center
	if at < 0 {
		at = -at
	}
	if at > radius {
		return 0
	}
	at /= radius
	return 1 - at*at*(3-2*at)
}
