//+build !noasm,!appengine

package numgo

import "fmt"

var (
	avxSupt, avx2Supt bool
)

func init() {
	avxSupt, avx2Supt = initasm()
	fmt.Println("AVX", avxSupt, "AVX2", avx2Supt)
}

func initasm() (a, a2 bool)

func addC(c float64, d []float64)

func subtrC(c float64, d []float64)

func multC(c float64, d []float64)

func divC(c float64, d []float64)

func add(a, b []float64)

func subtr(a, b []float64)

func mult(a, b []float64)

func div(a, b []float64)
