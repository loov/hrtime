package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/loov/hrtime"
	"github.com/loov/plot"
)

func main() {
	ConsoleHistogram()
	DensityPlot()
	TimingPlot()
}

// ConsoleHistogram demonstrates how to measure and print the output to console.
func ConsoleHistogram() {
	fmt.Println("Console Histogram")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

// TimingPlot demonstrates how to plot timing values based on the order.
func TimingPlot() {
	fmt.Println("Timing Plot (timing-plot.svg)")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}

	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = plot.R(5, 0, 5, 0)
	p.Add(stack)

	seconds := plot.DurationToSeconds(bench.Laps())

	stack.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewLine("", plot.Points(nil, seconds)),
		plot.NewTickLabels(),
	)

	svg := plot.NewSVG(800, 300)
	p.Draw(svg)
	ioutil.WriteFile("timing-plot.svg", svg.Bytes(), 0755)
}

// DensityPlot demonstrates how to create a density plot from the values.
func DensityPlot() {
	fmt.Println("Density Plot (density-plot.svg)")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(1000 * time.Nanosecond)
	}

	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = plot.R(5, 0, 5, 0)
	p.Add(stack)

	seconds := plot.DurationToSeconds(bench.Laps())

	stack.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewDensity("", seconds),
		plot.NewTickLabels(),
	)

	svg := plot.NewSVG(800, 300)
	p.Draw(svg)
	ioutil.WriteFile("density-plot.svg", svg.Bytes(), 0755)
}
