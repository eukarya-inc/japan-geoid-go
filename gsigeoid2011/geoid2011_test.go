package gsigeoid2011

import (
	"fmt"
	// "github.com/eukarya-inc/japan-geoid-go/gsigeoid2011"
)

func Example() {
	g, err := Load()
	// g, err := gsigeoid2011.Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Input: (lng: %.6f, lat: %.6f) -> Geoid height: %.6f\n", lng, lat, height)

	lng, lat = 10, 10
	height = g.GetHeight(lng, lat)
	fmt.Printf("Input: (lng: %.6f, lat: %.6f) -> Geoid height: %.6f\n", lng, lat, height)

	// Output:
	// Input: (lng: 138.283982, lat: 37.123786) -> Geoid height: 39.473871
	// Input: (lng: 10.000000, lat: 10.000000) -> Geoid height: NaN
}
