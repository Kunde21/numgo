package numgo

import (
	"fmt"
	"math"
	"runtime"
)

// Flatten reshapes the data to a 1-D array.
func (a *Array64) Flatten() *Array64 {
	if a.HasErr() {
		return a
	}
	return a.Reshape(a.strides[0])
}

// C will return a deep copy of the source array.
func (a *Array64) C() (b *Array64) {
	if a.HasErr() {
		return a
	}

	b = &Array64{
		nDimObject{
			strides: make([]int, len(a.strides)),
			shape:   make([]int, len(a.shape)),
			err:     nil,
			debug:   "",
			stack:   "",
			data:    make([]nDimElement, a.strides[0])},
	}

	copy(b.shape, a.shape)
	copy(b.strides, a.strides)
	copy(b.data, a.data)
	return b
}

// Shape returns a copy of the array shape
func (a *Array64) Shape() []int {
	if a.HasErr() {
		return nil
	}

	res := make([]int, len(a.shape), len(a.shape))
	copy(res, a.shape)
	return res
}

// At returns the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Array64) At(index ...int) float64 {
	idx := a.valIdx(index, "At")
	if a.HasErr() {
		return math.NaN()
	}

	return a.data[idx].(float64)
}

func (a *Array64) at(index []int) nDimElement {
	var idx int
	for i, v := range index {
		idx += v * a.strides[i+1]
	}
	return a.data[idx]
}

func (a *Array64) valIdx(index []int, mthd string) (idx int) {
	if a.HasErr() {
		return 0
	}
	if len(index) > len(a.shape) {
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by %s().  Shape: %v  Index: %v", mthd, a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return 0
	}
	for i, v := range index {
		if v >= a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by %s() does not exist shape: %v index: %v", mthd, a.shape, index)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return 0
		}
		idx += v * a.strides[i+1]
	}
	return
}

// SliceElement returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Array64) SliceElement(index ...int) (ret []nDimElement) {
	idx := a.valIdx(index, "SliceElement")
	switch {
	case a.HasErr():
		return nil
	case len(a.shape)-1 != len(index):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SliceElement().  Shape: %v  Index: %v", a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return nil
	}

	return append(ret, a.data[idx:idx+a.strides[len(a.strides)-2]]...)
}

// SubArr slices the array at a given index.
func (a *Array64) SubArr(index ...int) (ret *Array64) {
	idx := a.valIdx(index, "SubArr")
	if a.HasErr() {
		return a
	}

	ret = newArray64(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}

// Set sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Array64) Set(val nDimElement, index ...int) *Array64 {
	idx := a.valIdx(index, "Set")
	if a.HasErr() {
		return a
	}

	a.data[idx] = val
	return a
}

