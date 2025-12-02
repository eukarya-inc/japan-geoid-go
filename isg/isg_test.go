package isg

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseDMS(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"15°00'00\"", 15.0},
		{"50°00'00\"", 50.0},
		{"120°00'00\"", 120.0},
		{"160°00'00\"", 160.0},
		{"0°01'00\"", 1.0 / 60.0},            // 1 minute = 1/60 degree
		{"0°01'30\"", 1.5 / 60.0},            // 1 min 30 sec = 1.5/60 deg = 1/40 deg
		{"  15°00'00\"  ", 15.0},             // leading/trailing spaces
		{"-10°30'00\"", -10.5},               // negative value
		{"35°41'22.4\"", 35.0 + 41.0/60.0 + 22.4/3600.0}, // decimal seconds
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseDMS(tt.input)
			if err != nil {
				t.Fatalf("ParseDMS(%q) returned error: %v", tt.input, err)
			}
			// Use tolerance for floating point comparison
			if diff := result - tt.expected; diff < -1e-9 || diff > 1e-9 {
				t.Errorf("ParseDMS(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDMS_Invalid(t *testing.T) {
	tests := []string{
		"invalid",
		"15",
		"15°",
		"15°00'",
		"abc°00'00\"",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseDMS(input)
			if err == nil {
				t.Errorf("ParseDMS(%q) should return error", input)
			}
		})
	}
}

func TestParse(t *testing.T) {
	// Parse a simple ISG file
	isgContent := `begin_of_head ================================================
model name     : TestModel
lat min        =  35°00'00"
lat max        =  36°00'00"
lon min        = 139°00'00"
lon max        = 140°00'00"
delta lat      =   0°30'00"
delta lon      =   0°30'00"
nrows          =        3
ncols          =        3
nodata         =  -9999.0000
data ordering  : N-to-S, W-to-E
ISG format     =         2.0
end_of_head ==================================================
1.0 2.0 3.0
4.0 5.0 6.0
7.0 8.0 9.0
`

	grid, err := Parse(strings.NewReader(isgContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify header
	if grid.Header.ModelName != "TestModel" {
		t.Errorf("ModelName: got %q, want %q", grid.Header.ModelName, "TestModel")
	}
	if grid.Header.NRows != 3 {
		t.Errorf("NRows: got %d, want 3", grid.Header.NRows)
	}
	if grid.Header.NCols != 3 {
		t.Errorf("NCols: got %d, want 3", grid.Header.NCols)
	}
	if grid.Header.LatMin != 35.0 {
		t.Errorf("LatMin: got %f, want 35.0", grid.Header.LatMin)
	}
	if grid.Header.LonMin != 139.0 {
		t.Errorf("LonMin: got %f, want 139.0", grid.Header.LonMin)
	}
	if grid.Header.DataOrdering != "N-to-S, W-to-E" {
		t.Errorf("DataOrdering: got %q, want %q", grid.Header.DataOrdering, "N-to-S, W-to-E")
	}

	// Verify point data
	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}
	if len(grid.Points) != len(expected) {
		t.Fatalf("Points length: got %d, want %d", len(grid.Points), len(expected))
	}
	for i, v := range expected {
		if grid.Points[i] != v {
			t.Errorf("Points[%d]: got %f, want %f", i, grid.Points[i], v)
		}
	}
}

func Example() {
	isgContent := `begin_of_head ================================================
model name     : ExampleModel
lat min        =  35°00'00"
lat max        =  36°00'00"
lon min        = 139°00'00"
lon max        = 140°00'00"
delta lat      =   1°00'00"
delta lon      =   1°00'00"
nrows          =        2
ncols          =        2
nodata         =  -9999.0000
data ordering  : N-to-S, W-to-E
ISG format     =         2.0
end_of_head ==================================================
1.0 2.0
3.0 4.0
`

	grid, err := Parse(strings.NewReader(isgContent))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Model: %s\n", grid.Header.ModelName)
	fmt.Printf("Rows: %d, Cols: %d\n", grid.Header.NRows, grid.Header.NCols)
	fmt.Printf("Lat: %.1f - %.1f\n", grid.Header.LatMin, grid.Header.LatMax)
	fmt.Printf("Lon: %.1f - %.1f\n", grid.Header.LonMin, grid.Header.LonMax)

	// Output:
	// Model: ExampleModel
	// Rows: 2, Cols: 2
	// Lat: 35.0 - 36.0
	// Lon: 139.0 - 140.0
}

func ExampleParseDMS() {
	deg, err := ParseDMS("35°41'22.4\"")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.6f\n", deg)

	// Output:
	// 35.689556
}
