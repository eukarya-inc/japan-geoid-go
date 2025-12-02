package gsigeoid2024geoid

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"

	japangeoid "github.com/eukarya-inc/japan-geoid-go"
)

//go:embed JPGEO2024.bin.gz
var geoidRaw []byte

// Load はジオイド2024（日本とその周辺）のみを読み込みます。
// このパッケージをインポートすると約3.7MBのデータがバイナリに含まれます。
// 基準面補正パラメータとの合成版が必要な場合は gsigeoid2024.Load() を使用してください。
func Load() (*japangeoid.MemoryGrid, error) {
	b := bytes.NewReader(geoidRaw)
	g, err := gzip.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	return japangeoid.FromBinary(g)
}
