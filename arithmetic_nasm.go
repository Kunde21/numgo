//+build !amd64 noasm appengine

package numgo

var (
	sse3Supt, avxSupt, avx2Supt, fmaSupt bool
)

func initasm() {
}

func AddC(c nDimElement, d []nDimElement) {
	for i := range d {
		d[i] = d[i].(float64) + c.(float64)
	}
}

func SubtrC(c nDimElement, d []nDimElement) {
	for i := range d {
		d[i] = d[i].(float64) - c.(float64)
	}
}

func MultC(c nDimElement, d []nDimElement) {
	for i := range d {
		d[i] = d[i].(float64) * c.(float64)
	}
}

func DivC(c nDimElement, d []nDimElement) {
	for i := range d {
		d[i] = d[i].(float64) / c.(float64)
	}
}

func Add(a, b []nDimElement) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] = a[i].(float64) + b[j].(float64)
	}
}

func Vadd(a, b []nDimElement) {
	for i := range a {
		a[i] = a[i].(float64) + b[i].(float64)
	}
}

func Hadd(st uint64, a []nDimElement) {
	ln := uint64(len(a))
	for k := uint64(0); k < ln/st; k++ {
		a[k] = a[k*st]
		for i := uint64(1); i < st; i++ {
			a[k] = a[k].(float64) + a[k*st+i].(float64)
		}
	}
}

func Subtr(a, b []nDimElement) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] = a[i].(float64) - b[j].(float64)
	}
}

func Mult(a, b []nDimElement) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] = a[i].(float64) * b[j].(float64)
	}
}

func Div(a, b []nDimElement) {
	lna, lnb := len(a), len(b)
	for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		a[i] = a[i].(float64) / b[j].(float64)
	}
}

func Fma12(a nDimElement, x, b []nDimElement) {
	lnx, lnb := len(x), len(b)
	for i, j := 0, 0; i < lnx; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		x[i] = a.(float64)*x[i].(float64) + b[j].(float64)
	}
}

func Fma21(a nDimElement, x, b []nDimElement) {
	lnx, lnb := len(x), len(b)
	for i, j := 0, 0; i < lnx; i, j = i+1, j+1 {
		if j >= lnb {
			j = 0
		}
		x[i] = x[i].(float64)*b[j].(float64) + a.(float64)
	}
}
