package japangeoid

import (
	"io"
	"strings"

	"github.com/eukarya-inc/japan-geoid-go/isg"
)

// FromIsg はISG形式のジオイドモデルを読み込みます。
// ISG (International Service for the Geoid) 形式は、ヘッダーとグリッドデータで構成されます。
func FromIsg(r io.Reader) (*MemoryGrid, error) {
	grid, err := isg.Parse(r)
	if err != nil {
		return nil, err
	}

	header := grid.Header

	// XDenom, YDenom を計算
	// delta lon = 1/40度 → XDenom = 40
	// delta lat = 1/60度 → YDenom = 60
	xDenom := uint32(1.0 / header.DeltaLon)
	yDenom := uint32(1.0 / header.DeltaLat)

	// Version文字列は10文字以下に制限（バイナリ形式の制約）
	version := header.ModelName
	if len(version) > 10 {
		version = version[:10]
	}

	info := GridInfo{
		XNum:    uint32(header.NCols),
		YNum:    uint32(header.NRows),
		XDenom:  xDenom,
		YDenom:  yDenom,
		XMin:    float32(header.LonMin),
		YMin:    float32(header.LatMin),
		IKind:   0,
		Version: version,
	}

	// ポイントデータを内部形式に変換
	points := make([]int32, len(grid.Points))
	for i, val := range grid.Points {
		// nodata値の変換: -9999.0000 → 9990000 (内部nodata値)
		if val <= header.NoData {
			points[i] = 9990000
		} else {
			// 小数点4桁を整数に変換 (例: 12.3456 → 123456)
			points[i] = int32(val * 10000)
		}
	}

	// ISGはN→S（北から南）の順序だが、内部形式はS→N（南から北）なので行を反転する
	// ただし、data ordering を確認して判断
	if strings.Contains(header.DataOrdering, "N-to-S") {
		reversedPoints := make([]int32, len(points))
		xNum := int(info.XNum)
		yNum := int(info.YNum)
		for y := 0; y < yNum; y++ {
			srcRowStart := y * xNum
			dstRowStart := (yNum - 1 - y) * xNum
			copy(reversedPoints[dstRowStart:dstRowStart+xNum], points[srcRowStart:srcRowStart+xNum])
		}
		points = reversedPoints
	}

	return &MemoryGrid{
		Info:   info,
		Points: points,
	}, nil
}
