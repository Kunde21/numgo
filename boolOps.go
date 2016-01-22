package numgo

import (
	"fmt"
	"sort"
)

// Equals performs boolean '==' element-wise comparison
func (a *Array64) Equals(b *Array64) (r *Arrayb) {
	if debug {
		r = compValid(a, b, "Equals()")
	} else {
		r = compValid(a, b, "")
	}
	if r != nil {
		return r
	}

	r = comp(a, b, func(i, j float64) bool {
		if i == j {
			return true
		}
		return false
	})
	return
}

// Less performs boolean '<' element-wise comparison
func (a *Array64) Less(b *Array64) (r *Arrayb) {
	if debug {
		r = compValid(a, b, "Less()")
	} else {
		r = compValid(a, b, "")
	}
	if r != nil {
		return r
	}

	r = comp(a, b, func(i, j float64) bool {
		if i < j {
			return true
		}
		return false
	})
	return
}

// LessEq performs boolean '<=' element-wise comparison
func (a *Array64) LessEq(b *Array64) (r *Arrayb) {
	if debug {
		r = compValid(a, b, "LessEq()")
	} else {
		r = compValid(a, b, "")
	}
	if r != nil {
		return r
	}

	r = comp(a, b, func(i, j float64) bool {
		if i <= j {
			return true
		}
		return false
	})
	return
}

// Greater performs boolean '<' element-wise comparison
func (a *Array64) Greater(b *Array64) (r *Arrayb) {
	if debug {
		r = compValid(a, b, "Greater()")
	} else {
		r = compValid(a, b, "")
	}
	if r != nil {
		return r
	}

	r = comp(a, b, func(i, j float64) bool {
		if i > j {
			return true
		}
		return false
	})
	return
}

// GreaterEq performs boolean '<=' element-wise comparison
func (a *Array64) GreaterEq(b *Array64) (r *Arrayb) {
	if debug {
		r = compValid(a, b, "GreaterEq()")
	} else {
		r = compValid(a, b, "")
	}
	if r != nil {
		return r
	}

	r = comp(a, b, func(i, j float64) bool {
		if i >= j {
			return true
		}
		return false
	})
	return
}

func compValid(a, b *Array64, mthd string) (r *Arrayb) {

	switch {
	case a == nil:
		r = &Arrayb{err: NilError}
		if debug {
			r.debug = fmt.Sprintf("Nil pointer received by %s", mthd)
		}
		return r
	case b == nil:
		r = &Arrayb{err: NilError}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s is a Nil Pointer.", mthd)
		}
		return r
	case a.err != nil:
		r = &Arrayb{err: a.err}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
		}
		return r
	case b.err != nil:
		r = &Arrayb{err: b.err}
		if debug {
			r.debug = fmt.Sprintf("Error in %s arrays", mthd)
		}
		return r

	case len(a.shape) < len(b.shape):
		r = &Arrayb{err: ShapeError}
		if debug {
			r.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
		}
		return r
	}

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			r = &Arrayb{err: ShapeError}
			if debug {
				a.debug = fmt.Sprintf("Array received by %s can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
			}
			return r
		}
	}

	return nil
}

// Validation and error checks must be complete before calling comp
func comp(a, b *Array64, f func(i, j float64) bool) (r *Arrayb) {
	r = newArrayB(b.shape...)

	for i := range r.data {
		r.data[i] = f(a.data[i], b.data[i])
	}

	/*
		compChan := make(chan struct{})
		mul := len(a.data) / len(b.data)

		for k := 0; k < mul; k++ {
			go func(m int) {
				for i, v := range b.data {
					r.data[i] = f(a.data[i+m], v)
				}
				compChan <- struct{}{}
			}(k * len(b.data))
		}

		for k := 0; k < mul; k++ {
			<-compChan
		}
		close(compChan)
	*/
	return
}

// Any will return true if any element is non-zero, false otherwise.
func (a *Arrayb) Any(axis ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape) < len(axis):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Any().  Shape: %v  Axes: %v", a.shape, axis)
		}
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

	//Validate input
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Illegal axis received by Any().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}

	sort.IntSlice(axis).Sort()
	n := make([]uint64, len(a.shape)-len(axis))
	for i, t := 0, 0; i < len(a.shape); i++ {
		tmp := false
		for _, w := range axis {
			if i == w {
				tmp = true
				break
			}
		}
		if !tmp {
			n[t] = a.shape[i]
			t++
		}
	}

	t := a.data
	for i := 0; i < len(axis); i++ {

		maj, min := a.strides[axis[i]], a.strides[axis[i]+1]

		for j := uint64(0); j+maj <= uint64(len(t)); j += maj {
			for k := j; k < j+min; k += 1 {
				for z := k + min; z < j+maj; z += min {
					t[k] = t[k] || t[z]
				}
			}
		}

		j := uint64(1)
		for ; j < uint64(len(t))/maj; j++ {
			a := t[j*min : (j+1)*min]
			b := t[j*maj : j*maj+min]
			copy(a, b)
		}

		t = append(t[:0], t[0:j*min]...)
	}
	a.data = t
	a.shape = n

	tmp := uint64(1)
	for i := len(n); i > 0; i-- {
		a.strides[i] = tmp
		tmp *= n[i-1]
	}
	a.strides[0] = tmp
	a.strides = a.strides[0 : len(n)+1]
	return a
}

// Any will return true if all elements are non-zero, false otherwise.
func (a *Arrayb) All(axis ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape) < len(axis):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by All().  Shape: %v  Axes: %v", a.shape, axis)
		}
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

	//Validate input
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			return a
		}
	}

	sort.IntSlice(axis).Sort()
	n := make([]uint64, len(a.shape)-len(axis))
	for i, t := 0, 0; i < len(a.shape); i++ {
		tmp := false
		for _, w := range axis {
			if i == w {
				tmp = true
				break
			}
		}
		if !tmp {
			n[t] = a.shape[i]
			t++
		}
	}

	t := a.data
	for i := 0; i < len(axis); i++ {

		maj, min := a.strides[axis[i]], a.strides[axis[i]+1]

		for j := uint64(0); j+maj <= uint64(len(t)); j += maj {
			for k := j; k < j+min; k += 1 {
				for z := k + min; z < j+maj; z += min {
					t[k] = t[k] && t[z]
				}
			}
		}

		j := uint64(1)
		for ; j < uint64(len(t))/maj; j++ {
			a := t[j*min : (j+1)*min]
			b := t[j*maj : j*maj+min]
			copy(a, b)
		}

		t = append(t[:0], t[0:j*min]...)
	}
	a.data = t
	a.shape = n

	tmp := uint64(1)
	for i := len(n); i > 0; i-- {
		a.strides[i] = tmp
		tmp *= n[i-1]
	}
	a.strides[0] = tmp
	a.strides = append(a.strides[:0], a.strides[0:len(n)+1]...)
	return a
}

// Nonzero counts the number of non-zero elements are in the array
func (a *Array64) Nonzero() (c *uint64) {
	if a == nil || a.err != nil {
		return nil
	}

	*c = 0
	for _, v := range a.data {
		if v != float64(0) {
			(*c)++
		}
	}
	return
}
