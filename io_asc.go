package japangeoid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// FromAsc はASCII形式(asc)のジオイドモデルを読み込みます。
func FromAsc(r io.Reader) (*MemoryGrid, error) {
	scanner := bufio.NewScanner(r)

	info, err := gridInfoFromAsc(scanner)
	if err != nil {
		return nil, err
	}

	// read points
	var points []int32
	points = make([]int32, 0, info.XNum*info.YNum)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, f := range fields {
			// delete dot and convert to int
			s := strings.ReplaceAll(f, ".", "")
			val, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid data %q: %w", s, err)
			}

			points = append(points, int32(val))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// validdate points count
	expected := info.XNum * info.YNum
	if uint32(len(points)) != expected {
		return nil, fmt.Errorf("mismatch points count: got %d, want %d",
			len(points), expected)
	}

	return &MemoryGrid{
		Info:   info,
		Points: points,
	}, nil
}

func gridInfoFromAsc(s *bufio.Scanner) (GridInfo, error) {
	if !s.Scan() {
		return GridInfo{}, errors.New("failed to read the header line")
	}

	line := s.Text()
	parts := strings.Fields(line)
	if len(parts) != 8 {
		return GridInfo{}, fmt.Errorf("header line must have 8 values, got %d", len(parts))
	}

	//  c[0] -> y_min
	//  c[1] -> x_min
	//  c[2] -> x_denom
	//  c[3] -> y_denom
	//  c[4] -> y_num
	//  c[5] -> x_num
	//  c[6] -> ikind
	//  c[7] -> version

	if parts[2] != "0.016667" { // 1/60
		return GridInfo{}, fmt.Errorf("latitude interval must be 0.016667, got %s", parts[2])
	}

	if parts[3] != "0.025000" { // 1/40
		return GridInfo{}, fmt.Errorf("longitude interval must be 0.025000, got %s", parts[3])
	}

	yMin, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return GridInfo{}, fmt.Errorf("invalid y_min: %w", err)
	}

	xMin, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return GridInfo{}, fmt.Errorf("invalid x_min: %w", err)
	}

	yNum64, err := strconv.ParseUint(parts[4], 10, 32)
	if err != nil {
		return GridInfo{}, fmt.Errorf("invalid y_num: %w", err)
	}

	xNum64, err := strconv.ParseUint(parts[5], 10, 32)
	if err != nil {
		return GridInfo{}, fmt.Errorf("invalid x_num: %w", err)
	}

	iKind64, err := strconv.ParseUint(parts[6], 10, 16)
	if err != nil {
		return GridInfo{}, fmt.Errorf("invalid ikind: %w", err)
	}

	version := parts[7]

	return GridInfo{
		XNum:    uint32(xNum64),
		YNum:    uint32(yNum64),
		XDenom:  40,
		YDenom:  60,
		XMin:    float32(xMin),
		YMin:    float32(yMin),
		IKind:   uint16(iKind64),
		Version: version,
	}, nil
}
