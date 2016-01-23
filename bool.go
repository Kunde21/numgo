package numgo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Arrayb struct {
	shape   []uint64
	strides []uint64
	data    []bool
	err     error
	debug   string
}

// NewArrayB creates an Arrayb object with dimensions given in order from outer-most to inner-most
// All values will default to false
func NewArrayB(data []bool, shape ...int) (a *Arrayb) {
	if data != nil && len(shape) == 0 {
		shape = append(shape, len(data))
	}

	a = new(Arrayb)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Createb.  Shape: %v", shape)
			}
			return
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a.shape = sh
	a.data = make([]bool, sz)
	if data != nil {
		copy(a.data, data)
	}

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.err = nil
	return
}

// Internal function to create using the shape of another array
func newArrayB(shape ...uint64) (a *Arrayb) {
	a = new(Arrayb)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a.shape = sh
	a.data = make([]bool, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.err = nil
	return
}

// Full creates an Arrayb object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the returned array.
func Fullb(val bool, shape ...int) (a *Arrayb) {
	a = NewArrayB(nil, shape...)
	if a.err != nil {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

func fullb(val bool, shape ...uint64) (a *Arrayb) {
	a = newArrayB(shape...)
	if a.err != nil {
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
		return a.err.(*ngError).s
	case a.strides[0] == 0:
		return "[]"
	}

	stride := a.strides[len(a.strides)-2]
	for i, k := uint64(0), 0; i+stride < uint64(len(a.data)); i, k = i+stride, k+1 {

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
		if i+stride != uint64(len(a.data)) {
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
	if a == nil || a.err != nil {
		return a
	}

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Reshape().  Shape: %v", shape)
			}
			return a
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	if sz != uint64(len(a.data)) {
		a.err = ReshapeError
		if debug {
			a.debug = fmt.Sprintf("Reshape() can not change data size.  Dimensions: %v reshape: %v", a.shape, shape)
		}
		return a
	}

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
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
	if a == nil || a.err != nil {
		return a
	}

	b = newArrayB(a.shape...)
	for i, v := range a.data {
		b.data[i] = v
	}
	return
}

// At returns a copy of the element at the given index.
// Any errors will return a false value and record the error for the
// HasErr() and GetErr() functions.
func (a *Arrayb) At(index ...int) bool {
	switch {
	case a == nil || a.err != nil:
		return false
	case len(a.shape) != len(index):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Indexes E(%v) do not match array shape %v", index, a.shape)
		}
		return false
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index in E(%v) does not exist in array with shape %v", index, a.shape)
			}
			return false
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return a.data[idx]
}

// SliceElement returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Arrayb) SliceElement(index ...int) (ret []bool) {
	switch {
	case a == nil || a.err != nil:
		return nil
	case len(a.shape)-1 != len(index):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SliceElement().  Shape: %v  Index: %v", a.shape, index)
		}
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Index received by SliceElement() does not exist shape: %v index: %v", a.shape, index)
			}
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return append(ret, a.data[idx:idx+a.strides[len(a.strides)-2]]...)
}

// SubArr slices the array object by the index received on each corresponding axis.
//
// These are applied startig from the top axis.
// Intermediate slicing of axes is not available at this point.
func (a *Arrayb) SubArr(index ...int) (ret *Arrayb) {
	switch {
	case a == nil || a.err != nil:
		return nil
	case len(a.shape) < len(index):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Too many indicies received by SubArr().  Shape: %v Indicies: %v", a.shape, index)
		}
		return a
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			if debug {
				a.debug = fmt.Sprintf("Index received by SubArr() does not exist shape: %v index: %v", a.shape, index)
			}
			return
		}
		idx += uint64(v) * a.strides[i+1]
	}

	ret = newArrayB(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}

