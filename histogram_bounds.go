package hrtime

import "math"

func niceNumber(span float64, round bool) float64 {
	exp := math.Floor(math.Log10(span))
	frac := span / math.Pow(10, exp)

	var nice float64
	if round {
		nice = math.Round(frac)
	} else {
		nice = math.Ceil(frac)
	}

	return nice * math.Pow(10, exp)
}

func truncate(v float64, digits int) float64 {
	if digits == 0 {
		return 0
	}

	scale := math.Pow(10, math.Floor(math.Log10(v))+1-float64(digits))
	return scale * math.Trunc(v/scale)
}
