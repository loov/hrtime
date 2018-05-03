package benchplot

import (
	"image/color"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

type Density struct {
	Width  float64
	Height float64

	// when LogX > 0, it will use a logarithmic X axis
	LogX float64
	// when LogY > 0, it will use a logarithmic Y axis
	LogY float64

	// Data parameters
	Min float64 // NaN -> auto-detect
	Max float64 // NaN -> auto-detect

	// Lines to draw
	Lines []*Line
}

type Line struct {
	Label string
	Color color.Color
	Data  []float64 // sorted
}

func NewDensity() *Density {
	plot := &Density{}
	plot.LogX = 50.0
	plot.LogY = 50.0
	plot.Min = math.NaN()
	plot.Max = math.NaN()
	return plot
}

func (plot *Density) LineFloat64s(label string, values []float64) {
	line := &Line{}
	line.Label = label
	line.Data = values
	sort.Float64s(line.Data)
	plot.Lines = append(plot.Lines, line)
}

func (plot *Density) LineNanoseconds(label string, durations []time.Duration) {
	values := make([]float64, len(durations))
	for i, dur := range durations {
		values[i] = float64(dur.Nanoseconds())
	}
}

func (plot *Density) SVG() string {
	var w strings.Builder
	_ = plot.WriteSVG(&w)
	return w.String()
}

func (plot *Density) WriteSVG(w io.Writer) error {
	return nil
}

func (line *Line) density(count int, min, max, kernel float64, plot func(x, y float64)) {
	step := (max - min) / float64(count)
	index := sort.SearchFloat64s(line.Data, min-kernel)
	invkernel := 1 / kernel
	for at := min; at <= max; at += step {
		sample := 0.0
		low, high := at-kernel, at+kernel
		for ; index < len(line.Data); index++ {
			if line.Data[index] >= low {
				break
			}
		}
		for _, time := range line.Data[index:] {
			if time > high {
				break
			}
			sample += cubicPulse(at, kernel, invkernel, time)
		}
		plot(at, sample)
	}
}

func cubicPulse(center, radius, invradius, at float64) float64 {
	at = at - center
	if at < 0 {
		at = -at
	}
	if at > radius {
		return 0
	}
	at *= invradius
	return 1 - at*at*(3-2*at)
}
