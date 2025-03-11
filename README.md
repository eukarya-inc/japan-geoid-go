# japangeoid-go

日本のジオイド高を計算するための Go ライブラリ / Go library for calculating geoid height in Japan

- 本ライブラリは国土地理院が提供するものではありません。
- 本ライブラリは [ciscorn/japan-geoid](https://github.com/ciscorn/japan-geoid/tree/main) （MIT License）を参考に実装しています。

対応ジオイドモデル：

- 日本のジオイド 2011（Ver.2.2）（[出典](https://fgd.gsi.go.jp/download/geoid.php)）: `gsigeoid2011`
- 日本のジオイド 2025（[出典](https://www.gsi.go.jp/buturisokuchi/grageo_reference.html)）: 正式公開後対応予定

```go
package main

import (
	"fmt"
	"github.com/eukarya-inc/japan-geoid-go/gsigeoid2011"
)

func Example() {
	g, err := gsigeoid2011.Load()
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
```

## Lisence

[MIT Lisence](LICENSE)