// SetSliceElement sets the element group at one axis above the leaf elements.
// Source Array is returned, for function-chaining design.
func (a *Array64) SetSliceElement(vals []nDimElement, index ...int) *Array64 {
	idx := a.valIdx(index, "SetSliceElement")
	switch {
	case a.HasErr():
		return a
	case len(a.shape)-1 != len(index):
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SetSliceElement().  Shape: %v  Index: %v", a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		fallthrough
	case len(vals) != a.shape[len(a.shape)-1]:
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect slice length received by SetSliceElement().  Shape: %v  Index: %v", a.shape, len(index))
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	copy(a.data[idx:idx+a.strides[len(a.strides)-2]], vals[:a.strides[len(a.strides)-2]])
	return a
}

// SetSubArr sets the array below a given index to the values in vals.
// Values will be broadcast up multiple axes if the shapes match.
func (a *Array64) SetSubArr(vals *Array64, index ...int) *Array64 {
	idx := a.valIdx(index, "SetSubArr")
	switch {
	case a.HasErr():
		return a
	case vals.HasErr():
		a.err = vals.getErr()
		if debug {
			a.debug = "Array received by SetSubArr() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	case len(vals.shape)+len(index) > len(a.shape):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Array received by SetSubArr() cant be broadcast.  Shape: %v  Vals shape: %v index: %v", a.shape, vals.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	for i, j := len(a.shape)-1, len(vals.shape)-1; j >= 0; i, j = i-1, j-1 {
		if a.shape[i] != vals.shape[j] {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Shape of array recieved by SetSubArr() doesn't match receiver.  Shape: %v  Vals Shape: %v", a.shape, vals.shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
	}

	if len(a.shape)-len(index)-len(vals.shape) == 0 {
		copy(a.data[idx:idx+len(vals.data)], vals.data)
		return a
	}

	reps := 1
	for i := len(index); i < len(a.shape)-len(vals.shape); i++ {
		reps *= a.shape[i]
	}

	ln := len(vals.data)
	for i := 1; i <= reps; i++ {
		copy(a.data[idx+ln*(i-1):idx+ln*i], vals.data)
	}
	return a
}

// Resize will change the underlying array size.
//
// Make a copy C() if the original array needs to remain unchanged.
// Element location in the underlying slice will not be adjusted to the new shape.
func (a *Array64) Resize(shape ...int) *Array64 {
	switch {
	case a.HasErr():
		return a
	case len(shape) == 0:
		tmp := newArray64(0)
		a.shape, a.strides = tmp.shape, tmp.strides
		a.data = tmp.data
		return a
	}

	var sz = 1
	for _, v := range shape {
		if v >= 0 {
			sz *= v
			continue
		}

		a.err = NegativeAxis
		if debug {
			a.debug = fmt.Sprintf("Negative axis length received by Resize.  Shape: %v", shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	ln, cp := len(shape), cap(a.shape)
	if ln > cp {
		a.shape = append(a.shape[:cp], make([]int, ln-cp)...)
	} else {
		a.shape = a.shape[:ln]
	}

	ln, cp = ln+1, cap(a.strides)
	if ln > cp {
		a.strides = append(a.strides[:cp], make([]int, ln-cp)...)
	} else {
		a.strides = a.strides[:ln]
	}

	a.strides[ln-1] = 1
	for i := ln - 2; i >= 0; i-- {
		a.shape[i] = shape[i]
		a.strides[i] = a.shape[i] * a.strides[i+1]
	}

	cp = cap(a.data)
	if sz > cp {
		a.data = append(a.data[:cp], make([]nDimElement, sz-cp)...)
	} else {
		a.data = a.data[:sz]
	}

	return a
}

// Append will concatenate a and val at the given axis.
//
// Source array will be changed, so use C() if the original data is needed.
// All axes must be the same except the appending axis.
func (a *Array64) Append(val *Array64, axis int) *Array64 {
	switch {
	case a.HasErr():
		return a
	case axis >= len(a.shape), axis < 0:
		a.err = IndexError
		if debug {
			a.debug = fmt.Sprintf("Axis received by Append() out of range.  Shape: %v  Axis: %v", a.shape, axis)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	case val.HasErr():
		a.err = val.GetErr()
		if debug {
			a.debug = "Array received by Append() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	case len(a.shape) != len(val.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.shape, val.shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	for k, v := range a.shape {
		if v != val.shape[k] && k != axis {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.shape, val.shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
	}

	ln := len(a.data) + len(val.data)
	var dat []nDimElement
	cp := cap(a.data)
	if ln > cp {
		dat = make([]nDimElement, ln)
	} else {
		dat = a.data[:ln]
	}

	as, vs := a.strides[axis], val.strides[axis]
	for i, j := a.strides[0], val.strides[0]; i > 0; i, j = i-as, j-vs {
		copy(dat[i+j-vs:i+j], val.data[j-vs:j])
		copy(dat[i+j-as-vs:i+j-vs], a.data[i-as:i])
	}

	a.data = dat
	a.shape[axis] += val.shape[axis]

	for i := axis; i >= 0; i-- {
		a.strides[i] = a.strides[i+1] * a.shape[i]
	}

	return a
}
