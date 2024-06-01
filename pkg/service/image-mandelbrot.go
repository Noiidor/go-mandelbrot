package service

import (
	"fmt"
	colorhelp "go-mandelbrot/pkg/helpers/color"
	"go-mandelbrot/pkg/helpers/math"
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

const (
	itersPerRegion = 10

	defaultRegionsNum = 50
)

type colorRegion struct {
	// endIter  uint32
	nextRegion *colorRegion
	startColor color.RGBA
}

var savedRegions = generateRegions(defaultRegionsNum)

func GenerateMandelbrotImage(pointX, pointY float64, zoom uint64, maxIters uint32, width, height uint32) image.Image {
	pixelItersMap := generateItersMap(pointX, pointY, zoom, maxIters, width, height)

	checkRegionsDeficiency(maxIters)

	img := generateImage(pixelItersMap, width, height)

	return img
}

func checkRegionsDeficiency(maxIters uint32) {
	overflow := (int(maxIters) - len(savedRegions)*itersPerRegion)
	if overflow > 0 {
		// savedRegions = append(savedRegions, generateRegions(uint32(overflow/itersPerRegion)+1)...)
	}
}

func transformPixelToCartesian(point, pixelBounds uint32, axisMin, axisMax, offset float64, zoom uint64) float64 {
	// Scaling
	axisMin /= float64(zoom)
	axisMax /= float64(zoom)

	transformed := axisMin + ((float64(point) / (float64(pixelBounds - 1))) * (axisMax - axisMin))

	return transformed + offset
}

func generateItersMap(pointX, pointY float64, zoom uint64, maxIters uint32, width, height uint32) [][]uint32 {
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

func generateRegions(numOfRegions uint32) map[uint32]*colorRegion {
	regions := make(map[uint32]*colorRegion, numOfRegions)

	var prevRegion *colorRegion
	for i := range numOfRegions + 1 {
		region := &colorRegion{
			nextRegion: prevRegion,
			startColor: colorhelp.RandomRGBAColor(),
		}
		regionIter := uint32((numOfRegions * itersPerRegion) - (itersPerRegion * i))
		regions[regionIter] = region
		prevRegion = region
	}
	return regions
}

func generateImage(itersMap [][]uint32, width, height uint32) image.Image {
	defer timer("generate image")()
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	for x, xRow := range itersMap {
		for y, iter := range xRow {
			boundedIter := iter - (iter % itersPerRegion)
			region := savedRegions[boundedIter]

			leftColor := region.startColor
			rightColor := color.RGBA{}
			if region.nextRegion != nil {
				rightColor = region.nextRegion.startColor
			}

			ratio := math.RatioBetweenNums(int(boundedIter), int(boundedIter+itersPerRegion), int(iter))
			color := colorhelp.LerpColor(leftColor, rightColor, ratio)

			img.Set(y, x, color)
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
