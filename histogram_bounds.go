package hrtime

import "math"

func calculateSteps(min, max float64, bincount int) (minimum, spacing float64) {
	minimum = min
	spacing = (max - min) / float64(bincount)
	return minimum, spacing
}

func calculateNiceSteps(min, max float64, bincount int) (minimum, spacing float64) {
	span := niceNumber(max-min, false)
	spacing = niceNumber(span/float64(bincount-1), true)
	minimum = math.Floor(min/spacing) * spacing
	return minimum, spacing
}

func niceNumber(span float64, round bool) float64 {
	exp := math.Floor(math.Log10(span))
	frac := span / math.Pow(10, exp)

	var nice float64
	if round {
		switch {
		case frac < 1.5:
			nice = 1
		case frac < 3:
			nice = 2
		case frac < 7:
			nice = 5
		default:
			nice = 10
		}
	} else {
		switch {
		case frac <= 1:
			nice = 1
		case frac <= 2:
			nice = 2
		case frac <= 5:
			nice = 5
		default:
			nice = 10
		}
	}

	return nice * math.Pow(10, exp)
}

func truncate(v float64, digits int) float64 {
	if digits == 0 || v == 0 {
		return 0
	}

	scale := math.Pow(10, math.Floor(math.Log10(v))+1-float64(digits))
	return scale * math.Trunc(v/scale)
}

func round(v float64, digits int) float64 {
	if digits == 0 || v == 0 {
		return 0
	}

	scale := math.Pow(10, math.Floor(math.Log10(v))+1-float64(digits))
	return scale * math.Round(v/scale)
}
