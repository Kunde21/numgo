//+build !amd64 noasm appengine

package numgo

var (
	avxSupt, avx2Supt, fmaSupt bool
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
	ln := len(a)
	for k := 0; k < ln/st; k++ {
		a.data[k] = a.data[k*st]
		for _, v := range a.data[k*st : (k+1)*st] {
			a.data += v
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
