//+build !noasm,!appengine

package asm

var (
	sse3Supt, avxSupt, avx2Supt, fmaSupt bool
)

func init() {
	initasm()
}

func initasm()

func AddC(c float64, d []float64)

func SubtrC(c float64, d []float64)

func MultC(c float64, d []float64)

func DivC(c float64, d []float64)

func Add(a, b []float64)

func Vadd(a, b []float64)

func Hadd(st uint64, a []float64)

func Subtr(a, b []float64)

func Mult(a, b []float64)

func Div(a, b []float64)

func Fma12(a float64, x, b []float64)

func Fma21(a float64, x, b []float64)
