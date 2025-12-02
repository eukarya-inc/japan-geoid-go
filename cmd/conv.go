package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	japangeoid "github.com/eukarya-inc/japan-geoid-go"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: conv <file.asc|file.isg>")
		os.Exit(1)
	}

	p := os.Args[1]
	name := path.Base(p)

	f, err := loadFile(p)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// 拡張子に応じてパーサーを選択
	var grid *japangeoid.MemoryGrid
	switch {
	case strings.HasSuffix(strings.ToLower(name), ".isg"):
		grid, err = japangeoid.FromIsg(f)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		name = strings.TrimSuffix(name, ".isg")
	case strings.HasSuffix(strings.ToLower(name), ".asc"):
		grid, err = japangeoid.FromAsc(f)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		name = strings.TrimSuffix(name, ".asc")
	default:
		println("Unknown file extension. Supported: .asc, .isg")
		os.Exit(1)
	}

	outName := name + ".bin.gz"
	f2, err := os.Create(outName)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer f2.Close()

	w := gzip.NewWriter(f2)
	defer w.Close()

	if err := grid.ToBinary(w); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Done: %s (XNum=%d, YNum=%d, XMin=%.2f, YMin=%.2f)\n",
		outName, grid.Info.XNum, grid.Info.YNum, grid.Info.XMin, grid.Info.YMin)
}

func loadFile(url string) (io.ReadCloser, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		return res.Body, nil
	}

	return os.Open(url)
}
