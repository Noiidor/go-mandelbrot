package service

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"sync"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func GenerateMandelbrotImage(pointX, pointY float64, zoom uint64, maxIters uint32, width, height uint32) image.Image {
	pixelItersMap := generateItersMap2(pointX, pointY, zoom, maxIters, width, height)
	img := generateImage(pixelItersMap, width, height)

	return img
}

func transformPixelToCartesian(point, pixelBounds uint32, axisMin, axisMax, offset float64, zoom uint64) float64 {
	// Scaling
	axisMin /= float64(zoom)
	axisMax /= float64(zoom)

	transformed := axisMin + ((float64(point) / (float64(pixelBounds - 1))) * (axisMax - axisMin))

	return transformed + offset
}

func generateItersMap2(pointX, pointY float64, zoom uint64, maxIters uint32, width, height uint32) [][]uint32 {
	defer timer("generate iters map")()

	axisX := make([]float64, width)
	axisY := make([]float64, height)

	for i := range width {
		axisX[i] = transformPixelToCartesian(i, width, -2, 2, pointX, zoom)
	}

	for i := range height {
		axisY[i] = transformPixelToCartesian(i, height, -2, 2, pointY, zoom)
	}

	result := make([][]uint32, width)

	var wg sync.WaitGroup

	for x := range result {
		result[x] = make([]uint32, height)
		for y := range result[x] {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result[x][y] = iteratePoint(axisX[y], axisY[x], maxIters)
			}()
		}
	}

	wg.Wait()

	return result
}

func generateImage(itersMap [][]uint32, width, height uint32) image.Image {
	defer timer("generate image")()
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	for x, xRow := range itersMap {
		for y := range xRow {
			if itersMap[x][y] > 100 {
				img.Set(y, x, color.Black)
			}
		}
	}

	return img
}

func iteratePoint(x, y float64, maxIters uint32) uint32 {
	var z complex128
	c := complex(x, y)

	for n := range maxIters {
		z = z*z + c

		if cmplx.Abs(z) > 2 {
			return n
		}
	}
	return maxIters
}
