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
	density = flag.Int("density", 1e4, "points per line for graph")
	kernel  = flag.Float64("kernel", 0.2, "kernel width")
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
	rdtsc := make([]uint64, N)
	rdtscp := make([]uint64, N)

	var beforRDTSC, afterRDTSC, beforRDTSCP, afterRDTSCP time.Time

	debug.SetGCPercent(-1)
	{
		runtime.GC()

		for i := range times {
			times[i] = time.Now()
		}

		runtime.GC()

		for i := range timenanos {
			timenanos[i] = time.Now().UnixNano()
		}

		runtime.GC()

		start := time.Now()
		for i := range deltas {
			deltas[i] = time.Since(start)
		}

		runtime.GC()

		for i := range nanotimes {
			nanotimes[i] = nanotime()
		}

		runtime.GC()

		beforRDTSC = time.Now()
		for i := range rdtsc {
			rdtsc[i] = hrtime.RDTSC()
		}
		afterRDTSC = time.Now()

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
	ns_rdtsc := make([]float64, N-offset)
	ns_rdtscp := make([]float64, N-offset)

	for i := range ns_times {
		ns_times[i] = float64(times[offset+i].Sub(times[offset+i-1]).Nanoseconds())
		ns_timenanos[i] = float64(timenanos[offset+i] - timenanos[offset+i-1])
		ns_deltas[i] = float64((deltas[offset+i] - deltas[offset+i-1]).Nanoseconds())
		ns_nanotimes[i] = float64(nanotimes[offset+i] - nanotimes[offset+i-1])
		ns_rdtsc[i] = float64(rdtsc[offset+i]-rdtsc[offset+i-1]) * rdtscCalibration
		ns_rdtscp[i] = float64(rdtsc[offset+i]-rdtsc[offset+i-1]) * rdtscpCalibration
	}

	timings := []*Timing{
		{Name: "time.Now", Measured: ns_times},
		{Name: "time.UnixNano", Measured: ns_timenanos},
		{Name: "time.Since", Measured: ns_deltas},
		{Name: "nanotime", Measured: ns_nanotimes},
		{Name: "RDTSC", Measured: ns_rdtsc},
		{Name: "RDTSCP", Measured: ns_rdtscp},
	}

	out, err := os.Create(*svg)
	check(err)
	defer out.Close()

	buf := bufio.NewWriter(out)
	defer buf.Flush()

	plot(buf, timings)
}

func plot(w io.Writer, timings []*Timing) {
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
		text {
			font-family: monospace;
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
		/* ]]> */
	  </style>`)

	const nl = "\n"

	y := 0.0
	for _, timing := range timings {
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

			write(`<polyline class="line" stroke="#000" `)
			write(`points="`)
			for at := min; at <= max; at += pointstep {
				total := 0.0
				mul := 1.0 / float64(len(timing.Sanitized))

				low, high := at-kernel, at+kernel
				i := sort.SearchFloat64s(timing.Sanitized, low)
				for _, time := range timing.Sanitized[i:] {
					if time > high {
						break
					}
					total += cubicPulse(at, kernel, time) * mul
				}

				write(`%.2f,%.2f `, tox(at), toy(total))
			}
			write(`" />` + nl)

			write(`<rect x="10" y="10" width="%v" height="%v" style="fill:rgba(255,255,255,0.7);" />`+nl, legendwidth-20, height-20)
			write(`<text x="30" y="30" style="font-weight: bold;">%v</text>`, timing.Name)
			write(`<text x="30" y="45">measu=%v</text>`, len(timing.Measured))
			write(`<text x="30" y="60">valid=%v</text>`, len(timing.Sanitized))
			write(`<text x="30" y="75">zeros=%v</text>`, timing.Zero)
			write(`<text x="30" y="90">overs=%v</text>`, timing.Over)
		}()
	}
}

type Timing struct {
	Name      string
	Measured  []float64
	Sanitized []float64

	Zero int
	Over int
}

func (t *Timing) Prepare(max float64) {
	t.Sanitized = make([]float64, len(t.Measured))

	copy(t.Sanitized, t.Measured)
	sort.Float64s(t.Sanitized)
	tail := len(t.Sanitized) - 1
	for ; tail >= 0; tail-- {
		if t.Sanitized[tail] < max {
			break
		}
		t.Sanitized[tail] = max
		t.Over++
	}

	for _, v := range t.Sanitized {
		if v <= 0 {
			t.Zero++
		} else {
			break
		}
	}

	t.Sanitized = t.Sanitized[t.Zero:]
	if len(t.Sanitized) == 0 {
		t.Sanitized = []float64{0}
	}
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
