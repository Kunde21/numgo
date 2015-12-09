package numgo

import (
	"fmt"
	"strings"
	"sync"
)

type Arrayb struct {
	sync.RWMutex
	shape   []uint64
	strides []uint64
	data    []bool
	err     *ngError
	debug   []byte
}

// Create creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Createb(shape ...int) (a *Arrayb) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayb)
	a.shape = sh
	a.data = make([]bool, sz)

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
func createb(shape ...uint64) (a *Arrayb) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayb)
	a.shape = sh
	a.data = make([]bool, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	return
}

// Full creates an Arrayf object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the retuen
func Fullb(val bool, shape ...int) (a *Arrayb) {
	a = Createb(shape...)
	if a == nil {
		return nil
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

func fullb(val bool, shape ...uint64) (a *Arrayb) {
	a = createb(shape...)
	if a == nil {
		return nil
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrayb) String() (s string) {
	a.RLock()
	defer a.RUnlock()

	if a.err != nil {
		return a.err.s
	}
	if a.strides[0] == 0 {
		return "[]"
	}

	if a == nil {
		a.err = NilError
		return ""
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

// C will return a deep copy of the source array.
func (a *Arrayb) C() (b *Arrayb) {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}

	b = createb(a.shape...)
	for i, v := range a.data {
		b.data[i] = v
	}
	return
}

// E returns a pointer to a copy of the element at the given index.
// Any errors will return a nil value and record the error for the
// HasErr() and GetErr() functions.
func (a *Arrayb) E(index ...int) *bool {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}
	if len(a.shape) != len(index) {
		a.err = ShapeError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}
	b := a.data[idx]
	return &b
}

// Eslice returns the element group at one axis above the leaf elements.
// Data is returned as a copy  in a float slice.
func (a *Arrayb) SliceElement(index ...int) (ret []bool) {
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			a.err = ShapeError
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
	if a.err != nil {
		return nil
	}
	if a == nil {
		a.err = NilError
		return nil
	}

	if len(index) > len(a.shape) {
		a.err = ShapeError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			return nil
		}
		idx += uint64(v) * a.strides[i+1]
	}

	ret = createb(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}
