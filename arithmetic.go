package numgo

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

var nan float64

func init() {
	nan = math.NaN()
}

// Add performs element-wise addition
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Add(b nDimObject) *Array64 {
	if a.valRith(b, "Add") {
		return a
	}

	if b.fields().shape[len(b.fields().shape)-1] == a.shape[len(a.shape)-1] {
		Add(a.data, b.fields().data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.fields().data); i++ {
		AddC(b.fields().data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// AddC adds a constant to all elements of the array.
func (a *Array64) AddC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	AddC(b, a.data)
	return a
}

func (a Array64) values() *nDimFields {
	return &a.nDimFields
}

// Subtr performs element-wise subtraction.
// Arrays must be the same size or albe to broadcast.
// This will modify the source array.
func (a *Array64) Subtr(b nDimObject) *Array64 {
	if a.valRith(b, "Subtr") {
		return a
	}

	if b.fields().shape[len(b.fields().shape)-1] == a.shape[len(a.shape)-1] {
		Subtr(a.data, b.fields().data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.fields().data); i++ {
		SubtrC(b.fields().data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// SubtrC subtracts a constant from all elements of the array.
func (a *Array64) SubtrC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	SubtrC(b, a.data)
	return a
}

// Mult performs element-wise multiplication.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Mult(b nDimObject) *Array64 {
	if a.valRith(b, "Mult") {
		return a
	}

	if b.fields().shape[len(b.fields().shape)-1] == a.shape[len(a.shape)-1] {
		Mult(a.data, b.fields().data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.fields().data); i++ {
		MultC(b.fields().data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// MultC multiplies all elements of the array by a constant.
func (a *Array64) MultC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	MultC(b, a.data)
	return a
}

// Div performs element-wise division
// Arrays must be the same size or able to broadcast.
// Division by zero conforms to IEEE 754
// 0/0 = NaN, +x/0 = +Inf, -x/0 = -Inf
// This will modify the source array.
func (a *Array64) Div(b nDimObject) *Array64 {
	if a.valRith(b, "Div") {
		return a
	}

	if b.fields().shape[len(b.fields().shape)-1] == a.shape[len(a.shape)-1] {
		Div(a.data, b.fields().data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.fields().data); i++ {
		DivC(b.fields().data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// DivC divides all elements of the array by a constant.
// Division by zero conforms to IEEE 754
// 0/0 = NaN, +x/0 = +Inf, -x/0 = -Inf
func (a *Array64) DivC(b float64) *Array64 {
	switch {
	case a.HasErr():
		return a
	}

	DivC(b, a.data)
	return a
}

// Pow raises elements of a to the corresponding power in b.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Pow(b nDimObject) *Array64 {
	if a.valRith(b, "Pow") {
		return a
	}

	if b.fields().shape[len(b.fields().shape)-1] == a.shape[len(a.shape)-1] {
		lna, lnb := len(a.data), len(b.fields().data)
		for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
			if j >= lnb {
				j = 0
			}
			a.data[i] = math.Pow(a.data[i].(float64), b.fields().data[j].(float64))
		}
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.fields().data); i++ {
		for j := i * st; j < (i+1)*st; j++ {
			a.data[j] = math.Pow(a.data[j].(float64), b.fields().data[i].(float64))
		}
	}
	return a
}

// PowC raises all elements to a constant power.
// Negative powers will result in a math.NaN() values.
func (a *Array64) PowC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = math.Pow(a.data[i].(float64), b)
	}
	return a
}

// FMA12 is the fuse multiply add functionality.
// Array x will contain a[i] = x*a[i]+b[i]
func (a *Array64) FMA12(x float64, b *Array64) *Array64 {
	if a.valRith(b, "FMA") {
		return a
	}

	if b.strides[0] != a.strides[0] {
		cmp, mul := new(sync.WaitGroup), len(a.data)/len(b.data)
		cmp.Add(mul)
		for k := 0; k < mul; k++ {
			go func(m int) {
				Fma12(x, a.data[m:m+len(b.data)], b.data)
				cmp.Done()
			}(k * len(b.data))
		}
		cmp.Wait()
		return a
	}

	Fma12(x, a.data, b.data)
	return a
}

// FMA21 is the fuse multiply add functionality.
// Array x will contain a[i] = a[i]*b[i]+x
func (a *Array64) FMA21(x float64, b *Array64) *Array64 {
	if a.valRith(b, "FMA") {
		return a
	}
	if b.strides[0] != a.strides[0] {
		cmp, mul := new(sync.WaitGroup), len(a.data)/len(b.data)
		cmp.Add(mul)
		for k := 0; k < mul; k++ {
			go func(m int) {
				Fma21(x, a.data[m:m+len(b.data)], b.data)
				cmp.Done()
			}(k * len(b.data))
		}
		cmp.Wait()
		return a
	}

	Fma21(x, a.data, b.data)
	return a
}

// valAr needs to be called before
func (a *Array64) valRith(b nDimObject, mthd string) bool {
	var flag bool
	switch {
	case a.HasErr():
		return true
	case b == nil:
		a.err = NilError
		if debug {
			a.debug = "Array received by " + mthd + "() is a Nil pointer."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	case b.HasErr():
		a.err = b.fields().err
		if debug {
			a.debug = "Array received by " + mthd + "() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	case len(a.shape) < len(b.fields().shape):
		goto shape
	}

	for i, j := len(b.fields().shape)-1, len(a.fields().shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.fields().shape[i] {
			flag = true
			break
		}
	}
	if !flag {
		return false
	}
	if len(b.fields().shape) != len(a.shape) || b.fields().shape[len(b.fields().shape)-1] != 1 {
		goto shape
	}
	for i := 0; i < len(a.shape)-1; i++ {
		if a.shape[i] != b.fields().shape[i] {
			goto shape
		}
	}
	return false
shape:
	a.err = ShapeError
	if debug {
		a.debug = fmt.Sprintf("Array received by %s() can not be broadcast.  Shape: %v  Val shape: %v",
			mthd, a.shape, b.fields().shape)
		a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
	}
	return true
}
