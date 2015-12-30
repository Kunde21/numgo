package numgo

import (
	"fmt"
	"sort"
)

// Max will return the maximum along the given axes.
func (a *Arrayf) Max(axis ...int) (r *Arrayf) {
	max := func(d []float64) (r float64) {
		sort.Sort(sort.Reverse(sort.Float64Slice(d)))
		return d[0]
	}

	r = a.Map(max, axis...)

	return r
}

// Min will return the minimum along the given axes.
func (a *Arrayf) Min(axis ...int) (r *Arrayf) {
	min := func(d []float64) (r float64) {
		sort.Float64s(d)
		return d[0]
	}

	r = a.Map(min, axis...)

	return r
}

// MaxSet will return the element-wise maximum of arrays.
//
// All arrays must be the non-nil and the same shape.
func (a *Arrayf) MaxSet(arrSet ...*Arrayf) (b *Arrayf) {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		if debug {
			a.debug = "Nil pointer received by MaxSet()."
		}
		fallthrough
	case a.err != nil:
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

	b = create(a.shape...)

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
func (a *Arrayf) MinSet(arrSet ...*Arrayf) (b *Arrayf) {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		if debug {
			a.debug = "Nil pointer received by MinSet()."
		}
		fallthrough
	case a.err != nil:
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

	b = create(a.shape...)

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
