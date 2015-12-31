package numgo

import (
	"fmt"
	"math"
	"sort"
)

// Sum calculates the sum result array along a given axes.
// Empty call gives the grand sum of all elements.
func (a *Arrayf) Sum(axis ...int) (r *Arrayf) {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Sum().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	case len(axis) == 0:
		tot := float64(0)
		for _, v := range a.data {
			tot += v
		}
		return Full(tot, 1)
	}

	//Validate input
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Sum().  Shape: %v  Axes: %v", a.shape, axis)
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
					t[k] += t[z]
				}
			}
		}

		j := uint64(1)
		for ; j < uint64(len(t))/maj; j++ {
			copy(t[j*min:(j+1)*min], t[j*maj:j*maj+min])
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
	a.strides = a.strides[:len(n)+1]
	return a
}

// NaNSum calculates the sum result array along a given axes.
// All NaN values will be ignored in the Sum calculation.
// If all element values along the axis are NaN, NaN is in the return element.
//
// Empty call gives the grand sum of all elements.
func (a *Arrayf) NaNSum(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by NaNSum().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	}
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Sum().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}

	ns := func(d []float64) (r float64) {
		flag := false
		for _, v := range d {
			if !math.IsNaN(v) {
				flag = true
				r += v
			}
		}
		if flag {
			return r
		}
		return math.NaN()
	}

	return a.Map(ns, axis...)

}

// Count gives the number of elements along a set of axis.
// Value in the element is not tested, all elements are counted.
func (a *Arrayf) Count(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Count().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	case len(axis) == 0:
		return full(float64(a.strides[0]), 1)
	}

	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Sum().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}

	tAxis := make([]uint64, len(a.shape)-len(axis))
	cnt := uint64(1)
	for i, t := 0, 0; i < len(a.shape); i++ {
		tmp := false
		for _, w := range axis {
			if i == w {
				tmp = true
				break
			}
		}
		if !tmp {
			tAxis[t] = a.shape[i]
			t++
		} else {
			cnt *= a.shape[i]
		}
	}
	return full(float64(cnt), tAxis...)
}

// NaNCount calculates the number of values along a given axes.
// Empty call gives the total number of elements.
func (a *Arrayf) NaNCount(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Count().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	}
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Sum().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}

	nc := func(d []float64) (r float64) {
		for _, v := range d {
			if !math.IsNaN(v) {
				r++
			}
		}
		return r
	}

	return a.Map(nc, axis...)
}

// Mean calculates the mean across the given axes.
// NaN values in the dataa will result in NaN result elements.
func (a *Arrayf) Mean(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Mean().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	case len(axis) == 0:
		return a.C()
	}

	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Mean().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}
	return a.C().Sum(axis...).DivC(a.Count(axis...).data[0])
}

// NaNMean calculates the mean across the given axes.
// NaN values are ignored in this calculation.
func (a *Arrayf) NaNMean(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	switch {
	case a == nil || a.err != nil:
		return a
	case len(axis) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many axes received by Mean().  Shape: %v  Axes: %v", a.shape, axis)
		}
		return a
	case len(axis) == 0:
		return a.C()
	}

	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Axis out of range received by Mean().  Shape: %v  Axes: %v", a.shape, axis)
			}
			return a
		}
	}
	return a.NaNSum(axis...).Div(a.NaNCount(axis...))
}
