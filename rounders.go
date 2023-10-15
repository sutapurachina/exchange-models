package exchange_models

import "math"

func Round(number float64, precision int) float64 {
	return math.Round(number*math.Pow10(precision)) / math.Pow10(precision)
}

func Floor(number float64, precision int) float64 {
	return math.Floor(number*math.Pow10(precision)) / math.Pow10(precision)
}

func Ceil(number float64, precision int) float64 {
	return math.Ceil(number*math.Pow10(precision)) / math.Pow10(precision)
}

func Equals(n1, n2, eps float64) bool {
	return math.Abs(n1-n2) < eps
}
