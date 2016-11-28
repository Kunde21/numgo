package numgo

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/Kunde21/numgo/internal"
)

var nan float64

func init() {
	nan = math.NaN()
}

// Add performs element-wise addition
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Add(b *Array64) *Array64 {
	if a.valRith(b, "Add") {
		return a
	}

	if b.shape[len(b.shape)-1] == a.shape[len(a.shape)-1] {
		asm.Add(a.data, b.data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.data); i++ {
		asm.AddC(b.data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// AddC adds a constant to all elements of the array.
func (a *Array64) AddC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	asm.AddC(b, a.data)
	return a
}

// Subtr performs element-wise subtraction.
// Arrays must be the same size or albe to broadcast.
// This will modify the source array.
func (a *Array64) Subtr(b *Array64) *Array64 {
	if a.valRith(b, "Subtr") {
		return a
	}

	if b.shape[len(b.shape)-1] == a.shape[len(a.shape)-1] {
		asm.Subtr(a.data, b.data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.data); i++ {
		asm.SubtrC(b.data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// SubtrC subtracts a constant from all elements of the array.
func (a *Array64) SubtrC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	asm.SubtrC(b, a.data)
	return a
}

// Mult performs element-wise multiplication.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Mult(b *Array64) *Array64 {
	if a.valRith(b, "Mult") {
		return a
	}

	if b.shape[len(b.shape)-1] == a.shape[len(a.shape)-1] {
		asm.Mult(a.data, b.data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.data); i++ {
		asm.MultC(b.data[i], a.data[i*st:(i+1)*st])
	}
	return a
}

// MultC multiplies all elements of the array by a constant.
func (a *Array64) MultC(b float64) *Array64 {
	if a.HasErr() {
		return a
	}

	asm.MultC(b, a.data)
	return a
}

// Div performs element-wise division
// Arrays must be the same size or able to broadcast.
// Division by zero conforms to IEEE 754
// 0/0 = NaN, +x/0 = +Inf, -x/0 = -Inf
// This will modify the source array.
func (a *Array64) Div(b *Array64) *Array64 {
	if a.valRith(b, "Div") {
		return a
	}

	if b.shape[len(b.shape)-1] == a.shape[len(a.shape)-1] {
		asm.Div(a.data, b.data)
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.data); i++ {
		asm.DivC(b.data[i], a.data[i*st:(i+1)*st])
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

	asm.DivC(b, a.data)
	return a
}

// Pow raises elements of a to the corresponding power in b.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Pow(b *Array64) *Array64 {
	if a.valRith(b, "Pow") {
		return a
	}

	if b.shape[len(b.shape)-1] == a.shape[len(a.shape)-1] {
		lna, lnb := len(a.data), len(b.data)
		for i, j := 0, 0; i < lna; i, j = i+1, j+1 {
			if j >= lnb {
				j = 0
			}
			a.data[i] = math.Pow(a.data[i], b.data[j])
		}
		return a
	}

	st := a.strides[len(a.strides)-1] * a.shape[len(a.shape)-1]
	for i := 0; i < len(b.data); i++ {
		for j := i * st; j < (i+1)*st; j++ {
			a.data[j] = math.Pow(a.data[j], b.data[i])
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
		a.data[i] = math.Pow(a.data[i], b)
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
				asm.Fma12(x, a.data[m:m+len(b.data)], b.data)
				cmp.Done()
			}(k * len(b.data))
		}
		cmp.Wait()
		return a
	}

	asm.Fma12(x, a.data, b.data)
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
				asm.Fma21(x, a.data[m:m+len(b.data)], b.data)
				cmp.Done()
			}(k * len(b.data))
		}
		cmp.Wait()
		return a
	}

	asm.Fma21(x, a.data, b.data)
	return a
}

// valAr needs to be called before
func (a *Array64) valRith(b *Array64, mthd string) bool {
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
		a.err = b.err
		if debug {
			a.debug = "Array received by " + mthd + "() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	case len(a.shape) < len(b.shape):
		goto shape
	}

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			flag = true
			break
		}
	}
	if !flag {
		return false
	}
	if len(b.shape) != len(a.shape) || b.shape[len(b.shape)-1] != 1 {
		goto shape
	}
	for i := 0; i < len(a.shape)-1; i++ {
		if a.shape[i] != b.shape[i] {
			goto shape
		}
	}
	return false
shape:
	a.err = ShapeError
	if debug {
		a.debug = fmt.Sprintf("Array received by %s() can not be broadcast.  Shape: %v  Val shape: %v",
			mthd, a.shape, b.shape)
		a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
	}
	return true
}
