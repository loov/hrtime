package hrtime

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

type Histogram struct {
	Minimum time.Duration
	Average time.Duration
	Maximum time.Duration

	Bins []HistogramBin

	// for pretty printing
	Width int
}

type HistogramBin struct {
	Start time.Duration
	Count int
	Width float64
}

func NewHistogram(timing []time.Duration, binCount int) *Histogram {
	if binCount < 0 {
		panic("binCount must be larger than 0")
	}

	hist := &Histogram{}
	hist.Width = 40
	hist.Bins = make([]HistogramBin, binCount)
	if len(timing) == 0 {
		return hist
	}

	hist.Minimum = timing[0]
	hist.Average = timing[0] // TODO: fix potential overflow
	hist.Maximum = timing[0]

	for _, x := range timing {
		hist.Average += x
		if hist.Average < 0 {
			panic("average overflow")
		}
		if x < hist.Minimum {
			hist.Minimum = x
		}
		if x > hist.Maximum {
			hist.Maximum = x
		}
	}

	hist.Average /= time.Duration(len(timing))

	stepSize := float64(hist.Maximum-hist.Minimum) / float64(binCount)
	for i := range hist.Bins {
		hist.Bins[i].Start = time.Duration(stepSize*float64(i)) + hist.Minimum
	}
	for _, x := range timing {
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

func (hist *Histogram) WriteTo(w io.Writer) error {
	// TODO: use consistently single unit instead of multiple
	for _, bin := range hist.Bins {
		_, err := fmt.Fprintf(w, " %10v [%6v] ", bin.Start, bin.Count)
		if err != nil {
			return err
		}

		width := float64(hist.Width) * bin.Width
		frac := width - math.Trunc(width)

		if _, err = io.WriteString(w, strings.Repeat("█", int(width))); err != nil {
			return err
		}
		if frac > 0.5 {
			if _, err = io.WriteString(w, "▌"); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (hist *Histogram) String() string {
	var buffer strings.Builder
	_ = hist.WriteTo(&buffer)
	return buffer.String()
}
