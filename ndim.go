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
func (ndim *nDimMetadata) reshape(shape []int) {
	sh := make([]int, len(shape))
	strides := make([]int, len(shape)+1)
	strides[len(shape)] = 1
	for i := len(shape) - 1; i >= 0; i-- {
		if shape[i] < 0 {
			ndim.err = NegativeAxis
			if debug {
				ndim.debug = fmt.Sprintf("Negative axis length received by Create: %v", shape)
				ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
			}
			return
		}
		strides[i] = strides[i+1] * shape[i]
	}
	if strides[0] != ndim.strides[0] {
		ndim.err = ReshapeError
		if debug {
			ndim.debug = fmt.Sprintf("Reshape() can not change data size.  Dimensions: %v reshape: %v", ndim.shape, shape)
			ndim.stack = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return
	}
	copy(sh, shape)
	ndim.shape = sh
	ndim.strides = strides
}
