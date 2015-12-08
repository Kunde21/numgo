package numgo

import "sort"

type MapFunc func([]float64) float64

// cleanAxis removes any duplicate axes and returns the cleaned slice.
// only the first instance of an axis is retained.
func cleanAxis(axis ...int) []int {
	if len(axis) < 2 {
		return axis
	}
	length := len(axis) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if axis[i] == axis[j] {
				if j == length {
					axis = axis[:j]
				} else {
					axis = append(axis[:j], axis[j+1:]...)
				}
				length--
				j--
			}
		}
	}
	return axis
}

// collapse will reorganize data by putting element dataset in continuous sections of data slice.
// Returned Arrayf must be condensed with a summary calculation to create a valid array object.
func (a *Arrayf) collapse(axis ...int) (uint64, *Arrayf) {
	if a == nil {
		return 0, nil
	}

	if len(axis) == 0 {
		return 1, a.C()
	}
	axis = cleanAxis(axis...)

	span := uint64(1) // Span = size of "element" Mx = slicing
	mx := a.strides[len(a.strides)-1]
	steps := make([]uint64, len(axis)) // Element strides
	brks := make([]uint64, len(axis))  // Stride-ending breaks

	for i, v := range axis {
		span *= uint64(a.shape[v])
		steps[i], brks[i] = a.strides[v+1], a.strides[v]
		if brks[i] > mx {
			mx = brks[i]
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(axis)))
	asteps := make([]uint64, len(axis)) // Element strides
	abrks := make([]uint64, len(axis))  // Stride-ending breaks
	for i, v := range axis {
		asteps[i], abrks[i] = a.strides[v+1], a.strides[v]
	}

	// Axis-level offsets
	sort.IntSlice(axis).Sort()
	offsets := make([]uint64, len(a.strides))
	copy(offsets, a.strides)
	for _, v := range steps {
		for i, w := range offsets {
			if v == w && i == len(offsets)-1 {
				offsets = offsets[:i]
			} else if v == w {
				offsets = append(offsets[:i], offsets[i+1:]...)
			}
		}
	}

	// Reverse the offsets to make increment code cleaner
	for i, j := 0, len(offsets)-1; i < j; i, j = i+1, j-1 {
		offsets[i], offsets[j] = offsets[j], offsets[i]
	}

	for len(asteps) > 0 && offsets[0] > asteps[0] {
		asteps = asteps[1:]
	}

	tmp := make([]float64, a.strides[0]) // Holds re-arranged data for return
	inc := make([]uint64, len(axis))     // N-dimensional incrementor
	off := make([]uint64, len(axis))     // N-dimensional offset incrementor

	for sl := uint64(0); sl+mx <= a.strides[0]; sl += mx {

		// Inner loop might be made concurrent using slices
		// Unknown performance gains in doing so, tuning needed
		offset := uint64(0)
		for sp := uint64(0); sp+span <= mx; sp += span {

			//fmt.Println(off, offsets, asteps, steps, offset, span)

			for i, k := uint64(0), uint64(0); i < span; i++ {
				tmp[sl+i+sp] = a.data[sl+k+offset]
				k, inc[0] = k+steps[0], inc[0]+steps[0]

				// Incrementor loop to handle all dims
				for c, v := range brks {
					if uint64(i+1) == span {
						// Reset at end of loop
						inc[c] = 0
					}
					if inc[c] >= v {
						k = k - v + steps[c+1]
						inc[c] -= v
						inc[c+1] += steps[c+1]
					}
				}
			}

			//fmt.Println(tmp[sp : sp+span])

			// Increment to the next dimension
			offset, off[0] = offset+offsets[0], off[0]+1
			//fmt.Println(sp+span, mx)
			for c, v := range asteps {
				if sp+span == mx {
					// Reset at end of loop
					off[c] = 0
				}
				if off[c] >= v && c+1 < len(off) {
					offset = offset - v + offsets[c+1]
					off[c] -= v
					off[c+1] += offsets[c+1]
				}
			}
		}
	}

	// Create return object.  Data is invalid format until reform is called.
	b := new(Arrayf)
	b.shape = make([]uint64, len(a.shape)-len(axis))
	b.strides = make([]uint64, len(b.shape)+1)
	b.data = tmp

	for i, t := 0, 0; i < len(a.shape); i++ {
		tmp := false
		for _, w := range axis {
			if i == w {
				tmp = true
				break
			}
		}
		if !tmp {
			b.shape[t] = a.shape[i]
			t++
		}
	}

	t := uint64(1)
	for i := len(b.strides) - 1; i > 0; i-- {
		b.strides[i] = t
		t *= b.shape[i-1]
	}
	b.strides[0] = t

	return span, b
}

// MapCC applies function f along the given axes concurrently. Each call to f will launch a goroutine.
// In order to leverage this concurrency, MapCC should only be used for complex and CPU-heavy functions.
//
// Simple functions should use Map(f, axes...), as it's more performant.
func (a *Arrayf) MapCC(f MapFunc, axis ...int) (ret *Arrayf) {

	type rt struct {
		index uint64
		value float64
	}

	span, ret := a.collapse(axis...)

	retChan := make(chan rt)
	i := uint64(0)
	for i = uint64(0); i+span <= a.strides[0]; i += span {
		go func(i uint64) {
			retChan <- rt{i / span, f(ret.data[i : i+span])}
		}(i)
	}
	for i = uint64(0); i+span <= a.strides[0]; i += span {
		c := <-retChan
		ret.data[c.index] = c.value
	}

	close(retChan)
	ret.data = ret.data[:a.strides[0]]
	return ret
}

// Map applies function f along the given axes.
// Slice containing all data to be consolidated into an element will be passed to f.
// Return value will be the resulting element's value.
func (a *Arrayf) Map(f MapFunc, axis ...int) (ret *Arrayf) {
	span, ret := a.collapse(axis...)
	for i := uint64(0); i+span <= a.strides[0]; i += span {
		ret.data[i/span] = f(ret.data[i : i+span])
	}
	ret.data = ret.data[:a.strides[0]]
	return ret
}
