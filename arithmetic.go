package numgo

import (
	"fmt"
	"math"
	"runtime"
)

// Add performs element-wise addition
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Add(b *Array64) *Array64 {
	if a.valRith(b, "Add") {
		return a
	}
	/*	compChan := make(chan struct{})
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
		close(compChan)
	*/

	add(a.data, b.data)

	return a
}

// AddC adds a constant to all elements of the array.
func (a *Array64) AddC(b float64) *Array64 {
	if a == nil || a.err != nil {
		return a
	}

	addC(b, a.data)
	return a
}

// Subtr performs element-wise subtraction.
// Arrays must be the same size or albe to broadcast.
// This will modify the source array.
func (a *Array64) Subtr(b *Array64) *Array64 {
	if a.valRith(b, "Subtr") {
		return a
	}

	subtr(a.data, b.data)
	return a
}

// SubtrC subtracts a constant from all elements of the array.
func (a *Array64) SubtrC(b float64) *Array64 {
	if a == nil || a.err != nil {
		return a
	}

	subtrC(b, a.data)
	return a
}

// Mult performs element-wise multiplication.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Mult(b *Array64) *Array64 {
	if a.valRith(b, "Mult") {
		return a
	}

	mult(a.data, b.data)
	return a
}

// MultC multiplies all elements of the array by a constant.
func (a *Array64) MultC(b float64) *Array64 {
	if a == nil || a.err != nil {
		return a
	}

	multC(b, a.data)
	return a
}

// Div performs element-wise division
// Arrays must be the same size or able to broadcast.
// Division by zero will result in a math.NaN() values.
// This will modify the source array.
func (a *Array64) Div(b *Array64) *Array64 {
	if a.valRith(b, "Div") {
		return a
	}
	div(a.data, b.data)
	return a
}

// DivC divides all elements of the array by a constant.
// Division by zero will result in a math.NaN() values.
func (a *Array64) DivC(b float64) *Array64 {
	switch {
	case a == nil || a.err != nil:
		return a
	case b == 0:
		a.err = DivZeroError
		if debug {
			a.debug = "Division by zero encountered in DivC()"
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
	}

	for i := 0; i < len(a.data); i++ {
		if b == 0 {
			a.data[i] = math.NaN()
		} else {
			a.data[i] /= b
		}
	}
	return a
}

// Pow raises elements of a to the corresponding power in b.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Array64) Pow(b *Array64) *Array64 {
	if a.valRith(b, "Pow") {
		return a
	}

	compChan := make(chan struct{})
	mul := len(a.data) / len(b.data)
	for k := 0; k < mul; k++ {
		go func(m int) {
			for i, v := range b.data {
				a.data[i+m] = math.Pow(a.data[i+m], v)
			}
			compChan <- struct{}{}
		}(k * len(b.data))
	}

	for k := 0; k < mul; k++ {
		<-compChan
	}
	return a
}

// PowC raises all elements to a constant power.
// Negative powers will result in a math.NaN() values.
func (a *Array64) PowC(b float64) *Array64 {
	if a == nil || a.err != nil {
		return a
	}

	for i := 0; i < len(a.data); i++ {
		a.data[i] = math.Pow(a.data[i], b)
	}
	return a
}

// valAr needs to be called before
func (a *Array64) valRith(b *Array64, mthd string) bool {
	switch {
	case a == nil || a.err != nil:
		return true
	case b == nil:
		a.err = NilError
		if debug {
			a.debug = "Array received by " + mthd + "() is a Nil pointer."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	case b.err != nil:
		a.err = b.err
		if debug {
			a.debug = "Array received by " + mthd + "() is in error."
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
		if debug {
			a.debug = fmt.Sprintf("Array received by %s() can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
			a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return true
	}

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			a.err = ShapeError
			if debug {
				a.debug = fmt.Sprintf("Array received by %s() can not be broadcast.  Shape: %v  Val shape: %v", mthd, a.shape, b.shape)
				a.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return true
		}
	}
	return false
}
