package numgo

import (
	"fmt"
	"math"
	"runtime"
)

// Flatten reshapes the data to a 1-D array.
func (a *Array64) Flatten() *Array64 {
	if a == nil || a.err != nil {
		return a
	}
	return a.Reshape(int(a.strides[0]))
}

// C will return a deep copy of the source array.
func (a *Array64) C() (b *Array64) {
	if a == nil || a.err != nil {
		return a
	}

	b = newArray64(a.shape...)
	copy(b.data, a.data)
	return
}

// At returns the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Array64) At(index ...int) float64 {
	switch {
	case a == nil || a.err != nil:
		return math.NaN()
	case len(a.shape) != len(index):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Indexes E(%v) do not match array shape %v", index, a.shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return math.NaN()
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index in E(%v) does not exist in array with shape %v", index, a.shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return math.NaN()
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return a.data[idx]
}

// SliceElement returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Array64) SliceElement(index ...int) (ret []float64) {
	switch {
	case a == nil || a.err != nil:
		return nil
	case len(a.shape)-1 != len(index):
		a.err = IndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SliceElement().  Shape: %v  Index: %v", a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return nil
	}
	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SliceElement() does not exist shape: %v index: %v", a.shape, index)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return append(ret, a.data[idx:idx+a.strides[len(a.strides)-2]]...)
}

// SubArr slices the array at a given index.
func (a *Array64) SubArr(index ...int) (ret *Array64) {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(index) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many indicies received by SubArr().  Shape: %v Indicies: %v", a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SubArr() does not exist shape: %v index: %v", a.shape, index)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
		idx += uint64(v) * a.strides[i+1]
	}

	ret = newArray64(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}

// Set sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Array64) Set(val float64, index ...int) *Array64 {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape) != len(index):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SetE().  Shape: %v Index: %v", a.shape, index)
		}
		return a
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SetE() does not exist shape: %v index: %v", a.shape, index)
			}
			return a
		}
		idx += uint64(v) * a.strides[i+1]
	}
	a.data[idx] = val
	return a
}

// SetSliceElement sets the element group at one axis above the leaf elements.
// Source Array is returned, for function-chaining design.
func (a *Array64) SetSliceElement(vals []float64, index ...int) *Array64 {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape)-1 != len(index):
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SetSliceElement().  Shape: %v  Index: %v", a.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		fallthrough
	case uint64(len(vals)) != a.shape[len(a.shape)-1]:
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Incorrect slice length received by SetSliceElement().  Shape: %v  Index: %v", a.shape, len(index))
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}
	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SetSliceElement() does not exist shape: %v index: %v", a.shape, index)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}

			return a
		}
		idx += uint64(v) * a.strides[i+1]
	}

	copy(a.data[idx:idx+a.strides[len(a.strides)-2]], vals)
	return a
}

// SetSubArr sets the array below a given index to the values in vals.
// Values will be broadcast up multiple axes if the shapes match.
func (a *Array64) SetSubArr(vals *Array64, index ...int) *Array64 {
	switch {
	case a == nil || vals == nil:
		a.err = NilError
		if debug {
			a.debug = "Input array value received by SetE is a Nil pointer."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		fallthrough
	case a.err != nil:
		return a
	case vals.err != nil:
		a.err = vals.err
		if debug {
			a.debug = "Array received by SetSubArr() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
	case len(vals.shape)+len(index) > len(a.shape):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Array received by SetSubArr() cant be broadcast.  Shape: %v  Vals shape: %v index: %v", a.shape, vals.shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	for i, j := len(a.shape)-1, len(vals.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[i] != vals.shape[j] {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Shape of array recieved by SetSubArr() doesn't match receiver.  Shape: %v  Vals Shape: %v", a.shape, vals.shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SetSubArr() out of range.  Shape: %v  Index: %v", a.shape, index)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
		idx += uint64(v) * a.strides[i+1]
	}

	if len(vals.shape)-len(index)-len(a.shape) == 0 {
		copy(a.data[idx:idx+uint64(len(vals.data))], vals.data)
		return a
	}

	reps := uint64(1)
	for i := len(index); i < len(a.shape)-len(vals.shape); i++ {
		reps *= a.shape[i]
	}

	ln := uint64(len(vals.data))
	for i := uint64(1); i <= reps; i++ {
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
	case a == nil || a.err != nil:
		return a
	case len(shape) == 0:
		return newArray64(0)
	}

	var sz uint64 = 1
	a.shape = make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Resize.  Shape: %v", shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
		sz *= uint64(v)
		a.shape[i] = uint64(v)
	}

	if sz > a.strides[0] {
		a.data = append(a.data, make([]float64, a.strides[0]-sz)...)
	} else {
		a.data = a.data[:sz]
	}

	a.strides = make([]uint64, len(shape)+1)
	a.strides[0] = sz
	sz = 1
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = sz
		sz *= a.shape[i-1]
	}
	return a
}

// Append will concatenate a and val at the given axis.
//
// Source array will be changed, so use C() if the original data is needed.
// All axes must be the same except the appending axis.
func (a *Array64) Append(val *Array64, axis int) *Array64 {
	switch {
	case a == nil || a.err != nil:
		return a
	case axis >= len(a.shape) || axis < 0:
		a.err = IndexError
		if debug {
			a.debug = fmt.Sprintf("Axis received by Append() out of range.  Shape: %v  Axis: %v", a.shape, axis)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	case val.err != nil:
		a.err = val.err
		if debug {
			a.debug = "Array received by Append() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
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

	a.data = append(a.data, val.data...)

	as, vs := a.strides[axis], val.strides[axis+1]
	for i, j := a.strides[0]-as, val.strides[0]-vs; i < a.strides[0]; i, j = i-as, j-vs {
		copy(a.data[i+j+as:i+j+as+vs], val.data[j:j+vs])
		copy(a.data[i+j:i+j+as], a.data[i:i+as])
	}

	a.shape[axis] += val.shape[axis]

	tmp := a.strides[axis+1]
	for i := axis; i >= 0; i-- {
		tmp *= a.shape[i]
		a.strides[i] = tmp
	}

	return a
}
