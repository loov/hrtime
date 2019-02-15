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
	TimingAndDensityPlot()
}

// ConsoleHistogram demonstrates how to measure and print the output to console.
func ConsoleHistogram() {
	fmt.Println("Console Histogram")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(5000 * time.Nanosecond)
	}
	fmt.Println(bench.Histogram(10))
}

// TimingPlot demonstrates how to plot timing values based on the order.
func TimingPlot() {
	fmt.Println("Timing Plot (timing-plot.svg)")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(5000 * time.Nanosecond)
	}

	seconds := plot.DurationToSeconds(bench.Laps())

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.AddGroup(
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
		time.Sleep(5000 * time.Nanosecond)
	}

	seconds := plot.DurationToSeconds(bench.Laps())

	p := plot.New()
	p.Margin = plot.R(5, 0, 0, 5)
	p.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewDensity("", seconds),
		plot.NewTickLabels(),
	)

	svg := plot.NewSVG(800, 300)
	p.Draw(svg)
	ioutil.WriteFile("density-plot.svg", svg.Bytes(), 0755)
}

// TimingAndDensityPlot demonstrates how to combine both plots
func TimingAndDensityPlot() {
	fmt.Println("Stacked Plot (stacked-plot.svg)")

	bench := hrtime.NewBenchmark(4 << 10)
	for bench.Next() {
		time.Sleep(5000 * time.Nanosecond)
	}

	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = plot.R(5, 5, 5, 5)
	p.Add(stack)

	seconds := plot.DurationToSeconds(bench.Laps())

	lineplot := plot.NewAxisGroup()
	stack.Add(lineplot)
	lineplot.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewLine("", plot.Points(nil, seconds)),
		plot.NewTickLabels(),
	)

	densityplot := plot.NewAxisGroup()
	stack.Add(densityplot)
	densityplot.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		plot.NewDensity("", seconds),
		plot.NewTickLabels(),
	)

	svg := plot.NewSVG(800, 600)
	p.Draw(svg)
	ioutil.WriteFile("stacked-plot.svg", svg.Bytes(), 0755)
}
