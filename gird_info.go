package japangeoid

// GridInfo はジオイドモデルのメタデータを表します。
type GridInfo struct {
	XNum    uint32  // x方向グリッド数
	YNum    uint32  // y方向グリッド数
	XDenom  uint32  // x方向の分母 (0.025→1/0.025=40)
	YDenom  uint32  // y方向の分母 (0.016667→1/0.016667≈60)
	XMin    float32 // 経度の最小値
	YMin    float32 // 緯度の最小値
	IKind   uint16  // 未使用(ヘッダのikind)
	Version string  // バージョン文字列
}
