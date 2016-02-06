//+build !amd64 noasm appengine

package numgo

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
