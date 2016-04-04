package numgo

//Matrix and vector multiplication implementations
// - Dot Product
// - Cross Product
// - Tensor Product

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
