package benchplot

import (
	"image/color"
	"math"
)

type Length float64

type Point struct{ X, Y Length }

type Plot struct {
	// Size is the visual size of this plot
	Size Point
	// X, Y are the axis information
	X, Y Axis
	// Elements
	Elements []Element
	// DefaultStyle
	Line Style
	Font Style
	Fill Style
}

type Axis struct {
	// Min value of the axis (in value space)
	Min float64
	// Max value of the axis (in value space)
	Max float64
	// Transform transform [0..1] -> float64
	Transform func(p float64) float64
}

func (axis *Axis) IsAutomatic() bool {
	return math.IsNaN(axis.Min) || math.IsNaN(axis.Max)
}

// Element is a drawable plot element
type Element interface {
	Draw(plot *Plot, canvas Canvas)
}

// Dataset represents an Element that contains data
type Dataset interface {
	Element
	Stats(precentile float64) (min, max, avg, plow, p50, phigh Point)
}

func New() *Plot {
	return &Plot{
		Size: Point{800, 600},
		X: Axis{
			Min: math.NaN(),
			Max: math.NaN(),
		},
		Y: Axis{
			Min: math.NaN(),
			Max: math.NaN(),
		},
		Line: Style{
			Color: color.NRGBA{0, 0, 0, 255},
			Fill:  nil,
			Size:  1.0,
		},
		Font: Style{
			Color: color.NRGBA{0, 0, 0, 255},
			Fill:  nil,
			Size:  1.0,
		},
		Fill: Style{
			Color: nil,
			Fill:  color.NRGBA{128, 128, 128, 255},
			Size:  1.0,
		},
	}
}

func AutomaticAxis(elements []Element) (X, Y Axis) {
	// TODO:
	return Axis{}, Axis{}
}

func (plot *Plot) Draw(canvas Canvas) {
	if plot.X.IsAutomatic() || plot.Y.IsAutomatic() {
		tmpplot := &Plot{}
		*tmpplot = *plot
		plot = tmpplot
		plot.X, plot.Y = AutomaticAxis(plot.Elements)
	}

	for _, element := range plot.Elements {
		element.Draw(plot, canvas)
	}
}

func (plot *Plot) ToCanvas(x, y float64) Point {
	var p Point
	px := (x - plot.X.Min) / (plot.X.Max - plot.X.Min)
	if plot.X.Transform != nil {
		px = plot.X.Transform(px)
	}
	p.X = Length(px) * plot.Size.X

	py := (y - plot.Y.Min) / (plot.Y.Max - plot.Y.Min)
	if plot.Y.Transform != nil {
		py = plot.Y.Transform(py)
	}
	p.Y = Length(py) * plot.Size.Y

	return p
}

type Canvas interface {
	Translate(p Point) Canvas
	Text(glyph string, at Point, style *Style)
	Line(points []Point, style *Style)
	Fill(points []Point, style *Style)
}

type Style struct {
	Color color.Color
	Fill  color.Color
	Size  Length

	// line only
	Dash       []Length
	DashOffset []Length

	// text only
	Font     string
	Rotation float64
	Origin   Point // {-1..1, -1..1}
}
