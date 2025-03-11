package japangeoid

import (
	"bufio"
	"strings"
	"testing"
)

const data = `  20.00000 120.00000 0.016667 0.025000 2 3 1 ver2.2
12.3456 23.4567 34.5678
45.6789 56.7890 67.8901
`

func TestFromAsc(t *testing.T) {
	r := strings.NewReader(data)
	g, err := FromAsc(r)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(g.Points) != 6 {
		t.Errorf("len(Points): %v", len(g.Points))
	}

	if g.Points[0] != 123456 {
		t.Errorf("Points[0]: %v", g.Points[0])
	}

	if g.Points[1] != 234567 {
		t.Errorf("Points[1]: %v", g.Points[1])
	}

	if g.Points[2] != 345678 {
		t.Errorf("Points[2]: %v", g.Points[2])
	}

	if g.Points[3] != 456789 {
		t.Errorf("Points[3]: %v", g.Points[3])
	}

	if g.Points[4] != 567890 {
		t.Errorf("Points[4]: %v", g.Points[4])
	}

	if g.Points[5] != 678901 {
		t.Errorf("Points[5]: %v", g.Points[5])
	}
}

func TestGridInfoFromAsc(t *testing.T) {
	s := bufio.NewScanner(strings.NewReader(data))
	info, err := gridInfoFromAsc(s)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if info.XMin != 120.0 {
		t.Errorf("XMin: %v", info.XMin)
	}

	if info.YMin != 20.0 {
		t.Errorf("YMin: %v", info.YMin)
	}

	if info.XNum != 3 {
		t.Errorf("XNum: %v", info.XNum)
	}

	if info.YNum != 2 {
		t.Errorf("YNum: %v", info.YNum)
	}

	if info.XDenom != 40 {
		t.Errorf("XDenom: %v", info.XDenom)
	}

	if info.YDenom != 60 {
		t.Errorf("YDenom: %v", info.YDenom)
	}

	if info.IKind != 1 {
		t.Errorf("IKind: %v", info.IKind)
	}

	if info.Version != "ver2.2" {
		t.Errorf("Version: %v", info.Version)
	}
}
