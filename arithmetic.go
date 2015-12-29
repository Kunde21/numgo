package numgo

import (
	"fmt"
	"math"
)

// Add performs element-wise addition
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Arrayf) Add(b *Arrayf) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	case b == nil:
		a.err = NilError
		fallthrough
	case b.err != nil:
		return a
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			a.err = ShapeError
			return a
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
	close(compChan)
	return a
}

// AddC adds a constant to all elements of the array.
func (a *Arrayf) AddC(b float64) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] += b
	}
	return a
}

// Subtr performs element-wise subtraction.
// Arrays must be the same size or albe to broadcast.
// This will modify the source array.
func (a *Arrayf) Subtr(b *Arrayf) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	case b == nil:
		a.err = NilError
		fallthrough
	case b.err != nil:
		return a
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			a.err = ShapeError
			return a
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
	close(compChan)
	return a
}

// SubtrC subtracts a constant from all elements of the array.
func (a *Arrayf) SubtrC(b float64) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] -= b
	}
	return a
}

// Mult performs element-wise multiplication.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Arrayf) Mult(b *Arrayf) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	case b == nil:
		a.err = NilError
		fallthrough
	case b.err != nil:
		return a
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i, j := len(b.shape)-1, len(a.shape)-1; i >= 0; i, j = i-1, j-1 {
		if a.shape[j] != b.shape[i] {
			a.err = ShapeError
			return a
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
	close(compChan)
	return a
}

// MultC multiplies all elements of the array by a constant.
func (a *Arrayf) MultC(b float64) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] *= b
	}
	return a
}

// Mult performs element-wise division
// Arrays must be the same size or able to broadcast.
// Division by zero will result in a math.NaN() values.
// This will modify the source array.
func (a *Arrayf) Div(b *Arrayf) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	case b == nil:
		a.err = NilError
		fallthrough
	case b.err != nil:
		return a
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
		return a
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
	close(compChan)
	if flg {
		a.err = DivZeroError
	}
	return a
}

// MultC divides all elements of the array by a constant.
// Division by zero will result in a math.NaN() values.
func (a *Arrayf) DivC(b float64) *Arrayf {
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

// Pow raises elements of a to the corresponding power in b.
// Arrays must be the same size or able to broadcast.
// This will modify the source array.
func (a *Arrayf) Pow(b *Arrayf) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		return a
	case a.err != nil || b.err != nil:
		return a
	case b == nil:
		a.err = NilError
		return a
	case len(a.shape) < len(b.shape):
		a.err = ShapeError
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

// MultC divides all elements of the array by a constant.
// Division by zero will result in a math.NaN() values.
func (a *Arrayf) PowC(b float64) *Arrayf {
	switch {
	case a == nil:
		a = new(Arrayf)
		a.err = NilError
		fallthrough
	case a.err != nil:
		return a
	}

	a.Lock()
	defer a.Unlock()

	for i := 0; i < len(a.data); i++ {
		a.data[i] = math.Pow(a.data[i], b)
	}
	return a
}
