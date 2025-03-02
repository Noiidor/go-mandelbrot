package bigcomplex

import (
	"github.com/ericlagergren/decimal"
)

// BigComplex represents a complex number with decimal.Big components
type BigComplex struct {
	real, imag *decimal.Big
}

// NewBigComplex creates a new BigComplex
func NewBigComplex(real, imag *decimal.Big) *BigComplex {
	return &BigComplex{
		real: real,
		imag: imag,
	}
}

// Add adds two BigComplex numbers and stores the result in the receiver
func (z *BigComplex) Add(w *BigComplex) *BigComplex {
	z.real.Add(z.real, w.real)
	z.imag.Add(z.imag, w.imag)
	return z
}

// Mul multiplies two BigComplex numbers and stores the result in the receiver
func (z *BigComplex) Mul(w *BigComplex) *BigComplex {
	real := new(decimal.Big).Sub(
		new(decimal.Big).Mul(z.real, w.real),
		new(decimal.Big).Mul(z.imag, w.imag),
	)
	imag := new(decimal.Big).Add(
		new(decimal.Big).Mul(z.real, w.imag),
		new(decimal.Big).Mul(z.imag, w.real),
	)
	z.real, z.imag = real, imag
	return z
}

// Abs returns the absolute value (magnitude) of the BigComplex number
func (z *BigComplex) Abs() *decimal.Big {
	realSq := new(decimal.Big).Mul(z.real, z.real)
	imagSq := new(decimal.Big).Mul(z.imag, z.imag)
	sum := new(decimal.Big).Add(realSq, imagSq)
	return sum.Context.Sqrt(sum, sum)
}