// Set sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Arrayb) Set(val bool, index ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape) != len(index):
		a.err = ShapeError
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
func (a *Arrayb) SetSliceElement(vals []bool, index ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(a.shape)-1 != len(index):
		if debug {
			a.debug = fmt.Sprintf("Incorrect number of indicies received by SetSliceElement().  Shape: %v  Index: %v", a.shape, index)
		}
		fallthrough
	case uint64(len(vals)) != a.shape[len(a.shape)-1]:
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Incorrect slice length received by SetSliceElement().  Shape: %v  Index: %v", a.shape, len(index))
		}
		return a
	}
	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
			if debug {
				a.debug = fmt.Sprintf("Index received by SetSliceElement() does not exist shape: %v index: %v", a.shape, index)
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
func (a *Arrayb) SetSubArr(vals *Arrayb, index ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case vals == nil:
		a.err = NilError
		if debug {
			a.debug = "Input array value received by SetE is a Nil pointer."
		}
		return a
	case vals.err != nil:
		a.err = vals.err
		if debug {
			a.debug = "Array received by SetSubArr() is in error."
		}
	case len(vals.shape)+len(index) > len(a.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Array received by SetSubArr() cant be broadcast.  Shape: %v  Vals shape: %v index: %v", a.shape, vals.shape, index)
		}
		return a
	}

	for i, j := len(a.shape)-1, len(vals.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[i] != vals.shape[j] {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Shape of array recieved by SetSubArr() doesn't match receiver.  Shape: %v  Vals Shape: %v", a.shape, vals.shape)
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
func (a *Arrayb) Resize(shape ...int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case len(shape) == 0:
		return newArrayB(0)
	}

	var sz uint64 = 1
	a.shape = make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Resize.  Shape: %v", shape)
			}
			return a
		}
		sz *= uint64(v)
		a.shape[i] = uint64(v)
	}

	if sz > a.strides[0] {
		a.data = append(a.data, make([]bool, a.strides[0]-sz)...)
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
func (a *Arrayb) Append(val *Arrayb, axis int) *Arrayb {
	switch {
	case a == nil || a.err != nil:
		return a
	case axis >= len(a.shape) || axis < 0:
		a.err = IndexError
		if debug {
			a.debug = fmt.Sprintf("Axis received by Append() out of range.  Shape: %v  Axis: %v", a.shape, axis)
		}
		return a
	case val.err != nil:
		a.err = val.err
		if debug {
			a.debug = "Array received by Append() is in error."
		}
	case len(a.shape) != len(val.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.shape, val.shape)
		}
		return a
	}

	for k, v := range a.shape {
		if v != val.shape[k] && k != axis {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.shape, val.shape)
			}
			return a
		}
	}

	a.data = append(a.data, val.data...)

	as, vs := a.strides[axis], val.strides[axis+1]
	for i, j := a.strides[0]-as, val.strides[0]-vs; i >= 0; i, j = i-as, j-vs {
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

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *Arrayb) MarshalJSON() ([]byte, error) {
	if a == nil {
		return nil, NilError
	}
	return json.Marshal(struct {
		Shape []uint64 `json:"shape"`
		Data  []bool   `json:"data"`
		Err   int8     `json:"err,omitempty"`
	}{
		Shape: a.shape,
		Data:  a.data,
		Err:   encodeErr(a.err),
	})
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for decoding data.
// Custom Unmarshaler is needed to load/decode unexported values and build strides.
func (a *Arrayb) UnmarshalJSON(b []byte) error {

	if a == nil {
		return NilError
	}

	tmpA := new(struct {
		Shape []uint64 `json:"shape"`
		Data  []bool   `json:"data"`
		Err   int8     `json:"err,omitempty"`
	})

	err := json.Unmarshal(b, tmpA)
	if err != nil {
		return err
	}

	a.shape = tmpA.Shape
	a.data = tmpA.Data
	a.err = decodeErr(tmpA.Err)

	a.strides = make([]uint64, len(a.shape)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= a.shape[i-1]
	}
	a.strides[0] = tmp

	return nil
}
