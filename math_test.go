package japangeoid

import (
	"math"
	"testing"
)

func TestBilinear(t *testing.T) {
	cases := []struct {
		name               string
		xFrac, yFrac       float64
		v00, v01, v10, v11 float64
		want               float64
	}{
		{
			name:  "All zero fraction => v00",
			xFrac: 0.0, yFrac: 0.0,
			v00: 1.0, v01: 2.0, v10: 3.0, v11: 4.0,
			want: 1.0,
		},
		{
			name:  "x=1.0, y=0.0 => v01",
			xFrac: 1.0, yFrac: 0.0,
			v00: 1.0, v01: 2.0, v10: 3.0, v11: 4.0,
			want: 2.0,
		},
		{
			name:  "x=0.0, y=1.0 => v10",
			xFrac: 0.0, yFrac: 1.0,
			v00: 1.0, v01: 2.0, v10: 3.0, v11: 4.0,
			want: 3.0,
		},
		{
			name:  "x=1.0, y=1.0 => v11",
			xFrac: 1.0, yFrac: 1.0,
			v00: 1.0, v01: 2.0, v10: 3.0, v11: 4.0,
			want: 4.0,
		},
		{
			name:  "x=0.5, y=0.0 => midpoint between v00 and v01",
			xFrac: 0.5, yFrac: 0.0,
			v00: 10.0, v01: 20.0, v10: 30.0, v11: 40.0,
			want: 15.0,
		},
		{
			name:  "x=0.0, y=0.5 => midpoint between v00 and v10",
			xFrac: 0.0, yFrac: 0.5,
			v00: 10.0, v01: 20.0, v10: 30.0, v11: 40.0,
			want: 20.0,
		},
		{
			name:  "x=0.5, y=0.5 => center (average) of v00=10, v01=20, v10=30, v11=40 => 25",
			xFrac: 0.5, yFrac: 0.5,
			v00: 10.0, v01: 20.0, v10: 30.0, v11: 40.0,
			want: 25.0,
		},
		{
			name:  "If some corner is NaN => result depends on fraction",
			xFrac: 1.0, yFrac: 0.5,
			v00: 10.0, v01: math.NaN(), v10: 30.0, v11: 40.0,
			// xFrac=1.0 => v00*(1−xFrac) is 0. v01*xFracがNaN => result NaN?
			// yFrac=0.5 => mix with v10, v11
			// Actually, (NaN*(1−yFrac) + v11*yFrac) => NaN + 20 => NaN
			// => want NaN
			want: math.NaN(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bilinear(c.xFrac, c.yFrac, c.v00, c.v01, c.v10, c.v11)
			if math.IsNaN(c.want) {
				if !math.IsNaN(got) {
					t.Errorf("got %.6f, want NaN", got)
				}
			} else if math.Abs(got-c.want) > 1e-9 {
				t.Errorf("got %.6f, want %.6f", got, c.want)
			}
		})
	}
}
