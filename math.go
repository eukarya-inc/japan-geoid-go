package japangeoid

func bilinear(xFrac, yFrac, v00, v01, v10, v11 float64) float64 {
	if xFrac == 0.0 && yFrac == 0.0 {
		return v00
	} else if xFrac == 0.0 {
		return v00*(1.0-yFrac) + v10*yFrac
	} else if yFrac == 0.0 {
		return v00*(1.0-xFrac) + v01*xFrac
	}

	return v00*(1.0-xFrac)*(1.0-yFrac) +
		v01*xFrac*(1.0-yFrac) +
		v10*(1.0-xFrac)*yFrac +
		v11*xFrac*yFrac
}
