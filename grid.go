package japangeoid

import "math"

// Grid はジオイドデータを表すインターフェースです。
type Grid interface {
	// GetHeight は指定した経度(lng)・緯度(lat)におけるジオイドの高さを返します。
	GetHeight(lng, lat float64) float64
}

// MemoryGrid はジオイドデータを保持する構造体です。
type MemoryGrid struct {
	Info   GridInfo
	Points []int32 // グリッド点の高さ: [y_num * x_num]
}

// GetHeight は指定した経度(lng)・緯度(lat)におけるジオイドの高さを返します。無効な座標の場合は NaN を返します。
func (g *MemoryGrid) GetHeight(lng, lat float64) float64 {
	x, y, info := lng, lat, g.Info

	// 1. convert lng, lat to grid index
	gridX := (x - float64(info.XMin)) * float64(info.XDenom)
	gridY := (y - float64(info.YMin)) * float64(info.YDenom)
	if gridX < 0.0 || gridY < 0.0 {
		return math.NaN()
	}

	ix := int(math.Floor(gridX))
	iy := int(math.Floor(gridY))
	if ix < 0 || iy < 0 || ix >= int(info.XNum) || iy >= int(info.YNum) {
		return math.NaN()
	}

	// 2. get values of 4 points around
	v00 := g.lookupGridPoints(ix, iy)

	var v01 float64
	if ix < int(info.XNum)-1 {
		v01 = g.lookupGridPoints(ix+1, iy)
	} else {
		v01 = math.NaN()
	}

	var v10 float64
	if iy < int(info.YNum)-1 {
		v10 = g.lookupGridPoints(ix, iy+1)
	} else {
		v10 = math.NaN()
	}

	var v11 float64
	if ix < int(info.XNum)-1 && iy < int(info.YNum)-1 {
		v11 = g.lookupGridPoints(ix+1, iy+1)
	} else {
		v11 = math.NaN()
	}

	// 3. apply bilinear interpolation
	xFrac := gridX - float64(ix)
	yFrac := gridY - float64(iy)
	return bilinear(xFrac, yFrac, v00, v01, v10, v11)
}

func (g *MemoryGrid) lookupGridPoints(ix, iy int) float64 {
	pos := iy*int(g.Info.XNum) + ix
	v := g.Points[pos]

	// 9990000 is NaN
	if v == 9990000 {
		return math.NaN()
	}

	return float64(v) / 10000.0
}
