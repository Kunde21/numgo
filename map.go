package numgo

import (
	"fmt"
	"runtime"
	"sort"
)

// FoldFunc can be received by Fold and FoldCC to apply as a summary function
// across one or multiple axes.
type FoldFunc func([]float64) float64

// MapFunc can be received by Map to modify each element in an array.
type MapFunc func(float64) float64

// cleanAxis removes any duplicate axes and returns the cleaned slice.
// only the first instance of an axis is retained.
func cleanAxis(axis *[]int) *[]int {

	if len(*axis) < 2 {
		return axis
	}
	length := len(*axis) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if (*axis)[i] == (*axis)[j] {
				if j == length {
					*axis = (*axis)[:j]
				} else {
					*axis = append((*axis)[:j], (*axis)[j+1:]...)
				}
				length--
				j--
			}
		}
	}
	return axis
}

func (a *Array64) valAxis(axis *[]int, mthd string) bool {
	axis = cleanAxis(axis)
	switch {
	case a.HasErr():
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
	if len(*axis) == len(a.shape) {
		*axis = (*axis)[:0]
	}
	return false
}

// collapse will reorganize data by putting element dataset in continuous sections of data slice.
// Returned Arrayf must be condensed with a summary calculation to create a valid array object.
func (a *Array64) collapse(axis []int) (uint64, *Array64) {
	if len(axis) == 0 {
		r := newArray64(1)
		r.data = append(r.data[:0], a.data...)
		return a.strides[0], r
	}

	span := uint64(1) // Span = size of "element" Mx = slicing
	mx := a.strides[len(a.strides)-1]
	steps := make([]uint64, len(axis)) // Element strides
	brks := make([]uint64, len(axis))  // Stride-ending breaks

	for i, v := range axis {
		span *= uint64(a.shape[v])
		steps[i], brks[i] = a.strides[v+1], a.strides[v]
		if brks[i] > mx {
			mx = brks[i]
		}
	}

	ln := len(a.shape) - len(axis)
	asteps := make([]uint64, ln) // Element strides
	abrks := make([]uint64, ln)  // Stride-ending breaks
	newShape := make([]uint64, ln)
	sort.Ints(axis)

shape:
	for i, j := 0, len(asteps)-1; i < len(a.shape); i++ {
		for _, v := range axis {
			if i == v {
				continue shape
			}
		}
		newShape[ln-j-1] = a.shape[i]
		asteps[j], abrks[j] = a.strides[i+1], a.strides[i]
		j--
	}

	tmp := make([]float64, a.strides[0]) // Holds re-arranged data for return
	retChan, compChan := make(chan struct{}), make(chan struct{})
	defer close(retChan)
	defer close(compChan)

	go func() {
		for sl := uint64(0); sl+mx <= a.strides[0]; sl += mx {
			<-retChan
		}
		compChan <- struct{}{}
	}()

	for sl := uint64(0); sl+mx <= a.strides[0]; sl += mx {

		go func(sl uint64) {
			inc := make([]uint64, len(axis))              // N-dimensional incrementor
			off := make([]uint64, len(a.shape)-len(axis)) // N-dimensional offset incrementor

			// Inner loop might be made concurrent using slices
			// Unknown performance gains in doing so, tuning needed
			offset := uint64(0)
			for sp := uint64(0); sp+span <= mx; sp += span {

				for i, k := uint64(0), uint64(0); i < span; i++ {
					tmp[sl+i+sp] = a.data[sl+k+offset]

					k, inc[0] = k+steps[0], inc[0]+steps[0]

					// Incrementor loop to handle all dims
					for c, v := range brks {
						if uint64(i+1) == span {
							// Reset at end of loop
							inc[c] = 0
						}
						if inc[c] >= v {
							k = k - v + steps[c+1]
							inc[c] -= v
							inc[c+1] += steps[c+1]
						}
					}
				}

				// Increment to the next dimension
				offset, off[0] = offset+asteps[0], off[0]+1

				for c, v := range abrks {
					if sp+span == mx {
						// Reset at end of loop
						off[c] = 0
					}
					if off[c] >= v && c+1 < len(off) {
						offset = offset - v + asteps[c+1]
						off[c] -= v
						off[c+1] += asteps[c+1]
					}

				}
			}
			retChan <- struct{}{}
		}(sl)
	}

	<-compChan

	// Create return object.  Data is invalid format until reform is called.
	b := new(Array64)
	b.shape = newShape
	b.strides = make([]uint64, len(b.shape)+1)
	b.data = tmp

	t := uint64(1)
	for i := len(b.strides) - 1; i > 0; i-- {
		b.strides[i] = t
		t *= b.shape[i-1]
	}
	b.strides[0] = t

	return span, b
}

// FoldCC applies function f along the given axes concurrently. Each call to f will launch a goroutine.
// In order to leverage this concurrency, MapCC should only be used for complex and CPU-heavy functions.
//
// Simple functions should use Fold(f, axes...), as it's more performant on small functions.
func (a *Array64) FoldCC(f FoldFunc, axis ...int) (ret *Array64) {
	if a.valAxis(&axis, "FoldCC") {
		return a
	}

	type rt struct {
		index uint64
		value float64
	}

	rfunc := func(c chan rt, i uint64) {
		if r := recover(); r != nil {
			ret = a
			ret.err = FoldMapError
			ret.debug = fmt.Sprint(r)
			if debug {
				ret.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			c <- rt{i, 0}
		}
	}

	span, ret := a.collapse(axis)

	retChan, compChan := make(chan rt), make(chan struct{})
	defer close(retChan)
	defer close(compChan)
	go func() {
		d := make([]float64, ret.strides[0])
		for i := uint64(0); i+span <= a.strides[0]; i += span {
			c := <-retChan
			d[c.index] = c.value
		}
		ret.data = d
		compChan <- struct{}{}
	}()

	for i := uint64(0); i+span <= a.strides[0]; i += span {
		go func(i uint64) {
			defer rfunc(retChan, i/span)
			retChan <- rt{i / span, f(ret.data[i : i+span])}
		}(i)
	}
	<-compChan
	ret.data = ret.data[:ret.strides[0]]
	return ret
}

// Fold applies function f along the given axes.
// Slice containing all data to be consolidated into an element will be passed to f.
// Return value will be the resulting element's value.
func (a *Array64) Fold(f FoldFunc, axis ...int) (ret *Array64) {
	if a.valAxis(&axis, "Fold") {
		return a
	}

	defer func() {
		if r := recover(); r != nil {
			ret = a
			ret.err = FoldMapError
			ret.debug = fmt.Sprint(r)
			if debug {
				ret.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
		}
	}()

	span, ret := a.collapse(axis)
	for i := uint64(0); i+span <= a.strides[0]; i += span {
		ret.data[i/span] = f(ret.data[i : i+span])
	}
	ret.data = ret.data[:ret.strides[0]]
	return ret
}

// Map applies function f to each element in the array.
func (a *Array64) Map(f MapFunc) (ret *Array64) {
	if a == nil || a.err != nil {
		return a
	}

	defer func() {
		if r := recover(); r != nil {
			ret = a
			ret.err = FoldMapError
			ret.debug = fmt.Sprint(r)
			if debug {
				ret.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
		}
	}()

	ret = newArray64(a.shape...)
	for i := uint64(0); i < a.strides[0]; i++ {
		ret.data[i] = f(a.data[i])
	}
	return
}
