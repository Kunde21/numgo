package numgo

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strings"
)

type nDimFields struct {
	shape        []int
	strides      []int
	err          error
	debug, stack string
	data         []nDimElement
}

type nDimObject interface {
	/* TODO: add methods */
	fields() nDimFields
	Equals(nDimObject) *Arrayb
	Shape() []int
	C() nDimObject
	Count(...int) *Array64
	Flatten() nDimObject
	GetErr() error
	getErr() error //TODO is this nessesasy?
	GetDebug() (err error, debugStr, stackTrace string)
	At(...int) nDimElement
	Set(nDimElement, ...int) *nDimFields
	SetSliceElement(vals []nDimElement, index ...int) *nDimFields
	Reshape(...int) nDimObject
	HasErr() bool
	at([]int) nDimElement //TODO is this nessesary?
	SliceElement(...int) []nDimElement
	SubArr(...int) *nDimFields
}
type nDimElement interface {
}

func (a nDimFields) fields() nDimFields {
	return a
}

// Flatten reshapes the data to a 1-D array.
func (a nDimFields) Flatten() nDimObject {
	if a.HasErr() {
		return a
	}
	b := (a.Reshape(a.strides[0]))
	return b
}

func (a nDimFields) C() nDimObject {
	return a
}

// c will return a deep copy of the source array.
func (a *nDimFields) c() (b *nDimFields) {
	if a.HasErr() {
		return a
	}

	b =
		&nDimFields{
			strides: make([]int, len(a.strides)),
			shape:   make([]int, len(a.shape)),
			err:     nil,
			debug:   "",
			stack:   "",
			data:    make([]nDimElement, a.strides[0])}

	copy(b.shape, a.shape)
	copy(b.strides, a.strides)
	copy(b.data, a.data)
	return b
}

// Shape returns a copy of the array shape
func (a nDimFields) Shape() []int {
	if a.HasErr() {
		return nil
	}

	res := make([]int, len(a.shape), len(a.shape))
	copy(res, a.shape)
	return res
}

// At returns the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a nDimFields) At(index ...int) nDimElement {
	idx := a.valIdx(index, "At")
	if a.HasErr() {
		return math.NaN()
	}

	return a.data[idx]
}

func (a nDimFields) at(index []int) nDimElement {
	var idx int
	for i, v := range index {
		idx += v * a.strides[i+1]
	}
	return a.data[idx]
}

