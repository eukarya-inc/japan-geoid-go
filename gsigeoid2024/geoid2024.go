package gsigeoid2024

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"

	japangeoid "github.com/eukarya-inc/japan-geoid-go"
)

//go:embed JPGEO2024+Hrefconv2024.bin.gz
var combinedRaw []byte

//go:embed Hrefconv2024.bin.gz
var hrefconvRaw []byte

// Load はジオイド2024と基準面補正パラメータを合成したデータを読み込みます。
// 一部の離島において衛星測位によって標高を求める際には、この合成データを使用してください。
func Load() (*japangeoid.MemoryGrid, error) {
	return load(combinedRaw)
}

// LoadHrefconv は基準面補正パラメータのみを読み込みます。
// 東京湾平均海面と離島独自の平均海面の差（基準面補正量）を表したデータです。
func LoadHrefconv() (*japangeoid.MemoryGrid, error) {
	return load(hrefconvRaw)
}

func load(data []byte) (*japangeoid.MemoryGrid, error) {
	b := bytes.NewReader(data)
	g, err := gzip.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	return japangeoid.FromBinary(g)
}
