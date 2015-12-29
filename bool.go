package numgo

import (
	"encoding/json"
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
	debug   string
}

// Create creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Createb(shape ...int) (a *Arrayb) {
	a = new(Arrayb)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			a.err = NegativeAxis
			return
		}
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
	return
}

// Internal function to create using the shape of another array
func createb(shape ...uint64) (a *Arrayb) {
	a = new(Arrayb)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			a.err = NegativeAxis
			return
		}
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
	return
}

// Full creates an Arrayf object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the retuen
func Fullb(val bool, shape ...int) (a *Arrayb) {
	a = Createb(shape...)
	if a.err != nil {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = val
	}
	return
}

func fullb(val bool, shape ...uint64) (a *Arrayb) {
	a = createb(shape...)
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
		a = new(Arrayb)
		a.err = NilError
		return ""
	case a.err != nil:
		return a.err.s
	case a.strides[0] == 0:
		return "[]"
	}

	a.RLock()
	defer a.RUnlock()

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
	switch {
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	}

	a.Lock()
	defer a.Unlock()

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			return a
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	if sz != uint64(len(a.data)) {
		a.err = ReshapeError
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
	switch {
	case a == nil:
		b = new(Arrayb)
		b.err = NilError
		return b
	case a.err != nil:
		return a
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
	switch {
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return nil
	case len(a.shape) != len(index):
		a.err = ShapeError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
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
	switch {
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return nil
	case len(a.shape)-1 != len(index):
		a.err = ShapeError
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
	switch {
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return nil
	case len(a.shape) < len(index):
		a.err = ShapeError
		return nil
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {

			return
		}
		idx += uint64(v) * a.strides[i+1]
	}

	ret = createb(a.shape[len(index):]...)
	copy(ret.data, a.data[idx:idx+a.strides[len(index)]])

	return
}

// SetE sets the element at the given index.
// There should be one index per axis.  Generates a ShapeError if incorrect index.
func (a *Arrayb) SetE(val bool, index ...int) *Arrayb {
	switch {
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		return a
	case a.err != nil:
		return a
	case len(a.shape) != len(index):
		a.err = ShapeError
		return a
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
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
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		return a
	case a.err != nil:
		return a
	case len(a.shape)-1 != len(index) || uint64(len(vals)) != a.shape[len(a.shape)-1]:
		a.err = ShapeError
		return a
	}
	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] || v < 0 {
			a.err = IndexError
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
	case a == nil:
		a = new(Arrayb)
		a.err = NilError
		return a
	case a.err != nil:
		return a
	case len(vals.shape)+len(index) > len(a.shape):
		a.err = ShapeError
		return a
	}

	for i, j := len(a.shape)-1, len(vals.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[i] != vals.shape[j] {
			a.err = ShapeError
			return a
		}
	}

	idx := uint64(0)
	for i, v := range index {
		if uint64(v) > a.shape[i] {
			a.err = IndexError
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
		a = new(Arrayb)
		a.err = NilError
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
	a.err.decodeErr(tmpA.Err)

	a.strides = make([]uint64, len(a.shape)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= a.shape[i-1]
	}
	a.strides[0] = tmp

	return nil
}
