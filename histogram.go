package hrtime

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

// HistogramOptions is configuration.
type HistogramOptions struct {
	BinCount int
	// NiceRange will try to round the bucket sizes to have a nicer output.
	NiceRange bool
	// Clamp values to either percentile or to a specific ns value.
	ClampMaximum    float64
	ClampPercentile float64
}

var defaultOptions = HistogramOptions{
	BinCount:        10,
	NiceRange:       true,
	ClampMaximum:    0,
	ClampPercentile: 0.999,
}

// Histogram is a binned historgram with different statistics.
type Histogram struct {
	Minimum float64
	Average float64
	Maximum float64

	P50, P90, P99, P999, P9999 float64

	Bins []HistogramBin

	// for pretty printing
	Width int
}

// HistogramBin is a single bin in histogram
type HistogramBin struct {
	Start    float64
	Count    int
	Width    float64
	andAbove bool
}

// NewDurationHistogram creates a histogram from time.Duration-s.
func NewDurationHistogram(durations []time.Duration, opts *HistogramOptions) *Histogram {
	nanos := make([]float64, len(durations))
	for i, d := range durations {
		nanos[i] = float64(d.Nanoseconds())
	}
	return NewHistogram(nanos, opts)
}

// NewHistogram creates a new histogram from the specified nanosecond values.
func NewHistogram(nanoseconds []float64, opts *HistogramOptions) *Histogram {
	if opts.BinCount <= 0 {
		panic("binCount must be larger than 0")
	}

	hist := &Histogram{}
	hist.Width = 40
	hist.Bins = make([]HistogramBin, opts.BinCount)
	if len(nanoseconds) == 0 {
		return hist
	}

	nanoseconds = append(nanoseconds[:0:0], nanoseconds...)
	sort.Float64s(nanoseconds)

	hist.Minimum = nanoseconds[0]
	hist.Maximum = nanoseconds[len(nanoseconds)-1]

	hist.Average = float64(0)
	for _, x := range nanoseconds {
		hist.Average += x
	}
	hist.Average /= float64(len(nanoseconds))

	p := func(p float64) float64 {
		i := int(math.Round(p * float64(len(nanoseconds))))
		if i < 0 {
			i = 0
		}
		if i >= len(nanoseconds) {
			i = len(nanoseconds) - 1
		}
		return nanoseconds[i]
	}

	hist.P50, hist.P90, hist.P99, hist.P999, hist.P9999 = p(0.50), p(0.90), p(0.99), p(0.999), p(0.9999)

	clampMaximum := hist.Maximum
	if opts.ClampPercentile > 0 {
		clampMaximum = p(opts.ClampPercentile)
	}
	if opts.ClampMaximum > 0 {
		clampMaximum = opts.ClampMaximum
	}

	var minimum, spacing float64

	if opts.NiceRange {
		minimum, spacing = calculateNiceSteps(hist.Minimum, clampMaximum, opts.BinCount)
	} else {
		minimum, spacing = calculateSteps(hist.Minimum, clampMaximum, opts.BinCount)
	}

	for i := range hist.Bins {
		hist.Bins[i].Start = spacing*float64(i) + minimum
	}
	hist.Bins[0].Start = hist.Minimum

	for _, x := range nanoseconds {
		k := int(float64(x-minimum) / spacing)
		if k < 0 {
			k = 0
		}
		if k >= opts.BinCount {
			k = opts.BinCount - 1
			hist.Bins[k].andAbove = true
		}
		hist.Bins[k].Count++
	}

	maxBin := 0
	for _, bin := range hist.Bins {
		if bin.Count > maxBin {
			maxBin = bin.Count
		}
	}

	for k := range hist.Bins {
		bin := &hist.Bins[k]
		bin.Width = float64(bin.Count) / float64(maxBin)
	}

	return hist
}

// Divide divides histogram by number of repetitions for the tests.
func (hist *Histogram) Divide(n int) {
	hist.Minimum /= float64(n)
	hist.Average /= float64(n)
	hist.Maximum /= float64(n)

	hist.P50 /= float64(n)
	hist.P90 /= float64(n)
	hist.P99 /= float64(n)
	hist.P999 /= float64(n)
	hist.P9999 /= float64(n)

	for i := range hist.Bins {
		hist.Bins[i].Start /= float64(n)
	}
}

// WriteStatsTo writes formatted statistics to w.
func (hist *Histogram) WriteStatsTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "  avg %v;  min %v;  p50 %v;  max %v;\n  p90 %v;  p99 %v;  p999 %v;  p9999 %v;\n",
		time.Duration(truncate(hist.Average, 3)),
		time.Duration(truncate(hist.Minimum, 3)),
		time.Duration(truncate(hist.P50, 3)),
		time.Duration(truncate(hist.Maximum, 3)),

		time.Duration(truncate(hist.P90, 3)),
		time.Duration(truncate(hist.P99, 3)),
		time.Duration(truncate(hist.P999, 3)),
		time.Duration(truncate(hist.P9999, 3)),
	)
	return int64(n), err
}

// WriteTo writes formatted statistics and histogram to w.
func (hist *Histogram) WriteTo(w io.Writer) (int64, error) {
	written, err := hist.WriteStatsTo(w)
	if err != nil {
		return written, err
	}

	// TODO: use consistently single unit instead of multiple
	maxCountLength := 3
	for i := range hist.Bins {
		x := (int)(math.Ceil(math.Log10(float64(hist.Bins[i].Count + 1))))
		if x > maxCountLength {
			maxCountLength = x
		}
	}

	var n int
	for _, bin := range hist.Bins {
		if bin.andAbove {
			n, err = fmt.Fprintf(w, " %10v+[%[2]*[3]v] ", time.Duration(round(bin.Start, 3)), maxCountLength, bin.Count)
		} else {
			n, err = fmt.Fprintf(w, " %10v [%[2]*[3]v] ", time.Duration(round(bin.Start, 3)), maxCountLength, bin.Count)
		}

		written += int64(n)
		if err != nil {
			return written, err
		}

		width := float64(hist.Width) * bin.Width
		frac := width - math.Trunc(width)

		n, err = io.WriteString(w, strings.Repeat("█", int(width)))
		written += int64(n)
		if err != nil {
			return written, err
		}

		if frac > 0.5 {
			n, err = io.WriteString(w, `▌`)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}

		n, err = fmt.Fprintf(w, "\n")
		written += int64(n)
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

// StringStats returns a string representation of the histogram stats.
func (hist *Histogram) StringStats() string {
	var buffer strings.Builder
	_, _ = hist.WriteStatsTo(&buffer)
	return buffer.String()
}

// String returns a string representation of the histogram.
func (hist *Histogram) String() string {
	var buffer strings.Builder
	_, _ = hist.WriteTo(&buffer)
	return buffer.String()
}
