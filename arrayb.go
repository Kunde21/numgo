package numgo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Arrayb is an n-dimensional array of boolean values
type Arrayb struct {
	nDimMetadata
	data []bool
}

// NewArrayB creates an Arrayb object with dimensions given in order from outer-most to inner-most
// All values will default to false
func NewArrayB(data []bool, shape ...int) (a *Arrayb) {
	if data != nil && len(shape) == 0 {
		shape = append(shape, len(data))
	}

	a = new(Arrayb)
	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		if v <= 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Createb.  Shape: %v", shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		sz *= v
	}
	copy(sh, shape)

	a.shape = sh
	a.data = make([]bool, sz)
	if data != nil {
		copy(a.data, data)
	}

	a.strides = make([]int, len(sh)+1)
	tmp := 1
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.err = nil
	return
}

// Internal function to create using the shape of another array
func newArrayB(shape ...int) (a *Arrayb) {
	a = new(Arrayb)
	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		sz *= v
	}
	copy(sh, shape)

	a.shape = sh
	a.data = make([]bool, sz)

	a.strides = make([]int, len(sh)+1)
	tmp := 1
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.err = nil
	return
}

// Fullb creates an Arrayb object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the returned array.
func Fullb(val bool, shape ...int) (a *Arrayb) {
	a = NewArrayB(nil, shape...)
	if a.HasErr() || !val {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

func fullb(val bool, shape ...int) (a *Arrayb) {
	a = newArrayB(shape...)
	if a.HasErr() || !val {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrayb) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.err != nil:
		return "Error: " + a.err.Error()
	case a.shape == nil || a.strides == nil || a.data == nil:
		return "<nil>"
	case a.strides[0] == 0:
		return "[]"
	}

	stride := a.strides[len(a.strides)-2]
	for i, k := 0, 0; i+stride <= len(a.data); i, k = i+stride, k+1 {

		t := ""
		for j, v := range a.strides {
			if i%v == 0 && j < len(a.strides)-2 {
				t += "["
			}
		}

		s += strings.Repeat(" ", len(a.shape)-len(t)-1) + t
		s += fmt.Sprint(a.data[i : i+stride])

		t = ""
		for j, v := range a.strides {
			if (i+stride)%v == 0 && j < len(a.strides)-2 {
				t += "]"
			}
		}

		s += t + strings.Repeat(" ", len(a.shape)-len(t)-1)
		if i+stride != len(a.data) {
			s += "\n"
			if len(t) > 0 {
				s += "\n"
			}
		}
	}
	return
}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a *Arrayb) Reshape(shape ...int) *Arrayb {
	if a.HasErr() {
		return a
	}

	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Reshape().  Shape: %v", shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
		sz *= v
	}
	copy(sh, shape)

	if sz != len(a.data) {
		a.err = ReshapeError
		if debug {
			a.debug = fmt.Sprintf("Reshape() can not change data size.  Dimensions: %v reshape: %v", a.shape, shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	a.strides = make([]int, len(sh)+1)
	tmp := 1
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.shape = sh

	return a
}

// C will return a deep copy of the source array.
func (a *Arrayb) C() (b *Arrayb) {
	if a.HasErr() {
		return a
	}

	b = newArrayB(a.shape...)
	copy(b.data, a.data)
	return
}

// At returns a copy of the element at the given index.
// Any errors will return a false value and record the error for the
// HasErr() and GetErr() functions.
func (a *Arrayb) At(index ...int) bool {
	idx := a.valIdx(index, "At")
	if a.HasErr() {
		return false
	}
	return a.data[idx]
}

// SliceElement returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Arrayb) SliceElement(index ...int) (ret []bool) {
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

// SubArr slices the array object by the index received on each corresponding axis.
//
// These are applied startig from the top axis.
// Intermediate slicing of axes is not available at this point.
func (a *Arrayb) SubArr(index ...int) (ret *Arrayb) {
	idx := a.valIdx(index, "SubArr")
	if a.HasErr() {
		return nil
	}

	ret = newArrayB(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])
	return
}

// Set sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Arrayb) Set(val bool, index ...int) *Arrayb {
	idx := a.valIdx(index, "Set")
	if a.HasErr() {
		return a
	}

	a.data[idx] = val
	return a
}

// SetSliceElement sets the element group at one axis above the leaf elements.
// Source Array is returned, for function-chaining design.
func (a *Arrayb) SetSliceElement(vals []bool, index ...int) *Arrayb {
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
func (a *Arrayb) SetSubArr(vals *Arrayb, index ...int) *Arrayb {
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
func (a *Arrayb) Resize(shape ...int) *Arrayb {
	switch {
	case a.HasErr():
		return a
	case len(shape) == 0:
		tmp := newArrayB(0)
		a.shape, a.strides = tmp.shape, tmp.strides
		a.data = tmp.data
		return a
	}

	var sz int = 1
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
	copy(a.shape, shape)
	for i := ln - 2; i >= 0; i-- {
		a.strides[i] = a.shape[i] * a.strides[i+1]
	}

	cp = cap(a.data)
	if sz > cp {
		a.data = append(a.data[:cp], make([]bool, sz-cp)...)
	} else {
		a.data = a.data[:sz]
	}

	return a
}

// Append will concatenate a and val at the given axis.
//
// Source array will be changed, so use C() if the original data is needed.
// All axes must be the same except the appending axis.
func (a *Arrayb) Append(val *Arrayb, axis int) *Arrayb {
	switch {
	case a.HasErr():
		return a
	case axis >= len(a.shape) || axis < 0:
		a.err = IndexError
		if debug {
			a.debug = fmt.Sprintf("Axis received by Append() out of range.  Shape: %v  Axis: %v", a.shape, axis)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	case val.HasErr():
		a.err = val.getErr()
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
	var dat []bool
	cp := cap(a.data)
	if ln > cp {
		dat = append(a.data, make([]bool, ln-cp)...)
	} else {
		dat = a.data[:ln]
	}

	as, vs := a.strides[axis+1], val.strides[axis+1]
	for i, j := a.strides[0], val.strides[0]; i > 0; i, j = i-as, j-vs {
		copy(dat[i+j-vs:i+j], val.data[j-vs:j])
		copy(dat[i+j-as-vs:i+j-vs], a.data[i-as:i])
	}

	a.data = dat
	a.shape[axis] += val.shape[axis]

	for i := axis; i >= 0; i-- {
		a.strides[i] = a.strides[axis+1] * a.shape[i]
	}

	return a
}

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *Arrayb) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Shape []int  `json:"shape"`
		Data  []bool `json:"data"`
		Err   int8   `json:"err,omitempty"`
	}{
		Shape: a.shape,
		Data:  a.data,
		Err:   encodeErr(a.err),
	})
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for decoding data.
// Custom Unmarshaler is needed to load/decode unexported values and build strides.
func (a *Arrayb) UnmarshalJSON(b []byte) error {

	tmpA := new(struct {
		Shape []int  `json:"shape"`
		Data  []bool `json:"data"`
		Err   int8   `json:"err,omitempty"`
	})

	err := json.Unmarshal(b, tmpA)

	a.shape = tmpA.Shape
	a.data = tmpA.Data
	a.err = decodeErr(tmpA.Err)
	if a.data == nil && a.err == nil {
		a.err = NilError
		a.strides = nil
		return nil
	}

	a.strides = make([]int, len(a.shape)+1)
	tmp := 1
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= a.shape[i-1]
	}
	a.strides[0] = tmp

	return err
}

// helper function to validate index inputs
func (a *Arrayb) valIdx(index []int, mthd string) (idx int) {
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
