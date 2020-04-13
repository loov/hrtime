// Package hrplot provides convenience functions for easily creating plots out of benchmark results.

package hrplot

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"

	"github.com/loov/plot"
	"github.com/loov/plot/plotsvg"
)

// Benchmark declares interface for benchmarks it can plot.
type Benchmark interface {
	Name() string
	Unit() string
	Float64s() []float64
}

// Option is for declaring options to plotting.
type Option interface{ apply(*plotOptions) }

type plotOptions struct {
	Width  float64
	Height float64

	LineClip       float64
	DensityClip    float64
	PercentileClip float64
}

func applyAll(opts ...Option) plotOptions {
	result := defaultOptions
	for _, opt := range opts {
		opt.apply(&result)
	}
	return result
}

var defaultOptions = plotOptions{
	Width:  800,
	Height: 300,

	LineClip:       0.9999,
	DensityClip:    0.99,
	PercentileClip: 0.9995,
}

type optionFunc func(*plotOptions)

func (fn optionFunc) apply(options *plotOptions) { fn(options) }

// ClipPercentile specifies how to clip the values.
func ClipPercentile(percentile float64) Option {
	return optionFunc(func(options *plotOptions) {
		options.LineClip = percentile
		options.DensityClip = percentile
		options.PercentileClip = percentile
	})
}

// label formats the label to be used on plots.
func label(kind string, b Benchmark) string {
	return fmt.Sprintf("%s [%s] %s", b.Name(), b.Unit(), kind)
}

// All plots line, density and percentiles plot on a single image.
func All(svgfile string, b Benchmark, opts ...Option) error {
	measurements := b.Float64s()
	if len(measurements) == 0 {
		return nil
	}
	options := applyAll(opts...)

	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = plot.R(5, 5, 5, 5)
	p.Add(stack)

	line := plot.NewAxisGroup(lineOptimized(b, measurements)...)
	line.Y.Max = percentile(measurements, options.LineClip)
	stack.Add(line)

	density := plot.NewAxisGroup(density(b, measurements)...)
	density.X.Max = percentile(measurements, options.DensityClip)
	stack.Add(density)

	percentiles := plot.NewAxisGroup(percentiles(b, measurements)...)
	percentiles.X = plot.NewPercentilesAxis()
	percentiles.X.Transform = plot.NewPercentileTransform(4)
	percentiles.Y.Min, percentiles.Y.Max = 0, percentile(measurements, options.PercentileClip)
	stack.Add(percentiles)

	svg := plotsvg.New(options.Width, options.Height*3)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Line draws a line graph in timing order.
func Line(svgfile string, b Benchmark, opts ...Option) error {
	measurements := b.Float64s()
	if len(measurements) == 0 {
		return nil
	}
	options := applyAll(opts...)

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.Y.Max = percentile(measurements, options.LineClip)
	p.AddGroup(lineOptimized(b, measurements)...)

	svg := plotsvg.New(options.Width, options.Height)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Density draws a density plot out of benchmark measurements.
func Density(svgfile string, b Benchmark, opts ...Option) error {
	measurements := b.Float64s()
	if len(measurements) == 0 {
		return nil
	}
	options := applyAll(opts...)

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.X.Max = percentile(measurements, options.DensityClip)
	p.AddGroup(density(b, measurements)...)

	svg := plotsvg.New(options.Width, options.Height)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Percentiles draws a percentiles plot out of benchmark measurements.
func Percentiles(svgfile string, b Benchmark, opts ...Option) error {
	measurements := b.Float64s()
	if len(measurements) == 0 {
		return nil
	}
	options := applyAll(opts...)

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.X = plot.NewPercentilesAxis()
	p.X.Transform = plot.NewPercentileTransform(4)
	p.Y.Min, p.Y.Max = 0, percentile(measurements, options.PercentileClip)
	p.AddGroup(percentiles(b, measurements)...)

	svg := plotsvg.New(options.Width, options.Height)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

func percentile(measurements []float64, p float64) float64 {
	sorted := append(measurements[:0:0], measurements...)
	sort.Float64s(sorted)
	k := int(math.Ceil(p * float64(len(sorted))))
	if k >= len(sorted) {
		k = len(sorted) - 1
	}
	return sorted[k]
}

func line(b Benchmark, measurements []float64) []plot.Element {
	return []plot.Element{
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewLine(b.Unit(), plot.Points(nil, measurements)),
		plot.NewTickLabels(),
		plot.NewTextbox(label("line", b)),
	}
}

func lineOptimized(b Benchmark, measurements []float64) []plot.Element {
	return []plot.Element{
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewOptimizedLine(b.Unit(), plot.Points(nil, measurements), 2),
		plot.NewTickLabels(),
		plot.NewTextbox(label("line", b)),
	}
}

func density(b Benchmark, measurements []float64) []plot.Element {
	return []plot.Element{
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewDensity(b.Unit(), measurements),
		plot.NewTickLabels(),
		plot.NewTextbox(label("density", b)),
	}
}

func percentiles(b Benchmark, measurements []float64) []plot.Element {
	return []plot.Element{
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewPercentiles(b.Unit(), measurements),
		plot.NewTickLabels(),
		plot.NewTextbox(label("percentiles", b)),
	}
}
