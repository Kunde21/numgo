package numgo

import "fmt"

// Max will return the maximum along the given axes.
func (a *Array64) Max(axis ...int) (r *Array64) {
	if valAxis(a, axis, "Max") {
		return a
	}

	max := func(d []float64) (r float64) {
		r = d[0]
		for _, v := range d {
			if v > r {
				r = v
			}
		}
		return r
	}

	r = a.Fold(max, axis...)

	return r
}

func valAxis(a *Array64, axis []int, mthd string) bool {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return true
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by %s().  Shape: %v  Axes: %v", mthd, a.shape, axis)
		}
		return true
	}
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by %s().  Shape: %v  Axes: %v", mthd, a.shape, axis)
			}
			return true
		}
	}
	return false

}

// Min will return the minimum along the given axes.
func (a *Array64) Min(axis ...int) (r *Array64) {
	if valAxis(a, axis, "Max") {
		return a
	}

	min := func(d []float64) (r float64) {
		r = d[0]
		for _, v := range d {
			if v < r {
				r = v
			}
		}
		return r
	}

	r = a.Fold(min, axis...)

	return r
}

// MaxSet will return the element-wise maximum of arrays.
//
// All arrays must be the non-nil and the same shape.
func (a *Array64) MaxSet(arrSet ...*Array64) (b *Array64) {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(arrSet) == 0:
		return a.C()
	}
	for _, v := range arrSet {
		if v == nil {
			a.err = NilError
			if debug {
				a.debug = "MaxSet() received a Nil pointer array as an argument."
			}
			return a
		}
		for k, s := range v.shape {
			if s != a.shape[k] {
				a.err = ShapeError
				if debug {
					a.debug = fmt.Sprintf("Array received by MaxSet() does not match shape.  Shape: %v  Val shape: %v", a.shape, v.shape)
				}
				return a
			}
		}
	}

	b = newArray64(a.shape...)

	for i, v := range a.data {
		b.data[i] = v
		for _, s := range arrSet {
			if s.data[i] > b.data[i] {
				b.data[i] = s.data[i]
			}
		}
	}

	return
}

// MinSet will return the element-wise maximum of arrays.
//
// All arrays must be the non-nil and the same shape.
func (a *Array64) MinSet(arrSet ...*Array64) (b *Array64) {
	if a == nil || a.err != nil {
		return a
	}
	for _, v := range arrSet {
		if v == nil {
			a.err = NilError
			if debug {
				a.debug = "MinSet() received a Nil pointer array as an argument."
			}
			return a
		}
		for k, s := range v.shape {
			if s != a.shape[k] {
				a.err = ShapeError
				if debug {
					a.debug = fmt.Sprintf("Array received by MinSet() does not match shape.  Shape: %v  Val shape: %v", a.shape, v.shape)
				}
				return a
			}
		}
	}

	b = newArray64(a.shape...)

	for i, v := range a.data {
		b.data[i] = v
		for _, s := range arrSet {
			if s.data[i] < b.data[i] {
				b.data[i] = s.data[i]
			}
		}
	}

	return
}
