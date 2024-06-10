package service

import (
	"fmt"
	bigcomplex "go-mandelbrot/pkg/helpers/big-complex"
	colorhelp "go-mandelbrot/pkg/helpers/color"
	"go-mandelbrot/pkg/helpers/math"
	"image"
	"image/color"
	"image/draw"
	"math/big"
	"math/cmplx"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

const (
	itersPerRegion = 20

	defaultRegionsNum = 500
)

type colorRegion struct {
	nextRegion *colorRegion
	startColor color.RGBA
}

// Note: all "state" is client-side, except color palette
// Better to change this
var savedRegions = generateRegions(defaultRegionsNum)

func GenerateMandelbrotImage(pointX, pointY float64, zoom uint64, maxIters uint32, width, height uint32) image.Image {
	pixelItersMap := generateItersMap(pointX, pointY, zoom, maxIters, width, height)

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
	invZoom := float64(1.0) / float64(zoom)
	axisMin *= invZoom
	axisMax *= invZoom

	axisRange := axisMax - axisMin

	transformed := axisMin + (float64(point)/float64(pixelBounds-1))*axisRange

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
		axisY[i] = transformPixelToCartesian(i, height, -2, 2, -pointY, zoom)
	}

	result := make([][]uint32, width)

	var wg sync.WaitGroup

	for x := range result {
		result[x] = make([]uint32, height)
		for y := range result[x] {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result[x][y] = iteratePointRaw(axisX[y], axisY[x], maxIters)
			}()
		}
	}

	wg.Wait()

	return result
}

// TODO: make better and more consistent coloring(probably histogram)
func generateRegions(numOfRegions uint32) map[uint32]*colorRegion {
	regions := make(map[uint32]*colorRegion, numOfRegions)

	var prevRegion *colorRegion
	prevColor := colorhelp.RandomRGBAColor()
	for i := range numOfRegions + 1 {
		region := &colorRegion{
			nextRegion: prevRegion,
			// startColor: colorhelp.RandomRGBAColor(),
		}
		region.startColor = colorhelp.ShuffleColor(colorhelp.InvertColor(prevColor))
		regionIter := uint32((numOfRegions * itersPerRegion) - (itersPerRegion * i))
		regions[regionIter] = region
		prevRegion = region
	}
	return regions
}

func RegenerateRegions() {
	for _, v := range savedRegions {
		v.startColor = colorhelp.RandomRGBAColor()
	}
}

func generateImage(itersMap [][]uint32, width, height, maxIter uint32) image.Image {
	defer timer("generate image")()
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	for x, xRow := range itersMap {
		for y, iter := range xRow {
			if iter == maxIter {
				img.Set(y, x, colorhelp.BlackRGBA)
				continue
			}
			boundedIter := iter - (iter % itersPerRegion)
			region, ok := savedRegions[boundedIter]
			if !ok {
				region = &colorRegion{startColor: colorhelp.BlackRGBA}
			}

			leftColor := region.startColor
			rightColor := colorhelp.BlackRGBA
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

// Not efficient
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

func iteratePointRaw(x0, y0 float64, maxIters uint32) uint32 {
	x2 := float64(0)
	y2 := float64(0)

	x := float64(0)
	y := float64(0)

	i := uint32(0)
	for ; x2+y2 <= 4 && i < maxIters; i++ {
		y = (x+x)*y + y0
		x = x2 - y2 + x0

		x2 = x * x
		y2 = y * y
	}

	return i
}

// UNUSED
func transformPixelToCartesianBig(point, pixelBounds uint32, axisMin, axisMax, offset float64, zoom uint64) *big.Float {
	bigPoint := new(big.Float).SetUint64(uint64(point))
	bigPixelBounds := new(big.Float).SetUint64(uint64(pixelBounds))
	bigAxisMin := new(big.Float).SetFloat64(axisMin)
	bigAxisMax := new(big.Float).SetFloat64(axisMax)
	bigOffset := new(big.Float).SetFloat64(offset)
	bigZoom := new(big.Float).SetUint64(zoom)

	bigOne := new(big.Float).SetFloat64(1.0)
	invZoom := new(big.Float).Quo(bigOne, bigZoom)

	bigAxisMin.Mul(bigAxisMin, invZoom)
	bigAxisMax.Mul(bigAxisMax, invZoom)

	axisRange := new(big.Float).Sub(bigAxisMax, bigAxisMin)

	bigPixelBounds.Sub(bigPixelBounds, bigOne)
	pointRatio := new(big.Float).Quo(bigPoint, bigPixelBounds)
	transformed := new(big.Float).Mul(pointRatio, axisRange)
	transformed.Add(bigAxisMin, transformed)

	transformed.Add(transformed, bigOffset)

	return transformed
}

func transformPixelToCartesianDecimal(point, pixelBounds uint32, axisMin, axisMax, offset float64, zoom uint64) *decimal.Big {
	bigPoint := new(decimal.Big).SetUint64(uint64(point))
	bigPixelBounds := new(decimal.Big).SetUint64(uint64(pixelBounds))
	bigAxisMin := new(decimal.Big).SetFloat64(axisMin)
	bigAxisMax := new(decimal.Big).SetFloat64(axisMax)
	bigOffset := new(decimal.Big).SetFloat64(offset)
	bigZoom := new(decimal.Big).SetUint64(zoom)

	bigOne := new(decimal.Big).SetFloat64(1.0)
	invZoom := new(decimal.Big).Quo(bigOne, bigZoom)

	bigAxisMin.Mul(bigAxisMin, invZoom)
	bigAxisMax.Mul(bigAxisMax, invZoom)

	bigAxisMin.Mul(bigAxisMin, invZoom)
	bigAxisMax.Mul(bigAxisMax, invZoom)

	axisRange := new(decimal.Big).Sub(bigAxisMax, bigAxisMin)

	bigPixelBounds.Sub(bigPixelBounds, bigOne)
	pointRatio := new(decimal.Big).Quo(bigPoint, bigPixelBounds)
	transformed := new(decimal.Big).Mul(pointRatio, axisRange)
	transformed.Add(bigAxisMin, transformed)

	transformed.Add(transformed, bigOffset)

	return transformed
}

// UNUSED
func iteratePointDecimal(x, y *decimal.Big, maxIters uint32) uint32 {
	z := bigcomplex.NewBigComplex(new(decimal.Big).SetFloat64(0), new(decimal.Big).SetFloat64(0))
	c := bigcomplex.NewBigComplex(x, y)
	two := new(decimal.Big).SetFloat64(2.0)

	for n := uint32(0); n < maxIters; n++ {
		z = z.Mul(z).Add(c)

		if z.Abs().Cmp(two) > 0 {
			return n
		}
	}

	return maxIters
}
