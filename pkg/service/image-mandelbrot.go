package service

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"math/rand/v2"
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
	pixelItersMap := generateItersMap2(pointX, pointY, zoom, maxIters, width, height)

	checkRegionsDeficiency(maxIters)

	img := generateImage(pixelItersMap, width, height, maxIters)

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

func generateRegions(numOfRegions uint32) map[uint32]*colorRegion {
	regions := make(map[uint32]*colorRegion, numOfRegions)
	var prevRegion *colorRegion
	for i := range numOfRegions + 1 {
		regionIter := uint32((numOfRegions * itersPerRegion) - (itersPerRegion * i))
		region := &colorRegion{
			nextRegion: prevRegion,
			startColor: randomRGBAColor(),
		}
		regions[regionIter] = region
		prevRegion = region
	}
	return regions
}

func randomRGBAColor() color.RGBA {
	min := 50
	max := 255
	return color.RGBA{
		uint8(rand.IntN(max-min) + min),
		uint8(rand.IntN(max-min) + min),
		uint8(rand.IntN(max-min) + min), 255}
}

func generateImage(itersMap [][]uint32, width, height uint32) image.Image {
	defer timer("generate image")()
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	for x, xRow := range itersMap {
		for y := range xRow {
			iter := itersMap[x][y]
			boundedIter := iter
			iMod := iter % itersPerRegion
			if iMod > 0 {
				boundedIter = iter - iMod
			}
			region := savedRegions[boundedIter]
			leftColor := region.startColor
			rightColor := color.RGBA{}
			if region.nextRegion != nil {
				rightColor = region.nextRegion.startColor
			}

			ratio := ratioBetweenNums(int(boundedIter), int(boundedIter+itersPerRegion), int(iter))

			color := lerpColor(leftColor, rightColor, ratio)

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

func lerp(a, b, t float64) float64 {
	return (a * (1.0 - t)) + (b * t)
}

func lerpColor(a, b color.RGBA, t float64) color.Color {
	if t == 0 {
		return a
	}
	if t == 1.0 {
		return b
	}

	// Damn
	resultColor := color.RGBA{
		uint8(lerp(float64(a.R), float64(b.R), t)),
		uint8(lerp(float64(a.G), float64(b.G), t)),
		uint8(lerp(float64(a.B), float64(b.B), t)),
		uint8(lerp(float64(a.A), float64(b.A), t))}

	return resultColor
}

func ratioBetweenNums(a, b, x int) float64 {
	return 1.0 - (float64(b)-float64(x))/(float64(b)-float64(a))
}
