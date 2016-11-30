package numgo

import (
	"fmt"
	"runtime"
)

// Max will return the maximum along the given axes.
func (a *Array64) Max(axis ...int) (r *Array64) {
	if a.valAxis(&axis, "Max") {
		return a
	}

	max := func(d []nDimElement) (r nDimElement) {
		r = d[0]
		for _, v := range d {
			if v.(float64) > r.(float64) {
				r = v
			}
		}
		return r
	}

	r = a.Fold(max, axis...)

	return r
}

// Min will return the minimum along the given axes.
func (a *Array64) Min(axis ...int) (r *Array64) {
	if a.valAxis(&axis, "Max") {
		return a
	}

	min := func(d []nDimElement) (r nDimElement) {
		r = d[0]
		for _, v := range d {
			if v.(float64) < r.(float64) {
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
	if b = b.valSet(arrSet, "MaxSet"); b != nil {
		return b
	}

	b = arrSet[0].C()

	for j := 1; j < len(arrSet); j++ {
		for i := range b.data {
			if arrSet[j].data[i].(float64) > b.data[i].(float64) {
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
	if b = b.valSet(arrSet, "MaxSet"); b != nil {
		return b
	}

	b = arrSet[0].C()

	for j := 1; j < len(arrSet); j++ {
		for i := range b.data {
			if arrSet[j].data[i].(float64) < b.data[i].(float64) {
				b.data[i] = arrSet[j].data[i]
			}
		}
	}
	return
}

func (a *Array64) valSet(arrSet []*Array64, mthd string) (b *Array64) {

	if len(arrSet) == 0 {
		b = &Array64{nDimObject{err: NilError}}
		if debug {
			b.debug = mthd + "() called with no arrays"
			b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return b
	}

	a = arrSet[0]
	for _, v := range arrSet {
		if v == nil {
			b = &Array64{nDimObject{err: NilError}}
			if debug {
				b.debug = mthd + "() received a Nil pointer array as an argument."
				b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return b
		}
		if v.err != nil {
			b = &Array64{nDimObject{err: v.err}}
			if debug {
				b.debug = "Error in data passed to " + mthd + "()."
				b.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return b
		}

		for k, s := range v.shape {
			if s != a.shape[k] {
				b = &Array64{nDimObject{err: ShapeError}}
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
