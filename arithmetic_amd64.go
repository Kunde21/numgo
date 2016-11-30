//+build !noasm,!appengine

package numgo

import (
//"github.com/Kunde21/numgo"
)

var (
	Sse3Supt, AvxSupt, Avx2Supt, FmaSupt bool
)

func init() {
	initasm()
}

func initasm()

func AddC(c nDimElement, d []nDimElement)

func SubtrC(c nDimElement, d []nDimElement)

func MultC(c nDimElement, d []nDimElement)

func DivC(c nDimElement, d []nDimElement)

func Add(a, b []nDimElement)

func Vadd(a, b []nDimElement)

func Hadd(st uint64, a []nDimElement)

func Subtr(a, b []nDimElement)

func Mult(a, b []nDimElement)

func Div(a, b []nDimElement)

func Fma12(a nDimElement, x, b []nDimElement)

func Fma21(a nDimElement, x, b []nDimElement)
