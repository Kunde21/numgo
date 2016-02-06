package numgo

import (
	"fmt"
	"runtime"
)

type ngError struct {
	s string
}

func (n *ngError) Error() string {
	return n.s
}

var (
	// NilError flags any error where a nil pointer is received
	NilError = &ngError{"NilError: Nil pointer recieved."}
	// ShapeError flags any mismatching
	ShapeError = &ngError{"ShapeError: Array shapes don't match and can't be broadcast."}
	// ReshapeError flags incorrect use of Reshape() calls.
	// Resize() is the only call that can change the capacity of arrays.
	ReshapeError = &ngError{"ReshapeError: New shape cannot change the size of the array."}
	// NegativeAxis is forNew/Reshape/Resize calls with negative axis length: not allowed
	NegativeAxis = &ngError{"NegativeAxis: Negative axis length received."}
	// IndexError flags any attempt to index out of range
	IndexError = &ngError{"IndexError: Index or Axis out of range."}
	// DivZeroError flags division by zero
	DivZeroError = &ngError{"DivZeroError: Division by zero encountered."}
	// InvIndexError flags Negative or illegal indexes
	InvIndexError = &ngError{"InvIndexError: Invalid or illegal index received."}
	// FoldMapError catches panics within Fold/FoldCC/Map calls.
	// This will store the panic message in the debug string
	// when debugging is turned off, for proper errror reporting
	FoldMapError = &ngError{"FoldMapError: Fold/Map function panic encountered."}

	debug    bool
	stackBuf []byte
)

// Debug sets the error reporting level for the library.
// To get debugging data from the library, set this to true
// and use GetDebug() in place of GetErr().
//
// Debugging information includes the function call that generated
// the error, stack trace at the point the error was generated,
// and the values involved in that function call.
// This will add overhead to error reporting and handling, so
// use it for development and debugging purposes.
func Debug(set bool) bool {
	debug = set
	if debug && len(stackBuf) == 0 {
		stackBuf = make([]byte, 4096)
	}
	return debug
}

// HasErr tests for the existence of an error on the Array64 object.
//
// Errors will be maintained through a chain of function calls,
// so only the first error will be returned when GetErr() is called.
// Use HasErr() as a gate for the GetErr() or GetDebug() choice in
// error handling code.
func (a *Array64) HasErr() bool {
	if a == nil || (a.data == nil && a.err == nil) {
		return true
	}
	return a.err != nil
}

// GetErr returns the error object and clears the error from the array.
//
// This will only return an error value once per error instance.  Do not use
// it in the if statement to test for the existence of an error.  HasErr() is
// provided for that purpose.
func (a *Array64) GetErr() (err error) {
	if a == nil || (a.data == nil && a.err == nil) {
		return NilError
	}
	err = a.err
	a.err, a.debug, a.stack = nil, "", ""
	return
}

// GetDebug returns and clears the error object from the array object.  The returned debug string
// will include the function that generated the error and the arguments that caused it.
//
// This debug information will only be generated and returned if numgo.Debug is set to true
// before the function call that causes the error.
func (a *Array64) GetDebug() (err error, debugStr, stackTrace string) {
	if a == nil || (a.data == nil && a.err == nil) {
		err = NilError
		if debug {
			debugStr = "Nil pointer received by GetDebug().  Source array was not initialized."
			stackTrace = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return
	}
	err, debugStr, stackTrace = a.err, a.debug, a.stack
	a.err, a.debug, a.stack = nil, "", ""
	return
}

// encodeErr is a supporting function for MarshalJSON
func encodeErr(err error) int8 {
	e := err.(*ngError)
	switch e {
	case nil:
		return 0
	case NilError:
		return 1
	case ShapeError:
		return 2
	case ReshapeError:
		return 3
	case NegativeAxis:
		return 4
	case IndexError:
		return 5
	case DivZeroError:
		return 6
	case InvIndexError:
		return 7
	case FoldMapError:
		return 8
	}
	return -1
}

// decodeErr is a supporting method for UnmarshalJSON
func decodeErr(err int8) (a *ngError) {
	switch err {
	case 0:
		a = nil
	case 1:
		a = NilError
	case 2:
		a = ShapeError
	case 3:
		a = ReshapeError
	case 4:
		a = NegativeAxis
	case 5:
		a = IndexError
	case 6:
		a = DivZeroError
	case 7:
		a = InvIndexError
	case 8:
		a = FoldMapError
	default:
		a = &ngError{fmt.Sprintf("Unknown error Unmarshaled: %d", err)}
	}
	return
}

// HasErr tests for the existence of an error on the Arrayb object.
//
// Errors will be maintained through a chain of function calls,
// so only the first error will be returned when GetErr() is called.
// Use HasErr() as a gate for the GetErr() or GetDebug() choice in
// error handling code.
func (a *Arrayb) HasErr() bool {
	if a == nil || (a.data == nil && a.err == nil) {
		return true
	}
	return a.err != nil
}

// GetErr returns the error object and clears the error from the array.
//
// This will only return an error value once per error instance.  Do not use
// it in the if statement to test for the existence of an error.  HasErr() is
// provided for that purpose.
func (a *Arrayb) GetErr() (err error) {
	if a == nil || (a.data == nil && a.err == nil) {
		return NilError
	}
	err = a.err
	a.err, a.debug = nil, ""
	return
}

// GetDebug returns and clears the error object from the array object.  The returned debug string
// will include the function that generated the error and the arguments that caused it.
//
// This debug information will only be generated and returned if numgo.Debug is set to true
// before the function call that causes the error.
func (a *Arrayb) GetDebug() (err error, debugStr, stackTrace string) {
	if a == nil || (a.data == nil && a.err == nil) {
		err = NilError
		if debug {
			debugStr = "Nil pointer received in GetDebug().  Source array was not initialized."
			stackTrace = string(stackBuf[:runtime.Stack(stackBuf, false)])
		}
		return
	}
	err, debugStr, stackTrace = a.err, a.debug, a.stack
	a.err, a.debug, a.stack = nil, "", ""

	return
}
