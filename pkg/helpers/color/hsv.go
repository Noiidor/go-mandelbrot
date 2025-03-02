package color

import (
	"fmt"
	"image/color"
	"math"
)

var _ color.Color = (*HSV)(nil)

type HSV struct {
	H, S, V float64
}

func (c HSV) RGBA() (r, g, b, a uint32) {
	a = 255
	if c.H < 0 || c.H >= 360 ||
		c.S < 0 || c.S > 1 ||
		c.V < 0 || c.V > 1 {
		panic(fmt.Errorf("HSV out of bounds: %v, %v, %v", c.H, c.S, c.V))
	}
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ v ≤ 1:
	C := c.V * c.S
	X := C * (1 - math.Abs(math.Mod(c.H/60, 2)-1))
	m := c.V - C
	var Rnot, Gnot, Bnot float64
	switch {
	case 0 <= c.H && c.H < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= c.H && c.H < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= c.H && c.H < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= c.H && c.H < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= c.H && c.H < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= c.H && c.H < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = uint32(math.Round((Rnot + m) * 255))
	g = uint32(math.Round((Gnot + m) * 255))
	b = uint32(math.Round((Bnot + m) * 255))
	return r, g, b, a
}
