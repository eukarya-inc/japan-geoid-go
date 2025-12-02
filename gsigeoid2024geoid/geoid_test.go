package gsigeoid2024geoid

import (
	"fmt"
	"math"
	"testing"
)

func TestLoad(t *testing.T) {
	g, err := Load()
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	// グリッド情報の確認
	if g.Info.XNum != 1601 {
		t.Errorf("XNum: got %d, want 1601", g.Info.XNum)
	}
	if g.Info.YNum != 2101 {
		t.Errorf("YNum: got %d, want 2101", g.Info.YNum)
	}

	// 日本国内の座標でジオイド高を取得できることを確認
	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	if math.IsNaN(height) {
		t.Errorf("GetHeight(%f, %f): got NaN, want valid height", lng, lat)
	}
	t.Logf("LoadGeoid: GetHeight(%f, %f) = %f", lng, lat, height)
}

func Example() {
	g, err := Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Geoid height: %.6f\n", height)

	// Output:
	// Geoid height: 39.596702
}
