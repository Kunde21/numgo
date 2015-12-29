package numgo

import (
	"encoding/json"
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
	debug   string
}

// Create creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Create(shape ...int) (a *Arrayf) {
	a = new(Arrayf)
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v < 0 {
			a.err = NegativeAxis
			return
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

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
			a = new(Arrayf)
			a.err = NegativeAxis
			return
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

// FromSlice wraps a float64 slice as a 1-D Arrayf object.
func FromSlice(data []float64) (a *Arrayf) {
	if data == nil {
		a = new(Arrayf)
		a.err = NilError
		return nil
	}

	a = new(Arrayf)
	a.shape = []uint64{uint64(len(data))}
	a.strides = []uint64{a.shape[0], 1}
	a.data = make([]float64, len(data))
	copy(a.data, data)
	return a
}

// Full creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All elements will be set to the value passed in val.
func Full(val float64, shape ...int) (a *Arrayf) {
	a = Create(shape...)
	if a.err != nil {
		return
	}
	a.AddC(val)
	return
}

func full(val float64, shape ...uint64) (a *Arrayf) {
	a = create(shape...)
	if a.err != nil {
		return
	}
	a.AddC(val)
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
		return Create(0)
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
			a = new(Arrayf)
			a.err = ShapeError
			return
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
		r = new(Arrayf)
		r.err = ShapeError
		return
	}

	r = Create(size, size)
	for i := uint64(0); i < r.strides[0]; i = +r.strides[1] + r.strides[2] {
		r.data[i] = 1
	}
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrayf) String() (s string) {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		return "<nil>"
	case a.err != nil:
		return a.err.s
	case a.strides[0] == 0:
		return "[]"
	}

	a.RLock()
	defer a.RUnlock()

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
func (a *Arrayf) Reshape(shape ...int) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		return a
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

// encode is used to prepare data for MarshalJSON that isn't JSON defined.
func (a *Arrayf) encode() (inf, nan []int64, err int8) {
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
func (a *Arrayf) MarshalJSON() ([]byte, error) {
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
func (a *Arrayf) decode(i, n []int64, err int8) {
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
	a.err.decodeErr(err)
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for decoding data.
// Custom Unmarshaler is needed to load/decode unexported values and build strides.
func (a *Arrayf) UnmarshalJSON(b []byte) error {

	if a == nil {
		a = new(Arrayf)
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