func (a *nDimFields) valIdx(index []int, mthd string) (idx int) {
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
func (a nDimFields) SliceElement(index ...int) (ret []nDimElement) {
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
func (a nDimFields) SubArr(index ...int) (ret *nDimFields) {
	idx := a.valIdx(index, "SubArr")
	if a.HasErr() {
		return &a
	}

	ret = &newArray64(a.shape[len(index):]...).nDimFields
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}

// Set sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a nDimFields) Set(val nDimElement, index ...int) *nDimFields {
	idx := a.valIdx(index, "Set")
	if a.HasErr() {
		return &a
	}

	a.data[idx] = val
	return &a
}

// SetSliceElement sets the element group at one axis above the leaf elements.
// Source Array is returned, for function-chaining design.
func (a nDimFields) SetSliceElement(vals []nDimElement, index ...int) *nDimFields {
	idx := a.valIdx(index, "SetSliceElement")
	switch {
	case a.HasErr():
		return &a
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
		return &a
	}

	copy(a.data[idx:idx+a.strides[len(a.strides)-2]], vals[:a.strides[len(a.strides)-2]])
	return &a
}

// SetSubArr sets the array below a given index to the values in vals.
// Values will be broadcast up multiple axes if the shapes match.
func (a *nDimFields) SetSubArr(vals nDimObject, index ...int) *nDimFields {
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
	case len(vals.Shape())+len(index) > len(a.shape):
		a.err = InvIndexError
		if debug {
			a.debug = fmt.Sprintf("Array received by SetSubArr() cant be broadcast.  Shape: %v  Vals shape: %v index: %v", a.fields().shape, vals.fields().shape, index)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	for i, j := len(a.fields().shape)-1, len(vals.fields().shape)-1; j >= 0; i, j = i-1, j-1 {
		if a.fields().shape[i] != vals.fields().shape[j] {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Shape of array recieved by SetSubArr() doesn't match receiver.  Shape: %v  Vals Shape: %v", a.fields().shape, vals.fields().shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
	}

	if len(a.fields().shape)-len(index)-len(vals.fields().shape) == 0 {
		copy(a.fields().data[idx:idx+len(vals.fields().data)], vals.fields().data)
		return a
	}

	reps := 1
	for i := len(index); i < len(a.fields().shape)-len(vals.fields().shape); i++ {
		reps *= a.fields().shape[i]
	}

	ln := len(vals.fields().data)
	for i := 1; i <= reps; i++ {
		copy(a.data[idx+ln*(i-1):idx+ln*i], vals.fields().data)
	}
	return a
}

// Resize will change the underlying array size.
//
// Make a copy C() if the original array needs to remain unchanged.
// Element location in the underlying slice will not be adjusted to the new shape.
func (a *nDimFields) Resize(shape ...int) *nDimFields {
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
func (a *nDimFields) Append(val nDimObject, axis int) *nDimFields {
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
	case len(a.fields().shape) != len(val.fields().shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.fields().shape, val.fields().shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return a
	}

	for k, v := range a.shape {
		if v != val.fields().shape[k] && k != axis {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Array received by Append() can not be matched.  Shape: %v  Val shape: %v", a.fields().shape, val.fields().shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a
		}
	}

	ln := len(a.fields().data) + len(val.fields().data)
	var dat []nDimElement
	cp := cap(a.data)
	if ln > cp {
		dat = make([]nDimElement, ln)
	} else {
		dat = a.data[:ln]
	}

	as, vs := a.fields().strides[axis], val.fields().strides[axis]
	for i, j := a.fields().strides[0], val.fields().strides[0]; i > 0; i, j = i-as, j-vs {
		copy(dat[i+j-vs:i+j], val.fields().data[j-vs:j])
		copy(dat[i+j-as-vs:i+j-vs], a.fields().data[i-as:i])
	}

	a.data = dat
	a.shape[axis] += val.fields().shape[axis]

	for i := axis; i >= 0; i-- {
		a.strides[i] = a.strides[i+1] * a.shape[i]
	}

	return a
}

// String Satisfies the Stringer interface for fmt package
func (a *nDimFields) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.err != nil:
		return "Error: " + a.err.(*ngError).s
	case a.data == nil || a.shape == nil || a.strides == nil:
		return "<nil>"
	case a.strides[0] == 0:
		return "[]"
	case len(a.shape) == 1:
		return fmt.Sprint(a.data)
	}

	stride := a.shape[len(a.shape)-1]

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

func (a *nDimFields) Array64() *Array64 {
	return &Array64{*a}

}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a nDimFields) Reshape(shape ...int) nDimObject {
	if a.HasErr() || len(shape) == 0 {
		return a
	}

	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative dimension received by Reshape(): %v", shape)
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

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *nDimFields) MarshalJSON() ([]byte, error) {
	t := a.c()

	inf, nan, err := t.encode()
	return json.Marshal(struct {
		Shape []int         `json:"shape"`
		Data  []nDimElement `json:"data"`
		Inf   []int64       `json:"inf,omitempty"`
		Nan   []int64       `json:"nan,omitempty"`
		Err   int8          `json:"err,omitempty"`
	}{
		Shape: t.shape,
		Data:  t.data,
		Inf:   inf,
		Nan:   nan,
		Err:   err,
	})
}

// encode is used to prepare data for MarshalJSON that isn't JSON defined.
func (a *nDimFields) encode() (inf, nan []int64, err int8) {
	for k, v := range a.data {
		switch {
		case math.IsNaN(v.(float64)):
			a.data[k] = 0
			nan = append(nan, int64(k+1))
		case math.IsInf(v.(float64), 1):
			a.data[k] = 0
			inf = append(inf, int64(k+1))
		case math.IsInf(v.(float64), -1):
			a.data[k] = 0
			inf = append(inf, int64(-(k + 1)))
		}
	}

	err = encodeErr(a.err)
	return
}

// decode is used to build Array from UnmarshalJSON for values that aren't JSON defined.
func (a *nDimFields) decode(i, n []int64, err int8) {
	inf, nInf := math.Inf(1), math.Inf(-1)
	nan := math.NaN()

	for _, v := range n {
		a.data[v-1] = nan
	}

	for _, v := range i {
		if v-1 >= 0 {
			a.data[v-1] = inf
		} else {
			a.data[-v-1] = nInf
		}
	}
	a.err = decodeErr(err)
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for decoding data.
// Custom Unmarshaler is needed to load/decode unexported values and build strides.
func (a *nDimFields) UnmarshalJSON(b []byte) error {
	tmpA := new(struct {
		Shape []int         `json:"shape"`
		Data  []nDimElement `json:"data"`
		Inf   []int64       `json:"inf,omitempty"`
		Nan   []int64       `json:"nan,omitempty"`
		Err   int8          `json:"err,omitempty"`
	})

	err := json.Unmarshal(b, tmpA)

	a.shape = tmpA.Shape
	a.data = tmpA.Data
	a.decode(tmpA.Inf, tmpA.Nan, tmpA.Err)

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
