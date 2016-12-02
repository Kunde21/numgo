package numgo

import (
	"fmt"
	"math"
	"runtime"
	"sort"
)

// NotEq performs boolean '!=' element-wise comparison
func (a *nDimFields) NotEq(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "NotEq()")
	if r != nil {
		return r
	}

	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i != j && !(math.IsNaN(i.(float64)) && math.IsNaN(j.(float64)))
	})
	return
}

// Less performs boolean '<' element-wise comparison
func (a *nDimFields) Less(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "Less()")
	if r != nil {
		return r
	}

	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i.(float64) < j.(float64)
	})
	return
}

// LessEq performs boolean '<=' element-wise comparison
func (a *nDimFields) LessEq(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "LessEq()")
	if r != nil {
		return r
	}

	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i.(float64) <= j.(float64)
	})
	return
}

// Greater performs boolean '<' element-wise comparison
func (a *nDimFields) Greater(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "Greater()")
	if r != nil {
		return r
	}

	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i.(float64) > j.(float64)
	})
	return
}

// GreaterEq performs boolean '<=' element-wise comparison
func (a *nDimFields) GreaterEq(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "GreaterEq()")
	if r != nil {
		return r
	}

	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i.(float64) >= j.(float64)
	})
	return

}

func (a *nDimFields) compValid(b nDimObject, mthd string) (r *Arrayb) {
	c := b.fields()

	switch {
	case a == nil || a.data == nil && a.err == nil:
		r = &Arrayb{nDimFields{err: NilError}}
		if debug {
			r.debug = fmt.Sprintf("Nil pointer received by %s", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case b == nil || c.data == nil && c.err == nil:
		r = &Arrayb{nDimFields{err: NilError}}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s is a Nil Pointer.", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case a.err != nil:
		r = &Arrayb{nDimFields{err: a.err}}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	case c.err != nil:
		r = &Arrayb{nDimFields{err: c.err}}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r

	case len(a.shape) < len(c.shape):
		r = &Arrayb{nDimFields{err: ShapeError}}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, c.shape)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return r
	}

	for i, j := len(c.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != c.shape[i] {
			r = &Arrayb{nDimFields{err: ShapeError}}
			if debug {
				r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, c.shape)
				r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return r
		}
	}

	return nil
}

// Validation and error checks must be complete before calling comp
func (a *nDimFields) comp(b nDimFields, f func(i, j nDimElement) bool) (r *Arrayb) {
	r = newArrayB(b.Shape()...)

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
			if v.(bool) {
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
					t[k] = t[k].(bool) || t[z].(bool)
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
			if !v.(bool) {
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
					t[k] = t[k].(bool) && t[z].(bool)
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

// Equals performs boolean '==' element-wise comparison
func (a Array64) Equals(b nDimObject) (r *Arrayb) {
	return a.Equals(b)
}

// Equals performs boolean '==' element-wise comparison
func (a Arrayb) Equals(b nDimObject) (r *Arrayb) {
	return a.Equals(b)
}

// Equals performs boolean '==' element-wise comparison
func (a nDimFields) Equals(b nDimObject) (r *Arrayb) {
	r = a.compValid(b, "Equals()")
	if r != nil {
		return r
	}

	//r = a.comp(b, func(i, j nDimElement) bool {
	//return i == j || math.IsNaN(i.(float64)) && math.IsNaN(j.(float64))
	//})
	//return
	r = a.comp(b.fields(), func(i, j nDimElement) bool {
		return i == j
	})
	return
}
