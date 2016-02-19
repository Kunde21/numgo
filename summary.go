package numgo

import (
	"math"
	"sort"
)

// Sum calculates the sum result array along a given axes.
// Empty call gives the grand sum of all elements.
func (a *Array64) Sum(axis ...int) (r *Array64) {
	switch {
	case a.valAxis(&axis, "Sum"):
		return a
	case len(axis) == 0:
		tot := float64(0)
		for _, v := range a.data {
			tot += v
		}
		return Full(tot, 1)
	}

	sort.IntSlice(axis).Sort()
	n := make([]uint64, len(a.shape)-len(axis))

axisR:
	for i, t := 0, 0; i < len(a.shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axisR
			}
		}
		n[t] = a.shape[i]
		t++
	}

	ln := a.strides[0]
	for k := 0; k < len(axis); k++ {
		v, wd, st := a.shape[axis[k]], a.strides[axis[k]], a.strides[axis[k]+1]
		if st == 1 {
			for k := uint64(0); k < ln/wd; k++ {
				a.data[k] = a.data[k*wd]
				for i := uint64(1); i < wd; i++ {
					a.data[k] += a.data[i+k*wd]
				}
			}
			ln /= v
			continue
		}

		for w := uint64(0); w < ln; w += wd {
			for i := uint64(1); i*st+1 < wd; i++ {
				vadd(a.data[w:w+st], a.data[w+(i)*st:w+(i+1)*st])
			}
			copy(a.data[w/wd*st:(w/wd+1)*st], a.data[w:w+st])
		}
		ln /= v
	}
	a.shape = n

	tmp := uint64(1)
	for i := len(n); i > 0; i-- {
		a.strides[i] = tmp
		tmp *= n[i-1]
	}
	a.strides[0] = tmp
	a.data = a.data[:tmp]
	a.strides = a.strides[:len(n)+1]
	return a
}

// NaNSum calculates the sum result array along a given axes.
// All NaN values will be ignored in the Sum calculation.
// If all element values along the axis are NaN, NaN is in the return element.
//
// Empty call gives the grand sum of all elements.
func (a *Array64) NaNSum(axis ...int) *Array64 {
	if a.valAxis(&axis, "NaNSum") {
		return a
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

	return a.Fold(ns, axis...)

}

// Count gives the number of elements along a set of axis.
// Value in the element is not tested, all elements are counted.
func (a *Array64) Count(axis ...int) *Array64 {
	switch {
	case a.valAxis(&axis, "Count"):
		return a
	case len(axis) == 0:
		return full(float64(a.strides[0]), 1)
	}

	tAxis := make([]uint64, len(a.shape)-len(axis))
	cnt := uint64(1)
cntAx:
	for i, t := 0, 0; i < len(a.shape); i++ {
		for _, w := range axis {
			if i == w {
				cnt *= a.shape[i]
				continue cntAx
			}
		}
		tAxis[t] = a.shape[i]
		t++
	}
	return full(float64(cnt), tAxis...)
}

// NaNCount calculates the number of values along a given axes.
// Empty call gives the total number of elements.
func (a *Array64) NaNCount(axis ...int) *Array64 {
	if a.valAxis(&axis, "NaNCount") {
		return a
	}

	nc := func(d []float64) (r float64) {
		for _, v := range d {
			if !math.IsNaN(v) {
				r++
			}
		}
		return r
	}

	return a.Fold(nc, axis...)
}

// Mean calculates the mean across the given axes.
// NaN values in the dataa will result in NaN result elements.
func (a *Array64) Mean(axis ...int) *Array64 {
	switch {
	case a.valAxis(&axis, "Mean"):
		return a
	}
	c := a.Count(axis...).At(0)
	return a.C().Sum(axis...).DivC(c)
}

// NaNMean calculates the mean across the given axes.
// NaN values are ignored in this calculation.
func (a *Array64) NaNMean(axis ...int) *Array64 {
	switch {
	case a.valAxis(&axis, "Sum"):
		return a
	}
	return a.NaNSum(axis...).Div(a.NaNCount(axis...))
}

// Nonzero counts the number of non-zero elements are in the array
func (a *Array64) Nonzero(axis ...int) *Array64 {
	if a.valAxis(&axis, "Sum") {
		return a
	}

	cnz := func(d []float64) (r float64) {
		for _, v := range d {
			if v != 0 {
				r++
			}
		}
		return
	}
	return a.Fold(cnz, axis...)
}
