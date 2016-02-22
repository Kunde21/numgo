//+build !noasm,!appengine

package numgo

var (
	sse3Supt, avxSupt, avx2Supt, fmaSupt bool
)

func init() {
	initasm()
}

func initasm()

func addC(c float64, d []float64)

func subtrC(c float64, d []float64)

func multC(c float64, d []float64)

func divC(c float64, d []float64)

func add(a, b []float64)

func vadd(a, b []float64)

func hadd(st uint64, a []float64)

func subtr(a, b []float64)

func mult(a, b []float64)

func div(a, b []float64)

func fma12(a float64, x, b []float64)

func fma21(a float64, x, b []float64)
