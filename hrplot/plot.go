// Package hrplot provides convenience functions for easily creating plots out of benchmark results.

package hrplot

import (
	"fmt"
	"io/ioutil"

	"github.com/loov/plot"
)

const (
	defaultWidth  = 800
	defaultHeight = 300
)

// Benchmark declares interface for benchmarks it can plot.
type Benchmark interface {
	Name() string
	Unit() string
	Measurements() []float64
}

// Option is for declaring options to plotting.
type Option interface{ apply() }

// label formats the label to be used on plots.
func label(kind string, b Benchmark) string {
	return fmt.Sprintf("%s %s [%s]", kind, b.Name(), b.Unit())
}

// All plots line, density and percentiles plot on a single image.
func All(svgfile string, b Benchmark, options ...Option) error {
	measurements := b.Measurements()
	if len(measurements) == 0 {
		return nil
	}

	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = plot.R(5, 5, 5, 5)
	p.Add(stack)

	stack.Add(plot.NewAxisGroup(line(b, measurements)...))
	stack.Add(plot.NewAxisGroup(density(b, measurements)...))

	percentiles := plot.NewAxisGroup(percentiles(b, measurements)...)
	percentiles.X = plot.NewPercentilesAxis()
	stack.Add(percentiles)

	svg := plot.NewSVG(defaultWidth, defaultHeight*3)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Line draws a line graph in timing order.
func Line(svgfile string, b Benchmark, options ...Option) error {
	measurements := b.Measurements()
	if len(measurements) == 0 {
		return nil
	}

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.AddGroup(line(b, measurements)...)

	svg := plot.NewSVG(defaultWidth, defaultHeight)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Density draws a density plot out of benchmark measurements.
func Density(svgfile string, b Benchmark, options ...Option) error {
	measurements := b.Measurements()
	if len(measurements) == 0 {
		return nil
	}

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.AddGroup(density(b, measurements)...)

	svg := plot.NewSVG(defaultWidth, defaultHeight)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
}

// Percentiles draws a percentiles plot out of benchmark measurements.
func Percentiles(svgfile string, b Benchmark, options ...Option) error {
	measurements := b.Measurements()
	if len(measurements) == 0 {
		return nil
	}

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.X = plot.NewPercentilesAxis()
	p.AddGroup(percentiles(b, measurements)...)

	svg := plot.NewSVG(defaultWidth, defaultHeight)
	p.Draw(svg)

	return ioutil.WriteFile(svgfile, svg.Bytes(), 0755)
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
