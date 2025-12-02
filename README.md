# japan-geoid-go

日本のジオイド高を計算するための Go ライブラリ / Go library for calculating geoid height in Japan

- **外部依存ゼロ** - 標準ライブラリのみで動作します
- 本ライブラリは国土地理院が提供するものではありません
- 本ライブラリは [ciscorn/japan-geoid](https://github.com/ciscorn/japan-geoid/tree/main) （MIT License）を参考に実装しています

## 対応ジオイドモデル

| パッケージ | 関数 | 内容 | サイズ |
|-----------|------|------|--------|
| `gsigeoid2011` | `Load()` | 日本のジオイド2011（Ver.2.2） | 約300KB |
| `gsigeoid2024` | `Load()` | ジオイド2024＋基準面補正パラメータ（合成版、推奨） | 約400KB |
| `gsigeoid2024` | `LoadHrefconv()` | 基準面補正パラメータのみ | 約15KB |
| `gsigeoid2024geoid` | `Load()` | ジオイド2024のみ | 約3.7MB |

### 日本のジオイド2011（Ver.2.2）

- 出典: https://fgd.gsi.go.jp/download/geoid.php

### ジオイド2024（日本とその周辺）＋基準面補正パラメータ

- 出典: https://service.gsi.go.jp/kiban/app/geoid/

ジオイド・モデル「ジオイド2024日本とその周辺」は、世界測地系（日本測地系2024）における座標値（緯度と経度）で示された任意の位置でのジオイド高を表したデータです。

また、「基準面補正パラメータ」は、東京湾平均海面と離島独自の平均海面の差を基準面補正量と定義し、前述した任意の位置での基準面補正量を表したデータです。

一部の離島において衛星測位によって標高を求める際には、ジオイド高と基準面補正量を使用する必要があります。

ジオイドのみ（`gsigeoid2024geoid`）はサイズが大きいため、別パッケージとして分離しています。必要な場合のみインポートしてください。

※ ジオイドのみのデータが大きい理由：合成版は海域がnodata扱いですが、ジオイドのみのデータは海域にも値が含まれており、圧縮効率が低くなっています。

## 使い方

### gsigeoid2011

## 使い方

```go
package main

import (
	"fmt"
	"github.com/eukarya-inc/japan-geoid-go/gsigeoid2011"
)

func main() {
	g, err := gsigeoid2011.Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Geoid height: %.6f\n", height)
	// Output: Geoid height: 39.473871
}
```

### gsigeoid2024

```go
package main

import (
	"fmt"
	"github.com/eukarya-inc/japan-geoid-go/gsigeoid2024"
)

func main() {
	// ジオイド2024＋基準面補正パラメータ（合成版）
	g, err := gsigeoid2024.Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Geoid height: %.6f\n", height)
	// Output: Geoid height: 39.596702

	// 基準面補正パラメータのみ（離島で使用）
	gHref, _ := gsigeoid2024.LoadHrefconv()
	// 小笠原諸島（父島）付近
	fmt.Printf("Hrefconv (Ogasawara): %.6f\n", gHref.GetHeight(142.19, 27.09))
	// Output: Hrefconv (Ogasawara): 0.642000
}
```

### gsigeoid2024geoid（ジオイドのみ、約3.7MB）

```go
package main

import (
	"fmt"
	"github.com/eukarya-inc/japan-geoid-go/gsigeoid2024geoid"
)

func main() {
	g, err := gsigeoid2024geoid.Load()
	if err != nil {
		panic(err)
	}

	lng, lat := 138.2839817085188, 37.12378643088312
	height := g.GetHeight(lng, lat)
	fmt.Printf("Geoid height: %.6f\n", height)
	// Output: Geoid height: 39.596702
}
```

### ISGパーサー

ISG形式のファイルを直接パースする場合は、`isg` パッケージを使用できます。

```go
package main

import (
	"fmt"
	"os"

	"github.com/eukarya-inc/japan-geoid-go/isg"
)

func main() {
	f, _ := os.Open("geoid.isg")
	defer f.Close()

	grid, err := isg.Parse(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Model: %s\n", grid.Header.ModelName)
	fmt.Printf("Rows: %d, Cols: %d\n", grid.Header.NRows, grid.Header.NCols)
}
```

## 補間アルゴリズム

本ライブラリでは、グリッド点間のジオイド高を計算する際に **双線形補間（Bilinear Interpolation）** を使用しています。

指定された座標を囲む4つのグリッド点の値から、以下の式で補間値を計算します：

```
h = h00 * (1-x) * (1-y) + h01 * x * (1-y) + h10 * (1-x) * y + h11 * x * y
```

ここで：
- `h00`, `h01`, `h10`, `h11`: 4つのグリッド点のジオイド高
- `x`, `y`: グリッド内での相対位置（0.0〜1.0）

グリッド点がデータ範囲外または nodata の場合は `NaN` を返します。

## License

[MIT License](LICENSE)
