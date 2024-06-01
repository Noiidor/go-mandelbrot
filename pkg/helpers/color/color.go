package color

import (
	"go-mandelbrot/pkg/helpers/math"
	"image/color"
	"math/rand/v2"
)

func RandomRGBAColor() color.RGBA {
	min := 50
	max := 255
	return color.RGBA{
		uint8(rand.IntN(max-min) + min),
		uint8(rand.IntN(max-min) + min),
		uint8(rand.IntN(max-min) + min), 255}
}

func LerpColor(a, b color.RGBA, t float64) color.Color {
	if t == 0 {
		return a
	}
	if t == 1.0 {
		return b
	}

	// Damn
	resultColor := color.RGBA{
		uint8(math.Lerp(float64(a.R), float64(b.R), t)),
		uint8(math.Lerp(float64(a.G), float64(b.G), t)),
		uint8(math.Lerp(float64(a.B), float64(b.B), t)),
		uint8(math.Lerp(float64(a.A), float64(b.A), t))}

	return resultColor
}
