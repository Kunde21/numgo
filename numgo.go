package numgo

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

type Array struct {
	sync.RWMutex
	shape   []uint64
	strides []uint64
	data    []float64
}

// Create creates an Array object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Create(shape ...int) (a *Array) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Array)
	a.shape = sh
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(sh)-1)
	tmp := sh[len(sh)-1]

	for i := len(a.strides) - 1; i >= 0; i-- {
		tmp *= sh[i]
		a.strides[i] = tmp
	}
	return
}

// Full creates an Array object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the retuen
func Full(val float64, shape ...int) (a *Array) {
	a = Create(shape...)
	if a == nil {
		return nil
	}
	a.Addf(val)
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Array) String() string {

	a.RLock()
	defer a.RUnlock()

	stride := a.shape[len(a.shape)-1]

	var s string
	for i, k := uint64(0), 0; i+stride <= uint64(len(a.data)); i, k = i+stride, k+1 {
		t := ""
		for _, v := range a.strides {
			if i%v == 0 {
				t += "["
			}
		}
		s += strings.Repeat(" ", len(a.shape)-len(t)-1) + t
		s += fmt.Sprint(a.data[i : i+stride])
		t = ""
		for _, v := range a.strides {
			if (i+stride)%v == 0 {
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
	return s
}

// Arange Creates an array in one of three different ways, depending on input:
// One (stop):         Array from zero to positive value or negative value to zero
// Two (start,stop):   Array from start to stop, with increment of 1 or -1, depending on inputs
// Three (start, stop, step): Array from start to stop, with increment of step
//
// Any inputs beyond three values are ignored
func Arange(vals ...float64) (a *Array) {
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
	for i, v := 0, start; i < len(a.data) && v < stop; i, v = i+1, v+step {
		a.data[i] = v
	}
	return
}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a *Array) Reshape(shape ...int) *Array {
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

	a.strides = make([]uint64, len(sh)-1)
	tmp := sh[len(sh)-1]
	for i := len(a.strides) - 1; i >= 0; i-- {
		tmp *= sh[i]
		a.strides[i] = tmp
	}
	a.shape = sh

	return a
}

// Add will perform addition on two
func (a *Array) Add(b *Array) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	if len(a.shape) < len(b.shape) {
		fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
		return nil
	}
	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
			return nil
		}
	}
	compChan := make(chan struct{})
	mul := len(a.data) / len(b.data)
	for k := 0; k < mul; k++ {
		go func(m int) {
			for i, v := range b.data {
				a.data[i+m] += v
			}
			compChan <- struct{}{}
		}(k * len(b.data))
	}
	for k := 0; k < mul; k++ {
		<-compChan
	}
	return a
}

func (a *Array) Addf(b float64) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] += b
	}
	return a
}

func (a *Array) Subtr(b *Array) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	if len(a.shape) < len(b.shape) {
		fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
		return nil
	}
	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
			return nil
		}
	}
	compChan := make(chan struct{})
	mul := len(a.data) / len(b.data)
	for k := 0; k < mul; k++ {
		go func(m int) {
			for i, v := range b.data {
				a.data[i+m] -= v
			}
			compChan <- struct{}{}
		}(k * len(b.data))
	}
	for k := 0; k < mul; k++ {
		<-compChan
	}
	return a
}

func (a *Array) Subtrf(b float64) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] -= b
	}
	return a
}

func (a *Array) Mult(b *Array) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	if len(a.shape) < len(b.shape) {
		fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
		return nil
	}
	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
			return nil
		}
	}
	compChan := make(chan struct{})
	mul := len(a.data) / len(b.data)
	for k := 0; k < mul; k++ {
		go func(m int) {
			for i, v := range b.data {
				a.data[i+m] *= v
			}
			compChan <- struct{}{}
		}(k * len(b.data))
	}
	for k := 0; k < mul; k++ {
		<-compChan
	}
	return a
}

func (a *Array) Multf(b float64) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] *= b
	}
	return a
}

// Divides a by b
func (a *Array) Div(b *Array) *Array {
	if a == nil {
		return nil
	}
	if len(a.shape) < len(b.shape) {
		fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
		return nil
	}
	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			fmt.Println("Base array must have as many or more dimensions", a.shape, "<", b.shape)
			return nil
		}
	}
	compChan := make(chan bool)
	mul := len(a.data) / len(b.data)
	for k := 0; k < mul; k++ {
		go func(m int) {
			flg := false
			for i, v := range b.data {
				if v == 0 {
					flg = true
					a.data[i+m] = math.NaN()
				} else {
					a.data[i+m] /= v
				}
			}
			compChan <- flg
		}(k * len(b.data))
	}
	flg := false
	for k := 0; k < mul; k++ {
		flg = flg || <-compChan
	}
	if flg {
		fmt.Println("Division by zero encountered.")
	}
	return a
}

func (a *Array) Divc(b float64) *Array {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	flg := false
	for i := 0; i < len(a.data); i++ {
		if b == 0 {
			flg = true
			a.data[i] = math.NaN()
		} else {
			a.data[i] /= b
		}

	}
	if flg {
		fmt.Println("Division by zero encountered.")
	}
	return a
}
