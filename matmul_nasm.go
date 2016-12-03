//+build !amd64 noasm appengine

package numgo

func DotProd(a, b []nDimElement) float64 {
	var ret float64
	for i := range a {
		ret += a[i].(float64) * b[i].(float64)
	}
	return ret
}
