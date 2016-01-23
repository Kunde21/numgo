package numgo

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Array64 struct {
	shape   []uint64
	strides []uint64
	data    []float64
	err     error
	debug   string
}

// NewArray64 creates an Array64 object with dimensions given in order from outer-most to inner-most
// Passing a slice with no shape data will wrap the slice as a 1-D array.
// All values will default to zero.  Passing nil as the data parameter creates an empty array.
func NewArray64(data []float64, shape ...int) (a *Array64) {
	if data != nil && len(shape) == 0 {
		shape = append(shape, len(data))
	}

	a = new(Array64)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Create: %v", shape)
			}
			return
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a.shape = sh
	a.data = make([]float64, sz)
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
func newArray64(shape ...uint64) (a *Array64) {
	var sz uint64 = 1
	for _, v := range shape {
		sz *= uint64(v)
	}

	a = new(Array64)
	a.shape = shape
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(shape)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= shape[i-1]
	}
	//if sz == 0{

	a.strides[0] = tmp
	a.err = nil
	return
}

// Full creates an Array64 object with dimensions given in order from outer-most to inner-most
// All elements will be set to the value passed in val.
func Full(val float64, shape ...int) (a *Array64) {
	a = NewArray64(nil, shape...)
	if a.err != nil {
		return
	}
	a.AddC(val)
	return
}

func full(val float64, shape ...uint64) (a *Array64) {
	a = newArray64(shape...)
	if a.err != nil {
		return
	}
	a.AddC(val)
	return
}

// Arange Creates an array in one of three different ways, depending on input:
//  Arange(stop):              Array64 from zero to positive value or negative value to zero
//  Arange(start, stop):       Array64 from start to stop, with increment of 1 or -1, depending on inputs
//  Arange(start, stop, step): Array64 from start to stop, with increment of step
//
// Any inputs beyond three values are ignored
func Arange(vals ...float64) (a *Array64) {
	var start, stop, step float64 = 0, 0, 1

	switch len(vals) {
	case 0:
		return newArray64(0)
	case 1:
		if vals[0] <= 0 {
			start, stop, step = vals[0], 0, -1
		} else {
			stop = vals[0]
		}
	case 2:
		if vals[1] < vals[0] {
			step = -1
		}
		start, stop = vals[0], vals[1]
	default:
		if vals[1] < vals[0] && vals[2] >= 0 || vals[1] > vals[0] && vals[2] <= 0 {
			a = new(Array64)
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Arange received illegal values %v", vals)
			}
			return
		}
		start, stop, step = vals[0], vals[1], vals[2]
	}

	a = NewArray64(nil, int((stop-start)/step))
	for i, v := 0, start; i < len(a.data); i, v = i+1, v+step {
		a.data[i] = v
	}
	return
}

// Identity creates a size x size matrix with 1's on the main diagonal.
// All other values will be zero.
//
// Negative size values will generate an error and return a nil value.
func Identity(size int) (r *Array64) {
	if size < 0 {
		r = NewArray64(nil, 0)
		r.err = NegativeAxis
		if debug {
			r.debug = fmt.Sprintf("Negative dimension received by Identity: %d", size)
		}
		return
	}

	r = NewArray64(nil, size, size)
	for i := uint64(0); i < r.strides[0]; i = +r.strides[1] + r.strides[2] {
		r.data[i] = 1
	}
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Array64) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.err != nil:
		return "Error: " + a.err.(*ngError).s
	case a.strides[0] == 0:
		return "[]"
	}

	stride := a.shape[len(a.shape)-1]

	for i, k := uint64(0), 0; i+stride <= uint64(len(a.data)); i, k = i+stride, k+1 {

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
func (a *Array64) Reshape(shape ...int) *Array64 {
	if a == nil || a.err != nil {
		return a
	}

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			if debug {
				a.debug = fmt.Sprintf("Negative dimension received by Reshape(): %v", shape)
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

// encode is used to prepare data for MarshalJSON that isn't JSON defined.
func (a *Array64) encode() (inf, nan []int64, err int8) {
	for k, v := range a.data {
		a.data[k] = 0
		switch {
		case math.IsNaN(float64(v)):
			nan = append(nan, int64(k))
		case math.IsInf(float64(v), 1):
			inf = append(inf, int64(k))
		case math.IsInf(float64(v), -1):
			inf = append(inf, int64(-k))
		}
	}

	err = encodeErr(a.err)
	return
}

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *Array64) MarshalJSON() ([]byte, error) {
	if a == nil {
		return nil, NilError
	}
	t := a.C()
	inf, nan, err := t.encode()
	return json.Marshal(struct {
		Shape []uint64  `json:"shape"`
		Data  []float64 `json:"data"`
		Inf   []int64   `json:"inf,omitempty"`
		Nan   []int64   `json:"nan,omitempty"`
		Err   int8      `json:"err,omitempty"`
	}{
		Shape: t.shape,
		Data:  t.data,
		Inf:   inf,
		Nan:   nan,
		Err:   err,
	})
}

// decode is used to build Array from UnmarshalJSON for values that aren't JSON defined.
func (a *Array64) decode(i, n []int64, err int8) {
	inf, nInf := math.Inf(1), math.Inf(-1)
	nan := math.NaN()

	for _, v := range n {
		a.data[v] = nan
	}

	for _, v := range i {
		if v >= 0 {
			a.data[v] = inf
		} else {
			a.data[-v] = nInf
		}
	}
	a.err = decodeErr(err)
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for decoding data.
// Custom Unmarshaler is needed to load/decode unexported values and build strides.
func (a *Array64) UnmarshalJSON(b []byte) error {

	if a == nil {
		return NilError
	}

	tmpA := new(struct {
		Shape []uint64  `json:"shape"`
		Data  []float64 `json:"data"`
		Inf   []int64   `json:"inf,omitempty"`
		Nan   []int64   `json:"nan,omitempty"`
		Err   int8      `json:"err,omitempty"`
	})

	err := json.Unmarshal(b, tmpA)
	if err != nil {
		return err
	}

	a.shape = tmpA.Shape
	a.data = tmpA.Data
	a.decode(tmpA.Inf, tmpA.Nan, tmpA.Err)

	a.strides = make([]uint64, len(a.shape)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= a.shape[i-1]
	}
	a.strides[0] = tmp

	return nil
}
