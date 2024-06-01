package math

func Lerp(a, b, t float64) float64 {
	return (a * (1.0 - t)) + (b * t)
}

func RatioBetweenNums(a, b, x int) float64 {
	return 1.0 - (float64(b)-float64(x))/(float64(b)-float64(a))
}
