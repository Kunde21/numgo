package numgo

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

type Arrayf struct {
	sync.RWMutex
	shape   []uint64
	strides []uint64
	data    []float64
	err     *ngError
	debug   []byte
}

// Create creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Create(shape ...int) (a *Arrayf) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayf)
	a.shape = sh
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	return
}

// Internal function to create using the shape of another array
func create(shape ...uint64) (a *Arrayf) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayf)
	a.shape = sh
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	return
}

// Full creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All elements will be set to the value passed in val.
func Full(val float64, shape ...int) (a *Arrayf) {
	a = Create(shape...)
	if a == nil {
		return nil
	}
	a.AddC(val)
	return
}

func full(val float64, shape ...uint64) (a *Arrayf) {
	a = create(shape...)
	if a == nil {
		return nil
	}
	a.AddC(val)
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrayf) String() (s string) {

	a.RLock()
	defer a.RUnlock()

	if a.err != nil {
		return a.err.s
	}
	if a.strides[0] == 0 {
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

// Arange Creates an array in one of three different ways, depending on input:
// One (stop):         Arrayf from zero to positive value or negative value to zero
// Two (start, stop):   Arrayf from start to stop, with increment of 1 or -1, depending on inputs
// Three (start, stop, step): Arrayf from start to stop, with increment of step
//
// Any inputs beyond three values are ignored
func Arange(vals ...float64) (a *Arrayf) {
	var start, stop, step float64 = 0, 0, 1

	switch len(vals) {
	case 0:
		return nil
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
			return nil
		}
		start, stop, step = vals[0], vals[1], vals[2]
	}

	a = Create(int((stop - start) / step))
	for i, v := 0, start; i < len(a.data); i, v = i+1, v+step {
		a.data[i] = v
	}
	return
}

// Identity creates a size x size matrix with 1's on the main diagonal.
// All other values will be zero.
//
// Negative size values will generate an error and return a nil value.
func Identity(size int) (r *Arrayf) {
	if size < 0 {
		return nil
	}

	r = Create(size, size)
	for i := uint64(0); i < r.strides[0]; i = +r.strides[1] + r.strides[2] {
		r.data[i] = 1
	}
	return
}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a *Arrayf) Reshape(shape ...int) *Arrayf {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}

	a.Lock()
	defer a.Unlock()

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	if sz != uint64(len(a.data)) {
		a.err = ReshapeError
		return nil
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

// Flatten reshapes the data to a 1-D array.
func (a *Arrayf) Flatten() *Arrayf {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}
	a.shape[0] = a.strides[0]
	a.shape = a.shape[:1]
	fmt.Println(a.shape)
	return a.Reshape(int(a.strides[0]))
}

// C will return a deep copy of the source array.
func (a *Arrayf) C() (b *Arrayf) {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}

	b = create(a.shape...)
	for i, v := range a.data {
		b.data[i] = v
	}
	return
}

// E returns the element at the given index.
func (a *Arrayf) E(index ...int) float64 {
	if a.err != nil {
		return math.NaN()
	}
	if a == nil {
		a.err = NilError
		return math.NaN()
	}
	if len(a.shape) != len(index) {
		a.err = ShapeError
		return math.NaN()
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			return math.NaN()
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return a.data[idx]
}

// Eslice returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Arrayf) SliceElement(index ...int) (ret []float64) {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}
	if len(a.shape)-1 != len(index) {
		a.err = IndexError
		return nil
	}
	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			a.err = IndexError
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}
	return append(ret, a.data[idx:idx+a.strides[len(a.strides)-2]]...)
}

// SubArr slices the array at a given index.
func (a *Arrayf) SubArr(index ...int) (ret *Arrayf) {
	if len(index) > len(a.shape) {
		a.err = ShapeError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			a.err = IndexError
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}

	ret = create(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}
