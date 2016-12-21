package numgo

import (
	"fmt"
	"math"
	"runtime"
	"sort"
)

// Equals performs boolean '==' element-wise comparison
func (a *Array64) Equals(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "Equals()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i == j || math.IsNaN(i) && math.IsNaN(j)
	})
	return
}

// NotEq performs boolean '1=' element-wise comparison
func (a *Array64) NotEq(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "NotEq()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i != j && !(math.IsNaN(i) && math.IsNaN(j))
	})
	return
}

// Less performs boolean '<' element-wise comparison
func (a *Array64) Less(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "Less()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i < j
	})
	return
}

// LessEq performs boolean '<=' element-wise comparison
func (a *Array64) LessEq(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "LessEq()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i <= j
	})
	return
}

// Greater performs boolean '<' element-wise comparison
func (a *Array64) Greater(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "Greater()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i > j
	})
	return
}

// GreaterEq performs boolean '<=' element-wise comparison
func (a *Array64) GreaterEq(b *Array64) (r *Arrayb) {
	r = a.compValid(b, "GreaterEq()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j float64) bool {
		return i >= j
	})
	return

}

func (a *Array64) compValid(b *Array64, mthd string) (r *Arrayb) {

	switch {
	case a == nil || a.data == nil && a.err == nil:
		r = &Arrayb{nDimMetadata{err: NilError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Nil pointer received by %s", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case b == nil || b.data == nil && b.err == nil:
		r = &Arrayb{nDimMetadata{err: NilError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s is a Nil Pointer.", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case a.err != nil:
		r = &Arrayb{nDimMetadata{err: a.err}, nil}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case b.err != nil:
		r = &Arrayb{nDimMetadata{err: b.err}, nil}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r

	case len(a.shape) < len(b.shape):
		r = &Arrayb{nDimMetadata{err: ShapeError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	}

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			r = &Arrayb{nDimMetadata{err: ShapeError}, nil}
			if debug {
				r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
				r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return r
		}
	}

	return nil
}

// Validation and error checks must be complete before calling comp
func (a *Array64) comp(b *Array64, f func(i, j float64) bool) (r *Arrayb) {
	r = newArrayB(b.shape...)

	for i := range r.data {
		r.data[i] = f(a.data[i], b.data[i])
	}

	return
}

// Any will return true if any element is non-zero, false otherwise.
func (a *Arrayb) Any(axis ...int) *Arrayb {
	if a.valAxis(&axis, "All") {
		return a
	}

	if len(axis) == 0 {
		for _, v := range a.data {
			if v {
				return Fullb(true, 1)
			}
		}
		return Fullb(false, 1)
	}

	sort.IntSlice(axis).Sort()
	n := make([]int, len(a.shape)-len(axis))
axis:
	for i, t := 0, 0; i < len(a.shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axis
			}
		}
		n[t] = a.shape[i]
		t++
	}

	t := a.data
	for i := 0; i < len(axis); i++ {

		maj, min := a.strides[axis[i]], a.strides[axis[i]]/a.shape[axis[i]]

		for j := 0; j+maj <= len(t); j += maj {
			for k := j; k < j+min; k++ {
				for z := k + min; z < j+maj; z += min {
					t[k] = t[k] || t[z]
				}
			}
		}

		j := 1
		for ; j < len(t)/maj; j++ {
			copy(t[j*min:(j+1)*min], t[j*maj:j*maj+min])
		}

		t = append(t[:0], t[0:j*min]...)
	}
	a.data = t
	a.shape = n

	tmp := 1
	for i := len(n); i > 0; i-- {
		a.strides[i] = tmp
		tmp *= n[i-1]
	}
	a.strides[0] = tmp
	a.strides = a.strides[0 : len(n)+1]
	return a
}

// All will return true if all elements are non-zero, false otherwise.
func (a *Arrayb) All(axis ...int) *Arrayb {

	if a.valAxis(&axis, "All") {
		return a
	}

	if len(axis) == 0 {
		for _, v := range a.data {
			if !v {
				return Fullb(false, 1)
			}
		}
		return Fullb(true, 1)
	}

	sort.IntSlice(axis).Sort()
	n := make([]int, len(a.shape)-len(axis))
axis:
	for i, t := 0, 0; i < len(a.shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axis
			}
		}
		n[t] = a.shape[i]
		t++
	}

	t := a.data
	for i := 0; i < len(axis); i++ {

		maj, min := a.strides[axis[i]], a.strides[axis[i]]/a.shape[axis[i]]

		for j := 0; j+maj <= len(t); j += maj {
			for k := j; k < j+min; k++ {
				for z := k + min; z < j+maj; z += min {
					t[k] = t[k] && t[z]
				}
			}
		}

		j := 1
		for ; j < len(t)/maj; j++ {
			a := t[j*min : (j+1)*min]
			b := t[j*maj : j*maj+min]
			copy(a, b)
		}

		t = append(t[:0], t[0:j*min]...)
	}
	a.data = t
	a.shape = n

	tmp := 1
	for i := len(n); i > 0; i-- {
		a.strides[i] = tmp
		tmp *= n[i-1]
	}
	a.strides[0] = tmp
	a.strides = append(a.strides[:0], a.strides[0:len(n)+1]...)
	return a
}

func (a *Arrayb) valAxis(axis *[]int, mthd string) bool {
	axis = cleanAxis(axis)
	switch {
	case a == nil || a.err != nil:
		return true
	case len(*axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by %s().  Shape: %v  Axes: %v", mthd, a.shape, axis)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	}
	for _, v := range *axis {
		if v < 0 || v >= len(a.shape) {
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

// Equals performs boolean '==' element-wise comparison
func (a *Arrayb) Equals(b *Arrayb) (r *Arrayb) {
	r = a.compValid(b, "Equals()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j bool) bool {
		return i == j
	})
	return
}

// NotEq performs boolean '1=' element-wise comparison
func (a *Arrayb) NotEq(b *Arrayb) (r *Arrayb) {
	r = a.compValid(b, "NotEq()")
	if r != nil {
		return r
	}

	r = a.comp(b, func(i, j bool) bool {
		return i != j
	})
	return
}

func (a *Arrayb) compValid(b *Arrayb, mthd string) (r *Arrayb) {

	switch {
	case a == nil || a.data == nil && a.err == nil:
		r = &Arrayb{nDimMetadata{err: NilError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Nil pointer received by %s", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case b == nil || b.data == nil && b.err == nil:
		r = &Arrayb{nDimMetadata{err: NilError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s is a Nil Pointer.", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case a.err != nil:
		r = &Arrayb{nDimMetadata{err: a.err}, nil}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case b.err != nil:
		r = &Arrayb{nDimMetadata{err: b.err}, nil}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r

	case len(a.shape) < len(b.shape):
		r = &Arrayb{nDimMetadata{err: ShapeError}, nil}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	}

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			r = &Arrayb{nDimMetadata{err: ShapeError}, nil}
			if debug {
				r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
				r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return r
		}
	}

	return nil
}

// Validation and error checks must be complete before calling comp
func (a *Arrayb) comp(b *Arrayb, f func(i, j bool) bool) (r *Arrayb) {
	r = newArrayB(b.shape...)

	for i := range r.data {
		r.data[i] = f(a.data[i], b.data[i])
	}

	return
}
