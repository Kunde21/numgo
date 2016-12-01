package numgo

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// Arrayb is an n-dimensional array of boolean values
type Arrayb struct {
	nDimFields
}

// NewArrayB creates an Arrayb object with dimensions given in order from outer-most to inner-most
// All values will default to false
func NewArrayB(data []nDimElement, shape ...int) (a *Arrayb) {
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
	a.data = make([]nDimElement, sz)
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

func (a Arrayb) values() *nDimFields {
	return &a.nDimFields
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
	a.data = make([]nDimElement, sz)

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

// MarshalJSON fulfills the json.Marshaler Interface for encoding data.
// Custom Unmarshaler is needed to encode/send unexported values.
func (a *Arrayb) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Shape []int         `json:"shape"`
		Data  []nDimElement `json:"data"`
		Err   int8          `json:"err,omitempty"`
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
		Shape []int         `json:"shape"`
		Data  []nDimElement `json:"data"`
		Err   int8          `json:"err,omitempty"`
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
