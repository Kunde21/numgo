package numgo

import (
	"fmt"
	"runtime"
)

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
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	}
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by %s().  Shape: %v  Axes: %v", mthd, a.shape, axis)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
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
func MaxSet(arrSet ...*Array64) (b *Array64) {
	if b = valSet(arrSet, "MaxSet"); b != nil {
		return b
	}

	b = arrSet[0].C()

	for j := 1; j < len(arrSet); j++ {
		for i := range b.data {
			if arrSet[j].data[i] > b.data[i] {
				b.data[i] = arrSet[j].data[i]
			}
		}
	}
	return
}

// MinSet will return the element-wise maximum of arrays.
//
// All arrays must be the non-nil and the same shape.
func MinSet(arrSet ...*Array64) (b *Array64) {
	if b = valSet(arrSet, "MaxSet"); b != nil {
		return b
	}

	b = arrSet[0].C()

	for j := 1; j < len(arrSet); j++ {
		for i := range b.data {
			if arrSet[j].data[i] < b.data[i] {
				b.data[i] = arrSet[j].data[i]
			}
		}
	}
	return
}

func valSet(arrSet []*Array64, mthd string) (b *Array64) {

	if len(arrSet) == 0 {
		b = &Array64{err: NilError}
		if debug {
			b.debug = mthd + "() called with no arrays"
			b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return b
	}

	a := arrSet[0]
	for _, v := range arrSet {
		if v == nil {
			b = &Array64{err: NilError}
			if debug {
				b.debug = mthd + "() received a Nil pointer array as an argument."
				b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return b
		}
		if v.err != nil {
			b = &Array64{err: v.err}
			if debug {
				b.debug = "Error in data passed to " + mthd + "()."
				b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return b
		}

		for k, s := range v.shape {
			if s != a.shape[k] {
				b = &Array64{err: ShapeError}
				if debug {
					b.debug = fmt.Sprintf("Array received by %s() does not match shape.  Shape: %v  Val shape: %v", mthd, a.shape, v.shape)
					b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
				}
				return b
			}
		}
	}
	return nil
}
