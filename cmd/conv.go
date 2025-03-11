package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	japangeoid "github.com/eukarya-inc/japan-geoid-go"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: conv <file>")
		os.Exit(1)
	}

	p := os.Args[1]
	name := path.Base(p)

	f, err := loadAsc(p)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	defer f.Close()

	grid, err := japangeoid.FromAsc(f)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	f2, err := os.Create(strings.TrimSuffix(name, ".asc") + ".bin.gz")
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

	println("Done")
}

func loadAsc(url string) (io.ReadCloser, error) {
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
