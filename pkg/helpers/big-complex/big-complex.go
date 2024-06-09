package bigcomplex

import "math/big"

// Tried to use arbitrary precision floats for increased zoom
// Failed miserably(5m for single image at zoom 1.4*10^14)

type BigComplex struct {
	real, imag *big.Float
}

func NewBigComplex(real, imag *big.Float) *BigComplex {
	return &BigComplex{
		real: real,
		imag: imag,
	}
}

func (z *BigComplex) Add(w *BigComplex) *BigComplex {
	z.real.Add(z.real, w.real)
	z.imag.Add(z.imag, w.imag)
	return z
}

func (z *BigComplex) Mul(w *BigComplex) *BigComplex {
	real := new(big.Float).Sub(
		new(big.Float).Mul(z.real, w.real),
		new(big.Float).Mul(z.imag, w.imag),
	)
	imag := new(big.Float).Add(
		new(big.Float).Mul(z.real, w.imag),
		new(big.Float).Mul(z.imag, w.real),
	)
	z.real, z.imag = real, imag
	return z
}

func (z *BigComplex) Abs() *big.Float {
	realSq := new(big.Float).Mul(z.real, z.real)
	imagSq := new(big.Float).Mul(z.imag, z.imag)
	sum := new(big.Float).Add(realSq, imagSq)
	return new(big.Float).Sqrt(sum)
}
