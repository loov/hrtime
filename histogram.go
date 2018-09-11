package hrtime

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

type Histogram struct {
	Minimum float64
	Average float64
	Maximum float64

	P50, P90, P99, P999, P9999 float64

	Bins []HistogramBin

	// for pretty printing
	Width int
}

type HistogramBin struct {
	Start float64
	Count int
	Width float64
}

func NewDurationHistogram(durations []time.Duration, binCount int) *Histogram {
	nanos := make([]float64, len(durations))
	for i, d := range durations {
		nanos[i] = float64(d.Nanoseconds())
	}
	return NewHistogram(nanos, binCount)
}

func NewHistogram(nanoseconds []float64, binCount int) *Histogram {
	if binCount <= 1 {
		panic("binCount must be larger than 0")
	}

	hist := &Histogram{}
	hist.Width = 40
	hist.Bins = make([]HistogramBin, binCount)
	if len(nanoseconds) == 0 {
		return hist
	}

	nanoseconds = append(nanoseconds[:0:0], nanoseconds...)
	sort.Float64s(nanoseconds)

	hist.Minimum = nanoseconds[0]
	hist.Maximum = nanoseconds[len(nanoseconds)-1]

	hist.Average = nanoseconds[0]
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

	hist.P50 = p(0.50)
	hist.P90 = p(0.90)
	hist.P99 = p(0.99)
	hist.P999 = p(0.999)
	hist.P9999 = p(0.9999)

	stepSize := (hist.Maximum - hist.Minimum) / float64(binCount)

	for i := range hist.Bins {
		hist.Bins[i].Start = stepSize*float64(i) + hist.Minimum
	}

	for _, x := range nanoseconds {
		k := int(float64(x-hist.Minimum) / stepSize)
		if k < 0 {
			k = 0
		}
		if k >= binCount {
			k = binCount - 1
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
		n, err = fmt.Fprintf(w, " %10v [%[2]*[3]v] ", time.Duration(round(bin.Start, 3)), maxCountLength, bin.Count)
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
			n, err = io.WriteString(w, "▌")
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

func (hist *Histogram) String() string {
	var buffer strings.Builder
	hist.WriteTo(&buffer)
	return buffer.String()
}
