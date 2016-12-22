package numgo

import (
	"fmt"
	"runtime"
)

type nDimMetadata struct {
	shape        []int
	strides      []int
	err          error
	debug, stack string
}

func (ndim *nDimMetadata) hasErr() bool {
	if ndim == nil || (ndim.shape == nil && ndim.err == nil) {
		return true
	}
	return ndim.err != nil
}

// newNDim will create the metadata for an n-dimensional array.
func newNDim(shape []int) nDimMetadata {
	ndim := nDimMetadata{err: nil, debug: "", stack: ""}
	if len(shape) == 1 {
		if shape[0] == 0 {
			ndim.shape = []int{0}
			ndim.strides = []int{0, 0}
		} else {
			ndim.shape = []int{shape[0]}
			ndim.strides = []int{shape[0], 1}
		}
		return ndim
	}

	sh := make([]int, len(shape))
	ndim.strides = make([]int, len(shape)+1)
	ndim.strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		if shape[i] < 0 {
			ndim.err = NegativeAxis
			ndim.strides = ndim.strides[:1]
			if debug {
				ndim.debug = fmt.Sprintf("Negative axis length received by Create: %v", shape)
				ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return ndim
		}
		ndim.strides[i] = ndim.strides[i+1] * shape[i]
	}

	copy(sh, shape)
	ndim.shape = sh
	return ndim
}

// reshape will validate and adjust the array shape/strides and verify the size is unchanged.
// internal:  does not check for previous errors
func (ndim *nDimMetadata) reshape(shape []int) {
	sh := make([]int, len(shape))
	strides := make([]int, len(shape)+1)
	strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		if shape[i] < 0 {
			ndim.err = NegativeAxis
			if debug {
				ndim.debug = fmt.Sprintf("Negative axis length received by Reshape(): %v", shape)
				ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		strides[i] = strides[i+1] * shape[i]
	}
	if strides[0] != ndim.strides[0] {
		ndim.err = ReshapeError
		if debug {
			ndim.debug = fmt.Sprintf(
				"Reshape() can not change data size.  Dimensions: %v reshape: %v",
				ndim.shape, shape)
			ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return
	}
	copy(sh, shape)
	ndim.shape = sh
	ndim.strides = strides
}

// Resize will change the underlying array size.
//
// Make a copy C() if the original array needs to remain unchanged.
// Element location in the underlying slice will not be adjusted to the new shape.
func (ndim *nDimMetadata) resize(shape []int) {
	switch {
	case ndim.hasErr():
		return
	case len(shape) == 0:
		ndim.shape, ndim.strides = []int{0}, []int{0, 0}
		return
	}

	sh, strides := make([]int, len(shape)), make([]int, len(shape)+1)
	strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		if shape[i] < 0 {
			ndim.err = NegativeAxis
			if debug {
				ndim.debug = fmt.Sprintf(
					"Negative axis length received by Resize.  Shape: %v", shape)
				ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		strides[i] = strides[i+1] * shape[i]
	}
	copy(sh, shape)
	ndim.shape = sh
	ndim.strides = strides
}

// valIdx validates and calculates the flat index from multiple index
func (ndim *nDimMetadata) valIdx(index []int, mthd string) (idx int) {
	if ndim.hasErr() {
		return
	}
	if len(index) > len(ndim.shape) {
		ndim.err = InvIndexError
		if debug {
			ndim.debug = fmt.Sprintf(
				"Incorrect number of indicies received by %s().  Shape: %v  Index: %v",
				mthd, ndim.shape, index)
			ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return 0
	}
	for i, v := range index {
		if v >= ndim.shape[i] || v < 0 {
			ndim.err = IndexError
			if debug {
				ndim.debug = fmt.Sprintf(
					"Index received by %s() does not exist shape: %v index: %v",
					mthd, ndim.shape, index)
				ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return 0
		}
		idx += v * ndim.strides[i+1]
	}
	return idx
}
