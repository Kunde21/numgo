package numgo

//Matrix and vector multiplication implementations
// - Dot Product
// - Cross Product
// - Tensor Product

// DotProd calculates the dot (scalar) product of two vectors.
// NOTE: Only implemented on 1-D arrays, and other sizes are NOOP
func (a *Array64) DotProd(b *Array64) *Array64 {
	switch {
	case a.valRith(b, "DotProd"):
		return a
	case len(a.shape) == 1:
		return &Array64{
			shape:   []uint64{1},
			strides: []uint64{1, 1},
			data:    []float64{dotProd(a.data, b.data)},
			err:     nil,
			debug:   "",
			stack:   "",
		}
	}
	return a
}

func (a *Array64) MatProd(b *Array64) *Array64 {
    switch  {
    case a.valRith(b,"MatProd"):
            return a
    case len(a.shape) == 1 && len(b.shape) == 1:
		return &Array64{
			shape:   []uint64{1},
			strides: []uint64{1, 1},
			data:    []float64{dotProd(a.data, b.data)},
			err:     nil,
			debug:   "",
			stack:   "",
		}
    }
    return a
}