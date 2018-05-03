package benchplot

import (
	"image/color"
	"time"
)

type Density struct {
	// Visual parameters
	Width  float64
	Height float64
	// when LogX > 0, it will use a logarithmic X axis
	LogX float64
	// when LogY > 0, it will use a logarithmic Y axis
	LogY float64

	// Data parameters
	Min float64
	Max float64

	Lines []Density
}

type Line struct {
	Label  string
	Color  color.Color
	Values []float64
}

func (line *Line) densityCubic(count int, min, max, kernel float64) []float64 {
	result := make([]float64, count)
	return result
}

func NewDensity() *Density { return &Density{} }

func (plot *Density) LineDurations(label string, durations []time.Duration) {
}
