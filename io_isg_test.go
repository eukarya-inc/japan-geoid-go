package japangeoid

import (
	"strings"
	"testing"
)

func TestFromIsg(t *testing.T) {
	isgData := `begin_of_head ================================================
model name     : TestModel
lat min        =  35°00'00"
lat max        =  36°00'00"
lon min        = 139°00'00"
lon max        = 140°00'00"
delta lat      =   0°01'00"
delta lon      =   0°01'30"
nrows          =        2
ncols          =        2
nodata         =  -9999.0000
data ordering  : N-to-S, W-to-E
ISG format     =         2.0
end_of_head ==================================================
10.0000 20.0000
30.0000 40.0000
`

	grid, err := FromIsg(strings.NewReader(isgData))
	if err != nil {
		t.Fatalf("FromIsg failed: %v", err)
	}

	if grid.Info.XNum != 2 || grid.Info.YNum != 2 {
		t.Errorf("unexpected grid dimensions: %dx%d", grid.Info.XNum, grid.Info.YNum)
	}
	if grid.Info.XDenom != 40 || grid.Info.YDenom != 60 {
		t.Errorf("unexpected denoms: XDenom=%d, YDenom=%d", grid.Info.XDenom, grid.Info.YDenom)
	}
	if grid.Info.Version != "TestModel" {
		t.Errorf("unexpected version: %s", grid.Info.Version)
	}

	expected := []int32{300000, 400000, 100000, 200000}
	for i, exp := range expected {
		if grid.Points[i] != exp {
			t.Errorf("Points[%d]: got %d, want %d", i, grid.Points[i], exp)
		}
	}
}

func TestFromIsg_Precision(t *testing.T) {
	isgData := `begin_of_head ================================================
model name     : PrecisionTest
lat min        =  35°00'00"
lat max        =  36°00'00"
lon min        = 139°00'00"
lon max        = 140°00'00"
delta lat      =   0°30'00"
delta lon      =   0°30'00"
nrows          =        1
ncols          =        2
nodata         =  -9999.0000
data ordering  : S-to-N, W-to-E
ISG format     =         2.0
end_of_head ==================================================
0.0003 -9999.0000
`

	grid, err := FromIsg(strings.NewReader(isgData))
	if err != nil {
		t.Fatalf("FromIsg failed: %v", err)
	}

	if grid.Points[0] != 3 {
		t.Errorf("Points[0]: got %d, want 3", grid.Points[0])
	}
	if grid.Points[1] != 9990000 {
		t.Errorf("Points[1]: got %d, want 9990000 (nodata)", grid.Points[1])
	}
}
