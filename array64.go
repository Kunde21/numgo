package numgo

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strings"
)

type nDimMetadata struct {
	shape        []int
	strides      []int
	err          error
	debug, stack string
}

// Array64 is an n-dimensional array of float64 data
type Array64 struct {
	nDimMetadata
	data []float64
}

// NewArray64 creates an Array64 object with dimensions given in order from outer-most to inner-most
// Passing a slice with no shape data will wrap the slice as a 1-D array.
// All values will default to zero.  Passing nil as the data parameter creates an empty array.
func NewArray64(data []float64, shape ...int) (a *Array64) {
	if len(shape) == 0 {
		switch {
		case data != nil:
			return &Array64{
				nDimMetadata{
					shape:   []int{len(data)},
					strides: []int{len(data), 1},
					err:     nil,
					debug:   "",
					stack:   "",
				},

				data,
			}
		default:
			return &Array64{
				nDimMetadata{
					shape:   []int{0},
					strides: []int{0, 0},
					err:     nil,
					debug:   "",
					stack:   "",
				},
				[]float64{},
			}
		}
	}

	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		if v < 0 {
			a = &Array64{nDimMetadata{err: NegativeAxis}, nil}
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Create: %v", shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		sz *= v
	}
	copy(sh, shape)

	a = &Array64{

		nDimMetadata{
			shape:   sh,
			strides: make([]int, len(shape)+1),
			err:     nil,
			debug:   "",
			stack:   "",
		},
		make([]float64, sz),
	}

	if data != nil {
		copy(a.data, data)
	}

	a.strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		a.strides[i] = a.strides[i+1] * a.shape[i]
	}
	return
}

// Internal function to create using the shape of another array
func newArray64(shape ...int) (a *Array64) {
	var sz = 1
	for _, v := range shape {
		sz *= v
	}

	a = &Array64{
		nDimMetadata{
			shape:   shape,
			strides: make([]int, len(shape)+1),
			err:     nil,
			debug:   "",
			stack:   "",
		},

		make([]float64, sz),
	}

	a.strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		a.strides[i] = a.strides[i+1] * a.shape[i]
	}
	return
}

// FullArray64 creates an Array64 object with dimensions given in order from outer-most to inner-most
// All elements will be set to the value passed in val.
func FullArray64(val float64, shape ...int) (a *Array64) {
	a = NewArray64(nil, shape...)
	if a.HasErr() || val == 0 {
		return
	}

	return a.AddC(val)
}

func full(val float64, shape ...int) (a *Array64) {
	a = newArray64(shape...)
	if val == 0 {
		return
	}

	return a.AddC(val)
}

// RandArray64 creates an Arry64 object and fills it with random values from the default random source
// Use base and scale to adjust the default value range [0.0, 1.0)
// generated by formula rand * scale + default
func RandArray64(base, scale float64, shape ...int) (a *Array64) {
	a = NewArray64(nil, shape...)
	if a.HasErr() {
		return
	}

	for i := range a.data {
		a.data[i] = rand.Float64()*scale + base
	}
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
			start = vals[0]
		} else {
			stop = vals[0] - 1
		}
	case 2:
		if vals[1] < vals[0] {
			step = -1
		}
		start, stop = vals[0], vals[1]
	default:
		if vals[1] < vals[0] && vals[2] >= 0 || vals[1] > vals[0] && vals[2] <= 0 {
			a = &Array64{nDimMetadata{err: ShapeError}, nil}
			if debug {
				a.debug = fmt.Sprintf("Arange received illegal values %v", vals)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return a

		}
		start, stop, step = vals[0], vals[1], vals[2]
	}

	a = NewArray64(nil, int((stop-start)/(step))+1)
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
		r = &Array64{nDimMetadata{err: NegativeAxis}, nil}
		if debug {
			r.debug = fmt.Sprintf("Negative dimension received by Identity: %d", size)
			r.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return
	}

	r = NewArray64(nil, size, size)
	for i := 0; i < r.strides[0]; i += r.strides[1] + r.strides[2] {
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

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a *Array64) Reshape(shape ...int) *Array64 {
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

// encode is used to prepare data for MarshalJSON that isn't JSON defined.
func (a *Array64) encode() (inf, nan []int64, err int8) {
	for k, v := range a.data {
		switch {
		case math.IsNaN(float64(v)):
			a.data[k] = 0
			nan = append(nan, int64(k+1))
		case math.IsInf(float64(v), 1):
			a.data[k] = 0
			inf = append(inf, int64(k+1))
		case math.IsInf(float64(v), -1):
			a.data[k] = 0
			inf = append(inf, int64(-(k + 1)))
		}
	}

	err = encodeErr(a.err)
	return
}

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *Array64) MarshalJSON() ([]byte, error) {
	t := a.C()

	inf, nan, err := t.encode()
	return json.Marshal(struct {
		Shape []int     `json:"shape"`
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
func (a *Array64) UnmarshalJSON(b []byte) error {
	tmpA := new(struct {
		Shape []int     `json:"shape"`
		Data  []float64 `json:"data"`
		Inf   []int64   `json:"inf,omitempty"`
		Nan   []int64   `json:"nan,omitempty"`
		Err   int8      `json:"err,omitempty"`
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
