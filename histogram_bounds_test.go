package hrtime

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		in    float64
		trunc int
		exp   float64
	}{
		{0, 1, 0},
		{0, 2, 0},
		{0, 3, 0},

		{10, 1, 10},
		{10, 2, 10},
		{10, 3, 10},

		{1234, 1, 1000},
		{1234, 2, 1200},
		{1234, 3, 1230},
	}

	for _, test := range tests {
		got := truncate(test.in, test.trunc)
		if got != test.exp {
			t.Errorf("%f %d => %f expected %f", test.in, test.trunc, got, test.exp)
		}
	}
}
