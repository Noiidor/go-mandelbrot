package service

import (
	"image"
	"image/color"
	"math/rand"
)

func GenerateMandelbrotImage(pointX, pointY float64, zoom uint64, resolutionWidth, resolutionHeight uint32) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(resolutionWidth), int(resolutionHeight)))

	for i := range resolutionWidth {
		img.Set(int(i), rand.Intn(int(resolutionHeight)), color.Black)
	}

	return img
}
