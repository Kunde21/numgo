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
func (a *Arrayb) Reshape(shape ...int) *Arrayb {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	if sz != uint64(len(a.data)) {
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
	if a == nil {
		return nil
	}

	b = createb(a.shape...)
	for i, v := range a.data {
		b.data[i] = v
	}
	return
}
