package gsigeoid2011

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"

	japangeoid "github.com/eukarya-inc/japan-geoid-go"
)

//go:embed gsigeo2011_ver2_2.bin.gz
var geoid2011raw []byte

func Load() (*japangeoid.MemoryGrid, error) {
	b := bytes.NewReader(geoid2011raw)
	g, err := gzip.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	return japangeoid.FromBinary(g)
}
