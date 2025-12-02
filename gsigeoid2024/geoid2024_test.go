package gsigeoid2024

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
	t.Logf("Load: GetHeight(%f, %f) = %f", lng, lat, height)

	// 範囲外の座標はNaNを返すことを確認
	lng, lat = 10, 10
	height = g.GetHeight(lng, lat)
	if !math.IsNaN(height) {
		t.Errorf("GetHeight(%f, %f): got %f, want NaN", lng, lat, height)
	}
}

func TestLoadHrefconv(t *testing.T) {
	g, err := LoadHrefconv()
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

	// 日本国内の座標で値を取得できることを確認
	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	// Hrefconvは本土ではNaNになる可能性があるのでログ出力のみ
	t.Logf("LoadHrefconv: GetHeight(%f, %f) = %f", lng, lat, height)
}

func Example() {
	g, err := Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Input: (lng: %.6f, lat: %.6f) -> Geoid height: %.6f\n", lng, lat, height)

	lng, lat = 10, 10
	height = g.GetHeight(lng, lat)
	fmt.Printf("Input: (lng: %.6f, lat: %.6f) -> Geoid height: %.6f\n", lng, lat, height)

	// Output:
	// Input: (lng: 138.283982, lat: 37.123786) -> Geoid height: 39.596702
	// Input: (lng: 10.000000, lat: 10.000000) -> Geoid height: NaN
}

func ExampleLoadHrefconv() {
	g, err := LoadHrefconv()
	if err != nil {
		panic(err)
	}

	// 小笠原諸島（父島）付近 - 基準面補正量がある
	lng, lat := 142.19, 27.09
	height := g.GetHeight(lng, lat)
	fmt.Printf("Hrefconv: %.6f\n", height)

	// Output:
	// Hrefconv: 0.642000
}
