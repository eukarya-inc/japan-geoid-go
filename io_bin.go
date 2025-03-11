package japangeoid

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

// ToBinary は Grid をバイナリ形式で書き出します。
func (g *MemoryGrid) ToBinary(w io.Writer) error {
	if err := g.Info.toBinary(w); err != nil {
		return err
	}

	// write points with diff encoding to compress data
	var prevX1 int32 = 9990000
	var prevX1Y1 int32 = 9990000
	intBuf := make([]byte, 4)

	xNum := int(g.Info.XNum)
	yNum := int(g.Info.YNum)

	for pos := 0; pos < xNum*yNum; pos++ {
		curr := g.Points[pos]

		// use 9990000 if there is no element corresponding to the upper row (pos - xNum)
		var prevY1 int32 = 9990000
		if pos >= xNum {
			prevY1 = g.Points[pos-xNum]
		}

		predicted := prevX1 + prevY1 - prevX1Y1
		diff := curr - predicted

		// write diff as little-endian
		binary.LittleEndian.PutUint32(intBuf, uint32(int32(diff)))
		if _, err := w.Write(intBuf); err != nil {
			return err
		}

		// update
		prevX1 = curr
		prevX1Y1 = prevY1
	}

	return nil
}

// FromBinary はバイナリ形式から Grid を復元します。
func FromBinary(r io.Reader) (*MemoryGrid, error) {
	info, err := gridInfoFromBinary(r)
	if err != nil {
		return nil, err
	}

	// read points encoded with diff
	count := int(info.XNum * info.YNum)
	points := make([]int32, count)

	var prevX1 int32 = 9990000
	var prevX1Y1 int32 = 9990000

	buf := make([]byte, 4)
	for pos := 0; pos < count; pos++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		diff := int32(binary.LittleEndian.Uint32(buf))

		var prevY1 int32 = 9990000
		if pos >= int(info.XNum) {
			prevY1 = points[pos-int(info.XNum)]
		}

		predicted := prevX1 + prevY1 - prevX1Y1
		curr := predicted + diff
		points[pos] = curr

		// update
		prevX1 = curr
		prevX1Y1 = prevY1
	}

	mg := &MemoryGrid{
		Info:   info,
		Points: points,
	}
	return mg, nil
}

func (g *GridInfo) toBinary(w io.Writer) error {
	// x_num, y_num, x_denom, y_denom (uint16 x 4)
	// x_min, y_min (float32 x 2)
	// ikind (uint16)
	// version (10 byte ASCII, 0 padded)

	// x_num, y_num, x_denom, y_denom (uint16 x 4)
	var tmp16 [2]byte
	binary.LittleEndian.PutUint16(tmp16[:], uint16(g.XNum))
	if _, err := w.Write(tmp16[:]); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(tmp16[:], uint16(g.YNum))
	if _, err := w.Write(tmp16[:]); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(tmp16[:], uint16(g.XDenom))
	if _, err := w.Write(tmp16[:]); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(tmp16[:], uint16(g.YDenom))
	if _, err := w.Write(tmp16[:]); err != nil {
		return err
	}

	// x_min, y_min (float32 x 2)
	var tmp32 [4]byte
	// x_min
	math32bits := math.Float32bits(g.XMin) // float32 -> uint32
	binary.LittleEndian.PutUint32(tmp32[:], math32bits)
	if _, err := w.Write(tmp32[:]); err != nil {
		return err
	}
	// y_min
	math32bits = math.Float32bits(g.YMin)
	binary.LittleEndian.PutUint32(tmp32[:], math32bits)
	if _, err := w.Write(tmp32[:]); err != nil {
		return err
	}

	// ikind (uint16)
	binary.LittleEndian.PutUint16(tmp16[:], g.IKind)
	if _, err := w.Write(tmp16[:]); err != nil {
		return err
	}

	// version (10 byte ASCII, 0 padded)
	if len(g.Version) > 10 {
		return fmt.Errorf("version string must be <= 10 chars, got %d", len(g.Version))
	}

	verBuf := make([]byte, 10)
	copy(verBuf, []byte(g.Version))

	if _, err := w.Write(verBuf); err != nil {
		return err
	}

	return nil
}

func gridInfoFromBinary(r io.Reader) (GridInfo, error) {
	// x_num, y_num, x_denom, y_denom (uint16 x 4)
	// x_min, y_min (float32 x 2)
	// ikind (uint16)
	// version (10 byte ASCII, 0 padded)

	// read 28 bytes
	var head [28]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return GridInfo{}, err
	}

	// x_num, y_num, x_denom, y_denom
	xNum := binary.LittleEndian.Uint16(head[0:2])
	yNum := binary.LittleEndian.Uint16(head[2:4])
	xDenom := binary.LittleEndian.Uint16(head[4:6])
	yDenom := binary.LittleEndian.Uint16(head[6:8])

	// x_min, y_min
	xMinBits := binary.LittleEndian.Uint32(head[8:12])
	yMinBits := binary.LittleEndian.Uint32(head[12:16])
	xMin := math.Float32frombits(xMinBits)
	yMin := math.Float32frombits(yMinBits)

	// ikind
	iKind := binary.LittleEndian.Uint16(head[16:18])

	// version (10 bytes)
	verRaw := head[18:28]
	// treat as string, but ignore trailing 0s
	version := string(verRaw)
	version = strings.TrimRight(version, "\x00")

	return GridInfo{
		XNum:    uint32(xNum),
		YNum:    uint32(yNum),
		XDenom:  uint32(xDenom),
		YDenom:  uint32(yDenom),
		XMin:    xMin,
		YMin:    yMin,
		IKind:   iKind,
		Version: version,
	}, nil
}
