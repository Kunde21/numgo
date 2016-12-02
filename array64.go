package numgo

import (
	"fmt"
	"math/rand"
	"runtime"
)

// Array64 is an n-dimensional array of float64 data
type Array64 struct {
	nDimFields
}

// C will return a deep copy of the source array.
func (a Array64) C() nDimObject {
	b := Array64{a.nDimFields}
	return b
}

// NewArray64 creates an Array64 object with dimensions given in order from outer-most to inner-most
// Passing a slice with no shape data will wrap the slice as a 1-D array.
// All values will default to zero.  Passing nil as the data parameter creates an empty array.
func NewArray64(data []nDimElement, shape ...int) (a *Array64) {
	if len(shape) == 0 {
		switch {
		case data != nil:
			return &Array64{nDimFields{
				shape:   []int{len(data)},
				strides: []int{len(data), 1},
				data:    data,
				err:     nil,
				debug:   "",
				stack:   "",
			}}
		default:
			return &Array64{nDimFields{
				shape:   []int{0},
				strides: []int{0, 0},
				data:    []nDimElement{},
				err:     nil,
				debug:   "",
				stack:   "",
			}}
		}
	}

	var sz = 1
	sh := make([]int, len(shape))
	for _, v := range shape {
		if v < 0 {
			a = &Array64{nDimFields{err: NegativeAxis}}
			if debug {
				a.debug = fmt.Sprintf("Negative axis length received by Create: %v", shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		sz *= v
	}
	copy(sh, shape)

	a = &Array64{nDimFields{
		shape:   sh,
		strides: make([]int, len(shape)+1),
		data:    make([]nDimElement, sz),
		err:     nil,
		debug:   "",
		stack:   "",
	}}

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
func newnDimFields(shape ...int) (a *nDimFields) {
	var sz = 1
	for _, v := range shape {
		sz *= v
	}

	a = &nDimFields{
		shape:   shape,
		strides: make([]int, len(shape)+1),
		data:    make([]nDimElement, sz),
		err:     nil,
		debug:   "",
		stack:   "",
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

	a = &Array64{nDimFields{
		shape:   shape,
		strides: make([]int, len(shape)+1),
		data:    make([]nDimElement, sz),
		err:     nil,
		debug:   "",
		stack:   "",
	}}

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

//Reshape the array
func (a Array64) Reshape(shape ...int) nDimObject {
	b := Array64{*a.nDimFields.Reshape(shape...)}
	return b
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
			a = &Array64{nDimFields{err: ShapeError}}
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
		r = &Array64{nDimFields{err: NegativeAxis}}
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