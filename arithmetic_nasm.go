//+build !amd64 noasm appengine

package numgo

var (
	sse3Supt, avxSupt, avx2Supt, fmaSupt bool
)

func initasm() {
}

func addC(c float64, d []float64) {
	for i := range d {
		d[i] += c
	}
}

func subtrC(c float64, d []float64) {
	for i := range d {
		d[i] -= c
	}
}

func multC(c float64, d []float64) {
	for i := range d {
		d[i] *= c
	}
}

func divC(c float64, d []float64) {
	for i := range d {
		d[i] /= c
	}
}

func add(a, b []float64) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] += b[j]
	}
}

func vadd(a, b []float64) {
	for i := range a {
		a[i] += b[i]
	}
}

func hadd(st uint64, a []float64) {
	ln := uint64(len(a))
	for k := uint64(0); k < ln/st; k++ {
		a[k] = a[k*st]
		for i := uint64(1); i < st; i++ {
			a[k] += a[k*st+i]
		}
	}
}

func subtr(a, b []float64) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] -= b[j]
	}
}

func mult(a, b []float64) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] *= b[j]
	}
}

func div(a, b []float64) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] /= b[j]
	}
}

func fma12(a float64, x, b []float64) {
	lnx, lnb := len(x), len(b)
	for i, j := 0, 0; i < lnx; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		x[i] = a*x[i] + b[j]
	}
}

func fma21(a float64, x, b []float64) {
	lnx, lnb := len(x), len(b)
	for i, j := 0, 0; i < lnx; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		x[i] = x[i]*b[j] + a
	}
}
