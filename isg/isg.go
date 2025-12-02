// Package isg provides a parser for ISG (International Service for the Geoid) format files.
package isg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Header holds the header information of an ISG file.
type Header struct {
	ModelName    string
	ModelYear    string
	ModelType    string
	DataType     string
	DataUnits    string
	DataFormat   string
	DataOrdering string
	RefEllipsoid string
	RefFrame     string
	HeightDatum  string
	TideSystem   string
	CoordType    string
	CoordUnits   string
	LatMin       float64
	LatMax       float64
	LonMin       float64
	LonMax       float64
	DeltaLat     float64
	DeltaLon     float64
	NRows        int
	NCols        int
	NoData       float64
	CreationDate string
	ISGFormat    string
}

// Grid holds the grid data read from an ISG file.
type Grid struct {
	Header Header
	// Points holds the grid point values.
	// The data order follows the ISG file's data ordering (typically N-to-S, W-to-E).
	Points []float64
}

// Parse reads ISG format data and returns a Grid.
func Parse(r io.Reader) (*Grid, error) {
	scanner := bufio.NewScanner(r)
	// Extend buffer since ISG file lines can be very long
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	header, err := ParseHeader(scanner)
	if err != nil {
		return nil, err
	}

	// read grid data
	points := make([]float64, 0, header.NRows*header.NCols)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, f := range fields {
			val, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid data %q: %w", f, err)
			}
			points = append(points, val)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// validate points count
	expected := header.NRows * header.NCols
	if len(points) != expected {
		return nil, fmt.Errorf("mismatch points count: got %d, want %d", len(points), expected)
	}

	return &Grid{
		Header: *header,
		Points: points,
	}, nil
}

func ParseHeader(s *bufio.Scanner) (*Header, error) {
	header := &Header{}

	// Find begin_of_head
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "begin_of_head") {
			break
		}
	}

	// Parse header
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "end_of_head") {
			break
		}

		// Parse key = value or key : value
		var key, value string
		if idx := strings.Index(line, "="); idx != -1 {
			key = strings.TrimSpace(line[:idx])
			value = strings.TrimSpace(line[idx+1:])
		} else if idx := strings.Index(line, ":"); idx != -1 {
			key = strings.TrimSpace(line[:idx])
			value = strings.TrimSpace(line[idx+1:])
		} else {
			continue
		}

		switch key {
		case "model name":
			header.ModelName = value
		case "model year":
			header.ModelYear = value
		case "model type":
			header.ModelType = value
		case "data type":
			header.DataType = value
		case "data units":
			header.DataUnits = value
		case "data format":
			header.DataFormat = value
		case "data ordering":
			header.DataOrdering = value
		case "ref ellipsoid":
			header.RefEllipsoid = value
		case "ref frame":
			header.RefFrame = value
		case "height datum":
			header.HeightDatum = value
		case "tide system":
			header.TideSystem = value
		case "coord type":
			header.CoordType = value
		case "coord units":
			header.CoordUnits = value
		case "lat min":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid lat min %q: %w", value, err)
			}
			header.LatMin = v
		case "lat max":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid lat max %q: %w", value, err)
			}
			header.LatMax = v
		case "lon min":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid lon min %q: %w", value, err)
			}
			header.LonMin = v
		case "lon max":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid lon max %q: %w", value, err)
			}
			header.LonMax = v
		case "delta lat":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid delta lat %q: %w", value, err)
			}
			header.DeltaLat = v
		case "delta lon":
			v, err := ParseDMS(value)
			if err != nil {
				return nil, fmt.Errorf("invalid delta lon %q: %w", value, err)
			}
			header.DeltaLon = v
		case "nrows":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid nrows %q: %w", value, err)
			}
			header.NRows = v
		case "ncols":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid ncols %q: %w", value, err)
			}
			header.NCols = v
		case "nodata":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid nodata %q: %w", value, err)
			}
			header.NoData = v
		case "creation date":
			header.CreationDate = value
		case "ISG format":
			header.ISGFormat = value
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	// validate
	if header.NRows == 0 || header.NCols == 0 {
		return nil, errors.New("nrows or ncols is zero")
	}
	if header.DeltaLat == 0 || header.DeltaLon == 0 {
		return nil, errors.New("delta lat or delta lon is zero")
	}

	return header, nil
}

// dmsRegex is a regular expression for parsing DMS (degrees, minutes, seconds) format (e.g., 15°00'00")
var dmsRegex = regexp.MustCompile(`^\s*(-?\d+)°(\d+)'(\d+(?:\.\d+)?)"?\s*$`)

// ParseDMS converts a DMS (degrees, minutes, seconds) format string to decimal degrees.
// Examples:
//   - "15°00'00\"" → 15.0
//   - "0°01'30\"" → 0.025 (= 1.5/60 = 1/40)
func ParseDMS(s string) (float64, error) {
	matches := dmsRegex.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid DMS format: %q", s)
	}

	deg, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	min, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, err
	}

	sec, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return 0, err
	}

	// Handle negative degrees (minutes and seconds follow the sign)
	sign := 1.0
	if deg < 0 {
		sign = -1.0
		deg = -deg
	}

	return sign * (deg + min/60.0 + sec/3600.0), nil
}
