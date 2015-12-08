package numgo

import (
	"fmt"
	"math"
	"sort"
)

// Sum calculates the sum result array along a given axes.
// Empty call gives the grand sum of all elements.
func (a *Arrayf) Sum(axis ...int) *Arrayf {
	if a == nil {
		return nil
	}

	if len(axis) == 0 {
		tot := float64(0)
		for _, v := range a.data {
			tot += v
		}
		return Full(tot, 1)
	}

	axis = cleanAxis(axis...)

	//Validate input
	for _, v := range axis {
		if v < 0 || v > len(a.shape) {
			fmt.Println("Axis outside of range", v)
			return nil
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

// Count calculates the number of values along a given axes.
// Empty call gives the total number of elements.
func (a *Arrayf) NaNSum(axis ...int) *Arrayf {

	if a == nil {
		return nil
	}
	ns := func(d []float64) (r float64) {
		for _, v := range d {
			if !math.IsNaN(v) {
				r += v
			}
		}
		return r
	}

	return a.Map(ns, axis...)

}

// Any will return true if any element is non-zero, false otherwise.
func (a *Arrayf) Count(axis ...int) *Arrayf {
	if a == nil {
		return nil
	}
	if len(axis) == 0 {
		return full(float64(a.shape[0]), 1)
	}

	axis = cleanAxis(axis...)
	ret := full(1, a.shape...).Sum(axis...)

	return ret
}

// Count calculates the number of values along a given axes.
// Empty call gives the total number of elements.
func (a *Arrayf) NaNCount(axis ...int) *Arrayf {
	if a == nil {
		return nil
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

// Mean calculates the mean across the given axes
func (a *Arrayf) Mean(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	return a.C().Sum(axis...).Div(a.Count(axis...))
}

// NaNMean calculates the mean across the given axes
// NaN values are ignored in this calculation
func (a *Arrayf) NaNMean(axis ...int) *Arrayf {
	axis = cleanAxis(axis...)
	return a.C().NaNSum(axis...).Div(a.NaNCount(axis...))
}
