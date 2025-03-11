package japangeoid

import (
	"math"
	"testing"
)

var _ Grid = (*MemoryGrid)(nil)

func TestMemoryGrid_GetHeight(t *testing.T) {
	g := &MemoryGrid{
		Info: GridInfo{
			XNum:   3,
			YNum:   3,
			XDenom: 1,
			YDenom: 1,
			XMin:   0,
			YMin:   0,
		},
		Points: []int32{
			0, 1, 2,
			3, 4, 5,
			6, 7, 8,
		},
	}

	tests := [][]float64{
		{0.0, 0.0, 0.0},
		{0.0, 1.0, 0.0003},
		{1.0, 0.0, 0.0001},
		{1.0, 1.0, 0.0004},
		{2.0, 0.0, 0.0002},
		{0.0, 2.0, 0.0006},
		{2.0, 2.0, 0.0008},
		{1.5, 2.0, 0.00075},
		{-1.0, -1.0, math.NaN()},
		{3.0, 3.0, math.NaN()},
	}

	for _, test := range tests {
		x, y, expected := test[0], test[1], test[2]
		h := g.GetHeight(x, y)

		if math.IsNaN(expected) {
			if !math.IsNaN(h) {
				t.Errorf("invalid height: %v", h)
			}
			continue
		}

		if h != expected {
			t.Errorf("invalid height: %v", h)
		}
	}
}
