package math

import "math/rand/v2"

func Lerp(a, b, t float64) float64 {
	return (a * (1.0 - t)) + (b * t)
}

func RatioBetweenNums(a, b, x int) float64 {
	return 1.0 - (float64(b)-float64(x))/(float64(b)-float64(a))
}

func ClampInt(num, min, max int) int {
	if num > max {
		return max
	}
	if num < min {
		return num
	}
	return num
}

func RandomRange(min, max int) int {
	return rand.IntN(max-min) + min
}
