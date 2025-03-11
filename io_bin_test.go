package japangeoid

import (
	"bytes"
	"testing"
)

func TestGrid_ToBinary(t *testing.T) {
	g := MemoryGrid{
		Info: GridInfo{
			XNum:    3,
			YNum:    2,
			XDenom:  40,
			YDenom:  60,
			XMin:    120.0,
			YMin:    20.0,
			IKind:   1,
			Version: "ver2.2",
		},
		Points: []int32{
			1, 2, 3,
			4, 5, 6,
		},
	}

	buf := bytes.NewBuffer(nil)
	if err := g.ToBinary(buf); err != nil {
		t.Errorf("error (ToBinary): %v", err)
	}

	g2, err := FromBinary(buf)
	if err != nil {
		t.Errorf("error (FromBinary): %v", err)
	}

	if g.Info.XNum != g2.Info.XNum {
		t.Errorf("XNum: %v", g2.Info.XNum)
	}

	if g.Info.YNum != g2.Info.YNum {
		t.Errorf("YNum: %v", g2.Info.YNum)
	}

	if g.Info.XDenom != g2.Info.XDenom {
		t.Errorf("XDenom: %v", g2.Info.XDenom)
	}

	if g.Info.YDenom != g2.Info.YDenom {
		t.Errorf("YDenom: %v", g2.Info.YDenom)
	}

	if g.Info.XMin != g2.Info.XMin {
		t.Errorf("XMin: %v", g2.Info.XMin)
	}

	if g.Info.YMin != g2.Info.YMin {
		t.Errorf("YMin: %v", g2.Info.YMin)
	}

	if g.Info.IKind != g2.Info.IKind {
		t.Errorf("IKind: %v", g2.Info.IKind)
	}

	if g.Info.Version != g2.Info.Version {
		t.Errorf("Version: %v", g2.Info.Version)
	}

	if len(g.Points) != len(g2.Points) {
		t.Errorf("Points: %v", g2.Points)
	}

	for i := range g.Points {
		if g.Points[i] != g2.Points[i] {
			t.Errorf("Points[%d]: %v", i, g2.Points[i])
		}
	}
}
